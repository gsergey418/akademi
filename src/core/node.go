package core

import (
	crand "crypto/rand"
	"log"
	"math/bits"
	mrand "math/rand"
	"time"
)

// List of bootstrap nodes used for first connecting to
// the network.
var BootstrapHosts = [...]ListenAddr{
	"akademi_bootstrap:3856",
}

// Akademi uses 256-bit node and key IDs.
type BaseID [32]byte

// ListenAddress is used to identify node's IP address and
// port.
type ListenAddr string

// DataBytes is a type for values to be stored in akademi
// nodes.
type DataBytes []byte

// RoutingEntry is a structure that stores routing
// information about an akademi node.
type RoutingEntry struct {
	Host   ListenAddr
	NodeID BaseID
}

// AkademiNode is a structure containing the core kademlia
// logic.
type AkademiNode struct {
	Self          RoutingEntry
	KeyValueStore map[BaseID][]byte

	RoutingTable [256][20]RoutingEntry

	Dispatcher Dispatcher
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenAddr ListenAddr, bootstrap bool) {
	a.Self.Host = listenAddr
	_, err := crand.Read(a.Self.NodeID[:])
	if err != nil {
		log.Fatal(err)
	}
	a.Dispatcher = dispatcher
	a.Dispatcher.Initialize(&a.Self)
	if bootstrap {
		i := mrand.Intn(len(BootstrapHosts))
		var nodeID BaseID
		for nodeID, err = a.Dispatcher.Ping(BootstrapHosts[i]); err != nil; {
			log.Print(err)
			time.Sleep(5 * time.Second)
		}
		log.Print("Connected to bootstrap node \"", BootstrapHosts[i], "\". NodeID:", nodeID)
	}
}

// The function GetPrefixLength finds the length of the
// common prefix between two 256-bit Node/Key IDs.
func (id0 *BaseID) GetPrefixLength(id1 BaseID) int {
	for i := 0; i < 32; i++ {
		xor := id0[i] ^ id1[i]
		if xor != 0 {
			return i*8 + bits.LeadingZeros8(xor)
		}
	}
	return 0
}
