package rpc

import "github.com/gsergey418alt/akademi/core"

// BaseRPCRequest contains the base arguments for RPC calls
// sent out by nodes.
type BaseRPCRequest struct{}

// PingRequest contains the arguments for the Ping
// RPC call.
type PingRequest struct {
	BaseRPCRequest
}

// FindNodeRequest contains the arguments for the FindNode
// RPC call.
type FindNodeRequest struct {
	BaseRPCRequest
	nodeID core.NodeID
}

// FindKeyRequest contains the arguments for the FindKey
// RPC call.
type FindKeyRequest struct {
	BaseRPCRequest
	keyID core.KeyID
}

// StoreRequest contains the arguments for the Store
// RPC call.
type StoreRequest struct {
	BaseRPCRequest
	keyID core.KeyID
	data  core.DataBytes
}
