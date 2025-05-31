package main

import (
	"badies/router"
	"fmt"
)

func main() {
	// Step 1: Create a HashRing with 3 virtual replicas per physical node
	ring := router.NewHashRing(3)

	// Step 2: Add 5 nodes
	for i := 1; i <= 5; i++ {
		nodeID := "node" + fmt.Sprintf("%d", i)
		ring.AddNode(nodeID)
		fmt.Println("Added:", nodeID)
	}

	fmt.Println("\nAll nodes in ring:")
	fmt.Println(ring.GetAllNodes())

	// Step 3: Pick a sample key
	key := "apple"
	nodes := ring.GetNodes(key)

	fmt.Printf("\nFor key '%s', responsible nodes (1 primary + 2 replicas) are:\n", key)
	for i, n := range nodes {
		if i == 0 {
			fmt.Printf("🔵 Primary: %s\n", n)
		} else {
			fmt.Printf("🟡 Replica %d: %s\n", i, n)
		}
	}
}
