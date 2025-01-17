package listener

import (
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

// Parse listenAddrString and set AkademiNode
func (u *UDPListener) Initialize(listenAddrString string, a *core.AkademiNode) error {
	listenAddr, err := net.ResolveUDPAddr("udp", listenAddrString)
	if err != nil {
		return err
	}
	u.ListenAddr = listenAddr
	u.AkademiNode = a
	return nil
}

// Opens a UDP socket on UDPListener.ListenAddr
func (u *UDPListener) Listen() error {
	conn, err := net.ListenUDP("udp", u.ListenAddr)
	u.udpConn = conn
	if err != nil {
		return err
	}

	for {
		l, remoteAddr, err := u.udpConn.ReadFromUDP(u.udpReadBuffer[:])
		if err != nil {
			log.Print(err)
		}
		err = u.handleUDPMessage(remoteAddr, u.udpReadBuffer[:l])
		if err != nil {
			log.Print(err)
		}
	}
}

// Sends bytes to remoteAddr over UDP.
func (u *UDPListener) sendUDPBytes(remoteAddr *net.UDPAddr, buf []byte) error {
	log.Print("Writing response to ", remoteAddr, ": ", len(buf), " bytes.")
	_, err := u.udpConn.WriteTo(buf, remoteAddr)
	return err
}

// Sends pb.BaseMessage to remoteAddr
func (u *UDPListener) sendUDPMessage(remoteAddr *net.UDPAddr, res, req *pb.BaseMessage) error {
	u.populateDefaultResponse(res, req)
	res.Message = &pb.BaseMessage_PingResponse{}
	resBytes, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	err = u.sendUDPBytes(remoteAddr, resBytes)
	return err
}

// Populates the default response protobuf
func (u *UDPListener) populateDefaultResponse(res, req *pb.BaseMessage) {
	res.RequestID = req.RequestID
	res.ListenPort = uint32(u.ListenAddr.Port)
	res.NodeID = u.AkademiNode.NodeID[:]
}

// Multiplexer for the BaseMessage type.
func (u *UDPListener) reqMux(remoteAddr *net.UDPAddr, req *pb.BaseMessage) error {
	switch {
	case req.GetPingRequest() != nil:
		res := &pb.BaseMessage{}
		return u.sendUDPMessage(remoteAddr, res, req)
	}
	return nil
}

// Handle a slice of bytes as a UDP message
func (u *UDPListener) handleUDPMessage(remoteAddr *net.UDPAddr, buf []byte) error {
	log.Print("Message from ", remoteAddr, ": ", len(buf), " bytes.")
	req := &pb.BaseMessage{}
	err := proto.Unmarshal(buf, req)
	if err != nil {
		return err
	}
	err = u.reqMux(remoteAddr, req)
	if err != nil {
		return err
	}
	return nil
}
