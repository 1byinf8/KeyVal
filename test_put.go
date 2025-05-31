package main

import (
	"log"
	"path/filepath"
	"strings"

	"badies/router"
)

func main() {
	// Create NodeManager
	nodeManager := router.NewNodeManager()

	// Define node IDs and paths
	nodeIDs := []string{"node1", "node2", "node3", "node4", "node5"}
	for _, nodeID := range nodeIDs {
		dbPath := filepath.Join("dbs", nodeID)
		err := nodeManager.AddNode(nodeID, dbPath)
		if err != nil {
			log.Fatalf("Failed to add node %s: %v", nodeID, err)
		}
	}

	// Create a hash ring and add nodes
	ring := router.NewHashRing(3) // 3 replicas per key
	for _, nodeID := range nodeIDs {
		ring.AddNode(nodeID)
	}

	// Key-value to insert
	key := "hello"
	value := []byte("world")

	// Get target nodes from hashring
	targetNodes := ring.GetNodes(key)
	log.Printf("Storing key '%s' to nodes: %v", key, targetNodes)

	// Store key-value pair to each node
	for _, nodeID := range targetNodes {
		realNodeID := strings.Split(nodeID, "#")[0] // strip replica info

		db, err := nodeManager.GetDB(realNodeID)
		if err != nil {
			log.Printf("Failed to get DB for node %s: %v", realNodeID, err)
			continue
		}

		err = db.Put([]byte(key), value, nil)
		if err != nil {
			log.Printf("Error writing key '%s' to node %s: %v", key, realNodeID, err)
		} else {
			log.Printf("Successfully wrote key '%s' to node %s", key, realNodeID)
		}
	}

	// Close all DBs
	if err := nodeManager.Close(); err != nil {
		log.Fatalf("Failed to close databases: %v", err)
	}
}
