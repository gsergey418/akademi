package core

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	Initialize(listenPort IPPort) error
	Ping(host Host) (BaseID, error)
	FindNode(host Host, nodeID BaseID) (BaseID, []RoutingEntry, error)
	FindKey(host Host, keyID BaseID) (BaseID, BaseID, []RoutingEntry, error)
	Store(host Host, keyID BaseID, value DataBytes) (BaseID, error)
}
