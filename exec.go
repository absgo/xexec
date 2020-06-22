package xexec

import "os/exec"

// Exec defines the operations related to process execution.
type Exec interface {
	// RunCommand starts a process with the given configuration without waiting for it to finish.
	StartProcess(cmdConf *ProcessConf) (Process, error)

	// LookPath searches for an executable named file in the
	// directories named by the PATH environment variable.
	// If file contains a slash, it is tried directly and the PATH is not consulted.
	// LookPath also uses PATHEXT environment variable to match
	// a suitable candidate.
	// The result may be an absolute path or a path relative to the current directory.
	LookPath(file string) (string, error)
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

// LookPath searches for an executable named file in the
// directories named by the PATH environment variable.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// LookPath also uses PATHEXT environment variable to match
// a suitable candidate.
// The result may be an absolute path or a path relative to the current directory.
func (o *osExec) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
