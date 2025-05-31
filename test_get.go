package main

import (
	"context"
	"log"
	"time"

	pb "badies/proto/badiespb"

	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeyValClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Replace with any key you want to test
	key := "hello"

	// Send Get request
	resp, err := client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}

	if resp.GetFound() {
		log.Printf("Key: %s | Value: %s\n", key, resp.GetValue())
	} else {
		log.Printf("Key '%s' not found\n", key)
	}
}
