package test

import (
	"fmt"
	"testing"

	"github.com/gsergey418alt/akademi/core"
)

func TestGetPrefixLength(t *testing.T) {
	id0 := core.BaseID{0, 0, 0, 0, 0, 0, 127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	id1 := core.BaseID{0, 0, 0, 0, 0, 0, 63, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	fmt.Println(id0.GetPrefixLength(id1))
}
