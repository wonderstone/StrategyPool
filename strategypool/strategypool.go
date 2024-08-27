package strategypool

import (
	"fmt"
	"os"
	"sync"
	st "wonderstone/strategy_pool/strategytask"
)

type void struct{}

type Status int

const (
	Online  Status = iota
	Offline        // default for task that has not been started
	Offline_Done
	Offline_Terminated
	Offline_Other
)

type statusInfo struct {
	Status Status
	pid    interface{}
}

type taskInfo struct {
	Task    *st.StrategyTask
	Stsinfo statusInfo
	Args    []string
}
type stgStatusError struct {
	ID      string
	StsInfo statusInfo
	Err     error
}

// - StrategyPool is used to store strategytasks
type StrategyPool struct {
	allTaskInfo    map[string]taskInfo
	finalCheckPids map[string]int // tasks pid that should be stopped by method Stop
	StgErrorCh     chan stgStatusError
	// onLineTasks is the target tasks that should be running
	onLineTasks map[string]void
	// mu is a mutex to protect allTaskInfo and finalCheckPids
	m sync.Mutex
}

// & Function Section 0: Top Level
// - NewStrategyPool returns a new strategypool
func NewStrategyPool() *StrategyPool {
	return &StrategyPool{
		allTaskInfo:    make(map[string]taskInfo),
		finalCheckPids: make(map[string]int),
		StgErrorCh:     make(chan stgStatusError),
		onLineTasks:    make(map[string]void),
	}
}

// & Function Section 0 End

// & Method Section 1: Struct Level

// - Method Init will start a lasting goroutine to check stgErrorCh for errors
// and print them out to the console if there are any
// - Method Init to start a lasting goroutine to check ch for errors
// ~ it can also be used to update allTaskStatus and do some callback

func callback(saything string) {
	fmt.Println(saything)
}

func (sp *StrategyPool) Init() {
	// like a watcher
	// all sp status update will count on this goroutine
	go func() {
		for stgStatusErr := range sp.StgErrorCh {
			sp.m.Lock()

			// todo: delete following line
			// fmt.Println("stgStatusError info From channel:", stgStatusErr)
			// update allTaskStatus according to stgStatusErr.sts
			temp := sp.allTaskInfo[stgStatusErr.ID]
			temp.Stsinfo = stgStatusErr.StsInfo
			switch stgStatusErr.StsInfo.Status {
			case Online:
				// update the temp task status
				temp.Stsinfo.Status = Online
				// add pid to finalCheckPids
				sp.finalCheckPids[stgStatusErr.ID] = stgStatusErr.StsInfo.pid.(int)
				// fmt.Println("finalCheckPids-online:", sp.finalCheckPids)
				callback(fmt.Sprintf("The task %s started", stgStatusErr.ID))

			case Offline_Done, Offline_Terminated, Offline_Other:

				// update the temp task status
				temp.Stsinfo.Status = stgStatusErr.StsInfo.Status
				// change task pid to nil
				temp.Task.SetPid(nil)
				// remove pid from finalCheckPids
				delete(sp.finalCheckPids, stgStatusErr.ID)
				// fmt.Println("finalCheckPids-offline:", sp.finalCheckPids)
				callback(fmt.Sprintf("The task %s stopped", stgStatusErr.ID))

			}
			// update allTaskStatus
			sp.allTaskInfo[stgStatusErr.ID] = temp
		    sp.m.Unlock()

			callback(fmt.Sprintf("The task %s status updated", stgStatusErr.ID))
		}
	}()
}

// & Method Section 1 End

// & Method Section 2: All-Task Map Related
// ~ Method Register that add a strategytask to the strategypool allTasks
func (sp *StrategyPool) Register(task *st.StrategyTask, args []string) {
	sp.allTaskInfo[task.ID] = taskInfo{task, statusInfo{Offline, nil}, args}
}

// ~ Method UnRegister that remove a strategytask from the strategypool allTasks
func (sp *StrategyPool) UnRegister(taskID string) {
	delete(sp.allTaskInfo, taskID)
}

// ~ Method IfRegistered that check if a strategytask is in the strategypool allTasks
func (sp *StrategyPool) IfRegistered(taskID string) (*st.StrategyTask, bool) {
	taskInfo, ok := sp.allTaskInfo[taskID]
	if ok {
		return taskInfo.Task, true
	}
	return nil, false
}

// ~ Method ReloadArgs that reload the args of a strategytask
func (sp *StrategyPool) ReloadArgs(taskID string, args []string) error {
	// check if the task is in the strategypool
	_, ok := sp.IfRegistered(taskID)
	if !ok {
		return fmt.Errorf("task not found")
	}
	// update the args
	temp := sp.allTaskInfo[taskID]
	temp.Args = args
	sp.allTaskInfo[taskID] = temp
	return nil
}

