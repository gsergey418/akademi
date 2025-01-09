package main

import (
	"log"

	"github.com/gsergey418/kademlia/listener"
)

const (
	listenAddr = "0.0.0.0:3856"
)

func main() {
	log.Print("Starting Kademlia DHT node on address ", listenAddr)

	l := listener.Listener{}

	l.Listen(listenAddr)
}
