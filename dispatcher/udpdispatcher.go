package dispatcher

import (
	"crypto/rand"

	"github.com/gsergey418/akademi/core"
	"github.com/gsergey418/akademi/pb"
)

// Dispatcher constants
const (
	RequestIDLength = 32
)

// UDPDispatcher is an implementation of the dispatcher
// interface that dispatches RPC calls over UDP protobuf
// messages.
type UDPDispatcher struct {
	RoutingHeader core.RoutingHeader
}

// The Initialize functions sets the ListenPort of the
// UDPDispatcher.
func (u *UDPDispatcher) Initialize(h core.RoutingHeader) error {
	u.RoutingHeader = h
	return nil
}

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (u *UDPDispatcher) Ping(host core.Host) (core.RoutingHeader, error) {
	req := &pb.BaseMessage{}
	req.Message = &pb.BaseMessage_PingRequest{}
	res, err := u.dispatchUDPMessage(host, req)
	if err != nil {
		return core.RoutingHeader{}, err
	}
	return u.parseRoutingHeader(res), nil
}

// The FindNode function dispatches a FindNode RPC call
// to node located at host.
func (u *UDPDispatcher) FindNode(host core.Host, nodeID core.BaseID) (core.RoutingHeader, []core.RoutingEntry, error) {
	nodes := []core.RoutingEntry{}

	req := &pb.BaseMessage{}
	req.Message = &pb.BaseMessage_FindNodeRequest{FindNodeRequest: &pb.FindNodeRequest{NodeID: nodeID[:]}}
	res, err := u.dispatchUDPMessage(host, req)
	if err != nil {
		return core.RoutingHeader{}, nodes, err
	}

	for _, v := range res.GetFindNodeResponse().GetRoutingEntry() {
		nodes = append(nodes, core.RoutingEntry{
			Host:   core.Host(v.Address),
			NodeID: core.BaseID(v.NodeID),
		})
	}
	return u.parseRoutingHeader(res), nodes, nil
}

// The FindKey function dispatches a FindKey RPC call to
// node located at host.
func (u *UDPDispatcher) FindKey(host core.Host, keyID core.BaseID) (core.RoutingHeader, core.DataBytes, []core.RoutingEntry, error) {
	nodes := []core.RoutingEntry{}

	req := &pb.BaseMessage{}
	req.Message = &pb.BaseMessage_FindKeyRequest{FindKeyRequest: &pb.FindKeyRequest{KeyID: keyID[:]}}
	res, err := u.dispatchUDPMessage(host, req)
	if err != nil {
		return core.RoutingHeader{}, nil, nodes, err
	}

	data := core.DataBytes(res.GetFindKeyResponse().Data)
	for _, v := range res.GetFindKeyResponse().GetRoutingEntry() {
		nodes = append(nodes, core.RoutingEntry{
			Host:   core.Host(v.Address),
			NodeID: core.BaseID(v.NodeID),
		})
	}
	return u.parseRoutingHeader(res), data, nodes, nil
}

// The Store function dispatches a Store RPC call to node
// located at host.
func (u *UDPDispatcher) Store(host core.Host, data core.DataBytes) (core.RoutingHeader, error) {
	req := &pb.BaseMessage{}
	msg := &pb.BaseMessage_StoreRequest{StoreRequest: &pb.StoreRequest{Data: data}}
	req.Message = msg
	res, err := u.dispatchUDPMessage(host, req)
	if err != nil {
		return core.RoutingHeader{}, err
	}
	return u.parseRoutingHeader(res), nil
}

// Returns random RequestID.
func RandomRequestID() [RequestIDLength]byte {
	var o [RequestIDLength]byte
	rand.Read(o[:])
	return o
}
