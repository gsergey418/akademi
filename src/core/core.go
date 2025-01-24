package core

import (
	"fmt"
)

// AkademiNode constants.
const (
	IDLength           = 20
	BucketSize         = 40
	MaxDataLength      = 4096
	Replication        = 5
	ConcurrentRequests = 3
	Bootstraps         = 2
)

// Akademi uses node and key IDs, whose length is defined
// in bytes by IDLength.
type BaseID [IDLength]byte

// Separate IPPort type because the IP address is
// identified by receiving node.
type IPPort uint16

// Host is used to identify node's IP address and
// port.
type Host string

// DataBytes is a type for values to be stored in akademi
// nodes.
type DataBytes []byte

// RoutingEntry is a structure that stores routing
// information about an akademi node.
type RoutingEntry struct {
	Host   Host
	NodeID BaseID
}

// Type RoutingHeader contains routing information passed
// with every akademi request and response.
type RoutingHeader struct {
	NodeID     BaseID
	ListenPort IPPort
}

// Pretty-print core.RoutingEntry.
func (r RoutingEntry) String() string {
	return fmt.Sprintf("%s@%s", r.NodeID, r.Host)
}

// Pretty-print core.DataBytes.
func (d DataBytes) String() string {
	return string(d)
}
