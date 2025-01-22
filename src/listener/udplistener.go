package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/pb"
	"google.golang.org/protobuf/proto"
)

// The UDPListener struct receives protocol buffers
// encoded data via a UDP socket.
type UDPListener struct {
	ListenAddr  *net.UDPAddr
	AkademiNode *core.AkademiNode

	udpReadBuffer [65535]byte
	udpConn       *net.UDPConn
}

// Parse listenAddrString and set AkademiNode.
func (u *UDPListener) Initialize(listenAddrStr string, a *core.AkademiNode) error {
	listenAddr, err := net.ResolveUDPAddr("udp", listenAddrStr)
	if err != nil {
		return err
	}
	u.ListenAddr = listenAddr
	u.AkademiNode = a
	return nil
}

// Opens a UDP socket on UDPListener.ListenAddr.
func (u *UDPListener) Listen() error {
	conn, err := net.ListenUDP("udp", u.ListenAddr)
	u.udpConn = conn
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Print("UDP listener started on address ", u.ListenAddr, ".")

	for {
		l, host, err := u.udpConn.ReadFromUDP(u.udpReadBuffer[:])
		if err != nil {
			log.Print(err)
		}
		err = u.handleUDPMessage(host, u.udpReadBuffer[:l])
		if err != nil {
			log.Print(err)
		}
	}
}

// Populates the default response protobuf.
func (u *UDPListener) populateDefaultResponse(res, req *pb.BaseMessage) {
	res.RequestID = req.RequestID
	res.ListenPort = uint32(u.ListenAddr.Port)
	res.NodeID = u.AkademiNode.NodeID[:]
}

// Multiplexer for the BaseMessage type.
func (u *UDPListener) reqMux(host *net.UDPAddr, req *pb.BaseMessage) error {
	switch {
	case req.GetPingRequest() != nil:
		res := &pb.BaseMessage{}
		res.Message = &pb.BaseMessage_PingResponse{}
		return u.sendUDPMessage(host, res, req)
	case req.GetFindNodeRequest() != nil:
		res := &pb.BaseMessage{}
		msg := &pb.FindNodeResponse{}
		nodes, err := u.AkademiNode.GetClosestNodes(core.BaseID(req.GetFindNodeRequest().NodeID), core.BucketSize)
		if err != nil {
			return err
		}
		for _, v := range nodes {
			msg.RoutingEntry = append(msg.RoutingEntry, &pb.RoutingEntry{
				Address: string(v.Host),
				NodeID:  v.NodeID[:],
			})
		}
		res.Message = &pb.BaseMessage_FindNodeResponse{FindNodeResponse: msg}
		return u.sendUDPMessage(host, res, req)
	}
	return nil
}

// Handle a slice of bytes as a UDP message.
func (u *UDPListener) handleUDPMessage(host *net.UDPAddr, buf []byte) error {
	log.Print("Message from ", host, ": ", len(buf), " bytes.")
	req := &pb.BaseMessage{}
	err := proto.Unmarshal(buf, req)
	if err != nil {
		return err
	}

	if bytes.Equal(req.NodeID, u.AkademiNode.NodeID[:]) {
		return fmt.Errorf("Request from self, dropping.")
	}

	r := core.RoutingEntry{
		Host:   core.Host(fmt.Sprintf("%s:%d", host.IP, req.ListenPort)),
		NodeID: core.BaseID(req.NodeID),
	}
	err = u.AkademiNode.UpdateRoutingTable(r)
	if err != nil {
		return err
	}

	err = u.reqMux(host, req)
	if err != nil {
		return err
	}

	return nil
}
