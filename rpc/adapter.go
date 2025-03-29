package rpc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gsergey418/akademi/core"
	"github.com/gsergey418/akademi/node"
)

// Structure for exposing an RPC API for node control.
type AkademiNodeRPCServer struct {
	ListenAddr  string
	AkademiNode *node.AkademiNode
}

// Sets the underlying listen address and AkademiNode.
func (s *AkademiNodeRPCServer) Initialize(n *node.AkademiNode, listenAddr string) {
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
		return fmt.Errorf("could not lookup ID %s", args.ID)
	}
	reply.RoutingEntry = nodes[0]
	return nil
}

// Gets the node routing table as a string.
func (s *AkademiNodeRPCServer) RoutingTable(args *RoutingTableArgs, reply *RoutingTableReply) error {
	reply.RoutingTable = s.AkademiNode.RoutingTableString()
	return nil
}

// Gets the node datastore as a string.
func (s *AkademiNodeRPCServer) DataStore(args *DataStoreArgs, reply *DataStoreReply) error {
	reply.DataStore = s.AkademiNode.DataStoreString()
	return nil
}

// Gets the node information.
func (s *AkademiNodeRPCServer) NodeInfo(args *NodeInfoArgs, reply *NodeInfoReply) error {
	reply.NodeInfo = s.AkademiNode.NodeInfo()
	return nil
}

// Sends a bootstrap request (self-lookup) to Host. Useful
// if daemon was started with the "--no-bootstrap" flag.
func (s *AkademiNodeRPCServer) Bootstrap(args *BootstrapArgs, reply *BootstrapReply) error {
	_, _, err := s.AkademiNode.FindNode(args.Host, s.AkademiNode.NodeID)
	return err
}

// Publishes data to the DHT. Finds the best fitting nodes
// and replicates the data to them.
func (s *AkademiNodeRPCServer) Publish(args *PublishArgs, reply *PublishReply) error {
	keyID, err := s.AkademiNode.Publish(args.Data)
	reply.KeyID = keyID
	return err
}

// Gets the node datastore as a string.
func (s *AkademiNodeRPCServer) Get(args *GetArgs, reply *GetReply) error {
	data, err := s.AkademiNode.KeyLookup(args.KeyID)
	reply.Data = data
	return err
}
