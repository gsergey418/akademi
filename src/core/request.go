package core

// The responseHandler function manages the side effects
// of receiving an RPC response from the Dispatcher.
func (a *AkademiNode) responseHandler(host Host, header RoutingHeader) {
	r := RoutingEntry{
		Host:   host,
		NodeID: header.NodeID,
	}
	a.UpdateRoutingTable(r)
}

// Redefinitions of Dispatcher functions.

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (a *AkademiNode) Ping(host Host) (RoutingHeader, error) {
	header, err := a.dispatcher.Ping(host)
	a.responseHandler(host, header)
	return header, err
}

// The FindNode function dispatches a FindNode RPC call
// to node located at host.
func (a *AkademiNode) FindNode(host Host, nodeID BaseID) (RoutingHeader, []RoutingEntry, error) {
	header, nodes, err := a.dispatcher.FindNode(host, nodeID)
	a.responseHandler(host, header)
	for _, r := range nodes {
		a.UpdateRoutingTable(r)
	}
	return header, nodes, err
}

// The FindKey function dispatches a FindKey RPC call to
// node located at host.
func (a *AkademiNode) FindKey(host Host, keyID BaseID) (RoutingHeader, DataBytes, []RoutingEntry, error) {
	header, data, nodes, err := a.dispatcher.FindKey(host, keyID)
	a.responseHandler(host, header)
	for _, r := range nodes {
		a.UpdateRoutingTable(r)
	}
	return header, data, nodes, err
}

// The Store function dispatches a Store RPC call to node
// located at host.
func (a *AkademiNode) Store(host Host, keyID BaseID, value DataBytes) (RoutingHeader, error) {
	header, err := a.dispatcher.Store(host, keyID, value)
	a.responseHandler(host, header)
	return header, err
}
