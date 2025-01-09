package listener

import (
	"log"
	"net"
)

type Listener struct {
	udpReadBuffer [2048]byte
}

func (l *Listener) Listen(listenAddr string) {
	udpAddr := parseListenAddrToUDP(&listenAddr)

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
