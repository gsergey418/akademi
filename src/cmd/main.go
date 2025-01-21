package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gsergey418alt/akademi/daemon"
)

const (
	rpcListenAddr  = "127.0.0.1:3855"
	nodeListenAddr = "0.0.0.0:3865"
)

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() (cmd string, bootstrap bool) {
	bootstrap = true
	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments, please provide a command.\n")
		os.Exit(1)
	}
	cmd = os.Args[1]
	for _, arg := range os.Args[2:] {
		switch arg {
		case "--no-bootstrap":
			bootstrap = false
		}
	}
	return
}

// Akademi entrypoint.
func main() {
	cmd, bootstrap := parseArgs()
	switch cmd {
	case "daemon":
		log.Fatal(daemon.Daemon(nodeListenAddr, bootstrap, rpcListenAddr))
	default:
		fmt.Print("Command \"", cmd, "\" not found.\n")
		os.Exit(1)
	}
}
