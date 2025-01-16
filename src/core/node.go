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
var BootstrapHosts = [...]string{
	"akademi_bootstrap:3856",
}

// Akademi uses 256-bit node and key IDs.
type NodeID [32]byte
type KeyID [32]byte

// DataBytes is a type for values to be stored in akademi
// nodes.
type DataBytes []byte

// RoutingEntry is a structure that stores routing
// information about an akademi node.
type RoutingEntry struct {
	Host   string
	NodeID NodeID
}

// AkademiNode is a structure containing the core kademlia
// logic.
type AkademiNode struct {
	NodeID        NodeID
	KeyValueStore map[KeyID][]byte

	RoutingTable [256][20]RoutingEntry

	Dispatcher Dispatcher
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(bootstrap bool) {
	_, err := crand.Read(a.NodeID[:])
	if err != nil {
		log.Fatal(err)
	}
	if bootstrap {
		i := mrand.Intn(len(BootstrapHosts))
		var nodeID NodeID
		for nodeID, err = a.Dispatcher.Ping(BootstrapHosts[i]); err != nil; {
			log.Print(err)
			time.Sleep(5 * time.Second)
		}
		log.Print("Connected to bootstrap node \"", BootstrapHosts[i], "\". NodeID:", nodeID)
	}
}

// The function GetPrefixLength finds the length of the
// common prefix between two 256-bit Node/Key IDs.
func (a *AkademiNode) GetPrefixLength(id0, id1 [32]byte) int {
	for i := 0; i < 32; i++ {
		xor := id0[i] ^ id1[i]
		if xor != 0 {
			return i*8 + bits.LeadingZeros8(xor)
		}
	}
	return 0
}