// ~ Method GetTaskInfos returns all task infos
// ! Avoid using this method to manipulate the task
// = Just output the task info
func (sp *StrategyPool) GetTaskInfos() map[string]taskInfo {
	return sp.allTaskInfo
}

// ~ Method GetTaskInfo returns a task info
func (sp *StrategyPool) GetTaskInfo(taskID string) (taskInfo, error) {
	taskInfo, ok := sp.allTaskInfo[taskID]
	if ok {
		return taskInfo, nil
	}
	return taskInfo, fmt.Errorf("task not found")
}

// & Method Section 2 End

// & Method Section 3: Task Related
// ~ Run method:(in goroutine)execute the command and wait for it to finish.
// ~ Run method can do some callback after running
// ~ Start method: start the command and store the PID in the pid field but not wait for it to finish
// ~ Start method is really hard to do callback after starting
// - Method CheckRunning check if a strategytask is running
func (sp *StrategyPool) CheckRunning(taskID string) (interface{}, error) {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	if task != nil {
		pid, err := task.CheckRunning()
		return pid, err
	}
	return nil, fmt.Errorf("task is nil")
}

// - Method Run a strategytask
// + Task will get offline status info after running
// Run method runs a strategytask in the strategypool by taskID and task.Run()
func (sp *StrategyPool) Run(taskID string) error {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return fmt.Errorf("task not found")
	}

	if task != nil {
		go func() {
			// update status to online
			// no pid under run mode for blocking reason
			err := task.Start(sp.allTaskInfo[taskID].Args...)
			if err == nil {
				sp.StgErrorCh <- stgStatusError{taskID, statusInfo{Online, task.GetPid()}, nil}
			} else {
				sp.StgErrorCh <- stgStatusError{taskID, statusInfo{Offline_Other, nil}, err}
			}
			err2 := task.Wait4()
			sp.StgErrorCh <- stgStatusError{taskID, statusInfo{Offline_Done, nil}, err2}
		}()

	}
	return nil
}



// - Method Stop a strategytask
// Method Stop a strategytask
// Stop method stops a strategytask in the strategypool
// the fundemental difference task.Stop() will check if the pid is running
// therefore, Start and Run method both ok to use Stop method
//
func (sp *StrategyPool) Stop(taskID string) error {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return fmt.Errorf("task not found")
	}
	if task != nil {
		// add pid to finalCheckPids
		pid := task.GetPid()
		if pid != nil {
			sp.finalCheckPids[taskID] = pid.(int)
		}
		err := task.Stop()
		// update status to offline
		sp.StgErrorCh <- stgStatusError{taskID, statusInfo{Offline_Terminated, nil}, err}
		return err
	}
	return fmt.Errorf("task is nil")

}

// - Method StopAll stops all tasks in the strategypool
// StopAll method stops all tasks in the strategypool
func (sp *StrategyPool) StopAll() {
	for key := range sp.allTaskInfo {
		sp.Stop(key)
	}
}

// - Method GetTask returns a task
// GetTask method returns a specific task from the strategypool
func (sp *StrategyPool) GetTask(taskID string) (*st.StrategyTask, error) {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

// - Method GetTaskStatus returns a task status
// GetTaskStatus method returns a specific task status from the strategypool
func (sp *StrategyPool) GetTaskStatus(taskID string) (Status, error) {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return Offline_Other, fmt.Errorf("task not found")
	}
	if task != nil {
		// check task pid is really running and change the related status
		pid, err := task.CheckRunning()
		fmt.Println("pid:", pid, "err:", err)
		if err == nil {
			tmpstatusInfo := statusInfo{Online, pid}
			sp.allTaskInfo[taskID] = taskInfo{task, tmpstatusInfo, sp.allTaskInfo[taskID].Args}
		}
		return sp.allTaskInfo[taskID].Stsinfo.Status, nil
	}
	return Offline_Other, fmt.Errorf("task is nil")
}

// -Method GetOnlineTasks_ATI returns online tasks
// GetOnlineTasks_ATI method returns online tasks by iterating allTaskInfo
func (sp *StrategyPool) GetOnlineTasks_ATI() map[string]void {
	onlineTasks := make(map[string]void)
	for key, taskInfo := range sp.allTaskInfo {
		if taskInfo.Stsinfo.Status == Online {
			onlineTasks[key] = void{}
		}
	}
	return onlineTasks
}

