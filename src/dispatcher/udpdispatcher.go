package dispatcher

import (
	"log"
	"net"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/pb"
	"google.golang.org/protobuf/proto"
)

// UDPDispatcher is an implementation of the dispatcher
// interface that dispatches RPC calls over UDP protobuf
// messages.
type UDPDispatcher struct {
	BaseMessageHeader core.BaseMessageHeader
}

// The Initialize functions sets the ListenPort of the
// UDPDispatcher.
func (u *UDPDispatcher) Initialize(h core.BaseMessageHeader) error {
	u.BaseMessageHeader = h
	return nil
}

// dispatchUDPBytes is a function that manages sending
// raw bytes over a UDP connection.
func (u *UDPDispatcher) dispatchUDPBytes(host core.Host, buf []byte) ([]byte, error) {
	conn, err := net.Dial("udp", string(host))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	log.Print("Dispatching request to ", host, ": ", len(buf), " bytes.")
	_, err = conn.Write(buf)
	if err != nil {
		return nil, err
	}
	var udpReadBuffer [65535]byte
	l, err := conn.Read(udpReadBuffer[:])
	if err != nil {
		return nil, err
	}
	log.Print("Response from ", host, ": ", l, " bytes.")
	return udpReadBuffer[:l], nil
}

// dispatchUDPMessage is a function that wraps around
// dispatchUDPBytes, operating on pb.BaseMessage structs.
func (u *UDPDispatcher) dispatchUDPMessage(host core.Host, msg *pb.BaseMessage) (*pb.BaseMessage, error) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	resBytes, err := u.dispatchUDPBytes(host, buf)
	if err != nil {
		return nil, err
	}
	res := &pb.BaseMessage{}
	err = proto.Unmarshal(resBytes, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *UDPDispatcher) parseBaseMessageHeader(msg *pb.BaseMessage) core.BaseMessageHeader {
	return core.BaseMessageHeader{
		NodeID:     core.BaseID(msg.NodeID),
		ListenPort: core.IPPort(msg.ListenPort),
	}
}

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (u *UDPDispatcher) Ping(host core.Host) (core.BaseMessageHeader, error) {
	msg := &pb.BaseMessage{}
	msg.Message = &pb.BaseMessage_PingRequest{}
	res, err := u.dispatchUDPMessage(host, msg)
	if err != nil {
		return core.BaseMessageHeader{}, err
	}
	return u.parseBaseMessageHeader(res), nil
}

// The FindNode function dispatches a FindNode RPC call
// to node located at host.
func (u *UDPDispatcher) FindNode(host core.Host, nodeID core.BaseID) (core.BaseMessageHeader, []core.RoutingEntry, error) {
	panic("Function FindNode not implemented.")
}

// The FindKey function dispatches a FindKey RPC call to
// node located at host.
func (u *UDPDispatcher) FindKey(host core.Host, keyID core.BaseID) (core.BaseMessageHeader, core.BaseID, []core.RoutingEntry, error) {
	panic("Function FindKey not implemented.")
}

// The Store function dispatches a Store RPC call to node
// located at host.
func (u *UDPDispatcher) Store(host core.Host, keyID core.BaseID, value core.DataBytes) (core.BaseMessageHeader, error) {
	panic("Function Store not implemented.")
}
