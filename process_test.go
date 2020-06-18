package xexec

import (
	"context"
	"runtime"
	"testing"
)

func TestNewOsProcess_waitForExit(t *testing.T) {
	t.Parallel()
	proc, err := newOsProcess(&ProcessConf{
		Ctx:  context.Background(),
		Name: "go",
		Args: []string{"go", "help"},
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
	cmdPath := "sleep"
	cmdArgs := []string{"sleep", "100"}
	if runtime.GOOS == "windows" {
		cmdPath = "timeout"
		cmdArgs = []string{"timeout", "/t", "100"}
	}
	proc, err := newOsProcess(&ProcessConf{
		Ctx:  context.Background(),
		Name: cmdPath,
		Args: cmdArgs,
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
