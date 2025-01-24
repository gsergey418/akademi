package node

import (
	"fmt"
	"log"
	"sync"

	"math/rand"
	"time"

	"github.com/gsergey418alt/akademi/core"
)

// AkademiNode is a structure containing the core kademlia
// logic.
type AkademiNode struct {
	NodeID     core.BaseID
	ListenPort core.IPPort
	StartTime  time.Time
	dataStore  struct {
		data map[core.BaseID]core.DataBytes
		lock sync.Mutex
	}

	routingTable struct {
		data [core.IDLength*8 + 1][]core.RoutingEntry
		lock sync.Mutex
	}

	Dispatcher Dispatcher
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenPort core.IPPort, bootstrap bool, bootstrapList []core.Host) error {
	a.ListenPort = listenPort
	a.NodeID = core.RandomBaseID()
	a.StartTime = time.Now()
	log.Print("Initializing Akademi node. NodeID: ", a.NodeID)

	a.Dispatcher = dispatcher
	err := a.Dispatcher.Initialize(core.RoutingHeader{ListenPort: a.ListenPort, NodeID: a.NodeID})
	if err != nil {
		return err
	}

	a.dataStore.data = make(map[core.BaseID]core.DataBytes)

	if bootstrap {
		log.Print("Initiating node bootstrap. Hosts: ", bootstrapList, ".")
		bootstrapHosts := bootstrapList
		for bootstrapCount := 0; bootstrapCount < min(core.Bootstraps, len(bootstrapList)); bootstrapCount++ {
			var i int
			var err error
			var header core.RoutingHeader
			i = rand.Intn(len(bootstrapHosts))
			header, _, err = a.FindNode(bootstrapHosts[i], a.NodeID)
			for err != nil {
				log.Print(err)
				time.Sleep(5 * time.Second)
				i = rand.Intn(len(bootstrapHosts))
				header, _, err = a.FindNode(bootstrapHosts[i], a.NodeID)
			}
			log.Print("Connected to bootstrap node \"", bootstrapList[i], "\". NodeID: ", header.NodeID)
			bootstrapHosts = append(bootstrapHosts[:i], bootstrapHosts[i+1:]...)
		}
		log.Print("Bootstrapping process finished.")
	}
	return nil
}

// Returns time.Duration of the node's uptime.
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
