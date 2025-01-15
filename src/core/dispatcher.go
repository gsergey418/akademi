package core

// The Dispatcher interface is responsible for dispatching
// RPC requests to other nodes in the network.
type Dispatcher interface {
	Ping(addr string) error
	FindNode(addr string, nodeID NodeID) ([]RoutingEntry, error)
	FindKey(addr string, keyID KeyID) (KeyID, []RoutingEntry, error)
	Store(addr string, keyID KeyID, value DataBytes) error
}
