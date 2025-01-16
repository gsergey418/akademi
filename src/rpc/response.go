package rpc

import "github.com/gsergey418alt/akademi/core"

// BaseResponse is a structure embedded in all RPC responses of Akademi.
type BaseResponse struct {
	Self core.RoutingEntry
}

// PingResponse is a response structure for the Ping RPC.
type PingResponse struct {
	BaseResponse
}

// FindNodeResponse is a response structure for the FindNode RPC.
type FindNodeResponse struct {
	BaseResponse
	routingEntry []core.RoutingEntry
}

// FindKeyResponse is a response structure for the FindKey RPC.
type FindKeyResponse struct {
	BaseResponse
	keyID        core.BaseID
	routingEntry []core.RoutingEntry
}

// StoreResponse is a response structure for the Store RPC.
type StoreResponse struct {
	BaseResponse
}
