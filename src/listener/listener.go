package listener

import (
	"net"
	"net/http"
	"net/rpc"
)

// Listener is an interface represeting the module
// responsible for receiving RPC requests.
type Listener interface {
	Listen(string) error
}

// RPCListener is an implementation of the listener
// interface that opens a HTTP RPC api on a network
// address.
type RPCListener struct {
	RPCAdapter ListenerAdapter
}

// The Listen subroutine opens a HTTP RPC socket on the
// provided address in the format "127.0.0.1:443". This
// is the main loop of the program.
func (ul *RPCListener) Listen(listenAddr string) error {
	rpc.Register(ul.RPCAdapter)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	http.Serve(l, nil)
	return nil
}
