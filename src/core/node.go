package core

import (
	"log"

	mrand "math/rand"
	"time"
)

// AkademiNode constants.
const (
	IDLength           = 20
	BucketSize         = 20
	ConcurrentRequests = 3
)

// List of bootstrap nodes used for first connecting to
// the network.
var BootstrapHosts = [...]Host{
	"akademi_bootstrap_1:3865",
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

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenPort IPPort, bootstrap bool) error {
	a.ListenPort = listenPort
	a.NodeID = RandomBaseID()
	log.Print("Initializing Akademi node. NodeID: ", a.NodeID)

	a.Dispatcher = dispatcher
	err := a.Dispatcher.Initialize(RoutingHeader{ListenPort: a.ListenPort, NodeID: a.NodeID})
	if err != nil {
		return err
	}

	if bootstrap {
		i := mrand.Intn(len(BootstrapHosts))
		var header RoutingHeader
		var nodes []RoutingEntry
		for header, nodes, err = a.FindNode(BootstrapHosts[i], a.NodeID); err != nil; {
			log.Print(err)
			i = mrand.Intn(len(BootstrapHosts))
			time.Sleep(5 * time.Second)
		}
		log.Print("Connected to bootstrap node \"", BootstrapHosts[i], "\". NodeID: ", header.NodeID)
		log.Print("Neighbor nodes:")
		for _, v := range nodes {
			log.Print(v)
		}
		a.logRoutingTable()
	}
	return nil
}

// The responseHandler function manages the side effects
// of receiving an RPC response from the Dispatcher.
func (a *AkademiNode) responseHandler(host Host, header RoutingHeader) {
	r := RoutingEntry{
		Host:   host,
		NodeID: header.NodeID,
	}
	a.UpdateRoutingTable(r)
}

// Redefinitions of Dispatcher functions.

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (a *AkademiNode) Ping(host Host) (RoutingHeader, error) {
	header, err := a.Dispatcher.Ping(host)
	a.responseHandler(host, header)
	return header, err
}

// The FindNode function dispatches a FindNode RPC call
// to node located at host.
func (a *AkademiNode) FindNode(host Host, nodeID BaseID) (RoutingHeader, []RoutingEntry, error) {
	header, nodes, err := a.Dispatcher.FindNode(host, nodeID)
	a.responseHandler(host, header)
	for _, r := range nodes {
		a.UpdateRoutingTable(r)
	}
	return header, nodes, err
}

// The FindKey function dispatches a FindKey RPC call to
// node located at host.
func (a *AkademiNode) FindKey(host Host, keyID BaseID) (RoutingHeader, DataBytes, []RoutingEntry, error) {
	header, data, nodes, err := a.Dispatcher.FindKey(host, keyID)
	a.responseHandler(host, header)
	for _, r := range nodes {
		a.UpdateRoutingTable(r)
	}
	return header, data, nodes, err
}

// The Store function dispatches a Store RPC call to node
// located at host.
func (a *AkademiNode) Store(host Host, keyID BaseID, value DataBytes) (RoutingHeader, error) {
	header, err := a.Dispatcher.Store(host, keyID, value)
	a.responseHandler(host, header)
	return header, err
}
