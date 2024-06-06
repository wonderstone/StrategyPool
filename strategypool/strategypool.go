package strategypool

import (
	"fmt"
	"os"
	st "wonderstone/strategy_pool/strategytask"
)

type void struct{}

type status int

const (
	online status = iota
	offline
)

type statusInfo struct {
	status status
	pid    interface{}
}

type taskInfo struct {
	task    *st.StrategyTask
	stsinfo statusInfo
	args    []string
}
type stgStatusError struct {
	id      string
	stsInfo statusInfo
	err     error
}

// - StrategyPool is used to store strategytasks
type StrategyPool struct {
	allTaskInfo    map[string]taskInfo
	finalCheckPids map[string]int // tasks pid that should be stopped by method Stop
	stgErrorCh     chan stgStatusError
	// onLineTasks is the target tasks that should be running
	onLineTasks map[string]void
}

// & Function Section 0: Top Level
// - NewStrategyPool returns a new strategypool
func NewStrategyPool() *StrategyPool {
	return &StrategyPool{
		allTaskInfo:    make(map[string]taskInfo),
		finalCheckPids: make(map[string]int),
		stgErrorCh:     make(chan stgStatusError),
		onLineTasks:    make(map[string]void),
	}
}

// & Function Section 0 End

// & Method Section 1: Struct Level

// - Method Init will start a lasting goroutine to check stgErrorCh for errors
// and print them out to the console if there are any
// - Method Init to start a lasting goroutine to check ch for errors
func (sp *StrategyPool) Init() {
	// like a watcher
	// all sp status update will count on this goroutine
	go func() {
		for stgStatusErr := range sp.stgErrorCh {
			// todo: delete following line
			fmt.Println("stgStatusError info From channel:", stgStatusErr)
			// update allTaskStatus according to stgStatusErr.sts
			temp := sp.allTaskInfo[stgStatusErr.id]
			temp.stsinfo = stgStatusErr.stsInfo
			sp.allTaskInfo[stgStatusErr.id] = temp
		}
	}()
}

// & Method Section 1 End

// & Method Section 2: All-Task Map Related
// - Method Register that add a strategytask to the strategypool allTasks
func (sp *StrategyPool) Register(task *st.StrategyTask, args []string) {
	sp.allTaskInfo[task.ID] = taskInfo{task, statusInfo{offline, nil}, args}
}

// - Method UnRegister that remove a strategytask from the strategypool allTasks
func (sp *StrategyPool) UnRegister(taskID string) {
	delete(sp.allTaskInfo, taskID)
}

// - Method IfRegistered that check if a strategytask is in the strategypool allTasks
func (sp *StrategyPool) IfRegistered(taskID string) (*st.StrategyTask, bool) {
	taskInfo, ok := sp.allTaskInfo[taskID]
	if ok {
		return taskInfo.task, true
	}
	return nil, false
}

// -Method ReloadArgs that reload the args of a strategytask
func (sp *StrategyPool) ReloadArgs(taskID string, args []string) error {
	// check if the task is in the strategypool
	_, ok := sp.IfRegistered(taskID)
	if !ok {
		return fmt.Errorf("task not found")
	}
	// update the args
	temp := sp.allTaskInfo[taskID]
	temp.args = args
	sp.allTaskInfo[taskID] = temp
	return nil
}

// & Method Section 2 End

// & Method Section 3: Task Related
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
// + Really need this one?
// Run method runs a strategytask in the strategypool by taskID and task.Run()
func (sp *StrategyPool) Run(taskID string) error {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return fmt.Errorf("task not found")
	}

	if task != nil {
		// update status to online
		// no pid under run mode for blocking reason
		err := task.Start(sp.allTaskInfo[taskID].args...)
		if err == nil {
			sp.stgErrorCh <- stgStatusError{taskID, statusInfo{online, task.GetPid()}, nil}
		} else {
			sp.stgErrorCh <- stgStatusError{taskID, statusInfo{offline, nil}, err}
		}
		err2 := task.Wait4()
		sp.stgErrorCh <- stgStatusError{taskID, statusInfo{offline, nil}, err2}
		return err2
	}
	return fmt.Errorf("task is nil")
}

// - Method Start a strategytask
// Start method starts a strategytask in the strategypool
func (sp *StrategyPool) Start(taskID string) error {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return fmt.Errorf("task not found")
	}
	if task != nil {
		err := task.Start(sp.allTaskInfo[taskID].args...)

		if err != nil {
			// update status to offline by sending to channel
			sp.stgErrorCh <- stgStatusError{taskID, statusInfo{offline, nil}, err}
		} else {
			// update status to online by sending to channel
			sp.stgErrorCh <- stgStatusError{taskID, statusInfo{online, task.GetPid()}, nil}
		}
		return err
	}
	return fmt.Errorf("task is nil")
}

// - Method Stop a strategytask
// Method Stop a strategytask
// Stop method stops a strategytask in the strategypool
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
		sp.stgErrorCh <- stgStatusError{taskID, statusInfo{offline, nil}, err}
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
func (sp *StrategyPool) GetTaskStatus(taskID string) (status, error) {
	task, ok := sp.IfRegistered(taskID)
	if !ok {
		return offline, fmt.Errorf("task not found")
	}
	if task != nil {
		return sp.allTaskInfo[taskID].stsinfo.status, nil
	}
	return offline, fmt.Errorf("task is nil")
}

// -Method GetOnlineTasks returns online tasks
// GetOnlineTasks method returns online tasks by iterating allTaskInfo
func (sp *StrategyPool) GetOnlineTasksStatus() map[string]void {
	onlineTasks := make(map[string]void)
	for key, taskInfo := range sp.allTaskInfo {
		if taskInfo.stsinfo.status == online {
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
		if taskInfo.stsinfo.status == offline {
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

// - Method GetOnLineTasks returns onLineTasks
// GetOnLineTasks method returns onLineTasks
func (sp *StrategyPool) GetOnLineTasks() map[string]void {
	return sp.onLineTasks
}

// -Method MaintainOnLineTasks maintains onLineTasks
// MaintainOnLineTasks method maintains onLineTasks and actually onlinetasksStatus are equal
func (sp *StrategyPool) MaintainOnLineTasks() {
	// get real online tasks
	onlineTasksStatus := sp.GetOnlineTasksStatus()
	// task in onLineTasks should be in onlineTasksStatus
	// or it should be started
	for key := range sp.onLineTasks {
		if _, ok := onlineTasksStatus[key]; !ok {
			// start the task
			sp.Start(key)
		}
	}

	// task in onlineTasksStatus should be in onLineTasks
	// or it should be stopped
	for key := range onlineTasksStatus {
		if _, ok := sp.onLineTasks[key]; !ok {
			// stop the task
			sp.Stop(key)
		}
	}
}

// & Method Section 5 End


