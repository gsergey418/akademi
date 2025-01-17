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
func (u *UDPListener) sendUDPResponse(remoteAddr *net.UDPAddr, buf []byte) error {
	log.Print("Writing response to ", remoteAddr, ": ", len(buf), " bytes.")
	_, err := u.udpConn.WriteTo(buf, remoteAddr)
	return err
}

// Populates the default response protobuf
func (u *UDPListener) populateDefaultResponse(res, req *pb.BaseMessage) {
	res.RequestID = req.RequestID
	res.ListenPort = uint32(u.ListenAddr.Port)
	res.NodeID = u.AkademiNode.NodeID[:]
}

// Multiplexer for the BaseMessage type.
func (u *UDPListener) msgMux(remoteAddr *net.UDPAddr, msg *pb.BaseMessage) error {
	switch {
	case msg.GetPingRequest() != nil:
		res := &pb.BaseMessage{}
		u.populateDefaultResponse(res, msg)
		res.Message = &pb.BaseMessage_PingResponse{}
		resBytes, err := proto.Marshal(res)
		if err != nil {
			return err
		}
		err = u.sendUDPResponse(remoteAddr, resBytes)
		if err != nil {
			return &net.AddrError{}
		}
	}
	return nil
}

// Handle a slice of bytes as a UDP message
func (u *UDPListener) handleUDPMessage(remoteAddr *net.UDPAddr, buf []byte) error {
	log.Print("Message from ", remoteAddr, ": ", len(buf), " bytes.")
	msg := &pb.BaseMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		return err
	}
	err = u.msgMux(remoteAddr, msg)
	if err != nil {
		return err
	}
	return nil
}
