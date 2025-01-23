package node

import (
	"crypto/sha1"
	"fmt"

	"github.com/gsergey418alt/akademi/core"
)

// Hash function for stored data.
func (a *AkademiNode) getKeyID(data core.DataBytes) core.BaseID {
	h := sha1.New()
	return core.BaseID(h.Sum(data))
}

// Get all the entries in the dataStore as a string.
func (a *AkademiNode) DataStoreString() (table string) {
	for keyID, data := range a.dataStore.data {
		table += fmt.Sprintln(keyID, data)
	}
	if len(table) > 0 {
		return table[:len(table)-1]
	} else {
		return ""
	}
}
