package listener

import "github.com/gsergey418alt/akademi/core"

// RPCResponse is a structure embedded in all RPC responses of Akademi.
type RPCResponse struct {
	NodeID core.NodeID
}

// PingRPCResponse is a response structure for the Ping RPC.
type PingRPCResponse struct {
	RPCResponse
}

// FindNodeRPCResponse is a response structure for the FindNode RPC.
type FindNodeRPCResponse struct {
	RPCResponse
	routingEntry []core.RoutingEntry
}

// FindKeyRPCResponse is a response structure for the FindKey RPC.
type FindKeyRPCResponse struct {
	RPCResponse
	keyID        core.KeyID
	routingEntry []core.RoutingEntry
}

// StoreRPCResponse is a response structure for the Store RPC.
type StoreRPCResponse struct {
	RPCResponse
}
