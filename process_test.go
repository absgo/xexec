package xexec

import (
	"bytes"
	"context"
	"errors"
	"os"
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

func TestNewOsProcess(t *testing.T) {
	t.Parallel()
	cmdLine := getSleepCmdLine(100)
	proc, err := newOsProcess(&ProcessConf{
		Ctx:  context.Background(),
		Name: cmdLine[0],
		Args: cmdLine,
	})
	require.Nil(t, err)
	defer proc.Wait()
	osProc := proc.(*osProcess).proc.(*os.Process)
	require.Nil(t, osProc.Kill())
}

func TestProcess_Pid(t *testing.T) {
	t.Parallel()
	proc := &osProcess{
		pid: 123,
	}
	gotPid := proc.Pid()
	require.Equal(t, gotPid, 123)
}

func TestProcess_Kill(t *testing.T) {
	t.Parallel()
	procCtrl := new(mockOsProcessCtrl)
	defer procCtrl.AssertExpectations(t)
	proc := &osProcess{
		pid:  123,
		proc: procCtrl,
	}
	wantErr := errors.New("")
	procCtrl.On("Kill").
		Once().
		Return(wantErr)
	err := proc.Kill()
	require.Same(t, err, wantErr)
}

func TestProcess_Signal(t *testing.T) {
	t.Parallel()
	procCtrl := new(mockOsProcessCtrl)
	defer procCtrl.AssertExpectations(t)
	proc := &osProcess{
		pid:  123,
		proc: procCtrl,
	}
	wantErr := errors.New("")
	procCtrl.On("Signal", syscall.SIGINT).
		Once().
		Return(wantErr)
	err := proc.Signal(syscall.SIGINT)
	require.Same(t, err, wantErr)
}

func TestProcess_Release(t *testing.T) {
	t.Parallel()
	procCtrl := new(mockOsProcessCtrl)
	defer procCtrl.AssertExpectations(t)
	proc := &osProcess{
		pid:  123,
		proc: procCtrl,
	}
	wantErr := errors.New("")
	procCtrl.On("Release").
		Once().
		Return(wantErr)
	err := proc.Release()
	require.Same(t, err, wantErr)
}

func TestProcess_Wait(t *testing.T) {
	t.Parallel()
	procCtrl := new(mockOsProcessCtrl)
	defer procCtrl.AssertExpectations(t)
	proc := &osProcess{
		pid:  123,
		proc: procCtrl,
	}
	wantErr := errors.New("")
	procCtrl.On("Wait").
		Once().
		Return(new(os.ProcessState), wantErr)
	procState, err := proc.Wait()
	require.Nil(t, procState)
	require.Same(t, err, wantErr)

	wantProcState := new(os.ProcessState)
	procCtrl.On("Wait").
		Once().
		Return(wantProcState, nil)
	procState, err = proc.Wait()
	require.Nil(t, err)
	require.Same(t, procState.(*osProcessState).state, wantProcState)
}

func TestNewOsProcessWithStarter(t *testing.T) {
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
