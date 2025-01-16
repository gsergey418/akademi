package main

import (
	"log"
	"os"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
	"github.com/gsergey418alt/akademi/listener"
)

const (
	listenAddr = "0.0.0.0:3856"
)

// getDispatcher returns an instance of the Dispatcher
// interface.
func getDispatcher() core.Dispatcher {
	return &dispatcher.RPCDispatcher{}
}

// getAkademiNode creates and initializes an AkademiNode.
func getAkademiNode(bootstrap bool) *core.AkademiNode {
	a := &core.AkademiNode{Dispatcher: getDispatcher()}
	a.Initialize(bootstrap)
	return a
}

// getListener creates an instance of the Listener
// interface.
func getListener(bootstrap bool) Listener {
	l := &listener.RPCListener{}
	l.Initialize(getAkademiNode(bootstrap))
	return l
}

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() (bootstrap bool) {
	bootstrap = true
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--no-bootstrap":
			bootstrap = false
		}
	}
	return
}

// Akademi entrypoint.
func main() {
	log.Print("Starting Kademlia DHT node on address ", listenAddr)

	l := getListener(parseArgs())

	log.Fatal(l.Listen(listenAddr))
}
