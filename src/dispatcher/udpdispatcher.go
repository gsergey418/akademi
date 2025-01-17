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
	ListenPort core.IPPort
}

// The Initialize functions sets the ListenPort of the
// UDPDispatcher.
func (u *UDPDispatcher) Initialize(listenPort core.IPPort) error {
	u.ListenPort = listenPort
	return nil
}

// dispatchUDPMessage is a wrapper function that manages the
// UDP connection.
func (u *UDPDispatcher) dispatchUDPMessage(host core.Host, buf []byte) ([]byte, error) {
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

// The Ping function dispatches a Ping RPC call to node
// located at host.
func (u *UDPDispatcher) Ping(host core.Host) (core.BaseID, error) {
	msg := &pb.BaseMessage{}
	msg.Message = &pb.BaseMessage_PingRequest{}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return core.BaseID{}, err
	}
	resBytes, err := u.dispatchUDPMessage(host, buf)
	if err != nil {
		return core.BaseID{}, err
	}
	res := &pb.BaseMessage{}
	err = proto.Unmarshal(resBytes, res)
	if err != nil {
		return core.BaseID{}, err
	}
	return core.BaseID(res.NodeID), nil
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
