// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package stream

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm/network"
	streamTesting "github.com/ethereum/go-ethereum/swarm/network/stream/testing"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

var (
	deliveries map[discover.NodeID]*Delivery
	stores     map[discover.NodeID]storage.ChunkStore
	toAddr     func(discover.NodeID) *network.BzzAddr
	peerCount  func(discover.NodeID) int
)

func TestStreamerRetrieveRequest(t *testing.T) {
	tester, streamer, _, teardown, err := newStreamerTester(t)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	peerID := tester.IDs[0]

	streamer.delivery.RequestFromPeers(hash0[:], true)

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "RetrieveRequestMsg",
		Expects: []p2ptest.Expect{
			p2ptest.Expect{
				Code: 5,
				Msg: &RetrieveRequestMsg{
					Key:       hash0[:],
					SkipCheck: true,
				},
				Peer: peerID,
			},
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestStreamerUpstreamRetrieveRequestMsgExchangeWithoutStore(t *testing.T) {
	tester, streamer, _, teardown, err := newStreamerTester(t)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	peerID := tester.IDs[0]

	chunk := storage.NewChunk(storage.Key(hash0[:]), nil)

	peer := streamer.getPeer(peerID)

	peer.handleSubscribeMsg(&SubscribeMsg{
		Stream:   swarmChunkServerStreamName,
		Key:      nil,
		From:     0,
		To:       0,
		Priority: Top,
	})

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "RetrieveRequestMsg",
		Triggers: []p2ptest.Trigger{
			p2ptest.Trigger{
				Code: 5,
				Msg: &RetrieveRequestMsg{
					Key: chunk.Key[:],
				},
				Peer: peerID,
			},
		},
		Expects: []p2ptest.Expect{
			p2ptest.Expect{
				Code: 1,
				Msg: &OfferedHashesMsg{
					HandoverProof: nil,
					Hashes:        nil,
					From:          0,
					To:            0,
				},
				Peer: peerID,
			},
		},
	})

	expectedError := "exchange 0: 'RetrieveRequestMsg' timed out"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Expected error %v, got %v", expectedError, err)
	}
}

