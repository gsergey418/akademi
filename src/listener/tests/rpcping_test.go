package test

import (
	"fmt"
	"net/rpc"
	"testing"

	akademiRPC "github.com/gsergey418alt/akademi/rpc"
)

const (
	listenAddr = "127.0.0.1:3856"
)

var client *rpc.Client

func TestMain(m *testing.M) {
	fmt.Print("Connecting to RPC at ", listenAddr, ".\n")
	var err error
	client, err = rpc.DialHTTP("tcp", listenAddr)
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m.Run()
}

func TestRPCPing(t *testing.T) {
	args, reply := struct{}{}, akademiRPC.PingRPCResponse{}
	err := client.Call("AkademiNodeRPCAdapter.Ping", args, &reply)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Print("Ping Success! NodeID: ")
	for b := range reply.NodeID {
		fmt.Printf("%08b", b)
	}
	fmt.Println()
}
