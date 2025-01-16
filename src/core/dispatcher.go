package core

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	Initialize(listenPort IPPort)
	Ping(dispatchAddr DispatchAddr) (BaseID, error)
	FindNode(dispatchAddr DispatchAddr, nodeID BaseID) (BaseID, []RoutingEntry, error)
	FindKey(dispatchAddr DispatchAddr, keyID BaseID) (BaseID, BaseID, []RoutingEntry, error)
	Store(dispatchAddr DispatchAddr, keyID BaseID, value DataBytes) (BaseID, error)
}
