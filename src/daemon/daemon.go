package daemon

import (
	"log"
	"strconv"
	"strings"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
	"github.com/gsergey418alt/akademi/listener"
	"github.com/gsergey418alt/akademi/rpc"
)

// getDispatcher returns an instance of the Dispatcher
// interface.
func getDispatcher() core.Dispatcher {
	return &dispatcher.UDPDispatcher{}
}

// getAkademiNode creates and initializes an AkademiNode.
func getAkademiNode(d core.Dispatcher, listenPort core.IPPort, bootstrap bool) (*core.AkademiNode, error) {
	a := &core.AkademiNode{}
	err := a.Initialize(d, listenPort, bootstrap)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Extract port from IP address string
func parseIPPort(listenAddrString string) (core.IPPort, error) {
	listenPort, err := strconv.Atoi(strings.Split(listenAddrString, ":")[1])
	return core.IPPort(listenPort), err
}

// getListener creates an instance of the Listener
// interface.
func getListener(a *core.AkademiNode, listenAddrString string) Listener {
	l := &listener.UDPListener{}
	l.Initialize(listenAddrString, a)
	return l
}

// getRPCserver returns an RPC server instance.
func getRPCServer(a *core.AkademiNode) rpcServer {
	r := &rpc.AkademiNodeRPCServer{}
	r.Initialize(a)
	return r
}

// Wrapper that allows for subroutines to run asynchronously
// as goroutines.
func AsyncWrapper(c chan error, f func() error) {
	c <- f()
}

// Main loop of Akademi.
func Daemon(listenAddrString string, bootstrap bool, rpcListenAddr string) error {
	listenPort, err := parseIPPort(listenAddrString)
	if err != nil {
		return err
	}

	log.Print("Starting Kademlia DHT node on address ", listenAddrString)

	dis := getDispatcher()
	node, err := getAkademiNode(dis, listenPort, bootstrap)
	if err != nil {
		return err
	}
	listener := getListener(node, listenAddrString)

	c := make(chan error)
	go AsyncWrapper(c, listener.Listen)

	if rpcListenAddr != "" {
		rpcServer := getRPCServer(node)
		go AsyncWrapper(c, func() error {
			return rpcServer.Serve(rpcListenAddr)
		})
	}

	select {
	case err := <-c:
		return err
	}
}
