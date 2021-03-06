package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/pborman/uuid"

	cli "gopkg.in/urfave/cli.v1"
)

func cliUploadAndSync(c *cli.Context) error {
	defer func(now time.Time) { log.Info("total time", "time", time.Since(now), "size", filesize) }(time.Now())

	log.Info("uploading and syncing")

	f, cleanup := generateRandomFile(filesize * 1000000)
	defer cleanup()

	hash, err := upload(f, endpoints[0])
	if err != nil {
		return err
	}

	fhash, err := digest(f)
	if err != nil {
		return err
	}

	log.Info("uploaded successfully", "hash", hash, "digest", fmt.Sprintf("%x", fhash))

	wg := sync.WaitGroup{}
	for _, endpoint := range endpoints {
		endpoint := endpoint
		ruid := uuid.New()[:8]
		wg.Add(1)
		go func(endpoint string, ruid string) {
			for {
				err := fetch(hash, endpoint, fhash, ruid)
				if err != nil {
					continue
				}

				wg.Done()
				return
			}
		}(endpoint, ruid)
	}
	wg.Wait()
	log.Info("all endpoints synced random file successfully")

	return nil
}

// fetch is getting the requested `hash` from the `endpoint` and compares it with the `original` file
func fetch(hash string, endpoint string, original []byte, ruid string) error {
	log.Trace("sleeping", "ruid", ruid)
	time.Sleep(10 * time.Second)

	log.Trace("http get request", "ruid", ruid, "api", endpoint, "hash", hash)
	res, err := http.Get(endpoint + "/bzz:/" + hash)
	if err != nil {
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}
	log.Trace("http get response", "ruid", ruid, "api", endpoint, "hash", hash, "code", res.StatusCode, "len", res.ContentLength)

	if res.StatusCode != 200 {
		err := fmt.Errorf("expected status code %d, got %v", 200, res.StatusCode)
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}

	defer res.Body.Close()

	rdigest, err := digest(res.Body)
	if err != nil {
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}

	if !bytes.Equal(rdigest, original) {
		err := fmt.Errorf("downloaded imported file md5=%x is not the same as the generated one=%x", rdigest, original)
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}

	log.Trace("downloaded file matches random file", "ruid", ruid, "len", res.ContentLength)

	return nil
}

// upload is uploading a file `f` to `endpoint` via the `swarm up` cmd
func upload(f *os.File, endpoint string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("swarm", "--bzzapi", endpoint, "up", f.Name())
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	hash := strings.TrimRight(out.String(), "\r\n")
	return hash, nil
}

func digest(r io.Reader) ([]byte, error) {
	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// generateRandomFile is creating a temporary file with the requested byte size
func generateRandomFile(size int) (f *os.File, teardown func()) {
	// create a tmp file
	tmp, err := ioutil.TempFile("", "swarm-test")
	if err != nil {
		panic(err)
	}

	// callback for tmp file cleanup
	teardown = func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

	buf := make([]byte, size)
	_, err = rand.Read(buf)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(tmp.Name(), buf, 0755)

	return tmp, teardown
}
