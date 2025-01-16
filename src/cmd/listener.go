package main

import (
	"github.com/gsergey418alt/akademi/core"
)

// Listener is an interface represeting the module
// responsible for receiving requests.
type Listener interface {
	Initialize(string, *core.AkademiNode) error
	Listen() error
}
