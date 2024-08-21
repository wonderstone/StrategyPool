package main

import (
	"context"

	pb "wonderstone/strategy_pool/example"
	"wonderstone/strategy_pool/strategytask"
)

// InitStrategyPool(context.Context, *InitStrategyPoolRequest) (*InitStrategyPoolResponse, error)
func (s *server) InitStrategyPool(ctx context.Context, req *pb.InitStrategyRequest) (*pb.InitStrategyResponse, error) {

	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	// function part
	tmpsp.Init()

	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.InitStrategyResponse{InitStatus: status}, err
}

// Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	bin_loc := req.GetBinaryLocation()
	args := req.GetArgs()

	tmpTask := strategytask.NewStrategyTask(id, bin_loc)

	tmpsp.Register(tmpTask, args)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.RegisterResponse{RegisterStatus: status}, err
}
