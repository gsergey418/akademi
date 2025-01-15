package core

import (
	"crypto/rand"
	"log"
)

// List of bootstrap nodes used for first connecting to
// the network.
var BootstrapHosts = [...]string{
	"akademi_bootstrap:3306",
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
func (a *AkademiNode) Initialize() {
	_, err := rand.Read(a.NodeID[:])
	if err != nil {
		log.Fatal(err)
	}
}
