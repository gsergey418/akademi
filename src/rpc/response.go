package rpc

import "github.com/gsergey418alt/akademi/core"

// BaseRPCResponse is a structure embedded in all RPC responses of Akademi.
type BaseRPCResponse struct {
	NodeID core.NodeID
}

// PingRPCResponse is a response structure for the Ping RPC.
type PingRPCResponse struct {
	BaseRPCResponse
}

// FindNodeRPCResponse is a response structure for the FindNode RPC.
type FindNodeRPCResponse struct {
	BaseRPCResponse
	routingEntry []core.RoutingEntry
}

// FindKeyRPCResponse is a response structure for the FindKey RPC.
type FindKeyRPCResponse struct {
	BaseRPCResponse
	keyID        core.KeyID
	routingEntry []core.RoutingEntry
}

// StoreRPCResponse is a response structure for the Store RPC.
type StoreRPCResponse struct {
	BaseRPCResponse
}
