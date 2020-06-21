package xexec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOsExec(t *testing.T) {
	t.Parallel()
	e := NewOsExec()
	require.IsType(t, &osExec{}, e)
}

func TestExec_StartProcess(t *testing.T) {
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
