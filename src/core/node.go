package core

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/bits"
	mrand "math/rand"
	"time"
)

// AkademiNode constants.
const (
	IDLength   = 20
	BucketSize = 20
)

// List of bootstrap nodes used for first connecting to
// the network.
var BootstrapHosts = [...]Host{
	"akademi_bootstrap:3865",
}

// Akademi uses node and key IDs, whose length is defined
// in bytes by IDLength.
type BaseID [IDLength]byte

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

	RoutingTable [IDLength * 8][]RoutingEntry

	Dispatcher Dispatcher
}

// Update the routing table with a new entry
func (a *AkademiNode) UpdateRoutingTable(r RoutingEntry) error {
	prefix := r.NodeID.GetPrefixLength(a.NodeID)
	for i, v := range a.RoutingTable[prefix] {
		if v.NodeID == r.NodeID || v.Host == r.Host {
			a.RoutingTable[prefix] = append(a.RoutingTable[prefix][:i], a.RoutingTable[prefix][i+1:]...)
			break
		}
	}
	if len(a.RoutingTable[prefix]) >= BucketSize {
		return fmt.Errorf("Bucket already full.")
	}
	a.RoutingTable[prefix] = append(a.RoutingTable[prefix], r)
	return nil
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenPort IPPort, bootstrap bool) error {
	a.ListenPort = listenPort
	a.NodeID = RandomBaseID()
	a.Dispatcher = dispatcher
	err := a.Dispatcher.Initialize(BaseMessageHeader{ListenPort: a.ListenPort, NodeID: a.NodeID})
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
		log.Print("Connected to bootstrap node \"", BootstrapHosts[i], "\". NodeID: ", header.NodeID.Base64Str())
	}
	return nil
}

// Print all the entries in the routing table
func (a *AkademiNode) PrintRoutingTable() {
	fmt.Println("Node routing table:")
	for _, bucket := range a.RoutingTable {
		for _, r := range bucket {
			fmt.Println(r.Host, r.NodeID.Base64Str())
		}
	}
}

// The function GetPrefixLength finds the length of the
// common prefix between two Node/Key IDs.
func (id0 *BaseID) GetPrefixLength(id1 BaseID) int {
	for i := 0; i < IDLength; i++ {
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
	for i := 0; i < IDLength; i++ {
		out = out + fmt.Sprintf("%08b", id[i])
	}
	return out
}

// Returns base64 BaseID string.
func (id *BaseID) Base64Str() string {
	return base64.StdEncoding.EncodeToString(id[:])
}

// Returns random BaseID.
func RandomBaseID() BaseID {
	var o BaseID
	crand.Read(o[:])
	return o
}
