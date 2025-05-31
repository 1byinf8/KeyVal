package router

import (
	"fmt"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type NodeManager struct {
	mu        sync.RWMutex
	Instances map[string]*leveldb.DB
}

// NewNodeManager creates a new instance of NodeManager
func NewNodeManager() *NodeManager {
	return &NodeManager{
		Instances: make(map[string]*leveldb.DB),
	}
}

// AddNode adds a new node with the specified nodeID and database path
func (nm *NodeManager) AddNode(nodeID string, path string) error {
	return nm.AddNodeWithOptions(nodeID, path, nil)
}

// AddNodeWithOptions adds a new node with custom LevelDB options
func (nm *NodeManager) AddNodeWithOptions(nodeID string, path string, options *opt.Options) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Check if node already exists
	if _, exists := nm.Instances[nodeID]; exists {
		return fmt.Errorf("node %s already exists", nodeID)
	}

	// Open the database
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return fmt.Errorf("failed to open DB for node %s at path %s: %v", nodeID, path, err)
	}

	nm.Instances[nodeID] = db
	log.Printf("Successfully added node %s with database at %s", nodeID, path)
	return nil
}

// GetDB retrieves the database instance for the specified nodeID
func (nm *NodeManager) GetDB(nodeID string) (*leveldb.DB, error) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	db, exists := nm.Instances[nodeID]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}
	return db, nil
}

// RemoveNode removes a node and closes its database connection
func (nm *NodeManager) RemoveNode(nodeID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	db, exists := nm.Instances[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	// Close the database
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database for node %s: %v", nodeID, err)
	}

	// Remove from instances map
	delete(nm.Instances, nodeID)
	log.Printf("Successfully removed node %s", nodeID)
	return nil
}

// NodeExists checks if a node with the given nodeID exists
func (nm *NodeManager) NodeExists(nodeID string) bool {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	_, exists := nm.Instances[nodeID]
	return exists
}

// ListNodes returns a list of all node IDs
func (nm *NodeManager) ListNodes() []string {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	nodes := make([]string, 0, len(nm.Instances))
	for nodeID := range nm.Instances {
		nodes = append(nodes, nodeID)
	}
	return nodes
}

// NodeCount returns the number of managed nodes
func (nm *NodeManager) NodeCount() int {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	return len(nm.Instances)
}

// Close closes all database connections and cleans up resources
func (nm *NodeManager) Close() error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	var errs []error
	for nodeID, db := range nm.Instances {
		if err := db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close database for node %s: %v", nodeID, err))
		} else {
			log.Printf("Successfully closed database for node %s", nodeID)
		}
	}

	// Clear the instances map
	nm.Instances = make(map[string]*leveldb.DB)

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while closing databases: %v", errs)
	}

	log.Println("All node databases closed successfully")
	return nil
}

// CloseNode closes the database for a specific node without removing it from the manager
func (nm *NodeManager) CloseNode(nodeID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	db, exists := nm.Instances[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database for node %s: %v", nodeID, err)
	}

	delete(nm.Instances, nodeID)
	log.Printf("Successfully closed node %s", nodeID)
	return nil
}

// ReplaceNode replaces an existing node's database with a new one
func (nm *NodeManager) ReplaceNode(nodeID string, newPath string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Close existing database if it exists
	if existingDB, exists := nm.Instances[nodeID]; exists {
		if err := existingDB.Close(); err != nil {
			log.Printf("Warning: failed to close existing database for node %s: %v", nodeID, err)
		}
	}

	// Open new database
	db, err := leveldb.OpenFile(newPath, nil)
	if err != nil {
		return fmt.Errorf("failed to open new DB for node %s at path %s: %v", nodeID, newPath, err)
	}

	nm.Instances[nodeID] = db
	log.Printf("Successfully replaced database for node %s with new path %s", nodeID, newPath)
	return nil
}
