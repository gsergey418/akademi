package rpc

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/gsergey418alt/akademi/core"
)

// Structure for exposing an RPC API for node control.
type AkademiNodeRPCServer struct {
	ListenAddr  string
	AkademiNode *core.AkademiNode
}

// Sets the underlying listen address and AkademiNode.
func (a *AkademiNodeRPCServer) Initialize(n *core.AkademiNode, listenAddr string) {
	a.ListenAddr = listenAddr
	a.AkademiNode = n
}

// Main loop of the RPC server.
func (a *AkademiNodeRPCServer) Serve() error {
	rpc.Register(a)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", a.ListenAddr)
	if err != nil {
		return err
	}
	defer l.Close()
	return http.Serve(l, nil)
}

// Sends a ping request to args.Host.
func (a *AkademiNodeRPCServer) Ping(args *PingArgs, reply *PingReply) error {
	header, err := a.AkademiNode.Ping(args.Host)
	reply.Header = header
	return err
}

// Args for the Ping RPC.
type PingArgs struct {
	Host core.Host
}

// Reply for the Ping RPC.
type PingReply struct {
	Header core.RoutingHeader
}
