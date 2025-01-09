package listener

import (
	"log"
	"net"
	"strconv"
	"strings"
)

func parseListenAddrToUDP(listenAddr *string) net.UDPAddr {
	var octetList [4]byte
	splitAddr := strings.Split(*listenAddr, ":")
	ipAddr := splitAddr[0]
	ipPort := splitAddr[1]

	for i, octet := range strings.Split(ipAddr, ".") {
		octetValue, err := strconv.Atoi(octet)
		if err != nil {
			log.Fatal(err)
		}

		octetList[i] = byte(octetValue)
	}

	portNumber, err := strconv.Atoi(ipPort)
	if err != nil {
		log.Fatal(err)
	}

	udpAddr := net.UDPAddr{
		IP:   net.IPv4(octetList[0], octetList[1], octetList[2], octetList[3]),
		Port: int(portNumber),
		Zone: "",
	}

	return udpAddr
}
