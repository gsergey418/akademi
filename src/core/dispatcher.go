package core

// The Dispatcher interface is responsible for dispatching
// RPC requests to other nodes in the network.
type Dispatcher interface {
	Ping(addr string) (BaseID, error)
	FindNode(addr string, nodeID BaseID) (BaseID, []RoutingEntry, error)
	FindKey(addr string, keyID BaseID) (BaseID, BaseID, []RoutingEntry, error)
	Store(addr string, keyID BaseID, value DataBytes) (BaseID, error)
}
