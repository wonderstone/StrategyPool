package main

import (
	"context"
	"log"

	pb "wonderstone/strategy_pool/example"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Failed to call SayHello: %v", err)
	}

	log.Printf("Message: %s", resp.GetMessage())

	c2 := pb.NewStrategyPoolClient(conn)
	resp2, err := c2.InitStrategyPool(context.Background(), &pb.InitStrategyRequest{BinaryLocation: "./tasks/task"})
	if err != nil {
		log.Fatalf("Failed to call InitStrategyPool: %v", err)
	}

	log.Printf("Message: %s", resp2.GetID())

}
