package core

// The Dispatcher interface is responsible for dispatching
// RPC requests to other nodes in the network.
type Dispatcher interface {
	Ping(addr string) (NodeID, error)
	FindNode(addr string, nodeID NodeID) (NodeID, []RoutingEntry, error)
	FindKey(addr string, keyID KeyID) (NodeID, KeyID, []RoutingEntry, error)
	Store(addr string, keyID KeyID, value DataBytes) (NodeID, error)
}
