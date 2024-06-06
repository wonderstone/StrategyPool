package strategypool

import (
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
// Test Register  unRegister IfRegistered  ReloadArgs
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
	// check if tthe task is the same
	// Check if the task is the same
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

	// Unregister
	sp.UnRegister(id)
	if _, ifRegistered := sp.IfRegistered(id); ifRegistered {
		t.Errorf("Expected task to be unregistered, but it is not")
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
	// Run a task
	err := sp.Run(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	// Start a task
	err = sp.Start(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	// Stop a task
	err = sp.Stop(id)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
}