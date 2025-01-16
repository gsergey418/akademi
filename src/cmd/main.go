package main

import (
	"log"
	"os"

	"github.com/gsergey418alt/akademi/core"
)

const (
	listenAddr = "0.0.0.0:3865"
)

// getDispatcher returns an instance of the Dispatcher
// interface.
func getDispatcher() core.Dispatcher {
	panic("Function \"getDispatcher\" not implemented")
}

// getAkademiNode creates and initializes an AkademiNode.
func getAkademiNode(listenPort core.IPPort, bootstrap bool) *core.AkademiNode {
	a := &core.AkademiNode{}
	a.Initialize(getDispatcher(), listenPort, bootstrap)
	return a
}

// getListener creates an instance of the Listener
// interface.
func getListener(listenAddr string, bootstrap bool) Listener {
	panic("Function \"getListener\" not implemented")
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

	l := getListener(listenAddr, parseArgs())

	log.Fatal(l.Listen())
}
