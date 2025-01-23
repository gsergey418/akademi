package node

import (
	"crypto/sha1"
	"fmt"

	"github.com/gsergey418alt/akademi/core"
)

// Hash function for stored data.
func getKeyID(data core.DataBytes) core.BaseID {
	h := sha1.New()
	return core.BaseID(h.Sum(data))
}

// Write entry to the data storage
func (a *AkademiNode) Set(data core.DataBytes) {
	keyID := getKeyID(data)
	a.dataStore.lock.Lock()
	defer a.dataStore.lock.Unlock()
	a.dataStore.data[keyID] = data
}

// Read entry from data storage
func (a *AkademiNode) Get(keyID core.BaseID) core.DataBytes {
	a.dataStore.lock.Lock()
	defer a.dataStore.lock.Unlock()
	data, ok := a.dataStore.data[keyID]
	if ok {
		return data
	} else {
		return nil
	}
}

// Finds the best-fitting nodes and replicates the value
// to them.
func (a *AkademiNode) DHTSet(data core.DataBytes) error {
	keyID := getKeyID(data)
	nodes, err := a.Lookup(keyID, core.Replication)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("No suitable nodes found.")
	}
	var c chan error

	for _, node := range nodes {
		go func() {
			_, err := a.Store(node.Host, data)
			c <- err
		}()
	}
	for i := 0; i < len(nodes); i++ {
		err = <-c
		if err == nil {
			return nil
		}
	}
	return err
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