// -Method GetOfflineTasks returns offline tasks
// GetOfflineTasks method returns offline tasks by iterating allTaskInfo
func (sp *StrategyPool) GetOfflineTasks() map[string]void {
	offlineTasks := make(map[string]void)
	for key, taskInfo := range sp.allTaskInfo {
		// offline status includes offline, offline_Done, offline_Terminated, offline_Other
		if taskInfo.Stsinfo.Status == Offline || taskInfo.Stsinfo.Status == Offline_Done || taskInfo.Stsinfo.Status == Offline_Terminated || taskInfo.Stsinfo.Status == Offline_Other {
			offlineTasks[key] = void{}
		}
	}
	return offlineTasks
}

// & Method Section 3 End

// & Method Section 4: finalCheckPids
// - Method Check if Pids in finalCheckPids are still running
// CheckFinalCheckPids method checks if PIDs in finalCheckPids are still running
// if running, kill the process

func (sp *StrategyPool) CheckFinalCheckPids() (map[string]int, error) {
	var err error
	for key, pid := range sp.finalCheckPids {
		// use os package rather than syscall package in strategytask
		process, err := os.FindProcess(pid)
		if err != nil {
			// process not found
			delete(sp.finalCheckPids, key)
		} else {
			// kill the process
			err := process.Kill()
			if err == nil {
				delete(sp.finalCheckPids, key)
			}
		}
	}
	return sp.finalCheckPids, err
}

// - Method CheckFinalCheckPidEmpty check if finalCheckPids is empty
// CheckFinalCheckPidEmpty method checks if finalCheckPids is empty
func (sp *StrategyPool) CheckFinalCheckPidEmpty() bool {
	return len(sp.finalCheckPids) == 0
}

// & Method Section 4 End

// & Method Section 5: onLineTasks 
// & the inner onLineTasks is the target tasks that should be running
// - Method AddOnLineTasks adds tasks to onLineTasks
// AddOnLineTasks method adds tasks to onLineTasks
func (sp *StrategyPool) AddOnLineTasks(taskIDs ...string) {
	for _, taskID := range taskIDs {
		sp.onLineTasks[taskID] = void{}
	}
}

// - Method RemoveOnLineTasks removes tasks from onLineTasks
// RemoveOnLineTasks method removes tasks from onLineTasks
func (sp *StrategyPool) RemoveOnLineTasks(taskIDs ...string) {
	for _, taskID := range taskIDs {
		delete(sp.onLineTasks, taskID)
	}
}

// - Method GetonLineTasks_target returns onLineTasks
// GetonLineTasks_target method returns onLineTasks
func (sp *StrategyPool) GetonLineTasks_target() map[string]void {
	return sp.onLineTasks
}

// -Method StartonLineTasks
// StartOnLineTasks method starts onLineTasks
func (sp *StrategyPool) StartOnLineTasks() {
	// get real online tasks
	onlineTasks := sp.GetOnlineTasks_ATI()
	// task in sp.onLineTasks should be in onlineTasks
	// or it should be started
	for key := range sp.onLineTasks {
		if _, ok := onlineTasks[key]; !ok {
			// start the task
			sp.Run(key)
		}
	}

	// task in onlineTasks should be in sp.onLineTasks
	// or it should be stopped
	for key := range onlineTasks {
		if _, ok := sp.onLineTasks[key]; !ok {
			// stop the task
			sp.Stop(key)
		}
	}
}

//-Method RunOnLineTasks
// RunOnLineTasks method runs onLineTasks
func (sp *StrategyPool) RunOnLineTasks() {
	// get real online tasks
	onlineTasks := sp.GetOnlineTasks_ATI()
	// task in sp.onLineTasks should be in onlineTasks
	// or it should be started
	for key := range sp.onLineTasks {
		if _, ok := onlineTasks[key]; !ok {
			// run the task
			sp.Run(key)
		}
	}

	// task in onlineTasks should be in sp.onLineTasks
	// or it should be stopped
	for key := range onlineTasks {
		if _, ok := sp.onLineTasks[key]; !ok {
			// stop the task
			sp.Stop(key)
		}
	}
}



// - Method CheckOnLineTasks checks if onLineTasks are running
// CheckOnLineTasks method checks if onLineTasks are running
func (sp *StrategyPool) CheckOnLineTasks() (map[string]interface{}, error) {
	var err error
	onLineTasks := make(map[string]interface{})
	for key := range sp.onLineTasks {
		pid, err := sp.CheckRunning(key)
		if err == nil {
			onLineTasks[key] = pid
		}
	}
	return onLineTasks, err
}
// & Method Section 5 End
