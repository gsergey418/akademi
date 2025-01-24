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
func (a *AkademiNode) Set(data core.DataBytes) error {
	if len(data) > core.MaxDataLength {
		return fmt.Errorf("Data size too big. Max data length: %d bytes.", core.MaxDataLength)
	}
	keyID := getKeyID(data)
	a.dataStore.lock.Lock()
	defer a.dataStore.lock.Unlock()
	a.dataStore.data[keyID] = data
	return nil
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
func (a *AkademiNode) DHTStore(data core.DataBytes) (core.BaseID, error) {
	keyID := getKeyID(data)
	nodes, err := a.Lookup(keyID, core.Replication)
	if err != nil {
		return core.BaseID{}, err
	}
	if len(nodes) == 0 {
		return core.BaseID{}, fmt.Errorf("No suitable nodes found.")
	}
	c := make(chan error, core.Replication)

	for _, node := range nodes {
		if node.NodeID == a.NodeID {
			err := a.Set(data)
			c <- err
			continue
		}
		go func() {
			_, err := a.Store(node.Host, data)
			c <- err
		}()
	}
	for i := 0; i < len(nodes); i++ {
		err = <-c
		if err == nil {
			return keyID, nil
		}
	}
	return core.BaseID{}, err
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
