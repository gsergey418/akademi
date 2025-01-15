package main

import (
	"log"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
	"github.com/gsergey418alt/akademi/listener"
)

const (
	listenAddr = "127.0.0.1:3856"
)

func getDispatcher() core.Dispatcher {
	return &dispatcher.RPCDispatcher{}
}

func getAkademiNode() *core.AkademiNode {
	a := &core.AkademiNode{Dispatcher: getDispatcher()}
	a.Initialize()
	return a
}

func getListener() Listener {
	l := &listener.RPCListener{}
	l.Initialize(getAkademiNode())
	return l
}

func main() {
	log.Print("Starting Kademlia DHT node on address ", listenAddr)

	l := getListener()

	log.Fatal(l.Listen(listenAddr))
}
