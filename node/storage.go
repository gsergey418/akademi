package node

import (
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/gsergey418/akademi/core"
)

// Hash function for stored data.
func getKeyID(data core.DataBytes) core.BaseID {
	h := sha1.New()
	return core.BaseID(h.Sum(data))
}

// Write entry to the data storage
func (a *AkademiNode) Set(data core.DataBytes) error {
	if len(data) > core.MaxDataLength {
		return fmt.Errorf("data size too big. max data length: %d bytes", core.MaxDataLength)
	}
	keyID := getKeyID(data)
	a.dataStore.lock.Lock()
	defer a.dataStore.lock.Unlock()
	a.dataStore.data[keyID] = &core.DataContainer{Data: data, LastAccess: time.Now()}
	return nil
}

// Read entry from data storage
func (a *AkademiNode) Get(keyID core.BaseID) core.DataBytes {
	a.dataStore.lock.Lock()
	defer a.dataStore.lock.Unlock()
	data, ok := a.dataStore.data[keyID]
	if ok {
		a.dataStore.data[keyID].LastAccess = time.Now()
		return data.Data
	} else {
		return nil
	}
}

// Finds the best-fitting nodes and replicates the value
// to them.
func (a *AkademiNode) Publish(data core.DataBytes) (core.BaseID, error) {
	keyID := getKeyID(data)
	nodes, err := a.Lookup(keyID, core.Replication)
	if err != nil {
		return core.BaseID{}, err
	}
	if len(nodes) == 0 {
		return core.BaseID{}, fmt.Errorf("no suitable nodes found")
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

// Expires old records.
func (a *AkademiNode) ExpireOldData() {
	a.dataStore.lock.Lock()
	defer a.dataStore.lock.Unlock()
	expired := 0
	for k, v := range a.dataStore.data {
		if time.Since(v.LastAccess) > core.DataExpire {
			delete(a.dataStore.data, k)
			expired++
		}
	}
	if expired > 0 {
		log.Print("Expired ", expired, " records.")
	}
}
