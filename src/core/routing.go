package core

import (
	"fmt"
)

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
		return fmt.Errorf("Bucket already full.")
	}
	a.RoutingTable[prefix] = append(a.RoutingTable[prefix], r)
	return nil
}

// Print all the entries in the routing table.
func (a *AkademiNode) PrintRoutingTable() {
	fmt.Println("Node routing table:")
	for _, bucket := range a.RoutingTable {
		for _, r := range bucket {
			fmt.Println(r.Host, r.NodeID.Base64Str())
		}
	}
}
