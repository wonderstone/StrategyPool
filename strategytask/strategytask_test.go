package strategytask

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	// use google package to give uuid
	id := uuid.New().String()
	// turn id to string

	binaryLocation := "./tasks/task"
	// binaryLocation1 := "../tasks/task"
	st := NewStrategyTask(id, binaryLocation)
	// st1 := NewStrategyTask(id, binaryLocation1)

	if st.ID != id {
		t.Errorf("Expected ID %s, got %s", id, st.ID)
	}
	if st.BinaryLocation != binaryLocation {
		t.Errorf("Expected BinaryLocation %s, got %s", binaryLocation, st.BinaryLocation)
	}
	if st.pid != nil {
		t.Errorf("Expected pid to be nil, got %v", st.pid)
	}



	
	st.Run()
	// Test Start method with waitTag = false
	// it will not block the main thread
	st.Start()
	// st1.Start()
	// Test Check method before the process is done
	pid4st, err := st.CheckRunning()
	// pid4st should be a non-nil int interface and err should be nil
	if pid4st == nil || err != nil {
		t.Errorf("Expected pid to be non-nil and err to be nil, got %v and %v", pid4st, err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		// check if the process is running
		// only reasonable when Wait4 is called
		pid, err := st.CheckRunning()
		// pid should be a non-nil int interface and err should be nil
		if pid == nil || err != nil {
			t.Errorf("Expected pid to be non-nil and err to be nil, got %v and %v", pid, err)
		}
		err1 := st.Stop()
		// err should be nil: st1 should be stopped and return nil
		if err1 != nil {
			t.Errorf("Expected err to be nil, got %v", err1)
		}
	}()
	// Test Wait4 method
	err2 := st.Wait4()
	// err should be nil: st will be done and return nil
	if err2 != nil {
		t.Errorf("Expected err to be nil, got %v", err2)
	}


	// Check method after the process is done should be 
	// false no such process
	pid4st, err3 := st.CheckRunning()
	// pid4st should be nil and err should be sth not nil
	if pid4st != nil || err3 == nil {
		t.Errorf("Expected pid to be nil and err to be not nil, got %v and %v", pid4st, err3)
	}

}





// TestStart method with args
// Test Start method with args
func TestStart(t *testing.T) {
	id := uuid.New().String()
	binaryLocation := "./tasks/printArgs"
	st := NewStrategyTask(id, binaryLocation)
	// Test Start method with args
	st.Start( "arg1", "arg2")
	st.Start([]string{"arg3", "arg4"}...)
	st.Start()
	// test Run method with args
	// Test Run method with args
	st.Run("arg1", "arg2")
	st.Run([]string{"arg3", "arg4"}...)
	st.Run()

}