package core

import (
	"fmt"
	"log"
	"sync"

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
	StartTime     time.Time
	keyValueStore struct {
		data map[BaseID][]byte
		lock sync.Mutex
	}

	routingTable struct {
		data [IDLength * 8][]RoutingEntry
		lock sync.Mutex
	}

	dispatcher Dispatcher
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenPort IPPort, bootstrap bool) error {
	a.ListenPort = listenPort
	a.NodeID = RandomBaseID()
	a.StartTime = time.Now()
	log.Print("Initializing Akademi node. NodeID: ", a.NodeID)

	a.dispatcher = dispatcher
	err := a.dispatcher.Initialize(RoutingHeader{ListenPort: a.ListenPort, NodeID: a.NodeID})
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
		a.LogRoutingTable()
	}
	return nil
}

func (a *AkademiNode) Uptime() time.Duration {
	return time.Since(a.StartTime)
}

// Get node information string.
func (a *AkademiNode) NodeInfo() (nodeInfo string) {
	nodeInfo += fmt.Sprintf("NodeID: %s\n", a.NodeID)
	uptime := a.Uptime()
	nodeInfo += fmt.Sprintf("Uptime: %02d:%02d:%02d", int(uptime.Hours()), int(uptime.Minutes()), int(uptime.Seconds()))
	return
}
