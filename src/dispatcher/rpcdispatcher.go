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
type RPCDispatcher struct {
	Self *core.RoutingEntry
}

// The initializer for the Dispatcher sturct sets the self
func (d *RPCDispatcher) Initialize(self *core.RoutingEntry) {
	d.Self = self
}

// DispatchRPCCall is a wrapper function for makubg RPC
// calls to an address defined by dispatchAddr.
func (d *RPCDispatcher) DispatchRPCCall(dispatchAddr core.ListenAddr, f func(*rpc.Client) error) error {
	log.Print("Connecting to RPC at ", dispatchAddr, ".\n")
	client, err := rpc.DialHTTP("tcp", string(dispatchAddr))
	if err != nil {
		return err
	}
	defer client.Close()
	err = f(client)
	return err
}

// Populates the generic RPCRequest struct with values from
// the dispatcher.
func (d *RPCDispatcher) PopulateRPCRequest(r *akademiRPC.BaseRPCRequest) {
	r.Self = *d.Self
}

// The Ping function dispatches a Ping RPC to a node at the
// address dispatchAddr.
func (d *RPCDispatcher) Ping(dispatchAddr core.ListenAddr) (core.BaseID, error) {
	args, reply := akademiRPC.PingRequest{}, akademiRPC.PingResponse{}
	d.PopulateRPCRequest(&args.BaseRPCRequest)
	err := d.DispatchRPCCall(dispatchAddr, func(c *rpc.Client) error {
		return c.Call("AkademiNodeRPCAdapter.Ping", args, &reply)
	})
	return reply.Self.NodeID, err
}

// The FindNode function dispatches a FindNode RPC to a
// node at the address dispatchAddr.
func (d *RPCDispatcher) FindNode(dispatchAddr core.ListenAddr, nodeID core.BaseID) (core.BaseID, []core.RoutingEntry, error) {
	panic("Function	FindNode not implemented.")
}

// The FindKey function dispatches a FindKey RPC to a node
// at the address dispatchAddr.
func (d *RPCDispatcher) FindKey(dispatchAddr core.ListenAddr, keyID core.BaseID) (core.BaseID, core.BaseID, []core.RoutingEntry, error) {
	panic("Function	FindKey not implemented.")
}

// The Store function dispatches a Store RPC to a node at
// the address dispatchAddr.
func (d *RPCDispatcher) Store(dispatchAddr core.ListenAddr, keyID core.BaseID, value core.DataBytes) (core.BaseID, error) {
	panic("Function	Store not implemented.")
}
