package dispatcher

import (
	"github.com/gsergey418alt/akademi/core"
)

// UDPDispatcher is an implementation of the dispatcher
// interface that dispatches RPC calls over UDP protobuf
// messages.
type UDPDispatcher struct {
	ListenPort core.IPPort
}

// The Initialize functions sets the ListenPort of the
// UDPDispatcher.
func (u *UDPDispatcher) Initialize(listenPort core.IPPort) error {
	u.ListenPort = listenPort
	return nil
}

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (u *UDPDispatcher) Ping(host core.Host) (core.BaseID, error) {
	panic("Function Ping not implemented.")
}

// The FindNode function dispatches a FindNode RPC call
// to node located at host.
func (u *UDPDispatcher) FindNode(host core.Host, nodeID core.BaseID) (core.BaseID, []core.RoutingEntry, error) {
	panic("Function FindNode not implemented.")
}

// The FindKey function dispatches a FindKey RPC call to
// node located at host.
func (u *UDPDispatcher) FindKey(host core.Host, keyID core.BaseID) (core.BaseID, core.BaseID, []core.RoutingEntry, error) {
	panic("Function FindKey not implemented.")
}

// The Store function dispatches a Store RPC call to node
// located at host.
func (u *UDPDispatcher) Store(host core.Host, keyID core.BaseID, value core.DataBytes) (core.BaseID, error) {
	panic("Function Store not implemented.")
}
