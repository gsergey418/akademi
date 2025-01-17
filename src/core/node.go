package core

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"math/bits"
	mrand "math/rand"
	"time"
)

// List of bootstrap nodes used for first connecting to
// the network.
var BootstrapHosts = [...]Host{
	"akademi_bootstrap:3856",
}

// Akademi uses 256-bit node and key IDs.
type BaseID [32]byte

// Separate IPPort type because the IP address is
// identified by receiving node.
type IPPort uint16

// Host is used to identify node's IP address and
// port.
type Host string

// DataBytes is a type for values to be stored in akademi
// nodes.
type DataBytes []byte

// RoutingEntry is a structure that stores routing
// information about an akademi node.
type RoutingEntry struct {
	Host   Host
	NodeID BaseID
}

// AkademiNode is a structure containing the core kademlia
// logic.
type AkademiNode struct {
	NodeID        BaseID
	ListenPort    IPPort
	KeyValueStore map[BaseID][]byte

	RoutingTable [256][20]RoutingEntry

	Dispatcher Dispatcher
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenPort IPPort, bootstrap bool) error {
	a.ListenPort = listenPort
	nodeID, err := RandomBaseID()
	if err != nil {
		return err
	}
	a.NodeID = nodeID
	a.Dispatcher = dispatcher
	err = a.Dispatcher.Initialize(BaseMessageHeader{ListenPort: a.ListenPort, NodeID: a.NodeID})
	if err != nil {
		return err
	}
	if bootstrap {
		i := mrand.Intn(len(BootstrapHosts))
		var header BaseMessageHeader
		for header, err = a.Dispatcher.Ping(BootstrapHosts[i]); err != nil; {
			log.Print(err)
			time.Sleep(5 * time.Second)
		}
		log.Print("Connected to bootstrap node \"", BootstrapHosts[i], "\". NodeID:", header.NodeID)
	}
	return nil
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

// Returns BaseID string in binary.
func (id *BaseID) BinStr() string {
	out := ""
	for i := 0; i < 32; i++ {
		out = out + fmt.Sprintf("%08b", id[i])
	}
	return out
}

// Returns random BaseID.
func RandomBaseID() (BaseID, error) {
	var o BaseID
	_, err := crand.Read(o[:])
	return o, err
}
