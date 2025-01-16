package core

// The Dispatcher interface is responsible for dispatching
// requests to other nodes in the network.
type Dispatcher interface {
	Initialize(self *RoutingEntry)
	Ping(dispatchAddr ListenAddr) (BaseID, error)
	FindNode(dispatchAddr ListenAddr, nodeID BaseID) (BaseID, []RoutingEntry, error)
	FindKey(dispatchAddr ListenAddr, keyID BaseID) (BaseID, BaseID, []RoutingEntry, error)
	Store(dispatchAddr ListenAddr, keyID BaseID, value DataBytes) (BaseID, error)
}
