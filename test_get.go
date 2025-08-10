// main.go
//
// Example gRPC client for interacting with the KeyVal service.
// This script demonstrates how to connect to the KeyVal server
// and perform a Get operation for a specific key.
//
// Workflow:
//  1. Establish a gRPC connection to the server at localhost:50051.
//  2. Create a KeyVal client using the generated protobuf code.
//  3. Define a context with a timeout to avoid hanging requests.
//  4. Perform a Get request for a given key.
//  5. Print the value if found, otherwise print a "not found" message.
//
// Usage:
//  go run main.go
//
// Note:
//  - Ensure the KeyVal gRPC server is running before executing this client.
//  - Update the `key` variable to query different keys.
//
package main

import (
	"context"
	"log"
	"time"

	pb "badies/proto/badiespb"

	"google.golang.org/grpc"
)

func main() {
	// Step 1: Establish a connection to the gRPC server.
	// Using WithInsecure() for simplicity — not recommended for production.
	// WithBlock() ensures Dial waits until the connection is established.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Step 2: Create the KeyVal gRPC client from the protobuf definition.
	client := pb.NewKeyValClient(conn)

	// Step 3: Create a context with a 5-second timeout to avoid blocking indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Step 4: Specify the key to fetch from the server.
	key := "hello"

	// Step 5: Perform the Get request.
	resp, err := client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}

	// Step 6: Check if the key exists and handle accordingly.
	if resp.GetFound() {
		log.Printf("Key: %s | Value: %s\n", key, resp.GetValue())
	} else {
		log.Printf("Key '%s' not found\n", key)
	}
}
