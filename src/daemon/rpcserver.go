package daemon

type rpcServer interface {
	Serve() error
}
