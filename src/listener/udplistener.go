package listener

import (
	"log"
	"net"
	"runtime"

	"github.com/gsergey418alt/akademi/node"
)

// The UDPListener struct receives protocol buffers
// encoded data via a UDP socket.
type UDPListener struct {
	ListenAddr  *net.UDPAddr
	AkademiNode *node.AkademiNode

	udpReadBuffer [65535]byte
	udpConn       *net.UDPConn
}

// Parse listenAddrString and set AkademiNode.
func (u *UDPListener) Initialize(listenAddrStr string, a *node.AkademiNode) error {
	listenAddr, err := net.ResolveUDPAddr("udp", listenAddrStr)
	if err != nil {
		return err
	}
	u.ListenAddr = listenAddr
	u.AkademiNode = a
	return nil
}

// Sturcture that contains the source host and message bytes.
type UDPBytes struct {
	host *net.UDPAddr
	data []byte
}

// Main listener goroutine.
func (u *UDPListener) ListenUDP(msgChan chan UDPBytes, errChan chan error) {
	for {
		l, host, err := u.udpConn.ReadFromUDP(u.udpReadBuffer[:])
		if err != nil {
			log.Print(err)
		}
		msgChan <- UDPBytes{host, u.udpReadBuffer[:l]}
	}
}

// Handler goroutine.
func (u *UDPListener) UDPWorker(msgChan chan UDPBytes, errChan chan error) {
	for {
		udpBytes := <-msgChan
		err := u.handleUDPMessage(udpBytes.host, udpBytes.data)
		if err != nil {
			log.Print(err)
		}
	}
}

// Opens a UDP socket on UDPListener.ListenAddr.
func (u *UDPListener) Listen() error {
	conn, err := net.ListenUDP("udp", u.ListenAddr)
	u.udpConn = conn
	if err != nil {
		return err
	}
	defer conn.Close()

	cpuCount := runtime.NumCPU()
	log.Print("UDP listener started on address ", u.ListenAddr, ". Creating ", cpuCount, " workers.")

	errChan := make(chan error)
	msgChan := make(chan UDPBytes, 1024)

	go u.ListenUDP(msgChan, errChan)
	for i := 0; i < cpuCount; i++ {
		go u.UDPWorker(msgChan, errChan)
	}

	log.Fatal(<-errChan)
	return nil
}
