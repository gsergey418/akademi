package rpc

import (
	"log"
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
func (s *AkademiNodeRPCServer) Initialize(n *core.AkademiNode, listenAddr string) {
	s.ListenAddr = listenAddr
	s.AkademiNode = n
}

// Main loop of the RPC server.
func (s *AkademiNodeRPCServer) Serve() error {
	rpc.Register(s)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer l.Close()
	log.Print("RPC server started on address ", s.ListenAddr, ".")
	return http.Serve(l, nil)
}

// Sends a ping request to args.Host.
func (s *AkademiNodeRPCServer) Ping(args *PingArgs, reply *PingReply) error {
	header, err := s.AkademiNode.Ping(args.Host)
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
