package main

import (
	"context"
	"strconv"

	pb "wonderstone/strategy_pool/example"
	sp "wonderstone/strategy_pool/strategypool"
	"wonderstone/strategy_pool/strategytask"
)

// & Method Section 1: InitStrategyPool
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
// & Method Section 2: All-Task Map Related
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

// Unregister(context.Context, *UnregisterRequest) (*UnregisterResponse, error)
func (s *server) UnRegister(ctx context.Context, req *pb.UnregisterRequest) (*pb.UnregisterResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	tmpsp.UnRegister(id)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.UnregisterResponse{UnregisterStatus: status}, err
}

// IfRegistered(context.Context, *IfRegisteredRequest) (*IfRegisteredResponse, error)
func (s *server) IfRegistered(ctx context.Context, req *pb.IfRegisteredRequest) (*pb.IfRegisteredResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	var status string
	_, ok := tmpsp.IfRegistered(id)
	if ok {
		status = "registered"
	} else {
		status = "not registered"
	}
	// return part
	return &pb.IfRegisteredResponse{IfRegisteredStatus: status}, err
}

// ReloadArgs(context.Context, *ReloadArgsRequest) (*ReloadArgsResponse, error)
func (s *server) ReloadArgs(ctx context.Context, req *pb.ReloadArgsRequest) (*pb.ReloadArgsResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	args := req.GetArgs()
	tmpsp.ReloadArgs(id, args)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.ReloadArgsResponse{ReloadArgsStatus: status}, err
}

//GetTaskInfos(context.Context, *GetTaskInfosRequest) (*GetTaskInfosResponse, error)
func (s *server) GetTaskInfos(ctx context.Context, req *pb.GetTaskInfosRequest) (*pb.GetTaskInfosResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	taskInfos := tmpsp.GetTaskInfos()
	// return part: change taskInfos to pb.TaskInfo
	var taskInfos_res []*pb.TaskInfo
	for _, taskInfo := range taskInfos {
		taskInfos_res = append(taskInfos_res, &pb.TaskInfo{
			ID: taskInfo.Task.ID,
			BinaryLocation: taskInfo.Task.BinaryLocation,
			Args: taskInfo.Args,
		})
	}

	
	return &pb.GetTaskInfosResponse{TaskInfos: taskInfos_res}, err
}


// & Method Section 3: Task Related
// CheckRunning(context.Context, *CheckRunningRequest) (*CheckRunningResponse, error)
func (s *server) CheckRunning(ctx context.Context, req *pb.CheckRunningRequest) (*pb.CheckRunningResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	var status string
	pid, err := tmpsp.CheckRunning(id)

	if err != nil {
		status = "failed"
	} else {
		if pid == 0 {
			status = "not running"
		} else {
			status = "running"
		}
	}
	// return part
	return &pb.CheckRunningResponse{CheckRunningStatus: status}, err
}

// Run(context.Context, *RunRequest) (*RunResponse, error)
// only use Run method to call the task
func (s *server) Run(ctx context.Context, req *pb.RunRequest) (*pb.RunResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	// check if the task is already running
	_, err = tmpsp.CheckRunning(id)
	if err == nil {
		return &pb.RunResponse{RunStatus: "already running"}, err
	}
	err = tmpsp.Run(id)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.RunResponse{RunStatus: status}, err
}



// Stop(context.Context, *StopRequest) (*StopResponse, error)
func (s *server) Stop(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	err = tmpsp.Stop(id)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.StopResponse{StopStatus: status}, err
}

// StopAll(context.Context, *StopAllRequest) (*StopAllResponse, error)
func (s *server) StopAll(ctx context.Context, req *pb.StopAllRequest) (*pb.StopAllResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	tmpsp.StopAll()
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}

	return &pb.StopAllResponse{StopAllStatus: status}, err
}

// GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error)
func (s *server) GetTask(ctx context.Context, req *pb.GetTaskInfoRequest) (*pb.GetTaskInfoResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	// get taskinfo
	

	taskinfo,err := tmpsp.GetTaskInfo(id)

	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		switch taskinfo.Stsinfo.Status {
		case sp.Online:
			status = "running"
		case sp.Offline:
			status = "not started"
		case sp.Offline_Done:
			status = "Done"
		case sp.Offline_Terminated:
			status = "Terminated"
		case sp.Offline_Other:
			status = "Offline by Other reasons"
		}
	}

	return &pb.GetTaskInfoResponse{TaskInfo: &pb.TaskInfo{
		ID: taskinfo.Task.ID,
		BinaryLocation: taskinfo.Task.BinaryLocation,
		Args: taskinfo.Args,
		Status: status,
	}}, err
}

