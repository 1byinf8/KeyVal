package main

import (
	"context"
	"log"
	"net"
	"path/filepath"
	"strings"

	pb "badies/proto/badiespb"
	"badies/router"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKeyValServer
	nodeManager *router.NodeManager
	ring        *router.HashRing
}

// Put stores a key-value pair across the nodes determined by the hash ring
func (s *server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	key := req.GetKey()
	value := []byte(req.GetValue()) // Convert string to []byte for storage
	targetNodes := s.ring.GetNodes(key)
	log.Printf("Storing key '%s' to nodes: %v", key, targetNodes)

	for _, nodeID := range targetNodes {
		realNodeID := strings.Split(nodeID, "#")[0] // Strip replica info
		db, err := s.nodeManager.GetDB(realNodeID)
		if err != nil {
			log.Printf("Failed to get DB for node %s: %v", realNodeID, err)
			continue
		}
		err = db.Put([]byte(key), value, nil)
		if err != nil {
			log.Printf("Error writing key '%s' to node %s: %v", key, realNodeID, err)
			continue
		}
		log.Printf("Successfully wrote key '%s' to node %s", key, realNodeID)
	}
	return &pb.PutResponse{Success: true}, nil
}

// Get retrieves a value for a given key from the nodes in the hash ring
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	key := req.GetKey()
	targetNodes := s.ring.GetNodes(key)
	log.Printf("Retrieving key '%s' from nodes: %v", key, targetNodes)

	for _, nodeID := range targetNodes {
		realNodeID := strings.Split(nodeID, "#")[0] // Strip replica info
		db, err := s.nodeManager.GetDB(realNodeID)
		if err != nil {
			log.Printf("Failed to get DB for node %s: %v", realNodeID, err)
			continue
		}
		value, err := db.Get([]byte(key), nil)
		if err != nil {
			log.Printf("Error reading key '%s' from node %s: %v", key, realNodeID, err)
			continue
		}
		return &pb.GetResponse{Value: string(value), Found: true}, nil
	}
	return &pb.GetResponse{Value: "", Found: false}, nil
}

// Delete removes a key from the nodes in the hash ring
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	key := req.GetKey()
	targetNodes := s.ring.GetNodes(key)
	log.Printf("Deleting key '%s' from nodes: %v", key, targetNodes)

	success := false
	for _, nodeID := range targetNodes {
		realNodeID := strings.Split(nodeID, "#")[0]
		db, err := s.nodeManager.GetDB(realNodeID)
		if err != nil {
			log.Printf("Failed to get DB for node %s: %v", realNodeID, err)
			continue
		}
		err = db.Delete([]byte(key), nil)
		if err != nil {
			log.Printf("Error deleting key '%s' from node %s: %v", key, realNodeID, err)
			continue
		}
		log.Printf("Successfully deleted key '%s' from node %s", key, realNodeID)
		success = true
	}
	return &pb.DeleteResponse{Success: success}, nil
}

// UpdateKey renames a key across the nodes in the hash ring
func (s *server) UpdateKey(ctx context.Context, req *pb.UpdateKeyRequest) (*pb.UpdateKeyResponse, error) {
	oldKey := req.GetOldKey()
	newKey := req.GetNewKey()

	// Get value for oldKey
	resp, err := s.Get(ctx, &pb.GetRequest{Key: oldKey})
	if err != nil || !resp.Found {
		return &pb.UpdateKeyResponse{Success: false}, err
	}

	// Store value under newKey
	_, err = s.Put(ctx, &pb.PutRequest{Key: newKey, Value: resp.Value})
	if err != nil {
		return &pb.UpdateKeyResponse{Success: false}, err
	}

	// Delete oldKey
	_, err = s.Delete(ctx, &pb.DeleteRequest{Key: oldKey})
	if err != nil {
		return &pb.UpdateKeyResponse{Success: false}, err
	}

	return &pb.UpdateKeyResponse{Success: true}, nil
}

// UpdateValue updates the value for a key if the old value matches
func (s *server) UpdateValue(ctx context.Context, req *pb.UpdateValueRequest) (*pb.UpdateValueResponse, error) {
	key := req.GetKey()
	oldValue := req.GetOldValue()
	newValue := req.GetNewValue()

	// Verify old value
	resp, err := s.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil || !resp.Found || resp.Value != oldValue {
		return &pb.UpdateValueResponse{Success: false}, err
	}

	// Store new value
	_, err = s.Put(ctx, &pb.PutRequest{Key: key, Value: newValue})
	if err != nil {
		return &pb.UpdateValueResponse{Success: false}, err
	}

	return &pb.UpdateValueResponse{Success: true}, nil
}

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

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterKeyValServer(grpcServer, &server{nodeManager: nodeManager, ring: ring})

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Close all DBs on shutdown (optional, if server supports graceful shutdown)
	defer func() {
		if err := nodeManager.Close(); err != nil {
			log.Printf("Failed to close databases: %v", err)
		}
	}()
}
