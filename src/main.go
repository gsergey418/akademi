package main

import (
	"log"

	"github.com/gsergey418alt/akademi/listener"
)

const (
	listenAddr = "0.0.0.0:3856"
)

func getListener() listener.Listener {
	return &listener.UDPListener{}
}

func main() {
	log.Print("Starting Kademlia DHT node on address ", listenAddr)

	l := getListener()

	log.Fatal(l.Listen(listenAddr))
}
