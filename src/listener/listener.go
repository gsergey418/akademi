package listener

import (
	"net"
	"net/http"
	"net/rpc"
)

type Listener interface {
	Listen(string) error
}

type UDPListener struct {
	RPCAdapter    RPCAdapter
	udpReadBuffer [2048]byte
}

// Listen opens a UDP RPC socket on the provided address
// in the format "127.0.0.1:443". This is the main loop
// of the program.
func (ul *UDPListener) Listen(listenAddr string) error {
	rpc.Register(ul.RPCAdapter)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	http.Serve(l, nil)
	return nil
}
