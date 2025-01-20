package core

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	// To send out RPC requests, the dispatcher needs to
	// be initialized with a RoutingHeader.
	Initialize(h RoutingHeader) error

	// The Ping function dispatches a Ping RPC call to node
	// located at host.
	Ping(host Host) (RoutingHeader, error)

	// The FindNode function dispatches a FindNode RPC call
	// to node located at host.
	FindNode(host Host, nodeID BaseID) (RoutingHeader, []RoutingEntry, error)

	// The FindKey function dispatches a FindKey RPC call to
	// node located at host.
	FindKey(host Host, keyID BaseID) (RoutingHeader, DataBytes, []RoutingEntry, error)

	// The Store function dispatches a Store RPC call to node
	// located at host.
	Store(host Host, keyID BaseID, value DataBytes) (RoutingHeader, error)
}

// Type RoutingHeader contains routing information passed
// with every akademi request and response.
type RoutingHeader struct {
	NodeID     BaseID
	ListenPort IPPort
}
