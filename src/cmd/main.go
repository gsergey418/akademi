package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/daemon"
	akademiRPC "github.com/gsergey418alt/akademi/rpc"
)

const (
	// Never expose RPC to the public! For docker.
	defaultRpcListenAddr  = "127.0.0.1:3855"
	defaultNodeListenAddr = "0.0.0.0:3865"
)

// Settings populated by parseArgs()
type cmdOptions struct {
	cmd            string
	target         string
	rpcListenAddr  string
	nodeListenAddr string
	bootstrap      bool
	bootstrapList  []core.Host
}

// Global instance of cmdOptions
var opts cmdOptions

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() {
	// Commands with no positional arguments.
	noPosArgs := map[string]bool{
		"daemon":        true,
		"routing_table": true,
		"datastore":     true,
		"info":          true,
	}

	// Default values for command-line options.
	opts.bootstrap = true
	opts.bootstrapList = []core.Host{
		"akademi_bootstrap_1:3865",
		"akademi_bootstrap_2:3865",
		"akademi_bootstrap_3:3865",
	}
	opts.nodeListenAddr = defaultNodeListenAddr
	opts.rpcListenAddr = defaultRpcListenAddr

	// Start of argument parsing.
	argLen := len(os.Args)
	if argLen < 2 {
		fmt.Print("Not enough arguments, please provide a command.\n")
		os.Exit(1)
	}
	optStart, optStop := 2, argLen
	opts.cmd = os.Args[1]
	if _, ok := noPosArgs[opts.cmd]; !ok {
		if argLen < 3 {
			fmt.Print("Not enough arguments, please provide a positional argument.\n")
			os.Exit(1)
		}
		opts.target = os.Args[argLen-1]
		optStop--
	}
	for i := optStart; i < optStop; i++ {
		switch os.Args[i] {
		case "--no-bootstrap":
			opts.bootstrap = false
		case "--bootstrap-nodes":
			opts.bootstrapList = []core.Host{}
			for _, v := range strings.Split(os.Args[i+1], ",") {
				opts.bootstrapList = append(opts.bootstrapList, core.Host(v))
			}
			i++
		case "--rpc-addr":
			opts.rpcListenAddr = os.Args[i+1]
			i++
		case "--listen-addr":
			opts.nodeListenAddr = os.Args[i+1]
			i++
		default:
			fmt.Print("Unknown argument: \"", os.Args[i], "\".\n")
			os.Exit(1)
		}
	}
	return
}

// Wrapper for RPC calls.
func RPCSessionManager(f func(client *rpc.Client) error) {
	client, err := rpc.DialHTTP("tcp", opts.rpcListenAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()
	err = f(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Multiplexer for CLI commands.
func runCommand() {
	switch opts.cmd {
	case "daemon":
		log.Fatal(daemon.Daemon(opts.nodeListenAddr, opts.bootstrap, opts.bootstrapList, opts.rpcListenAddr))
	case "ping":
		args := akademiRPC.PingArgs{Host: core.Host(opts.target)}
		reply := akademiRPC.PingReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Ping", args, &reply)
		})
		fmt.Print("Received reply from ", opts.target, ". NodeID: ", reply.Header.NodeID, ".\n")
	case "lookup":
		id, err := core.B32ToID(opts.target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		args := akademiRPC.LookupArgs{ID: id}
		reply := akademiRPC.LookupReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Lookup", args, &reply)
		})
		fmt.Print("Node located successfully. Address: ", reply.RoutingEntry, ".\n")
	case "routing_table":
		args := akademiRPC.RoutingTableArgs{}
		reply := akademiRPC.RoutingTableReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.RoutingTable", args, &reply)
		})
		fmt.Print("Node routing table:\n", reply.RoutingTable, "\n")
	case "datastore":
		args := akademiRPC.DataStoreArgs{}
		reply := akademiRPC.DataStoreReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.DataStore", args, &reply)
		})
		fmt.Print("Node datastore:\n", reply.DataStore, "\n")
	case "info":
		args := akademiRPC.NodeInfoArgs{}
		reply := akademiRPC.NodeInfoReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.NodeInfo", args, &reply)
		})
		fmt.Print("Node information:\n", reply.NodeInfo, "\n")
	case "bootstrap":
		args := akademiRPC.BootstrapArgs{Host: core.Host(opts.target)}
		reply := akademiRPC.BootstrapReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Bootstrap", args, &reply)
		})
		fmt.Print("Successfully bootstrapped node with ", opts.target, ".\n")
	case "store":
		args := akademiRPC.StoreArgs{Data: core.DataBytes(opts.target)}
		reply := akademiRPC.StoreReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Store", args, &reply)
		})
		fmt.Print("Data stored on the DHT successfully. KeyID: ", reply.KeyID, ".\n")
	case "get":
		keyID, err := core.B32ToID(opts.target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		args := akademiRPC.GetArgs{KeyID: keyID}
		reply := akademiRPC.GetReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Get", args, &reply)
		})
		fmt.Print("Data retreived from the DHT successfully.\nContent:\n", reply.Data, "\n")
	default:
		fmt.Print("Command \"", opts.cmd, "\" not found.\n")
		os.Exit(1)
	}
}

// Akademi entrypoint.
func main() {
	parseArgs()
	runCommand()
}
