package core

import (
	"fmt"
	"time"
)

// AkademiNode constants.
const (
	IDLength           = 20
	BucketSize         = 20
	MaxDataLength      = 4096
	Replication        = 10
	ConcurrentRequests = 3
	Bootstraps         = 2
	DataExpire         = time.Hour
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

// DataContainer is a structure that represents an entry in
// the node's data storage.
type DataContainer struct {
	Data       DataBytes
	LastAccess time.Time
}

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

// Pretty-print core.DataBytes.
func (d DataContainer) String() string {
	return d.Data.String()
}
