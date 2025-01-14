package main

import (
	"log"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/listener"
)

const (
	listenAddr = "127.0.0.1:3856"
)

func getAkademiNode() *core.AkademiNode {
	return &core.AkademiNode{}
}

func getRPCAdapater() listener.ListenerAdapter {
	return &listener.AkademiNodeRPCAdapter{AkademiNode: getAkademiNode()}
}

func getListener() listener.Listener {
	return &listener.RPCListener{RPCAdapter: getRPCAdapater()}
}

func main() {
	log.Print("Starting Kademlia DHT node on address ", listenAddr)

	l := getListener()

	log.Fatal(l.Listen(listenAddr))
}
