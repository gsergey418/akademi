package test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/node"
)

type MockDispatcher struct {
	pingError     bool
	routingHeader core.RoutingHeader
}

func (d *MockDispatcher) Initialize(h core.RoutingHeader) error {
	return fmt.Errorf("Error: function not implemented.")
}

func (d *MockDispatcher) Ping(host core.Host) (core.RoutingHeader, error) {
	if d.pingError {
		return core.RoutingHeader{}, fmt.Errorf("Ping failed.")
	} else {
		return d.routingHeader, nil
	}
}

func (d *MockDispatcher) FindNode(host core.Host, nodeID core.BaseID) (core.RoutingHeader, []core.RoutingEntry, error) {
	return core.RoutingHeader{}, []core.RoutingEntry{}, fmt.Errorf("Error: function not implemented.")
}

func (d *MockDispatcher) FindKey(host core.Host, keyID core.BaseID) (core.RoutingHeader, core.DataBytes, []core.RoutingEntry, error) {
	return core.RoutingHeader{}, core.DataBytes{}, []core.RoutingEntry{}, fmt.Errorf("Error: function not implemented.")
}

func (d *MockDispatcher) Store(host core.Host, data core.DataBytes) (core.RoutingHeader, error) {
	return core.RoutingHeader{}, fmt.Errorf("Error: function not implemented.")
}

func TestUpdateRoutingTable(t *testing.T) {
	a := node.AkademiNode{}
	a.NodeID = core.RandomBaseID()
	dispatcher := &MockDispatcher{}
	a.Dispatcher = dispatcher
	r1 := core.RoutingEntry{
		NodeID: core.RandomBaseID(),
		Host:   "153.75.235.134:3865",
	}
	id2 := r1.NodeID
	id2[19] = 42
	r2 := core.RoutingEntry{
		NodeID: id2,
		Host:   "122.56.173.15:3865",
	}

	a.UpdateRoutingTable(r1)
	for i := 0; i < core.BucketSize+10; i++ {
		id := r1.NodeID
		id[19] = byte(rand.IntN(255))
		listenPort := rand.IntN(65535)
		r := core.RoutingEntry{
			NodeID: id,
			Host:   core.Host(fmt.Sprintf("%d.%d.%d.%d:%d", rand.IntN(255), rand.IntN(255), rand.IntN(255), rand.IntN(255), listenPort)),
		}
		dispatcher.pingError = i%2 == 0
		aliveID, _ := core.B32ToID(a.RoutingTableString()[:32])
		dispatcher.routingHeader = core.RoutingHeader{
			NodeID:     aliveID,
			ListenPort: core.IPPort(listenPort),
		}
		a.UpdateRoutingTable(r)
		if !dispatcher.pingError {
			nodes, err := a.GetClosestNodes(aliveID, 1)
			if nodes[0].NodeID != aliveID || err != nil {
				t.Fail()
			}
		} else {
			nodes, err := a.GetClosestNodes(id, 1)
			if nodes[0].NodeID != id || err != nil {
				t.Fail()
			}
		}
	}
	a.UpdateRoutingTable(r2)
	a.UpdateRoutingTable(r1)
}
