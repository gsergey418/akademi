package core

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	Initialize(h RoutingHeader) error
	Ping(host Host) (RoutingHeader, error)
	FindNode(host Host, nodeID BaseID) (RoutingHeader, []RoutingEntry, error)
	FindKey(host Host, keyID BaseID) (RoutingHeader, BaseID, []RoutingEntry, error)
	Store(host Host, keyID BaseID, value DataBytes) (RoutingHeader, error)
}

// Type RoutingHeader contains routing information passed
// with every akademi request and response.
type RoutingHeader struct {
	NodeID     BaseID
	ListenPort IPPort
}
