package core

import (
	"fmt"
	"log"
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

// Print all the entries in the routing table.
func (a *AkademiNode) PrintRoutingTable() {
	fmt.Println("Node routing table:")
	for _, bucket := range a.RoutingTable {
		for _, r := range bucket {
			fmt.Println(r)
		}
	}
}

// Log all the entries in the routing table.
func (a *AkademiNode) LogRoutingTable() {
	log.Print("Node routing table:")
	for _, bucket := range a.RoutingTable {
		for _, r := range bucket {
			log.Print(r)
		}
	}
}

// Gets the BucketSize closest nodes to the passed
// argument.
func (a *AkademiNode) GetClosestNodes(nodeID BaseID) []RoutingEntry {
	var nodes []RoutingEntry
	i := a.NodeID.GetPrefixLength(nodeID)
	for i >= 0 && len(nodes) < BucketSize {
		nodes = append(nodes, a.RoutingTable[i][:]...)
		i--
	}
	return nodes
}
