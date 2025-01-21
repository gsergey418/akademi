package core

import (
	crand "crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/bits"
)

// The function GetPrefixLength finds the length of the
// common prefix between two Node/Key IDs.
func (id0 BaseID) GetPrefixLength(id1 BaseID) int {
	for i := 0; i < IDLength; i++ {
		xor := id0[i] ^ id1[i]
		if xor != 0 {
			return i*8 + bits.LeadingZeros8(xor)
		}
	}
	return 0
}

// Gets the distance between two NodeIDs.
func (id0 BaseID) GetDistance(id1 BaseID) int {
	var xor BaseID
	for i := 0; i < IDLength; i++ {
		xor[i] = id0[i] ^ id1[i]
	}
	return int(binary.BigEndian.Uint64(xor[:]))
}

// Returns BaseID string in binary.
func (id BaseID) BinStr() string {
	out := ""
	for i := 0; i < IDLength; i++ {
		out = out + fmt.Sprintf("%08b", id[i])
	}
	return out
}

// Returns base32 BaseID string.
func (id BaseID) String() string {
	return base32.StdEncoding.EncodeToString(id[:])
}

// Returns random BaseID.
func RandomBaseID() BaseID {
	var o BaseID
	crand.Read(o[:])
	return o
}
