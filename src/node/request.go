package node

import "github.com/gsergey418alt/akademi/core"

// The responseHandler function manages the side effects
// of receiving an RPC response from the Dispatcher.
func (a *AkademiNode) responseHandler(host core.Host, header core.RoutingHeader) {
	r := core.RoutingEntry{
		Host:   host,
		NodeID: header.NodeID,
	}
	a.UpdateRoutingTable(r)
}

// Redefinitions of Dispatcher functions.

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (a *AkademiNode) Ping(host core.Host) (core.RoutingHeader, error) {
	header, err := a.dispatcher.Ping(host)
	if err != nil {
		return header, err
	}
	a.responseHandler(host, header)
	return header, err
}

// The FindNode function dispatches a FindNode RPC call
// to node located at host.
func (a *AkademiNode) FindNode(host core.Host, nodeID core.BaseID) (core.RoutingHeader, []core.RoutingEntry, error) {
	header, nodes, err := a.dispatcher.FindNode(host, nodeID)
	if err != nil {
		return header, nodes, err
	}
	a.responseHandler(host, header)
	for _, r := range nodes {
		a.UpdateRoutingTable(r)
	}
	return header, nodes, err
}

// The FindKey function dispatches a FindKey RPC call to
// node located at host.
func (a *AkademiNode) FindKey(host core.Host, keyID core.BaseID) (core.RoutingHeader, core.DataBytes, []core.RoutingEntry, error) {
	header, data, nodes, err := a.dispatcher.FindKey(host, keyID)
	if err != nil {
		return header, data, nodes, err
	}
	a.responseHandler(host, header)
	for _, r := range nodes {
		a.UpdateRoutingTable(r)
	}
	return header, data, nodes, err
}

// The Store function dispatches a Store RPC call to node
// located at host.
func (a *AkademiNode) Store(host core.Host, keyID core.BaseID, value core.DataBytes) (core.RoutingHeader, error) {
	header, err := a.dispatcher.Store(host, keyID, value)
	if err != nil {
		return header, err
	}
	a.responseHandler(host, header)
	return header, err
}
