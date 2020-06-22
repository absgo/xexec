package xexec

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOsExec(t *testing.T) {
	t.Parallel()
	e := NewOsExec()
	require.IsType(t, &osExec{}, e)
}

func TestOsExec_StartProcess(t *testing.T) {
	t.Parallel()
	e := NewOsExec()
	cmdLine := getSleepCmdLine(10)
	proc, err := e.StartProcess(&ProcessConf{
		Name: cmdLine[0],
		Args: cmdLine,
	})
	require.Nil(t, err)
	require.Nil(t, proc.Kill())
	_, err = proc.Wait()
	require.Nil(t, err)
}

func TestOsExec_LookPath(t *testing.T) {
	t.Parallel()
	e := NewOsExec()
	tests := []struct {
		file string
	}{
		// Unix
		{
			file: "ls",
		},
		{
			file: "echo",
		},
		// Windows
		{
			file: "notepad",
		},
		{
			file: "calc",
		},
	}
	for _, tt := range tests {
		gotPath, err := e.LookPath(tt.file)
		wantPath, wantErr := exec.LookPath(tt.file)
		require.Equal(t, gotPath, wantPath)
		require.Equal(t, err, wantErr)
	}
}
