package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gsergey418alt/akademi/core"
)

func TestGetPrefixLength(t *testing.T) {
	id0, err := core.B32ToID("RL4O5LSQ2RS5ZVZ7DX4R2VHVX5TTB2Y5")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	id1, err := core.B32ToID("TCFUKN4QEBMV6T7NV3WW3I7SHBF2EG5N")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	id2, err := core.B32ToID("RL4O5LSQ2RS5ZVZ7DX4R2VHVX5TTB2Y4")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if id0.GetPrefixLength(id1) != 3 {
		t.Fail()
	}
	if id0.GetPrefixLength(id2) != 159 {
		t.Fail()
	}
	if id0.GetPrefixLength(id0) != 160 {
		t.Fail()
	}
}
