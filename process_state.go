package xexec

// ProcessState stores information about a process.
type ProcessState interface {
	ExitCode() int
	String() string
	Success() bool
}

// processStateDelegater is a used for obtaining the state of a process, but the
// information is delegated to another process state instance.
type processStateDelegater struct {
	delegate ProcessState
}

func newProcessStateDelegater(delegate ProcessState) ProcessState {
	return &processStateDelegater{
		delegate: delegate,
	}
}

func (o *processStateDelegater) ExitCode() int {
	return o.delegate.ExitCode()
}

func (o *processStateDelegater) String() string {
	return o.delegate.String()
}

func (o *processStateDelegater) Success() bool {
	return o.delegate.Success()
}
