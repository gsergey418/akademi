package core

import (
	"fmt"
	"log"
	"sort"
)

// Pretty-print RoutingEntry.
func (r RoutingEntry) String() string {
	return fmt.Sprintf("%s@%s", r.NodeID, r.Host)
}

// Update the routing table with a new entry.
func (a *AkademiNode) UpdateRoutingTable(r RoutingEntry) error {
	prefix := r.NodeID.GetPrefixLength(a.NodeID)
	for i, v := range a.RoutingTable[prefix] {
		if v.NodeID == r.NodeID || v.Host == r.Host {
			a.RoutingTable[prefix] = append(a.RoutingTable[prefix][:i], a.RoutingTable[prefix][i+1:]...)
			break
		}
	}
	if len(a.RoutingTable[prefix]) >= BucketSize {
		return fmt.Errorf("RoutingError: Bucket already full.")
	}
	a.RoutingTable[prefix] = append(a.RoutingTable[prefix], r)
	return nil
}

// Gets the BucketSize closest nodes to the passed
// argument.
func (a *AkademiNode) GetClosestNodes(nodeID BaseID, amount int) []RoutingEntry {
	var nodes []RoutingEntry
	i := a.NodeID.GetPrefixLength(nodeID)
	for i >= 0 && len(nodes) < amount {
		nodes = append(nodes, a.RoutingTable[i][:]...)
		i--
	}
	sort.Sort(sortBucketByDistance{NodeID: nodeID, Bucket: &nodes})
	return nodes
}

// Locates a BaseID across the network.
func (a *AkademiNode) locateNode(nodeID BaseID) (RoutingEntry, error) {
	a.GetClosestNodes(nodeID, ConcurrentRequests)
	panic("Function locateNode not implemented.")
}

// Print all the entries in the routing table.
func (a *AkademiNode) printRoutingTable() {
	fmt.Println("Node routing table:")
	for _, bucket := range a.RoutingTable {
		for _, r := range bucket {
			fmt.Println(r)
		}
	}
}

// Log all the entries in the routing table.
func (a *AkademiNode) logRoutingTable() {
	log.Print("Node routing table:")
	for _, bucket := range a.RoutingTable {
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
