package daemon

import (
	"log"
	"strconv"
	"strings"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
	"github.com/gsergey418alt/akademi/listener"
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

// Wrapper that allows for subroutines to run asynchronously
// as goroutines.
func AsyncWrapper(c chan error, f func() error) {
	c <- f()
}

// Main loop of Akademi.
func Daemon(listenAddrString string, bootstrap bool) error {
	log.Print("Starting Kademlia DHT node on address ", listenAddrString)
	l := getListener(listenAddrString, bootstrap)

	c := make(chan error)
	go AsyncWrapper(c, l.Listen)

	select {
	case err := <-c:
		return err
	}
}