// upstream request server receives a retrieve Request and responds with
// offered hashes or delivery if skipHash is set to true
func TestStreamerUpstreamRetrieveRequestMsgExchange(t *testing.T) {
	tester, streamer, localStore, teardown, err := newStreamerTester(t)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	peerID := tester.IDs[0]
	peer := streamer.getPeer(peerID)

	peer.handleSubscribeMsg(&SubscribeMsg{
		Stream:   swarmChunkServerStreamName,
		Key:      nil,
		From:     0,
		To:       0,
		Priority: Top,
	})

	hash := storage.Key(hash0[:])
	chunk := storage.NewChunk(hash, nil)
	chunk.SData = hash
	localStore.Put(chunk)
	chunk.WaitToStore()

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "RetrieveRequestMsg",
		Triggers: []p2ptest.Trigger{
			p2ptest.Trigger{
				Code: 5,
				Msg: &RetrieveRequestMsg{
					Key: hash,
				},
				Peer: peerID,
			},
		},
		Expects: []p2ptest.Expect{
			p2ptest.Expect{
				Code: 1,
				Msg: &OfferedHashesMsg{
					HandoverProof: &HandoverProof{
						Handover: &Handover{},
					},
					Hashes: hash,
					From:   0,
					// TODO: why is this 32???
					To:     32,
					Key:    []byte{},
					Stream: swarmChunkServerStreamName,
				},
				Peer: peerID,
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	hash = storage.Key(hash1[:])
	chunk = storage.NewChunk(hash, nil)
	chunk.SData = hash1[:]
	localStore.Put(chunk)
	chunk.WaitToStore()

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "RetrieveRequestMsg",
		Triggers: []p2ptest.Trigger{
			p2ptest.Trigger{
				Code: 5,
				Msg: &RetrieveRequestMsg{
					Key:       hash,
					SkipCheck: true,
				},
				Peer: peerID,
			},
		},
		Expects: []p2ptest.Expect{
			p2ptest.Expect{
				Code: 6,
				Msg: &ChunkDeliveryMsg{
					Key:   hash,
					SData: hash,
				},
				Peer: peerID,
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestStreamerDownstreamChunkDeliveryMsgExchange(t *testing.T) {
	tester, streamer, localStore, teardown, err := newStreamerTester(t)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	streamer.RegisterClientFunc("foo", func(p *Peer, t []byte) (Client, error) {
		return &testClient{
			t: t,
		}, nil
	})

	peerID := tester.IDs[0]

	err = streamer.Subscribe(peerID, "foo", []byte{}, 5, 8, Top, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	chunkKey := hash0[:]
	chunkData := hash1[:]
	chunk, created := localStore.GetOrCreateRequest(chunkKey)

	if !created {
		t.Fatal("chunk already exists")
	}
	select {
	case <-chunk.ReqC:
		t.Fatal("chunk is already received")
	default:
	}

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "Subscribe message",
		Expects: []p2ptest.Expect{
			p2ptest.Expect{
				Code: 4,
				Msg: &SubscribeMsg{
					Stream:   "foo",
					Key:      []byte{},
					From:     5,
					To:       8,
					Priority: Top,
				},
				Peer: peerID,
			},
		},
	},
		p2ptest.Exchange{
			Label: "ChunkDeliveryRequest message",
			Triggers: []p2ptest.Trigger{
				p2ptest.Trigger{
					Code: 6,
					Msg: &ChunkDeliveryMsg{
						Key:   chunkKey,
						SData: chunkData,
					},
					Peer: peerID,
				},
			},
		})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	timeout := time.NewTimer(1 * time.Second)

	select {
	case <-timeout.C:
		t.Fatal("timeout receiving chunk")
	case <-chunk.ReqC:
	}

	storedChunk, err := localStore.Get(chunkKey)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(storedChunk.SData, chunkData) {
		t.Fatal("Retrieved chunk has different data than original")
	}

}

func TestDeliveryFromNodes(t *testing.T) {
	testDeliveryFromNodes(t, 2, 1, dataChunkCount, true)
	testDeliveryFromNodes(t, 2, 1, dataChunkCount, false)
	testDeliveryFromNodes(t, 4, 1, dataChunkCount, true)
	testDeliveryFromNodes(t, 4, 1, dataChunkCount, false)
	testDeliveryFromNodes(t, 8, 1, dataChunkCount, true)
	testDeliveryFromNodes(t, 8, 1, dataChunkCount, false)
	testDeliveryFromNodes(t, 16, 1, dataChunkCount, true)
	testDeliveryFromNodes(t, 16, 1, dataChunkCount, false)
}

func testDeliveryFromNodes(t *testing.T, nodes, conns, chunkCount int, skipCheck bool) {
	defaultSkipCheck = skipCheck
	toAddr = network.NewAddrFromNodeID
	conf := &streamTesting.RunConfig{
		Adapter:   *adapter,
		NodeCount: nodes,
		ConnLevel: conns,
		ToAddr:    toAddr,
		Services:  services,
	}

	sim, teardown, err := streamTesting.NewSimulation(conf)
	defer teardown()
	if err != nil {
		t.Fatal(err.Error())
	}
	stores = make(map[discover.NodeID]storage.ChunkStore)
	deliveries = make(map[discover.NodeID]*Delivery)
	for i, id := range sim.IDs {
		stores[id] = sim.Stores[i]
	}
	peerCount = func(id discover.NodeID) int {
		if sim.IDs[0] == id || sim.IDs[nodes-1] == id {
			return 1
		}
		return 2
	}

	// here we distribute chunks of a random file into Stores of nodes 1 to nodes
	rrdpa := storage.NewDPA(newRoundRobinStore(sim.Stores[1:]...), storage.NewChunkerParams())
	rrdpa.Start()
	size := chunkCount * chunkSize
	fileHash, wait, err := rrdpa.Store(io.LimitReader(crand.Reader, int64(size)), int64(size))
	// wait until all chunks stored
	wait()
	defer rrdpa.Stop()
	if err != nil {
		t.Fatal(err.Error())
	}
	errc := make(chan error, 1)
	waitPeerErrC = make(chan error)
	quitC := make(chan struct{})

	action := func(ctx context.Context) error {
		// each node Subscribes to each other's swarmChunkServerStreamName
		// need to wait till an aynchronous process registers the peers in streamer.peers
		// that is used by Subscribe
		// using a global err channel to share betweem action and node service
		i := 0
		for err := range waitPeerErrC {
			if err != nil {
				return fmt.Errorf("error waiting for peers: %s", err)
			}
			i++
			if i == nodes {
				break
			}
		}

		// each node subscribes to the upstream swarm chunk server stream
		// which responds to chunk retrieve requests all but the last node in the chain does not
		for j := 0; j < nodes-1; j++ {
			id := sim.IDs[j]
			err := sim.CallClient(id, func(client *rpc.Client) error {
				err := streamTesting.WatchDisconnections(sim.IDs[j], client, peerCount(sim.IDs[j]), errc, quitC)
				if err != nil {
					return err
				}
				ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()
				j++
				sid := sim.IDs[j]
				return client.CallContext(ctx, nil, "stream_subscribeStream", sid, swarmChunkServerStreamName, nil, 0, 0, Top, false)
			})
			if err != nil {
				return err
			}
		}
		// create a retriever dpa for the pivot node
		delivery := deliveries[sim.IDs[0]]
		retrieveFunc := func(chunk *storage.Chunk) error {
			return delivery.RequestFromPeers(chunk.Key[:], skipCheck)
		}
		netStore := storage.NewNetStore(sim.Stores[0].(*storage.LocalStore), retrieveFunc)
		dpa := storage.NewDPA(netStore, storage.NewChunkerParams())
		dpa.Start()

		go func() {
			defer dpa.Stop()
			// start the retrieval on the pivot node - this will spawn retrieve requests for missing chunks
			// we must wait for the peer connections to have started before requesting
			n, err := readAll(dpa, fileHash)
			log.Info(fmt.Sprintf("retrieved %v", fileHash), "read", n, "err", err)
			if err != nil {
				errc <- fmt.Errorf("requesting chunks action error: %v", err)
			}
		}()
		return nil
	}
	checkC := make(chan struct{})
	check := func(ctx context.Context, id discover.NodeID) (bool, error) {
		defer func() { checkC <- struct{}{} }()
		select {
		case err := <-errc:
			return false, err
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}
		var total int64
		err := sim.CallClient(id, func(client *rpc.Client) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return client.CallContext(ctx, &total, "stream_readAll", common.BytesToHash(fileHash))
		})
		log.Info(fmt.Sprintf("check if %08x is available locally: number of bytes read %v/%v (error: %v)", fileHash, total, size, err))
		if err != nil || total != int64(size) {
			return false, nil
		}
		return true, nil
	}

	conf.Step = &simulations.Step{
		Action:  action,
		Trigger: streamTesting.Trigger(10*time.Millisecond, quitC, sim.IDs[0]),
		// we are only testing the pivot node (net.Nodes[0])
		Expect: &simulations.Expectation{
			Nodes: sim.IDs[0:1],
			Check: check,
		},
	}
	startedAt := time.Now()
	timeout := 300 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	result, err := sim.Run(ctx, conf)
	finishedAt := time.Now()
	if err != nil {
		t.Fatalf("Setting up simulation failed: %v", err)
	}
	if result.Error != nil {
		t.Fatalf("Simulation failed: %s", result.Error)
	}
	streamTesting.CheckResult(t, result, startedAt, finishedAt)
}

func BenchmarkDeliveryFromNodesWithoutCheck(b *testing.B) {
	for chunks := 32; chunks <= 128; chunks *= 2 {
		for i := 2; i < 32; i *= 2 {
			b.Run(
				fmt.Sprintf("nodes=%v,chunks=%v", i, chunks),
				func(b *testing.B) {
					benchmarkDeliveryFromNodes(b, i, 1, chunks, true)
				},
			)
		}
	}
}

func BenchmarkDeliveryFromNodesWithCheck(b *testing.B) {
	for chunks := 32; chunks <= 128; chunks *= 2 {
		for i := 2; i < 32; i *= 2 {
			b.Run(
				fmt.Sprintf("nodes=%v,chunks=%v", i, chunks),
				func(b *testing.B) {
					benchmarkDeliveryFromNodes(b, i, 1, chunks, false)
				},
			)
		}
	}
}

func benchmarkDeliveryFromNodes(b *testing.B, nodes, conns, chunkCount int, skipCheck bool) {
	defaultSkipCheck = skipCheck
	toAddr = network.NewAddrFromNodeID
	timeout := 300 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conf := &streamTesting.RunConfig{
		Adapter:   *adapter,
		NodeCount: nodes,
		ConnLevel: conns,
		ToAddr:    toAddr,
		Services:  services,
	}
	sim, teardown, err := streamTesting.NewSimulation(conf)
	defer teardown()
	if err != nil {
		b.Fatal(err.Error())
	}

	stores = make(map[discover.NodeID]storage.ChunkStore)
	deliveries = make(map[discover.NodeID]*Delivery)
	for i, id := range sim.IDs {
		stores[id] = sim.Stores[i]
	}
	peerCount = func(id discover.NodeID) int {
		if sim.IDs[0] == id || sim.IDs[nodes-1] == id {
			return 1
		}
		return 2
	}
	// wait channel for all nodes all peer connections to set up
	waitPeerErrC = make(chan error)

	// create a dpa for the last node in the chain which we are gonna write to
	remoteDpa := storage.NewDPA(sim.Stores[nodes-1], storage.NewChunkerParams())
	remoteDpa.Start()
	defer remoteDpa.Stop()

	// channel to signal simulation initialisation with action call complete
	// or node disconnections
	simErrC := make(chan error)
	quitC := make(chan struct{})
	defer close(quitC)

	action := func(ctx context.Context) error {
		// each node Subscribes to each other's swarmChunkServerStreamName
		// need to wait till an aynchronous process registers the peers in streamer.peers
		// that is used by Subscribe
		// waitPeerErrC using a global err channel to share betweem action and node service
		i := 0
		for err := range waitPeerErrC {
			if err != nil {
				return fmt.Errorf("error waiting for peers: %s", err)
			}
			i++
			if i == nodes {
				break
			}
		}

		// each node except the last one subscribes to the upstream swarm chunk server stream
		// which responds to chunk retrieve requests
		for j := 0; j < nodes-1; j++ {
			id := sim.IDs[j]
			simErrC <- sim.CallClient(id, func(client *rpc.Client) error {
				err := streamTesting.WatchDisconnections(id, client, peerCount(id), simErrC, quitC)
				if err != nil {
					return err
				}
				ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()
				sid := sim.IDs[j+1] // the upstream peer's id
				return client.CallContext(ctx, nil, "stream_subscribeStream", sid, swarmChunkServerStreamName, nil, 0, 0, Top, false)
			})
		}
		// signal to the benchmark that setup is complete
		return err
	}

	// the check function is only triggered when the benchmark finishes
	checkC := make(chan error)
	trigger := make(chan discover.NodeID)
	check := func(ctx context.Context, id discover.NodeID) (_ bool, err error) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-checkC:
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}

	conf.Step = &simulations.Step{
		Action:  action,
		Trigger: trigger,
		// we are only testing the pivot node (net.Nodes[0])
		Expect: &simulations.Expectation{
			Nodes: sim.IDs[0:1],
			Check: check,
		},
	}

	// run the simulation in the background
	errc := make(chan error)
	go func() {
		_, err := sim.Run(ctx, conf)
		errc <- err
	}()

	// wait for simulation action to complete stream subscriptions
	err = <-simErrC
	if err != nil {
		b.Fatalf("simulation failed to initialise. expected no error. got %v", err)
	}
	go func() {
		for {
			var err error
			select {
			case err = <-simErrC:
			case <-quitC:
				return
			}
			trigger <- sim.IDs[0]
			checkC <- err
		}
	}()

	// create a retriever dpa for the pivot node
	// by now deliveries are set for each node by the streamer service
	delivery := deliveries[sim.IDs[0]]
	retrieveFunc := func(chunk *storage.Chunk) error {
		return delivery.RequestFromPeers(chunk.Key[:], skipCheck)
	}
	netStore := storage.NewNetStore(sim.Stores[0].(*storage.LocalStore), retrieveFunc)

	// benchmark loop
	b.ResetTimer()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		// uploading chunkCount random chunks to the last node
		hashes := make([]storage.Key, chunkCount)
		for i := 0; i < chunkCount; i++ {
			// create actual size real chunks
			hash, wait, err := remoteDpa.Store(io.LimitReader(crand.Reader, int64(chunkSize)), int64(chunkSize))
			// wait until all chunks stored
			wait()
			if err != nil {
				b.Fatalf("expected no error. got %v", err)
			}
			// collect the hashes
			hashes[i] = hash
		}
		// now benchmark the actual retrieval
		// netstore.Get is called for each hash in a go routine and errors are collected
		b.StartTimer()
		errs := make(chan error)
		for _, hash := range hashes {
			go func(h storage.Key) {
				_, err := netStore.Get(h)
				log.Warn("test check netstore get", "hash", h, "err", err)
				errs <- err
			}(hash)
		}
		// count and report retrieval errors
		// if there are misses then chunk timeout is too low for the distance and volume (?)
		var total, misses int
		for err := range errs {
			if err != nil {
				log.Warn(err.Error())
				misses++
			}
			total++
			if total == chunkCount {
				break
			}
		}
		b.StopTimer()
		if misses > 0 {
			simErrC <- fmt.Errorf("%v chunk not found out of %v", misses, total)
		}
	}
	// benchmark over, trigger the check function to conclude the simulation
	err = <-errc
	if err != nil {
		b.Fatalf("expected no error. got %v", err)
	}
}