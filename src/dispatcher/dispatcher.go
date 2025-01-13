package dispatcher

import (
	"github.com/gsergey418alt/akademi/core"
)

// The Dispatcher interface is responsible for dispatching
// RPC requests to other nodes in the network.
type Dispatcher interface {
	Ping(*core.RoutingEntry)
}
