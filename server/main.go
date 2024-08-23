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
		fmt.Println("StrategyPool is ",tmpsp)
		fmt.Println("====================================")
		fmt.Println("Online Tasks from SP_onLineTasks map: ")
		// tmp := tmpsp.GetOnLineTasks_OLT()
		tmp := tmpsp.GetOnlineTasks_ATI()
		for k := range tmp{
			fmt.Println("Task ID: ",k)
			tmptskinfo,_ := tmpsp.GetTaskInfo(k)
			fmt.Println("Task Info: ",tmptskinfo)
			info , err := tmpsp.GetTaskStatus(k)
			if err != nil {
				fmt.Println("Task Status: ",err)
			}else{
				fmt.Println("Task Status: ",info)
			}
		}
		fmt.Println("************************************")

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
