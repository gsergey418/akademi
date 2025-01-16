package listener

import (
	"log"
	"net"

	"github.com/gsergey418alt/akademi/core"
)

type UDPListener struct {
	ListenAddr  *net.UDPAddr
	AkademiNode *core.AkademiNode

	udpReadBuffer [65535]byte
}

func (u *UDPListener) Initialize(listenAddrString string, a *core.AkademiNode) error {
	listenAddr, err := net.ResolveUDPAddr("udp", listenAddrString)
	if err != nil {
		return err
	}
	u.ListenAddr = listenAddr
	u.AkademiNode = a
	return nil
}

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
		log.Print(u.udpReadBuffer[:l])
	}
}
