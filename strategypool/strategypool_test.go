package strategypool

import (
	"fmt"
	"testing"
	"time"
	st "wonderstone/strategy_pool/strategytask"

	"github.com/google/uuid"
)

func TestNewStrategyPool(t *testing.T) {
	// Test NewStrategyPool
	sp := NewStrategyPool()
	sp.Init()
	time.Sleep(1 * time.Second)
	if sp == nil {
		t.Errorf("Expected StrategyPool to be non-nil, got nil")
	}
}

// // Test Add、Run、Start、Stop
// Test Register  unRegister IfRegistered  ReloadArgs GetTaskInfos
func TestRegister(t *testing.T) {
	// Positive Test Case
	id := uuid.New().String()
	binaryLocation := "./tasks/task"
	tmpTask := st.NewStrategyTask(id, binaryLocation)
	sp := NewStrategyPool()
	sp.Init()
	sp.Register(tmpTask,[]string{})
	if sp.allTaskInfo[id].Task == nil {
		t.Errorf("Expected StrategyPool to be non-nil, got nil")
	}

	// ifRegistered
	task ,ifRegistered:= sp.IfRegistered(id)
	if !ifRegistered {
		t.Errorf("Expected task to be registered, but it is not")
	}
	// Check if the task is the same
	// this is a pointer comparison, or say map variable comparison
	fmt.Printf("   task:%v\ntmpTask:%v\n",task,tmpTask)
	if task != tmpTask {
		t.Errorf("Expected task to be the same, but it is different")
	}

	// ReloadArgs
	err := sp.ReloadArgs(id, []string{"arg1", "arg2"})
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	if sp.allTaskInfo[id].Args[0] != "arg1" {
		t.Errorf("Expected arg1, got %v", sp.allTaskInfo[id].Args[0])
	}

	// GetTaskInfos
	taskInfos := sp.GetTaskInfos()
	tmp:= taskInfos[id].Task
	if tmp != tmpTask {
		t.Errorf("Expected task to be the same, but it is different")
	}
	// ! Do sth not recommended
	tmp.Run()

	
	// Unregister
	sp.UnRegister(id)
	if _, ifRegistered := sp.IfRegistered(id); ifRegistered {
		t.Errorf("Expected task to be unregistered, but it is not")
	}

	// GetTaskInfos
	taskInfos = sp.GetTaskInfos()
	if taskInfos[id].Task != nil {
		t.Errorf("Expected task to be nil, got %v", taskInfos[id].Task)
	}

}


// test Task Related
func TestTaskRelated(t *testing.T) {
	id := uuid.New().String()
	binaryLocation := "./tasks/task"
	tmpTask := st.NewStrategyTask(id, binaryLocation)
	id1 := uuid.New().String()
	binaryLocation1 := "./tasks/task1"
	tmpTask1 := st.NewStrategyTask(id1, binaryLocation1)


	sp := NewStrategyPool()
	sp.Init()
	sp.Register(tmpTask,[]string{})
	sp.Register(tmpTask1,[]string{})
	
	// Run Start Stop
	// Run a task with status info
	err := sp.Run(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	time.Sleep(1 * time.Second)

	// GetOnlineTasks
	onlineTasks := sp.GetOnlineTasks_ATI()
	fmt.Println(onlineTasks)
	// GetOfflineTasks
	offlineTasks := sp.GetOfflineTasks()
	fmt.Println(offlineTasks)

	// check the state of the task
	fmt.Println("After Run the task: ",StatusStr(sp.allTaskInfo[id].Stsinfo.Status))

	time.Sleep(65 * time.Second)

	fmt.Println("sleep 65s")

	// check the state of the task
	fmt.Println("when the task should be done",StatusStr(sp.allTaskInfo[id].Stsinfo.Status))

	// Run a task with status info
	err = sp.Run(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	time.Sleep(5 * time.Second)

	// ~ stop the task
	err = sp.Stop(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	time.Sleep(5 * time.Second)

	fmt.Println("After Stop the task: ",StatusStr(sp.allTaskInfo[id].Stsinfo.Status))

	// // GetTask
	time.Sleep(5 * time.Second)
	task,err := sp.GetTask(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	fmt.Println("GetTask: ",task.BinaryLocation, task.ID)

	// GetTaskInfo
	taskInfo,err := sp.GetTaskInfo(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	fmt.Println("GetTaskInfo: ",taskInfo.Task.BinaryLocation, taskInfo.Task.ID)

	// GetTaskStatus
	taskStatus,err := sp.GetTaskStatus(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	fmt.Println("GetTaskStatus: ",StatusStr(taskStatus))


}

// test onLineTasks
func TestOnLineTasks(t *testing.T) {

	id := uuid.New().String()
	binaryLocation := "./tasks/task"
	tmpTask := st.NewStrategyTask(id, binaryLocation)
	id1 := uuid.New().String()
	binaryLocation1 := "./tasks/task"
	tmpTask1 := st.NewStrategyTask(id1, binaryLocation1)


	sp := NewStrategyPool()
	sp.Init()
	sp.Register(tmpTask,[]string{})
	sp.Register(tmpTask1,[]string{})
	
	// AddOnLineTasks
	sp.AddOnLineTasks(id,id1)

	// RemoveOnLineTasks
	sp.RemoveOnLineTasks(id)

	// GetonLineTasks_target

	tmpTasks := sp.GetonLineTasks_target()
	fmt.Println(tmpTasks)

	

	// RunOnLineTasks
	sp.RunOnLineTasks()

	time.Sleep(5 * time.Second)
	// CheckOnLineTasks
	res, err :=  sp.CheckOnLineTasks()
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	fmt.Println("CheckOnLineTasks: ",res)




}

// test 

// output status func:
func StatusStr(status Status) string{
	switch status {
	case Online:
		return "Online"
	case Offline:
		return "Offline"
	case Offline_Done:
		return "Offline_Done"
	case Offline_Terminated:
		return "Offline_Terminated"
	case Offline_Other:
		return "Offline_Other"
	default:
		return "Unknown"
	}
}