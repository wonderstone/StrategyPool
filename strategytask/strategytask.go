// Well! I know exec.CommandContext could be used here
// but it only adds more complexity to the code
// a map with stg name as key and *stgtask as value
// could be used to kill the target process when needed

package strategytask

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// - StrategyTask interface
type StrategyTaskInterface interface {
	// Run method: execute the command and wait for it to finish.
	// no pid stored
	// low flexibility in control
	// * not recommended
	Run(arg ... string) error
	// Start method: start the command and store the PID in the pid field
	// high flexibility in control
	// * recommended
	Start(arg ... string) error
	// CheckRunning method: if the StrategyTask is running in system by pid check.
	// - Attention: the pid field is an interface{}
	// pid is an int, but it could be nil if the process is not running
	CheckRunning() (interface{}, error)
	// Wait4 method: wait for the cmd change state and return error if any
	// + Is it really needed?
	Wait4() error
	// Stop method: stop the command and set the pid field to nil
	Stop() error
	// GetPid method: return the PID stored in the pid field
	GetPid() interface{}
}


// - StrategyTask struct
type StrategyTask struct {
	ID             string
	BinaryLocation string
	pid            interface{} // default nil
}

// - NewStrategyTask returns a pointer to a StrategyTask
func NewStrategyTask(id, binaryLocation string) *StrategyTask {
	return &StrategyTask{
		ID:             id,
		BinaryLocation: binaryLocation,
	}
}

// - Run method: execute the command and wait for it to finish.
// no pid stored
func (st *StrategyTask) Run(arg ...string) error {
	cmd := exec.Command(st.BinaryLocation, arg...)
	// the following two lines could be deleted
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	st.pid = cmd.Process.Pid
	err = cmd.Wait()
	st.pid = nil	
	return err
}

// - Start method: start the command and store the PID in the pid field
func (st *StrategyTask) Start(arg ...string) error {
	cmd := exec.Command(st.BinaryLocation, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	st.pid = cmd.Process.Pid
	return err
}

// - CheckRunning method: if the StrategyTask is running in system by pid check.
// - Attention: the pid field is an interface{} 
// pid is an int, but it could be nil if the process is not running
// but int type is not nilable in golang
// the default value of int is 0, which is a valid PID
// Do check the pid field is nil or not before type assertion
func (st *StrategyTask) CheckRunning() (interface{} , error) {
	// err is nil means the process is running
	// pid is nil means the process is not running
	pidInt, ok := st.pid.(int)
	if !ok {
		// maybe the pid field is nil for the process is not running
		return nil , fmt.Errorf("failed to convert pid to int")
	}
	// check if the process with the PID stored in the pid field is running
	// if running, return true, nil
	// if not running, return false, error
	err := syscall.Kill(pidInt, 0)
	// if err is error(syscall.Errno) ESRCH (3), the process is not running
	if err != nil {
		// if the process is not running, set the pid field to nil and the inProgress field to false
		st.pid = nil
		return nil, err
	} 
	return pidInt, err
}

// - Wait4 method: wait for the cmd change state and return error if any
func (st *StrategyTask) Wait4() error {
	// check st is running with CheckRunning method
	// if not running, return error
	pid, err := st.CheckRunning()
	if err != nil {
		return err
	}
	// should be pass the assertion
	pidInt, _ := pid.(int)
	// wait for the process with the PID stored in the pid field to finish
	_, err = syscall.Wait4(pidInt, nil, 0, nil)
	return err
}

// - Stop method: stop the command and set the pid field to nil
func (st *StrategyTask) Stop() error {
	// check st is running with CheckRunning method
	pid, err := st.CheckRunning()
	if err != nil {
		return err
	}
	// should be pass the assertion
	pidInt, _ := pid.(int)

	err = syscall.Kill(pidInt, syscall.SIGKILL)
	if err == nil {
		// err is nil means the process is killed
		st.pid = nil
	} 
	return err
}

// - GetPid method: return the PID stored in the pid field
func (st *StrategyTask) GetPid() (interface{}) {
	return st.pid
}

// - SetPid method: set the PID stored in the pid field
func (st *StrategyTask) SetPid(pid interface{}) {
	st.pid = pid
}