package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	sp "wonderstone/strategy_pool/strategypool"

	pb "wonderstone/strategy_pool/example"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
	pb.UnimplementedStrategyPoolServer
}

var tmpsp = sp.NewStrategyPool()


// SayHello(context.Context, *HelloRequest) (*HelloResponse, error)

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", req.GetName())
	return &pb.HelloResponse{Message: "Hello " + req.GetName()}, nil
}


// check the strategy pool 
func Printer(){
	for {
		time.Sleep(2 * time.Second)
		fmt.Println(tmpsp)
	}
}


func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go Printer()
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	pb.RegisterStrategyPoolServer(s, &server{})
	log.Printf("Server started on port 8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
