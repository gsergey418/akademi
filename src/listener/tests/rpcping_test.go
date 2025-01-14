package listener

import (
	"fmt"
	"net/rpc"
	"testing"

	"github.com/gsergey418alt/akademi/listener"
)

const (
	listenAddr = "127.0.0.1:3856"
)

var client *rpc.Client

func TestMain(m *testing.M) {
	fmt.Print("Connecting to RPC at ", listenAddr, ".\n")
	var err error
	client, err = rpc.DialHTTP("tcp", listenAddr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m.Run()
}

func TestRPCPing(t *testing.T) {
	args, reply := struct{}{}, struct{ listener.RPCResponse }{}
	err := client.Call("AkademiNodeRPCAdapter.Ping", args, &reply)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(reply.NodeID)
}
