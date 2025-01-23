package node

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/gsergey418alt/akademi/core"
)

// Update the routing table with a new entry.
func (a *AkademiNode) UpdateRoutingTable(r core.RoutingEntry) error {
	prefix := r.NodeID.GetPrefixLength(a.NodeID)
	a.routingTable.lock.Lock()
	for i, v := range a.routingTable.data[prefix] {
		if v.NodeID == r.NodeID || v.Host == r.Host {
			a.routingTable.data[prefix] = append(a.routingTable.data[prefix][:i], a.routingTable.data[prefix][i+1:]...)
			break
		}
	}
	if len(a.routingTable.data[prefix]) >= core.BucketSize {
		return fmt.Errorf("RoutingError: Bucket already full.")
	}
	a.routingTable.data[prefix] = append(a.routingTable.data[prefix], r)
	a.routingTable.lock.Unlock()
	return nil
}

// Returns an iterator for fetching routing entries right of
// the start point.
func (a *AkademiNode) rightRoutingTableFetcher(start int) func() (core.RoutingEntry, bool) {
	prefix := start
	i := 0
	return func() (core.RoutingEntry, bool) {
		for prefix < core.IDLength*8 && i >= len(a.routingTable.data[prefix]) {
			prefix++
			i = 0
		}
		if prefix < core.IDLength*8 || (prefix == core.IDLength*8 && i < len(a.routingTable.data[prefix])) {
			r := a.routingTable.data[prefix][i]
			i++
			return r, false
		} else {
			return core.RoutingEntry{}, true
		}
	}
}

// Returns an iterator for fetching routing entries left of
// the start point.
func (a *AkademiNode) leftRoutingTableFetcher(start int) func() (core.RoutingEntry, bool) {
	prefix := start
	i := 0
	return func() (core.RoutingEntry, bool) {
		for prefix > 0 && i >= len(a.routingTable.data[prefix]) {
			prefix--
			i = 0
		}
		if prefix > 0 || (prefix == 0 && i < len(a.routingTable.data[prefix])) {
			r := a.routingTable.data[prefix][i]
			i++
			return r, false
		} else {
			return core.RoutingEntry{}, true
		}
	}
}

// Gets the core.BucketSize closest nodes to the passed
// argument.
func (a *AkademiNode) GetClosestNodes(nodeID core.BaseID, amount int) ([]core.RoutingEntry, error) {
	var nodes []core.RoutingEntry
	prefix := a.NodeID.GetPrefixLength(nodeID)

	a.routingTable.lock.Lock()
	nodes = append(nodes, a.routingTable.data[prefix][:]...)

	lNext := a.leftRoutingTableFetcher(prefix - 1)
	rNext := a.rightRoutingTableFetcher(prefix + 1)
	var l, r core.RoutingEntry
	var lDone, rDone bool
	for len(nodes) < amount && (!lDone || !rDone) {
		l, lDone = lNext()
		if !lDone {
			nodes = append(nodes, l)
		}
		r, rDone = rNext()
		if !rDone {
			nodes = append(nodes, r)
		}
	}
	a.routingTable.lock.Unlock()

	if len(nodes) == 0 {
		return nodes, fmt.Errorf("Node doesn't exist.")
	}
	sort.Sort(sortBucketByDistance{NodeID: nodeID, Bucket: &nodes})
	return nodes, nil
}

// Locates a core.BaseID across the network.
func (a *AkademiNode) Lookup(nodeID core.BaseID, amount int) ([]core.RoutingEntry, error) {
	nodes, err := a.GetClosestNodes(nodeID, core.ConcurrentRequests)
	if err != nil {
		return nodes, err
	}
	if nodes[0].NodeID == nodeID {
		return nodes[:min(amount, len(nodes))], nil
	}

	var wg sync.WaitGroup
	queriedHosts := map[core.Host]bool{}
	prevClosestNode := core.RoutingEntry{}

	for nodes[0] != prevClosestNode {
		reqCounter := 0
		for i := 0; i < len(nodes) && reqCounter < core.ConcurrentRequests; i++ {
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
			return nodes[:min(amount, len(nodes))], nil
		}
		prevClosestNode = nodes[0]
		nodes, err = a.GetClosestNodes(nodeID, core.ConcurrentRequests)
		if err != nil {
			return nodes, err
		}
	}

	reqCounter := 0
	for i := 0; i < len(nodes) && reqCounter < core.BucketSize; i++ {
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
	nodes, err = a.GetClosestNodes(nodeID, core.ConcurrentRequests)
	if err != nil {
		return nodes, err
	}

	return nodes[:min(amount, len(nodes))], nil
}

// Locates a core.BaseID across the network with FindKey.
func (a *AkademiNode) KeyLookup(keyID core.BaseID) (core.DataBytes, error) {
	nodes, err := a.GetClosestNodes(keyID, core.ConcurrentRequests)
	if err != nil {
		return nil, err
	}

	queriedHosts := map[core.Host]bool{}
	prevClosestNode := core.RoutingEntry{}

	var c chan core.DataBytes

	for nodes[0] != prevClosestNode {
		reqCounter := 0
		for i := 0; i < len(nodes) && reqCounter < core.ConcurrentRequests; i++ {
			if _, ok := queriedHosts[nodes[i].Host]; ok == false {
				go func() {
					_, data, _, err := a.FindKey(nodes[i].Host, keyID)
					c <- data
					if err != nil {
						log.Print(err)
					}
				}()
				queriedHosts[nodes[i].Host] = true
				reqCounter++
			}
		}
		for i := 0; i < reqCounter; i++ {
			data := <-c
			if data != nil {
				return data, nil
			}
		}
		prevClosestNode = nodes[0]
		nodes, err = a.GetClosestNodes(keyID, core.ConcurrentRequests)
		if err != nil {
			return nil, err
		}
	}

	reqCounter := 0
	for i := 0; i < len(nodes) && reqCounter < core.BucketSize; i++ {
		if _, ok := queriedHosts[nodes[i].Host]; ok == false {
			go func() {
				_, data, _, err := a.FindKey(nodes[i].Host, keyID)
				c <- data
				if err != nil {
					log.Print(err)
				}
			}()
			queriedHosts[nodes[i].Host] = true
			reqCounter++
		}
	}
	for i := 0; i < reqCounter; i++ {
		data := <-c
		if data != nil {
			return data, nil
		}
	}
	nodes, err = a.GetClosestNodes(keyID, core.ConcurrentRequests)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Requested key not found.")
}

// Get all the entries in the routing table as a string.
func (a *AkademiNode) RoutingTableString() (table string) {
	for _, bucket := range a.routingTable.data {
		for _, r := range bucket {
			table += fmt.Sprintln(r)
		}
	}
	if len(table) > 0 {
		return table[:len(table)-1]
	} else {
		return ""
	}
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
	NodeID core.BaseID
	Bucket *[]core.RoutingEntry
}

// Finds bucket length in a SortBucketByDistance structure.
func (b sortBucketByDistance) Len() int {
	return len(*b.Bucket)
}

// Compares two entries in a SortBucketByDistance structure.
func (b sortBucketByDistance) Less(i, j int) bool {
	for o := 0; o < core.IDLength; o++ {
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
