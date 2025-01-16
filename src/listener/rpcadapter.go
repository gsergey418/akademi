package listener

import (
	"log"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/rpc"
)

// RPCListenerAdapter is an interface that enables connection
// of the core logic to the Listener.
type RPCListenerAdapter interface {
	Ping(args rpc.PingRequest, reply *rpc.PingResponse) error
	FindNode(args rpc.FindNodeRequest, reply *rpc.FindNodeResponse) error
	FindKey(args rpc.FindKeyRequest, reply *rpc.FindKeyResponse) error
	Store(args rpc.StoreRequest, reply *rpc.StoreResponse) error
}

// AkademiNodeRPCAdapter is an interface that enables
// connection of AkademiNode to the Listener.
type AkademiNodeRPCAdapter struct {
	AkademiNode *core.AkademiNode
}

// Populates the generic RPCResponse struct with values from AkademiCore.
func (a *AkademiNodeRPCAdapter) PopulateRPCResponse(r *rpc.BaseResponse) {
	r.Self = a.AkademiNode.Self
}

// The PING RPC probes a node to see if it is online.
// STORE instructs a node to store a 〈key, value〉 pair
// for later retrieval.
func (a *AkademiNodeRPCAdapter) Ping(args rpc.PingRequest, reply *rpc.PingResponse) error {
	a.PopulateRPCResponse(&reply.BaseResponse)
	log.Print("Received ping from address \"", args.Self.Host, "\".")
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
func (a *AkademiNodeRPCAdapter) FindNode(args rpc.FindNodeRequest, reply *rpc.FindNodeResponse) error {
	a.PopulateRPCResponse(&reply.BaseResponse)
	panic("Function \"FindNode\" not implemented.")
}

// FIND VALUE behaves like FIND NODE—returning
// 〈IP address, UDP port, Node ID〉 triples—with one
// exception. If the RPC recipient has received a STORE
// RPC for the key, it just returns the stored value.
func (a *AkademiNodeRPCAdapter) FindKey(args rpc.FindKeyRequest, reply *rpc.FindKeyResponse) error {
	a.PopulateRPCResponse(&reply.BaseResponse)
	panic("Function \"FindKey\" not implemented.")
}

// STORE instructs a node to store a 〈key, value〉 pair
// for later retrieval.
func (a *AkademiNodeRPCAdapter) Store(args rpc.StoreRequest, reply *rpc.StoreResponse) error {
	a.PopulateRPCResponse(&reply.BaseResponse)
	panic("Function \"Store\" not implemented.")
}
