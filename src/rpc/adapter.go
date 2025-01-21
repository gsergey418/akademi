package rpc

import "github.com/gsergey418alt/akademi/core"

type AkademiNodeRPCAdapter struct {
	AkademiNode *core.AkademiNode
}

func (a *AkademiNodeRPCAdapter) Ping(args *PingArgs, reply *PingReply) error {
	header, err := a.AkademiNode.Ping(args.host)
	reply.header = header
	return err
}

type PingArgs struct {
	host core.Host
}

type PingReply struct {
	header core.RoutingHeader
}
