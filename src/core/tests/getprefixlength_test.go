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
	fmt.Println(id0.BinStr())
	id1, err := core.B32ToID("TCFUKN4QEBMV6T7NV3WW3I7SHBF2EG5N")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(id1.BinStr())
	id2, err := core.B32ToID("RL4O5LSQ2RS5ZVZ7DX4R2VHVX5TTB2Y4")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(id2.BinStr())

	fmt.Println(id0.GetPrefixLength(id1))
	fmt.Println(id0.GetPrefixLength(id2))
	fmt.Println(id0.GetPrefixLength(id0))
}
