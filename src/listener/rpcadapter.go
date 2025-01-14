package listener

import (
	"github.com/gsergey418alt/akademi/core"
)

// RPCResponse is a structure embedded in all RPC responses of Akademi.
type RPCResponse struct {
	NodeID core.NodeID
}

// RPCAdapter is an interface that enables connection of
// the core logic to the Listener.
type RPCAdapter interface {
	Ping(args struct{}, reply *struct{ RPCResponse }) error
	FindNode(args struct{ nodeID core.NodeID }, reply *struct {
		RPCResponse
		routingEntry *[]core.RoutingEntry
	}) error
	FindKey(args struct{ keyID core.KeyID }, reply *struct {
		RPCResponse
		keyID        core.KeyID
		routingEntry []core.RoutingEntry
	}) error
	Store(args struct {
		keyID core.KeyID
		data  core.DataBytes
	}, reply *struct{ RPCResponse }) error
}

// AkademiNodeRPCAdapter is an interface that enables
// connection of AkademiNode to the Listener.
type AkademiNodeRPCAdapter struct {
	AkademiNode *core.AkademiNode
}

// The PING RPC probes a node to see if it is online.
// STORE instructs a node to store a 〈key, value〉 pair
// for later retrieval.
func (a *AkademiNodeRPCAdapter) Ping(args struct{}, reply *struct{ RPCResponse }) error {
	reply.NodeID = a.AkademiNode.NodeID
	return nil
}

// FIND NODE takes a 160-bit ID as an argu-
// ment. The recipient of a the RPC returns
// 〈IP address, UDP port, Node ID〉 triples for the k
// nodes it knows about closest to the target ID. These
// triples can come from a single k-bucket, or they may
// come from multiple k-buckets if the closest k-bucket
// is not full. In any case, the RPC recipient must return
// k items (unless there are fewer than k nodes in all its
// k-buckets combined, in which case it returns every
// node it knows about).
func (a *AkademiNodeRPCAdapter) FindNode(args struct{ nodeID core.NodeID }, reply *struct {
	RPCResponse
	routingEntry *[]core.RoutingEntry
}) error {
	reply.NodeID = a.AkademiNode.NodeID
	panic("Function \"FindNode\" not implemented.")
}

// FIND VALUE behaves like FIND NODE—returning
// 〈IP address, UDP port, Node ID〉 triples—with one
// exception. If the RPC recipient has received a STORE
// RPC for the key, it just returns the stored value.
func (a *AkademiNodeRPCAdapter) FindKey(args struct{ keyID core.KeyID }, reply *struct {
	RPCResponse
	keyID        core.KeyID
	routingEntry []core.RoutingEntry
}) error {
	reply.NodeID = a.AkademiNode.NodeID
	panic("Function \"FindKey\" not implemented.")
}

// STORE instructs a node to store a 〈key, value〉 pair
// for later retrieval.
func (a *AkademiNodeRPCAdapter) Store(args struct {
	keyID core.KeyID
	data  core.DataBytes
}, reply *struct{ RPCResponse }) error {
	reply.NodeID = a.AkademiNode.NodeID
	panic("Function \"Store\" not implemented.")
}
