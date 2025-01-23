package daemon

import (
	"log"
	"strconv"
	"strings"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
	"github.com/gsergey418alt/akademi/listener"
	"github.com/gsergey418alt/akademi/node"
	"github.com/gsergey418alt/akademi/rpc"
)

// getDispatcher returns an instance of the Dispatcher
// interface.
func getDispatcher() node.Dispatcher {
	return &dispatcher.UDPDispatcher{}
}

// getAkademiNode creates and initializes an AkademiNode.
func getAkademiNode() (*node.AkademiNode, error) {
	a := &node.AkademiNode{}
	return a, nil
}

// Extract port from IP address string
func parseIPPort(listenAddrString string) (core.IPPort, error) {
	listenPort, err := strconv.Atoi(strings.Split(listenAddrString, ":")[1])
	return core.IPPort(listenPort), err
}

// getListener creates an instance of the Listener
// interface.
func getListener(a *node.AkademiNode, listenAddr string) Listener {
	l := &listener.UDPListener{}
	l.Initialize(listenAddr, a)
	return l
}

// getRPCserver returns an RPC server instance.
func getRPCServer(a *node.AkademiNode, listenAddr string) rpcServer {
	r := &rpc.AkademiNodeRPCServer{}
	r.Initialize(a, listenAddr)
	return r
}

// Wrapper that allows for subroutines to run asynchronously
// as goroutines.
func AsyncWrapper(c chan error, f func() error) {
	c <- f()
}

// Main loop of Akademi.
func Daemon(listenAddr string, bootstrap bool, rpcListenAddr string) error {
	listenPort, err := parseIPPort(listenAddr)
	if err != nil {
		return err
	}

	log.Print("Starting Kademlia DHT node...")

	dis := getDispatcher()
	node, err := getAkademiNode()

	if err != nil {
		return err
	}
	listener := getListener(node, listenAddr)

	c := make(chan error)
	go AsyncWrapper(c, listener.Listen)

	if rpcListenAddr != "" {
		rpcServer := getRPCServer(node, rpcListenAddr)
		go AsyncWrapper(c, rpcServer.Serve)
	}

	err = node.Initialize(dis, listenPort, bootstrap)
	if err != nil {
		return err
	}

	select {
	case err := <-c:
		return err
	}
}
