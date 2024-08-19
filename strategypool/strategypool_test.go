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
	if sp.allTaskInfo[id].task == nil {
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
	if sp.allTaskInfo[id].args[0] != "arg1" {
		t.Errorf("Expected arg1, got %v", sp.allTaskInfo[id].args[0])
	}

	// GetTaskInfos
	taskInfos := sp.GetTaskInfos()
	tmp:= taskInfos[id].task
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
	if taskInfos[id].task != nil {
		t.Errorf("Expected task to be nil, got %v", taskInfos[id].task)
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

	time.Sleep(70 * time.Second)
	fmt.Println("sleep 60s")
	// // Start a task
	// err = sp.Start(id)
	// if err != nil {
	// 	t.Errorf("Expected error to be nil, got %v", err)
	// }

	// // Stop a task
	// err = sp.Stop(id)
	// if err != nil {
	// 	t.Errorf("Expected error to be nil, got %v", err)
	// }


	// // start 2 tasks
	// err = sp.Start(id)
	// if err != nil {
	// 	t.Errorf("Expected error to be nil, got %v", err)
	// }

	// err = sp.Start(id1)
	// if err != nil {
	// 	t.Errorf("Expected error to be nil, got %v", err)
	// }

	// // GetOnlineTasks
	// onlineTasks := sp.GetOnlineTasks()
	// fmt.Println(onlineTasks)

	// // 

}