package main

import (
	"context"
	"log"
	"time"

	pb "badies/proto/badiespb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to server with updated credentials
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeyValClient(conn)

	// Extended timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("=== KeyVal Client Test ===")

	// 1. PUT
	log.Println("\n1. Testing PUT operation...")
	putRes, err := client.Put(ctx, &pb.PutRequest{Key: "hello", Value: "world"})
	if err != nil {
		log.Fatalf("Put failed: %v", err)
	}
	log.Printf("Put Success: %v", putRes.GetSuccess())

	// 2. GET
	log.Println("\n2. Testing GET operation...")
	getRes, err := client.Get(ctx, &pb.GetRequest{Key: "hello"})
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	log.Printf("Get: Found=%v, Value=%s", getRes.GetFound(), getRes.GetValue())

	// 3. Test GET with non-existent key
	log.Println("\n3. Testing GET with non-existent key...")
	getRes2, err := client.Get(ctx, &pb.GetRequest{Key: "nonexistent"})
	if err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	log.Printf("Non-existent key: Found=%v, Value=%s", getRes2.GetFound(), getRes2.GetValue())

	// 4. UPDATE VALUE
	log.Println("\n4. Testing UPDATE VALUE operation...")
	updValRes, err := client.UpdateValue(ctx, &pb.UpdateValueRequest{
		Key:      "hello",
		OldValue: "world",
		NewValue: "universe",
	})
	if err != nil {
		log.Fatalf("UpdateValue failed: %v", err)
	}
	log.Printf("UpdateValue Success: %v", updValRes.GetSuccess())

	// Verify UPDATE VALUE worked
	log.Println("\n5. Verifying UPDATE VALUE...")
	getRes3, err := client.Get(ctx, &pb.GetRequest{Key: "hello"})
	if err != nil {
		log.Fatalf("Get after UpdateValue failed: %v", err)
	}
	log.Printf("After UpdateValue - Value: %s", getRes3.GetValue())

	// 6. UPDATE KEY
	log.Println("\n6. Testing UPDATE KEY operation...")
	updKeyRes, err := client.UpdateKey(ctx, &pb.UpdateKeyRequest{
		OldKey: "hello",
		NewKey: "hi",
	})
	if err != nil {
		log.Fatalf("UpdateKey failed: %v", err)
	}
	log.Printf("UpdateKey Success: %v", updKeyRes.GetSuccess())

	// Verify UPDATE KEY worked
	log.Println("\n7. Verifying UPDATE KEY...")
	getRes4, err := client.Get(ctx, &pb.GetRequest{Key: "hi"})
	if err != nil {
		log.Fatalf("Get with new key failed: %v", err)
	}
	log.Printf("After UpdateKey - New key 'hi' value: %s", getRes4.GetValue())

	// Verify old key is gone
	getRes5, err := client.Get(ctx, &pb.GetRequest{Key: "hello"})
	if err != nil {
		log.Fatalf("Get with old key failed: %v", err)
	}
	log.Printf("Old key 'hello' found: %v", getRes5.GetFound())

	// 8. Test UPDATE VALUE with wrong old value (should fail)
	log.Println("\n8. Testing UPDATE VALUE with wrong old value...")
	updValRes2, err := client.UpdateValue(ctx, &pb.UpdateValueRequest{
		Key:      "hi",
		OldValue: "wrong_value",
		NewValue: "galaxy",
	})
	if err != nil {
		log.Fatalf("UpdateValue failed: %v", err)
	}
	log.Printf("UpdateValue with wrong old value - Success: %v (should be false)", updValRes2.GetSuccess())

	// 9. Test UPDATE KEY with non-existent key (should fail)
	log.Println("\n9. Testing UPDATE KEY with non-existent key...")
	updKeyRes2, err := client.UpdateKey(ctx, &pb.UpdateKeyRequest{
		OldKey: "nonexistent",
		NewKey: "newkey",
	})
	if err != nil {
		log.Fatalf("UpdateKey failed: %v", err)
	}
	log.Printf("UpdateKey with non-existent key - Success: %v (should be false)", updKeyRes2.GetSuccess())

	// 10. DELETE
	log.Println("\n10. Testing DELETE operation...")
	delRes, err := client.Delete(ctx, &pb.DeleteRequest{Key: "hi"})
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete Success: %v", delRes.GetSuccess())

	// Verify DELETE worked
	log.Println("\n11. Verifying DELETE...")
	getRes6, err := client.Get(ctx, &pb.GetRequest{Key: "hi"})
	if err != nil {
		log.Fatalf("Get after delete failed: %v", err)
	}
	log.Printf("After Delete - Key 'hi' found: %v (should be false)", getRes6.GetFound())

	log.Println("\n=== All tests completed! ===")
}
