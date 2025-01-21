package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
	"github.com/gsergey418alt/akademi/listener"
)

const (
	listenAddrString = "0.0.0.0:3865"
)

// getDispatcher returns an instance of the Dispatcher
// interface.
func getDispatcher() core.Dispatcher {
	return &dispatcher.UDPDispatcher{}
}

// getAkademiNode creates and initializes an AkademiNode.
func getAkademiNode(listenPort core.IPPort, bootstrap bool) *core.AkademiNode {
	a := &core.AkademiNode{}
	err := a.Initialize(getDispatcher(), listenPort, bootstrap)
	if err != nil {
		log.Fatal(err)
	}
	return a
}

// Extract port from IP address string
func parseIPPort(listenAddrString string) (core.IPPort, error) {
	listenPort, err := strconv.Atoi(strings.Split(listenAddrString, ":")[1])
	return core.IPPort(listenPort), err
}

// getListener creates an instance of the Listener
// interface.
func getListener(listenAddrString string, bootstrap bool) Listener {
	l := &listener.UDPListener{}
	listenPort, err := parseIPPort(listenAddrString)
	if err != nil {
		log.Fatal(err)
	}
	l.Initialize(listenAddrString, getAkademiNode(listenPort, bootstrap))
	return l
}

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() (cmd string, bootstrap bool) {
	bootstrap = true
	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments, please provide a command.\n")
		os.Exit(1)
	}
	cmd = os.Args[1]
	for _, arg := range os.Args[2:] {
		switch arg {
		case "--no-bootstrap":
			bootstrap = false
		}
	}
	return
}

// Akademi entrypoint.
func main() {
	cmd, bootstrap := parseArgs()
	switch cmd {
	case "daemon":
		log.Print("Starting Kademlia DHT node on address ", listenAddrString)
		l := getListener(listenAddrString, bootstrap)

		log.Fatal(l.Listen())
	default:
		fmt.Print("Command \"", cmd, "\" not found.\n")
		os.Exit(1)
	}
}
