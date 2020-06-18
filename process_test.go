package xexec

import (
	"context"
	"strconv"
	"testing"
)

func getSleepCmdLine(sleepSec int) []string {
	return []string{"./bin/sleep", strconv.Itoa(sleepSec)}
}

func TestNewOsProcess_waitForExit(t *testing.T) {
	t.Parallel()
	cmdLine := getSleepCmdLine(0)
	proc, err := newOsProcess(&ProcessConf{
		Ctx:  context.Background(),
		Name: cmdLine[0],
		Args: cmdLine,
	})
	if err != nil {
		t.Fatalf("failed to create process: %v", err)
	}
	procState, err := proc.Wait()
	if err != nil {
		t.Fatalf("failed to wait: %v", err)
	}
	if procState.ExitCode() != 0 {
		t.Fatalf("got exit code: %d, want: %d", procState.ExitCode(), 0)
	}
}

func TestNewOsProcess_killProcess(t *testing.T) {
	t.Parallel()
	cmdLine := getSleepCmdLine(100)
	proc, err := newOsProcess(&ProcessConf{
		Ctx:  context.Background(),
		Name: cmdLine[0],
		Args: cmdLine,
	})
	if err != nil {
		t.Fatalf("failed to create process: %v", err)
	}
	defer func() {
		if err := proc.Release(); err != nil {
			t.Fatalf("failed to release: %v", err)
		}
	}()
	if err := proc.Kill(); err != nil {
		t.Fatalf("failed to kill process: %v", err)
	}
}
