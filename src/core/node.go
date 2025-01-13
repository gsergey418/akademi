package core

type NodeID [256]byte
type KeyID [256]byte
type DataBytes []byte

// RoutingEntry is a structure that stores routing
// information about an akademi node.
type RoutingEntry struct {
	IPAddress string
	IPPort    uint8
	NodeID    NodeID
}

// AkademiNode is a structure containing the core kademlia
// logic.
type AkademiNode struct {
	NodeID        NodeID
	KeyValueStore map[KeyID][]byte

	RoutingTable [256][20]RoutingEntry
}