// GetTaskStatus(context.Context, *GetTaskStatusRequest) (*GetTaskStatusResponse, error)
func (s *server) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusRequest) (*pb.GetTaskStatusResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	id := req.GetID()
	
	status, err := tmpsp.GetTaskStatus(id)
	// return part
	if err != nil {
		return &pb.GetTaskStatusResponse{Status: "failed"}, err
	}
	// change status to string
	var status_str string
	switch status {
	case sp.Online:
		status_str = "running"
	case sp.Offline:
		status_str = "not started"
	case sp.Offline_Done:
		status_str = "Done"
	case sp.Offline_Terminated:
		status_str = "Terminated"
	case sp.Offline_Other:
		status_str = "Offline by Other reasons"
	}
	return &pb.GetTaskStatusResponse{Status: status_str}, err
}

// GetOnlineTasks(context.Context, *GetOnlineTasksRequest) (*GetOnlineTasksResponse, error)
func (s *server) GetOnlineTasks(ctx context.Context, req *pb.GetOnlineTasksRequest) (*pb.GetOnlineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	onlineTasks := tmpsp.GetOnlineTasks_ATI()
	// return part
	var onlineTasks_res []string
	for id := range onlineTasks {
		onlineTasks_res = append(onlineTasks_res, id)
	}
	return &pb.GetOnlineTasksResponse{OnlineTasks: onlineTasks_res}, err
}

// GetOfflineTasks(context.Context, *GetOfflineTasksRequest) (*GetOfflineTasksResponse, error)
func (s *server) GetOfflineTasks(ctx context.Context, req *pb.GetOfflineTasksRequest) (*pb.GetOfflineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	offlineTasks := tmpsp.GetOfflineTasks()
	// return part
	var offlineTasks_res []string
	for id := range offlineTasks {
		offlineTasks_res = append(offlineTasks_res, id)
	}
	return &pb.GetOfflineTasksResponse{OfflineTasks: offlineTasks_res}, err
}

// AddOnLineTasks(context.Context, *AddOnLineTasksRequest) (*AddOnLineTasksResponse, error)
func (s *server) AddOnLineTasks(ctx context.Context, req *pb.AddOnLineTasksRequest) (*pb.AddOnLineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	ids := req.GetIDs()
	tmpsp.AddOnLineTasks(ids...)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}
	return &pb.AddOnLineTasksResponse{AddOnLineTasksStatus: status}, err
}

// RemoveOnLineTasks(context.Context, *RemoveOnLineTasksRequest) (*RemoveOnLineTasksResponse, error)
func (s *server) RemoveOnLineTasks(ctx context.Context, req *pb.RemoveOnLineTasksRequest) (*pb.RemoveOnLineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	ids := req.GetIDs()
	tmpsp.RemoveOnLineTasks(ids...)
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}
	return &pb.RemoveOnLineTasksResponse{RemoveOnLineTasksStatus: status}, err
}

// GetOnLineTasks(context.Context, *GetOnLineTasksRequest) (*GetOnLineTasksResponse, error)
func (s *server) GetOnLineTasks(ctx context.Context, req *pb.GetOnLineTasksRequest) (*pb.GetOnLineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	onlineTasks := tmpsp.GetOnlineTasks_ATI()
	// return part
	var onlineTasks_res []string
	for id := range onlineTasks {
		onlineTasks_res = append(onlineTasks_res, id)
	}
	return &pb.GetOnLineTasksResponse{OnLineTasks: onlineTasks_res}, err
}

// StartOnLineTasks(context.Context, *StartOnLineTasksRequest) (*StartOnLineTasksResponse, error)
func (s *server) StartOnLineTasks(ctx context.Context, req *pb.StartOnLineTasksRequest) (*pb.StartOnLineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	tmpsp.StartOnLineTasks()
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}
	return &pb.StartOnLineTasksResponse{StartOnLineTasksStatus: status}, err
}

// RunOnLineTasks(context.Context, *RunOnLineTasksRequest) (*RunOnLineTasksResponse, error)
func (s *server) RunOnLineTasks(ctx context.Context, req *pb.RunOnLineTasksRequest) (*pb.RunOnLineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part
	tmpsp.RunOnLineTasks()
	// return part
	var status string
	if err != nil {
		status = "failed"
	} else {
		status = "success"
	}
	return &pb.RunOnLineTasksResponse{RunOnLineTasksStatus: status}, err
}

// CheckOnLineTasks(context.Context, *CheckOnLineTasksRequest) (*CheckOnLineTasksResponse, error)
func (s *server) CheckOnLineTasks(ctx context.Context, req *pb.CheckOnLineTasksRequest) (*pb.CheckOnLineTasksResponse, error) {
	// catch error part
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// function part

	tmpMap, err := tmpsp.CheckOnLineTasks()

	resMap := make(map[string]string)
	for id, pid := range tmpMap {
		if pid == 0 {
			resMap[id] = "not running"
		} else {
			// type assertion

			switch typedPid := pid.(type) {
			case int:
				resMap[id] = strconv.Itoa(typedPid)
			case nil:
				resMap[id] = "no PID"
			default:
				resMap[id] = "Not running by other reasons"
			}
		}
	}

	
	// return part
	return &pb.CheckOnLineTasksResponse{CheckOnLineTasksStatus: resMap}, err
}