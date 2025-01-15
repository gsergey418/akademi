package listener

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/gsergey418alt/akademi/core"
)

// RPCListener is an implementation of the listener
// interface that opens a HTTP RPC api on a network
// address.
type RPCListener struct {
	RPCAdapter RPCListenerAdapter
}

// The Initialize method on RPCListener assigns an
// RPCAdapter to it.
func (rl *RPCListener) Initialize(a *core.AkademiNode) {
	rl.RPCAdapter = &AkademiNodeRPCAdapter{AkademiNode: a}
}

// The Listen subroutine opens a HTTP RPC socket on the
// provided address in the format "127.0.0.1:443". This
// is the main loop of the program.
func (rl *RPCListener) Listen(listenAddr string) error {
	rpc.Register(rl.RPCAdapter)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	http.Serve(l, nil)
	return nil
}
