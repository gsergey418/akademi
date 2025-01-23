package tests

import (
	"fmt"
	"testing"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
)

func TestPing(t *testing.T) {
	d := &dispatcher.UDPDispatcher{}
	d.Initialize(core.RoutingHeader{NodeID: core.RandomBaseID(), ListenPort: core.IPPort(1337)})
	header, err := d.Ping(core.Host("127.0.0.1:3865"))
	if err != nil {
		panic(err)
	}
	fmt.Println("NodeID: ", header.NodeID)
}
