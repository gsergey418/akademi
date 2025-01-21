package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

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
	targetHost     string
	rpcListenAddr  string
	nodeListenAddr string
	bootstrap      bool
}

// Global instance of cmdOptions
var opts cmdOptions

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() {
	opts.bootstrap = true
	opts.nodeListenAddr = defaultNodeListenAddr
	opts.rpcListenAddr = defaultRpcListenAddr

	argLen := len(os.Args)
	if argLen < 2 {
		fmt.Print("Not enough arguments, please provide a command.\n")
		os.Exit(1)
	}
	optStart, optStop := 2, argLen
	opts.cmd = os.Args[1]
	if opts.cmd != "daemon" {
		opts.targetHost = os.Args[argLen-1]
		optStop--
	}
	for i := optStart; i < optStop; i++ {
		switch os.Args[i] {
		case "--no-bootstrap":
			opts.bootstrap = false
		case "--rpc-addr":
			opts.rpcListenAddr = os.Args[i+1]
			i++
		default:
			fmt.Print("Unknown argument: \"", os.Args[i], "\".\n")
		}
	}
	return
}

// Akademi entrypoint.
func main() {
	parseArgs()
	switch opts.cmd {
	case "daemon":
		log.Fatal(daemon.Daemon(opts.nodeListenAddr, opts.bootstrap, opts.rpcListenAddr))
	case "ping":
		client, err := rpc.DialHTTP("tcp", opts.rpcListenAddr)
		if err != nil {
			fmt.Println(err)
		}
		defer client.Close()
		args := akademiRPC.PingArgs{Host: core.Host(opts.targetHost)}
		reply := akademiRPC.PingReply{}
		err = client.Call("AkademiNodeRPCServer.Ping", args, &reply)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("Received reply from ", opts.targetHost, ". NodeID: ", reply.Header.NodeID, ".\n")
	default:
		fmt.Print("Command \"", opts.cmd, "\" not found.\n")
		os.Exit(1)
	}
}
