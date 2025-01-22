package core

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
)

// Pretty-print RoutingEntry.
func (r RoutingEntry) String() string {
	return fmt.Sprintf("%s@%s", r.NodeID, r.Host)
}

// Update the routing table with a new entry.
func (a *AkademiNode) UpdateRoutingTable(r RoutingEntry) error {
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
	prefix := a.NodeID.GetPrefixLength(nodeID)

	a.routingTable.lock.Lock()
	nodes = append(nodes, a.routingTable.data[prefix][:]...)
	for i := 1; (i <= prefix || i < IDLength*8+1-prefix) && len(nodes) < amount; i++ {
		if i <= prefix {
			nodes = append(nodes, a.routingTable.data[prefix-i][:]...)
		}
		if i < IDLength*8+1 {
			nodes = append(nodes, a.routingTable.data[prefix+i][:]...)
		}
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
		if nodes[0].NodeID == nodeID {
			return nodes[:amount], nil
		}
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
	nodes, err = a.GetClosestNodes(nodeID, ConcurrentRequests)
	if err != nil {
		return nodes, err
	}

	return nodes[:amount], nil
}

// Get all the entries in the routing table as a string.
func (a *AkademiNode) RoutingTableString() (table string) {
	for _, bucket := range a.routingTable.data {
		for _, r := range bucket {
			table += fmt.Sprintln(r)
		}
	}
	return table[:len(table)-1]
}

// Print the routing table.
func (a *AkademiNode) PrintRoutingTable() {
	fmt.Println("Node routing table:")
	fmt.Println(a.RoutingTableString())
}

// Log all the entries in the routing table.
func (a *AkademiNode) LogRoutingTable() {
	log.Print("Node routing table:")
	for _, line := range strings.Split(a.RoutingTableString(), "\n") {
		log.Print(line)
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
