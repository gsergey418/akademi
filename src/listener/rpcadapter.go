package listener

import (
	"github.com/gsergey418alt/akademi/core"
)

// RPCAdapter is an interface that enables connection of the core logic to the Listener.
type RPCAdapter interface {
	Ping(args struct{}, reply *struct{}) error
	FindNode(args struct{ nodeID core.NodeID }, reply *struct{ routingEntry *[]core.RoutingEntry }) error
	FindKey(args struct{ keyID core.KeyID }, reply *struct {
		keyID        core.KeyID
		routingEntry []core.RoutingEntry
	}) error
	Store(args struct {
		keyID core.KeyID
		data  core.DataBytes
	}, reply *struct{}) error
}

// AkademiNodeRPCAdapter is an interface that enables connection of AkademiNode to the Listener.
type AkademiNodeRPCAdapter struct {
	AkademiNode *core.AkademiNode
}

func (a *AkademiNodeRPCAdapter) Ping(args struct{}, reply *struct{}) error {
	return nil
}

func (a *AkademiNodeRPCAdapter) FindNode(args struct{ nodeID core.NodeID }, reply *struct{ routingEntry *[]core.RoutingEntry }) error {
	panic("Function \"FindNode\" not implemented.")
}

func (a *AkademiNodeRPCAdapter) FindKey(args struct{ keyID core.KeyID }, reply *struct {
	keyID        core.KeyID
	routingEntry []core.RoutingEntry
}) error {
	panic("Function \"FindKey\" not implemented.")
}

func (a *AkademiNodeRPCAdapter) Store(args struct {
	keyID core.KeyID
	data  core.DataBytes
}, reply *struct{}) error {
	panic("Function \"Store\" not implemented.")
}
