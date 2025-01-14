package listener

import (
	"github.com/gsergey418alt/akademi/core"
)

// ListenerAdapter is an interface that enables connection
// of the core logic to the Listener.
type ListenerAdapter interface {
	Ping(args struct{}, reply *PingRPCResponse) error
	FindNode(args struct{ nodeID core.NodeID }, reply *FindNodeRPCResponse) error
	FindKey(args struct{ keyID core.KeyID }, reply *FindKeyRPCResponse) error
	Store(args struct {
		keyID core.KeyID
		data  core.DataBytes
	}, reply *StoreRPCResponse) error
}

// AkademiNodeRPCAdapter is an interface that enables
// connection of AkademiNode to the Listener.
type AkademiNodeRPCAdapter struct {
	AkademiNode *core.AkademiNode
}

// Populates the generic RPCResponse struct with values from AkademiCore.
func (a *AkademiNodeRPCAdapter) PopulateRPCResponse(r *RPCResponse) {
	r.NodeID = a.AkademiNode.NodeID
}

// The PING RPC probes a node to see if it is online.
// STORE instructs a node to store a 〈key, value〉 pair
// for later retrieval.
func (a *AkademiNodeRPCAdapter) Ping(args struct{}, reply *PingRPCResponse) error {
	a.PopulateRPCResponse(&reply.RPCResponse)
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
func (a *AkademiNodeRPCAdapter) FindNode(args struct{ nodeID core.NodeID }, reply *FindNodeRPCResponse) error {
	a.PopulateRPCResponse(&reply.RPCResponse)
	panic("Function \"FindNode\" not implemented.")
}

// FIND VALUE behaves like FIND NODE—returning
// 〈IP address, UDP port, Node ID〉 triples—with one
// exception. If the RPC recipient has received a STORE
// RPC for the key, it just returns the stored value.
func (a *AkademiNodeRPCAdapter) FindKey(args struct{ keyID core.KeyID }, reply *FindKeyRPCResponse) error {
	a.PopulateRPCResponse(&reply.RPCResponse)
	panic("Function \"FindKey\" not implemented.")
}

// STORE instructs a node to store a 〈key, value〉 pair
// for later retrieval.
func (a *AkademiNodeRPCAdapter) Store(args struct {
	keyID core.KeyID
	data  core.DataBytes
}, reply *StoreRPCResponse) error {
	a.PopulateRPCResponse(&reply.RPCResponse)
	panic("Function \"Store\" not implemented.")
}
