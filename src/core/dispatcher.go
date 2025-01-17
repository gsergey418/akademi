package core

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	Initialize(h BaseMessageHeader) error
	Ping(host Host) (BaseMessageHeader, error)
	FindNode(host Host, nodeID BaseID) (BaseMessageHeader, []RoutingEntry, error)
	FindKey(host Host, keyID BaseID) (BaseMessageHeader, BaseID, []RoutingEntry, error)
	Store(host Host, keyID BaseID, value DataBytes) (BaseMessageHeader, error)
}

// Type BaseMessageHeader contains routing information passed
// with every akademi request and response.
type BaseMessageHeader struct {
	NodeID     BaseID
	ListenPort IPPort
}
