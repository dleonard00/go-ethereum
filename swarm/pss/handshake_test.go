package pss

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// asymmetrical key exchange between two directly connected peers
// full address, partial address (8 bytes) and empty address
func TestHandshake(t *testing.T) {
	t.Run("32", testHandshake)
	t.Run("8", testHandshake)
	t.Run("0", testHandshake)
}

func testHandshake(t *testing.T) {

	// how much of the address we will use
	useHandshake = true
	var addrsize int64
	var err error
	addrsizestring := strings.Split(t.Name(), "/")
	addrsize, _ = strconv.ParseInt(addrsizestring[1], 10, 0)

	// set up two nodes directly connected
	// (we are not testing pss routing here)
	topic := BytesToTopic([]byte("foo:42"))

	clients, err := setupNetwork(2)
	if err != nil {
		t.Fatal(err)
	}

	var loaddr []byte
	err = clients[0].Call(&loaddr, "pss_baseAddr")
	if err != nil {
		t.Fatalf("rpc get node 1 baseaddr fail: %v", err)
	}
	loaddr = loaddr[:addrsize]
	var roaddr []byte
	err = clients[1].Call(&roaddr, "pss_baseAddr")
	if err != nil {
		t.Fatalf("rpc get node 2 baseaddr fail: %v", err)
	}
	roaddr = roaddr[:addrsize]
	log.Debug("addresses", "left", loaddr, "right", roaddr)

	// retrieve public key from pss instance
	// set this public key reciprocally
	lpubkey := make([]byte, 32)
	err = clients[0].Call(&lpubkey, "pss_getPublicKey")
	if err != nil {
		t.Fatalf("rpc get node 1 pubkey fail: %v", err)
	}
	rpubkey := make([]byte, 32)
	err = clients[1].Call(&rpubkey, "pss_getPublicKey")
	if err != nil {
		t.Fatalf("rpc get node 2 pubkey fail: %v", err)
	}

	time.Sleep(time.Millisecond * 1000) // replace with hive healthy code

	// give each node its peer's public key
	err = clients[0].Call(nil, "pss_setPeerPublicKey", rpubkey, topic, roaddr)
	if err != nil {
		t.Fatal(err)
	}
	err = clients[1].Call(nil, "pss_setPeerPublicKey", lpubkey, topic, loaddr)
	if err != nil {
		t.Fatal(err)
	}

	// perform the handshake
	// after this each side will have defaultSymKeyBufferCapacity symkeys each for in- and outgoing messages:
	// L -> request 4 keys -> R
	// L <- send 4 keys, request 4 keys <- R
	// L -> send 4 keys -> R
	// the call will fill the array with symkeys L needs for sending to R
	err = clients[0].Call(nil, "pss_addHandshake", topic)
	if err != nil {
		t.Fatal(err)
	}
	err = clients[1].Call(nil, "pss_addHandshake", topic)
	if err != nil {
		t.Fatal(err)
	}

	var lhsendsymkeyids []string
	err = clients[0].Call(&lhsendsymkeyids, "pss_handshake", common.ToHex(rpubkey), topic, true, true)
	if err != nil {
		t.Fatal(err)
	}

	// make sure the r-node gets its keys
	time.Sleep(time.Second)

	// check if we have 6 outgoing keys stored, and they match what was received from R
	var lsendsymkeyids []string
	err = clients[0].Call(&lsendsymkeyids, "pss_getHandshakeKeys", common.ToHex(rpubkey), topic, false, true)
	if err != nil {
		t.Fatal(err)
	}
	m := 0
	for _, hid := range lhsendsymkeyids {
		for _, lid := range lsendsymkeyids {
			if lid == hid {
				m++
			}
		}
	}
	if m != defaultSymKeyCapacity {
		t.Fatalf("buffer size mismatch, expected %d, have %d: %v", defaultSymKeyCapacity, m, lsendsymkeyids)
	}

	// check if in- and outgoing keys on l-node and r-node match up and are in opposite categories (l recv = r send, l send = r recv)
	var rsendsymkeyids []string
	err = clients[1].Call(&rsendsymkeyids, "pss_getHandshakeKeys", common.ToHex(lpubkey), topic, false, true)
	if err != nil {
		t.Fatal(err)
	}
	var lrecvsymkeyids []string
	err = clients[0].Call(&lrecvsymkeyids, "pss_getHandshakeKeys", common.ToHex(rpubkey), topic, true, false)
	if err != nil {
		t.Fatal(err)
	}
	var rrecvsymkeyids []string
	err = clients[1].Call(&rrecvsymkeyids, "pss_getHandshakeKeys", common.ToHex(lpubkey), topic, true, false)
	if err != nil {
		t.Fatal(err)
	}

	// get outgoing symkeys in byte form from both sides
	var lsendsymkeys [][]byte
	for _, id := range lsendsymkeyids {
		var key []byte
		err = clients[0].Call(&key, "pss_getSymmetricKey", id)
		if err != nil {
			t.Fatal(err)
		}
		lsendsymkeys = append(lsendsymkeys, key)
	}
	var rsendsymkeys [][]byte
	for _, id := range rsendsymkeyids {
		var key []byte
		err = clients[1].Call(&key, "pss_getSymmetricKey", id)
		if err != nil {
			t.Fatal(err)
		}
		rsendsymkeys = append(rsendsymkeys, key)
	}

	// get incoming symkeys in byte form from both sides and compare
	var lrecvsymkeys [][]byte
	for _, id := range lrecvsymkeyids {
		var key []byte
		err = clients[0].Call(&key, "pss_getSymmetricKey", id)
		if err != nil {
			t.Fatal(err)
		}
		match := false
		for _, otherkey := range rsendsymkeys {
			if bytes.Equal(otherkey, key) {
				match = true
			}
		}
		if !match {
			t.Fatalf("no match right send for left recv key %s", id)
		}
		lrecvsymkeys = append(lrecvsymkeys, key)
	}
	var rrecvsymkeys [][]byte
	for _, id := range rrecvsymkeyids {
		var key []byte
		err = clients[1].Call(&key, "pss_getSymmetricKey", id)
		if err != nil {
			t.Fatal(err)
		}
		match := false
		for _, otherkey := range lsendsymkeys {
			if bytes.Equal(otherkey, key) {
				match = true
			}
		}
		if !match {
			t.Fatalf("no match left send for right recv key %s", id)
		}
		rrecvsymkeys = append(rrecvsymkeys, key)
	}

	// send new handshake request, should send no keys
	err = clients[0].Call(nil, "pss_handshake", common.ToHex(rpubkey), topic, false)
	if err == nil {
		t.Fatal("expected full symkey buffer error")
	}

	// expire one key, send new handshake request
	err = clients[0].Call(nil, "pss_releaseHandshakeKey", common.ToHex(rpubkey), topic, lsendsymkeyids[0], true)
	if err != nil {
		t.Fatalf("release left send key %s fail: %v", lsendsymkeyids[0], err)
	}

	var newlhsendkeyids []string

	// send new handshake request, should now receive one key
	// check that it is not in previous right recv key array
	err = clients[0].Call(&newlhsendkeyids, "pss_handshake", common.ToHex(rpubkey), topic, true, false)
	if err != nil {
		t.Fatalf("handshake send fail: %v", err)
	} else if len(newlhsendkeyids) != defaultSymKeyCapacity {
		t.Fatalf("wrong receive count, expected 1, got %d", len(newlhsendkeyids))
	}

	var newlrecvsymkey []byte
	err = clients[0].Call(&newlrecvsymkey, "pss_getSymmetricKey", newlhsendkeyids[0])
	if err != nil {
		t.Fatal(err)
	}
	var rmatchsymkeyid *string
	for i, id := range rrecvsymkeyids {
		var key []byte
		err = clients[1].Call(&key, "pss_getSymmetricKey", id)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(newlrecvsymkey, key) {
			rmatchsymkeyid = &rrecvsymkeyids[i]
		}
	}
	if rmatchsymkeyid != nil {
		t.Fatalf("right sent old key id %s in second handshake", *rmatchsymkeyid)
	}

	// clean the pss core keystore. Should clean the key released earlier
	var cleancount int
	clients[0].Call(&cleancount, "psstest_clean")
	if cleancount > 1 {
		t.Fatalf("pss clean count mismatch; expected 1, got %d", cleancount)
	}
}