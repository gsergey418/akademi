package dispatcher

import "github.com/gsergey418alt/akademi/core"

// RPCDispatcher is an implementation of the Dispatcher
// interface that interacts with other peers through HTTP
// RPC.
type RPCDispatcher struct{}

// The Ping function dispatches a Ping RPC to a node at the
// address addr.
func (d *RPCDispatcher) Ping(addr string) error {
	panic("Function	Ping not implemented.")
}

// The FindNode function dispatches a FindNode RPC to a
// node at the address addr.
func (d *RPCDispatcher) FindNode(addr string, nodeID core.NodeID) ([]core.RoutingEntry, error) {
	panic("Function	FindNode not implemented.")
}

// The FindKey function dispatches a FindKey RPC to a node
// at the address addr.
func (d *RPCDispatcher) FindKey(addr string, keyID core.KeyID) (core.KeyID, []core.RoutingEntry, error) {
	panic("Function	FindKey not implemented.")
}

// The Store function dispatches a Store RPC to a node at
// the address addr.
func (d *RPCDispatcher) Store(addr string, keyID core.KeyID, value core.DataBytes) error {
	panic("Function	Store not implemented.")
}
