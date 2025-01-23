package node

import "github.com/gsergey418alt/akademi/core"

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	// To send out RPC requests, the dispatcher needs to
	// be initialized with a core.RoutingHeader.
	Initialize(h core.RoutingHeader) error

	// The Ping function dispatches a Ping RPC call to node
	// located at host.
	Ping(host core.Host) (core.RoutingHeader, error)

	// The FindNode function dispatches a FindNode RPC call
	// to node located at host.
	FindNode(host core.Host, nodeID core.BaseID) (core.RoutingHeader, []core.RoutingEntry, error)

	// The FindKey function dispatches a FindKey RPC call to
	// node located at host.
	FindKey(host core.Host, keyID core.BaseID) (core.RoutingHeader, core.DataBytes, []core.RoutingEntry, error)

	// The Store function dispatches a Store RPC call to node
	// located at host.
	Store(host core.Host, keyID core.BaseID, value core.DataBytes) (core.RoutingHeader, error)
}
