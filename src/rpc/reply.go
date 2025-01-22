package rpc

import "github.com/gsergey418alt/akademi/core"

// Reply for the Ping RPC.
type PingReply struct {
	Header core.RoutingHeader
}

// Reply for the Lookup RPC.
type LookupReply struct {
	Header       core.RoutingHeader
	RoutingEntry core.RoutingEntry
}

// Reply for the RoutingTable RPC.
type RoutingTableReply struct {
	RoutingTable string
}

// Reply for the NodeInfo RPC.
type NodeInfoReply struct {
	NodeInfo string
}

// Reply for the Bootstrap RPC.
type BootstrapReply struct{}
