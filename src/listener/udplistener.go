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
	_, err := u.udpConn.WriteTo(buf, remoteAddr)
	return err
}

// Handle a slice of bytes as a UDP message
func (u *UDPListener) handleUDPMessage(remoteAddr *net.UDPAddr, buf []byte) error {
	msg := &pb.BaseMessage{}
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		log.Print(err)
	}
	switch {
	case msg.GetPingRequest() != nil:
		res := &pb.BaseMessage{}
		res.Message = &pb.BaseMessage_PingResponse{}
		resBytes, err := proto.Marshal(res)
		if err != nil {
			log.Print(err)
		}
		err = u.sendUDPResponse(remoteAddr, resBytes)
		if err != nil {
			log.Print(err)
		}
	}
	log.Print("Message from ", remoteAddr, ": ", msg)
	return nil
}
