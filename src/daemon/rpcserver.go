package daemon

type rpcServer interface {
	Serve(listenAddr string) error
}
