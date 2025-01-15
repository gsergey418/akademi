package test

import (
	"fmt"
	"testing"

	"github.com/gsergey418alt/akademi/core"
)

func TestGetPrefixLength(t *testing.T) {
	a := core.AkademiNode{}
	id0 := [32]byte{0, 0, 0, 0, 0, 0, 127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	id1 := [32]byte{0, 0, 0, 0, 0, 0, 63, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	fmt.Println(a.GetPrefixLength(id0, id1))
}
