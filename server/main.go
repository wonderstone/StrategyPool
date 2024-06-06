package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	sp "wonderstone/strategy_pool/strategypool"
	st "wonderstone/strategy_pool/strategytask"

	pb "wonderstone/strategy_pool/example"

	"github.com/google/uuid"
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

// InitStrategyPool(context.Context, *InitStrategyPoolRequest) (*InitStrategyPoolResponse, error)
func (s *server) InitStrategyPool(ctx context.Context, req *pb.InitStrategyRequest) (*pb.InitStrategyResponse, error) {
	id := uuid.New().String()
	binaryLocation := "./tasks/task"
	tmpTask := st.NewStrategyTask(id, binaryLocation)
	tmpsp.Init()
	tmpsp.Register(tmpTask,[]string{})
	fmt.Println(tmpsp)

	return &pb.InitStrategyResponse{Message: "Initialized with "}, nil
}


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
