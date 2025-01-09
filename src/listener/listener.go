package listener

import (
	"log"
	"net"
)

type Listener struct {
	udpReadBuffer [2048]byte
}

// Listen opens a UDP socket on the provided address in the format "127.0.0.1:443". This is the main loop of the program.
func (l *Listener) Listen(listenAddr string) {
	udpAddr := ParseListenAddrToUDP(&listenAddr)

	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		n, addr, err := conn.ReadFromUDP(l.udpReadBuffer[:])
		if err != nil {
			continue
		}
		log.Print("Message from ", addr, ": ", string(l.udpReadBuffer[:n]))
	}
}
