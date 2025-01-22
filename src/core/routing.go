package core

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

// Pretty-print RoutingEntry.
func (r RoutingEntry) String() string {
	return fmt.Sprintf("%s@%s", r.NodeID, r.Host)
}

// Update the routing table with a new entry.
func (a *AkademiNode) UpdateRoutingTable(r RoutingEntry) error {
	if r.NodeID == a.NodeID {
		return fmt.Errorf("Can't put your own ID into the routing table.")
	}
	prefix := r.NodeID.GetPrefixLength(a.NodeID)
	a.routingTable.lock.Lock()
	for i, v := range a.routingTable.data[prefix] {
		if v.NodeID == r.NodeID || v.Host == r.Host {
			a.routingTable.data[prefix] = append(a.routingTable.data[prefix][:i], a.routingTable.data[prefix][i+1:]...)
			break
		}
	}
	if len(a.routingTable.data[prefix]) >= BucketSize {
		return fmt.Errorf("RoutingError: Bucket already full.")
	}
	a.routingTable.data[prefix] = append(a.routingTable.data[prefix], r)
	a.routingTable.lock.Unlock()
	return nil
}

// Gets the BucketSize closest nodes to the passed
// argument.
func (a *AkademiNode) GetClosestNodes(nodeID BaseID, amount int) ([]RoutingEntry, error) {
	var nodes []RoutingEntry
	i := a.NodeID.GetPrefixLength(nodeID)
	a.routingTable.lock.Lock()
	for i >= 0 && len(nodes) < amount {
		nodes = append(nodes, a.routingTable.data[i][:]...)
		i--
	}
	a.routingTable.lock.Unlock()
	if len(nodes) == 0 {
		return nodes, fmt.Errorf("Node doesn't exist.")
	}
	sort.Sort(sortBucketByDistance{NodeID: nodeID, Bucket: &nodes})
	return nodes, nil
}

// Locates a BaseID across the network.
func (a *AkademiNode) Lookup(nodeID BaseID, amount int) ([]RoutingEntry, error) {
	nodes, err := a.GetClosestNodes(nodeID, ConcurrentRequests)
	if err != nil {
		return nodes, err
	}
	if nodeID == a.NodeID {
		return []RoutingEntry{}, fmt.Errorf("Can't do a lookup on your own ID.")
	}
	if nodes[0].NodeID == nodeID {
		return nodes[:amount], nil
	}

	var wg sync.WaitGroup
	queriedHosts := map[Host]bool{}
	prevClosestNode := RoutingEntry{}

	for nodes[0] != prevClosestNode {
		reqCounter := 0
		if nodes[0].NodeID == nodeID {
			return nodes[:amount], nil
		}
		for i := 0; i < len(nodes) && reqCounter < ConcurrentRequests; i++ {
			if _, ok := queriedHosts[nodes[i].Host]; ok == false {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, _, err := a.FindNode(nodes[i].Host, nodeID)
					if err != nil {
						log.Print(err)
					}
				}()
				queriedHosts[nodes[i].Host] = true
				reqCounter++
			}
		}
		wg.Wait()
		prevClosestNode = nodes[0]
		nodes, err = a.GetClosestNodes(nodeID, ConcurrentRequests)
		if err != nil {
			return nodes, err
		}
	}

	reqCounter := 0
	for i := 0; i < len(nodes) && reqCounter < BucketSize; i++ {
		if _, ok := queriedHosts[nodes[i].Host]; ok == false {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _, err := a.FindNode(nodes[i].Host, nodeID)
				if err != nil {
					log.Print(err)
				}
			}()
			queriedHosts[nodes[i].Host] = true
			reqCounter++
		}
	}
	wg.Wait()

	return nodes[:amount], nil
}

// Print all the entries in the routing table.
func (a *AkademiNode) PrintRoutingTable() {
	fmt.Println("Node routing table:")
	for _, bucket := range a.routingTable.data {
		for _, r := range bucket {
			fmt.Println(r)
		}
	}
}

// Log all the entries in the routing table.
func (a *AkademiNode) LogRoutingTable() {
	log.Print("Node routing table:")
	for _, bucket := range a.routingTable.data {
		for _, r := range bucket {
			log.Print(r)
		}
	}
}

// Utility structure for sorting that implements the
// sort.Interface interface.
type sortBucketByDistance struct {
	NodeID BaseID
	Bucket *[]RoutingEntry
}

// Finds bucket length in a SortBucketByDistance structure.
func (b sortBucketByDistance) Len() int {
	return len(*b.Bucket)
}

// Compares two entries in a SortBucketByDistance structure.
func (b sortBucketByDistance) Less(i, j int) bool {
	for o := 0; o < IDLength; o++ {
		xor := int((*b.Bucket)[j].NodeID[o]^b.NodeID[o]) - int((*b.Bucket)[i].NodeID[o]^b.NodeID[o])
		if xor != 0 {
			return xor > 0
		}
	}
	return false
}

// Swaps two entries in a SortBucketByDistance structure.
func (b sortBucketByDistance) Swap(i, j int) {
	(*b.Bucket)[i], (*b.Bucket)[j] = (*b.Bucket)[j], (*b.Bucket)[i]
}
