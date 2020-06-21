package xexec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProcessStateDelegater(t *testing.T) {
	t.Parallel()
	mockProcState := &MockProcessState{}
	procState := newProcessStateDelegater(mockProcState)

	mockProcState.On("ExitCode").
		Once().
		Return(111)
	require.Equal(t, procState.ExitCode(), 111)
	mockProcState.AssertExpectations(t)

	mockProcState.On("String").
		Once().
		Return("hello world")
	require.Equal(t, procState.String(), "hello world")
	mockProcState.AssertExpectations(t)

	mockProcState.On("Success").
		Once().
		Return(true)
	require.Equal(t, procState.Success(), true)
	mockProcState.AssertExpectations(t)

	mockProcState.On("Success").
		Once().
		Return(false)
	require.Equal(t, procState.Success(), false)
	mockProcState.AssertExpectations(t)
}
