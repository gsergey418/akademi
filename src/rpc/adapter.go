package rpc

import (
	"fmt"
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

// Sends a lookup request for args.ID.
func (s *AkademiNodeRPCServer) Lookup(args *LookupArgs, reply *LookupReply) error {
	nodes, err := s.AkademiNode.Lookup(args.ID, 1)
	empty := core.RoutingEntry{}
	if err != nil {
		return err
	}
	if len(nodes) < 1 || nodes[0] == empty {
		return fmt.Errorf("Could not lookup ID %s.", args.ID)
	}
	reply.RoutingEntry = nodes[0]
	return nil
}

// Gets the node routing table as a string.
func (s *AkademiNodeRPCServer) RoutingTable(args *RoutingTableArgs, reply *RoutingTableReply) error {
	reply.RoutingTable = s.AkademiNode.RoutingTableString()
	return nil
}

// Args for the Ping RPC.
type PingArgs struct {
	Host core.Host
}

// Reply for the Ping RPC.
type PingReply struct {
	Header core.RoutingHeader
}

// Args for the Lookup RPC.
type LookupArgs struct {
	ID core.BaseID
}

// Reply for the Lookup RPC.
type LookupReply struct {
	Header       core.RoutingHeader
	RoutingEntry core.RoutingEntry
}

// Args for the RoutingTable RPC.
type RoutingTableArgs struct{}

// Reply for the RoutingTable RPC.
type RoutingTableReply struct {
	RoutingTable string
}
