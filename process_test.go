package xexec

import (
	"bytes"
	"context"
	"errors"
	os "os"
	"os/exec"
	"reflect"
	"strconv"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func getSleepCmdLine(sleepSec int) []string {
	return []string{"./bin/sleep", strconv.Itoa(sleepSec)}
}

func TestProcess_Wait(t *testing.T) {
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

func TestProcess_Kill(t *testing.T) {
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

func TestNewOsProcessWithStarter_fail(t *testing.T) {
	t.Parallel()
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	buf3 := &bytes.Buffer{}
	files := []*os.File{}
	tests := map[string]struct {
		procConf     *ProcessConf
		startProcErr error
		wantCmd      *exec.Cmd
		wantErr      bool
	}{
		"fail to start process": {
			procConf: &ProcessConf{
				Name: "myprogram",
			},
			startProcErr: errors.New(""),
			wantCmd: &exec.Cmd{
				Path: "myprogram",
			},
			wantErr: true,
		},
		"start process with all settings": {
			procConf: &ProcessConf{
				Ctx:         context.Background(),
				Name:        "myprogram",
				Args:        []string{"arg1", "arg2"},
				Env:         []string{"env1=val1", "env2=val2"},
				Dir:         "testdir",
				Stdin:       buf1,
				Stdout:      buf2,
				Stderr:      buf3,
				ExtraFiles:  files,
				SysProcAttr: &syscall.SysProcAttr{},
			},
			wantCmd: &exec.Cmd{
				Path:        "myprogram",
				Args:        []string{"arg1", "arg2"},
				Env:         []string{"env1=val1", "env2=val2"},
				Dir:         "testdir",
				Stdin:       buf1,
				Stdout:      buf2,
				Stderr:      buf3,
				ExtraFiles:  files,
				SysProcAttr: &syscall.SysProcAttr{},
			},
		},
	}

	wantPid := 111
	osProcInternal := &os.Process{Pid: wantPid}
	for label, tt := range tests {
		tt := tt
		t.Run(label, func(t *testing.T) {
			t.Parallel()
			proc, err := newOsProcessWithStarter(tt.procConf, func(cmd *exec.Cmd) error {
				require.Equal(t, cmd.Path, tt.wantCmd.Path)
				require.Equal(t, cmd.Args, tt.wantCmd.Args)
				require.Equal(t, cmd.Dir, tt.wantCmd.Dir)
				require.Equal(t, cmd.Env, tt.wantCmd.Env)
				if !cmp.Equal(cmd.ExtraFiles, tt.wantCmd.ExtraFiles) {
					t.Fatalf("not same extra files. got: %v, want: %v", cmd.ExtraFiles, tt.wantCmd.ExtraFiles)
				}
				requireSameOrNil(t, cmd.Stdin, tt.wantCmd.Stdin)
				requireSameOrNil(t, cmd.Stdout, tt.wantCmd.Stdout)
				requireSameOrNil(t, cmd.Stderr, tt.wantCmd.Stderr)
				require.Equal(t, cmd.SysProcAttr, tt.wantCmd.SysProcAttr)
				cmd.Process = osProcInternal
				return tt.startProcErr
			})
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error when fail to create process")
				}
				if proc != nil {
					t.Fatalf("want nil process, see %v", proc)
				}
			} else {
				if err != nil {
					t.Fatalf("fail to start process: %v", err)
				}
				osProc := proc.(*osProcess)
				wantOsProc := &osProcess{
					proc: osProcInternal,
					pid:  wantPid,
				}
				if !reflect.DeepEqual(osProc, wantOsProc) {
					t.Fatalf("incorrect process. got: %v, want: %v", osProc, wantOsProc)
				}
			}
		})
	}
}

func requireSameOrNil(t *testing.T, gotVal interface{}, wantVal interface{}) {
	t.Helper()
	if gotVal == wantVal {
		return
	}
	if wantVal == nil {
		require.Nil(t, gotVal)
	} else {
		require.Same(t, gotVal, wantVal)
	}
}
