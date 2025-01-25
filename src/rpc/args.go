package rpc

import (
	"github.com/gsergey418alt/akademi/core"
)

// Arguments for the Ping RPC.
type PingArgs struct {
	Host core.Host
}

// Arguments for the Lookup RPC.
type LookupArgs struct {
	ID core.BaseID
}

// Arguments for the RoutingTable RPC.
type RoutingTableArgs struct{}

// Arguments for the DataStore RPC.
type DataStoreArgs struct{}

// Arguments for the NodeInfo RPC.
type NodeInfoArgs struct{}

// Arguments for the Bootstrap RPC.
type BootstrapArgs struct {
	Host core.Host
}

// Arguments for the Store RPC.
type PublishArgs struct {
	Data core.DataBytes
}

// Arguments for the Get RPC.
type GetArgs struct {
	KeyID core.BaseID
}
