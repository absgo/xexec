package xexec

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// Process represents a process that is executed.
type Process interface {
	// Pid returns the ID of the process.
	Pid() int

	// Wait waits for the Process to exit, and then returns a
	// ProcessState describing its status and an error, if any.
	// Wait releases any resources associated with the Process.
	// On most operating systems, the Process must be a child
	// of the current process or an error will be returned.
	Wait() (ProcessState, error)

	// Kill causes the Process to exit immediately. Kill does not wait until
	// the Process has actually exited. This only kills the Process itself,
	// not any other processes it may have started.
	Kill() error

	// Release releases any resources associated with the Process p,
	// rendering it unusable in the future.
	// Release only needs to be called if Wait is not.
	Release() error
}

type osProcess struct {
	proc *os.Process
}

// Pid returns the ID of the process.
func (o *osProcess) Pid() int {
	return o.proc.Pid
}

// Wait waits for the Process to exit, and then returns a
// ProcessState describing its status and an error, if any.
// Wait releases any resources associated with the Process.
// On most operating systems, the Process must be a child
// of the current process or an error will be returned.
func (o *osProcess) Wait() (ProcessState, error) {
	state, err := o.proc.Wait()
	if err != nil {
		return nil, err
	}
	return newProcessState(state), nil
}

// Kill causes the Process to exit immediately. Kill does not wait until
// the Process has actually exited. This only kills the Process itself,
// not any other processes it may have started.
func (o *osProcess) Kill() error {
	return o.proc.Kill()
}

// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.
func (o *osProcess) Release() error {
	return o.proc.Release()
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (o *osProcess) Signal(sig Signal) error {
	return o.proc.Signal(sig)
}

// newOsProcess starts a new process using the default implementation in the
// "os" package.
func newOsProcess(procConf *ProcessConf) (Process, error) {
	cmd := exec.CommandContext(procConf.Ctx, procConf.Name)
	// Set the path after creating the command so that we are able to control the first
	// argument.
	cmd.Args = procConf.Args
	cmd.Dir = procConf.Dir
	cmd.Env = procConf.Env
	cmd.ExtraFiles = procConf.ExtraFiles
	cmd.Stdin = procConf.Stdin
	cmd.Stderr = procConf.Stderr
	cmd.Stdout = procConf.Stdout
	cmd.SysProcAttr = procConf.SysProcAttr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &osProcess{
		proc: cmd.Process,
	}, nil
}

// ProcessConf is the config used for running a process.
type ProcessConf struct {

	// Ctx is the context used to kill the process (by calling
	// os.Process.Kill) if the context becomes done before the command
	// completes on its own.
	Ctx context.Context

	// Name is the name of the program to run.
	// If name contains no path separators, {Exec.LookPath} is used to
	// resolve name to a complete path if possible. Otherwise it uses name
	// directly as the program path.
	Name string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, {Name} is used.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// ExtraFiles specifies additional open files to be inherited by the
	// new process. It does not include standard input, standard output, or
	// standard error. If non-nil, entry i becomes file descriptor 3+i.
	//
	// ExtraFiles is not supported on Windows.
	ExtraFiles []*os.File

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr
}
