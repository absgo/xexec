package xexec

// Exec defines the operations related to process execution.
type Exec interface {
	StartProcess(cmdConf *ProcessConf) (Process, error)
}

// osExec is an implementation of "Exec" using the "os" package.
type osExec struct {
}

// NewOsExec creates a new "Exec" that is implemented using the "os" package.
func NewOsExec() Exec {
	return &osExec{}
}

// RunCommand starts a process with the given configuration without waiting for it to finish.
func (o *osExec) StartProcess(cmdConf *ProcessConf) (Process, error) {
	return newOsProcess(cmdConf)
}
