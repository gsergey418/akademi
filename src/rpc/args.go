package rpc

import (
	"github.com/gsergey418alt/akademi/core"
)

// Args for the Ping RPC.
type PingArgs struct {
	Host core.Host
}

// Args for the Lookup RPC.
type LookupArgs struct {
	ID core.BaseID
}

// Args for the RoutingTable RPC.
type RoutingTableArgs struct{}

// Args for the DataStore RPC.
type DataStoreArgs struct{}

// Args for the NodeInfo RPC.
type NodeInfoArgs struct{}

// Args for the Bootstrap RPC.
type BootstrapArgs struct {
	Host core.Host
}
