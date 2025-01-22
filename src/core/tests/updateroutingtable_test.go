package test

import (
	"testing"

	"github.com/gsergey418alt/akademi/core"
)

func TestUpdateRoutingTable(t *testing.T) {
	a := core.AkademiNode{}
	a.NodeID = core.RandomBaseID()
	r1 := core.RoutingEntry{
		NodeID: core.RandomBaseID(),
		Host:   "153.75.235.134:3865",
	}
	id2 := r1.NodeID
	id2[19] = 94
	r2 := core.RoutingEntry{
		NodeID: id2,
		Host:   "127.0.0.1:3865",
	}
	id3 := r1.NodeID
	id3[19] = 42
	r3 := core.RoutingEntry{
		NodeID: id3,
		Host:   "122.56.173.15:3865",
	}

	a.UpdateRoutingTable(r1)
	a.PrintRoutingTable()
	a.UpdateRoutingTable(r2)
	a.PrintRoutingTable()
	a.UpdateRoutingTable(r3)
	a.PrintRoutingTable()
	a.UpdateRoutingTable(r1)
	a.PrintRoutingTable()
}
