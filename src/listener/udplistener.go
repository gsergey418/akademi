package listener

import (
	"log"
	"net"

	"github.com/gsergey418alt/akademi/core"
)

// The UDPListener struct receives protocol buffers
// encoded data via a UDP socket.
type UDPListener struct {
	ListenAddr  *net.UDPAddr
	AkademiNode *core.AkademiNode

	udpReadBuffer [65535]byte
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
	if err != nil {
		return err
	}

	for {
		l, err := conn.Read(u.udpReadBuffer[:])
		if err != nil {
			log.Print(err)
		}
		err = u.handleUDPMessage(u.udpReadBuffer[:l])
		if err != nil {
			log.Print(err)
		}
	}
}

// Handle a slice of bytes as a UDP message
func (u *UDPListener) handleUDPMessage(msg []byte) error {
	log.Print(msg)
	return nil
}
