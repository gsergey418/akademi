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
	rpcListenAddr  = "127.0.0.1:3855"
	nodeListenAddr = "0.0.0.0:3865"
)

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() (cmd string, dest string, bootstrap bool) {
	bootstrap = true
	argLen := len(os.Args)
	if argLen < 2 {
		fmt.Print("Not enough arguments, please provide a command.\n")
		os.Exit(1)
	}
	optStart, optStop := 2, argLen
	cmd = os.Args[1]
	if cmd != "daemon" {
		dest = os.Args[argLen-1]
		optStop--
	}
	for _, arg := range os.Args[optStart:optStop] {
		switch arg {
		case "--no-bootstrap":
			bootstrap = false
		}
	}
	return
}

// Akademi entrypoint.
func main() {
	cmd, dest, bootstrap := parseArgs()
	switch cmd {
	case "daemon":
		log.Fatal(daemon.Daemon(nodeListenAddr, bootstrap, rpcListenAddr))
	case "ping":
		client, err := rpc.DialHTTP("tcp", rpcListenAddr)
		if err != nil {
			fmt.Println(err)
		}
		defer client.Close()
		args := akademiRPC.PingArgs{Host: core.Host(dest)}
		reply := akademiRPC.PingReply{}
		err = client.Call("AkademiNodeRPCServer.Ping", args, &reply)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("Received reply from ", dest, ". NodeID: ", reply.Header.NodeID, ".\n")
	default:
		fmt.Print("Command \"", cmd, "\" not found.\n")
		os.Exit(1)
	}
}
