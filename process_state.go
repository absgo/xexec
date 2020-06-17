package xexec

import "os"

// ProcessState stores information about a process.
type ProcessState interface {
	ExitCode() int
	String() string
	Success() bool
}

type osProcessState struct {
	state *os.ProcessState
}

func newProcessState(state *os.ProcessState) ProcessState {
	return &osProcessState{
		state: state,
	}
}

func (o *osProcessState) ExitCode() int {

	return o.state.ExitCode()
}

func (o *osProcessState) String() string {
	return o.state.String()
}

func (o *osProcessState) Success() bool {
	return o.state.Success()
}
