package core

import (
	crand "crypto/rand"
	"encoding/base32"
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
	return IDLength * 8
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

// Create a BaseID from base32.
func B32ToID(s string) (BaseID, error) {
	bytes, err := base32.StdEncoding.DecodeString(s)
	if err != nil || len(bytes) != IDLength {
		return BaseID{}, fmt.Errorf("Wrong ID format, use %d-byte base32 string.\n", IDLength)
	}
	return BaseID(bytes), nil
}

// Returns random BaseID.
func RandomBaseID() BaseID {
	var o BaseID
	crand.Read(o[:])
	return o
}
