package rpc

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/gsergey418alt/akademi/core"
)

type AkademiNodeRPCServer struct {
	AkademiNode *core.AkademiNode
}

func (a *AkademiNodeRPCServer) Initialize(n *core.AkademiNode) {
	a.AkademiNode = n
}

func (a *AkademiNodeRPCServer) Serve(listenAddr string) error {
	rpc.Register(a)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	return http.Serve(l, nil)
}

func (a *AkademiNodeRPCServer) Ping(args *PingArgs, reply *PingReply) error {
	header, err := a.AkademiNode.Ping(args.host)
	reply.header = header
	return err
}

type PingArgs struct {
	host core.Host
}

type PingReply struct {
	header core.RoutingHeader
}
