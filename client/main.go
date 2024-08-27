package main

import (
	"context"
	"log"
	"time"

	// "time"

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

	

	// ~ test Greeter Section
	c := pb.NewGreeterClient(conn)
	
	// @ SayHello method
	resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Failed to call SayHello: %v", err)
	}

	log.Printf("Message: %s", resp.GetMessage())
	
	// ~ test StrategyPool
	c2 := pb.NewStrategyPoolClient(conn)
	
	// @ InitStrategyPool method
	resp2, err := c2.InitStrategyPool(context.Background(), &pb.InitStrategyRequest{})
	if err != nil {
		log.Fatalf("Failed to call InitStrategyPool: %v", err)
	}

	log.Printf("Message: %s", resp2.GetInitStatus())
	// @ Register method
	resp3, err := c2.Register(context.Background(), &pb.RegisterRequest{ID: "1", BinaryLocation: "./tasks/task", Args: []string{"arg1", "arg2"}})
	resp3, err = c2.Register(context.Background(), &pb.RegisterRequest{ID: "2", BinaryLocation: "./tasks/task", Args: []string{"arg3", "arg4"}})
	resp3, err = c2.Register(context.Background(), &pb.RegisterRequest{ID: "3", BinaryLocation: "./tasks/task", Args: []string{"arg5", "arg6"}})
	if err != nil {
		log.Fatalf("Failed to call Register: %v", err)
	}

	log.Printf("Message: %s", resp3.GetRegisterStatus())

	// @ UnRegister method	
	resp4, err := c2.UnRegister(context.Background(), &pb.UnregisterRequest{ID: "3"})
	if err != nil {
		log.Fatalf("Failed to call UnRegister: %v", err)
	}

	log.Printf("Message: %s", resp4.GetUnregisterStatus())

	// @ IfRegistered method
	resp5, err := c2.IfRegistered(context.Background(), &pb.IfRegisteredRequest{ID: "1"})
	if err != nil {
		log.Fatalf("Failed to call IfRegistered: %v", err)
	}

	log.Printf("Message: %s", resp5.GetIfRegisteredStatus())

	// @ ReloadArgs method
	resp6, err := c2.ReloadArgs(context.Background(), &pb.ReloadArgsRequest{ID: "1", Args: []string{}})
	resp6, err = c2.ReloadArgs(context.Background(), &pb.ReloadArgsRequest{ID: "2", Args: []string{}})
	if err != nil {
		log.Fatalf("Failed to call ReloadArgs: %v", err)
	}

	log.Printf("Message: %s", resp6.GetReloadArgsStatus())

	// @ GetTaskInfos method
	resp7, err := c2.GetTaskInfos(context.Background(), &pb.GetTaskInfosRequest{})
	if err != nil {
		log.Fatalf("Failed to call GetTaskInfos: %v", err)
	}

	log.Printf("Message: %s", resp7.GetTaskInfos())

	// ~ Method Section 3: Task Related	

	// @ Run method
	// - Recommended to use Run method to call the task
	resp9, err := c2.Run(context.Background(), &pb.RunRequest{ID: "1"})
	if err != nil {
		log.Fatalf("Failed to call Run: %v", err)
	}

	log.Printf("Message: %s", resp9.GetRunStatus())

	// $ Sleep for 5 seconds
	// $ wait for the task to start and update the pid
	// $ then call the Stop method
	time.Sleep(5 * time.Second)
	// @ Stop method
	resp10, err := c2.Stop(context.Background(), &pb.StopRequest{ID: "1"})
	if err != nil {
		log.Fatalf("Failed to call Stop: %v", err)
	}

	log.Printf("Message: %s", resp10.GetStopStatus())

	// @ GetTask method
	resp11, err := c2.GetTask(context.Background(), &pb.GetTaskInfoRequest{ID: "1"})
	if err != nil {
		log.Fatalf("Failed to call GetTask: %v", err)
	}

	log.Printf("Message: %s", resp11.GetTaskInfo())










}
