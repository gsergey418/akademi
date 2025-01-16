package dispatcher

import (
	"log"
	"net/rpc"

	"github.com/gsergey418alt/akademi/core"
	akademiRPC "github.com/gsergey418alt/akademi/rpc"
)

// RPCDispatcher is an implementation of the Dispatcher
// interface that interacts with other peers through HTTP
// RPC.
type RPCDispatcher struct{}

func (d *RPCDispatcher) DispatchRPCCall(addr string, f func(*rpc.Client) error) error {
	log.Print("Connecting to RPC at ", addr, ".\n")
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return err
	}
	defer client.Close()
	err = f(client)
	return err
}

// The Ping function dispatches a Ping RPC to a node at the
// address addr.
func (d *RPCDispatcher) Ping(addr string) (core.BaseID, error) {
	args, reply := struct{}{}, akademiRPC.PingRPCResponse{}
	err := d.DispatchRPCCall(addr, func(c *rpc.Client) error {
		return c.Call("AkademiNodeRPCAdapter.Ping", args, &reply)
	})
	return reply.NodeID, err
}

// The FindNode function dispatches a FindNode RPC to a
// node at the address addr.
func (d *RPCDispatcher) FindNode(addr string, nodeID core.BaseID) (core.BaseID, []core.RoutingEntry, error) {
	panic("Function	FindNode not implemented.")
}

// The FindKey function dispatches a FindKey RPC to a node
// at the address addr.
func (d *RPCDispatcher) FindKey(addr string, keyID core.BaseID) (core.BaseID, core.BaseID, []core.RoutingEntry, error) {
	panic("Function	FindKey not implemented.")
}

// The Store function dispatches a Store RPC to a node at
// the address addr.
func (d *RPCDispatcher) Store(addr string, keyID core.BaseID, value core.DataBytes) (core.BaseID, error) {
	panic("Function	Store not implemented.")
}
