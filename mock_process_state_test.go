// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package xexec

import mock "github.com/stretchr/testify/mock"

// MockProcessState is an autogenerated mock type for the ProcessState type
type MockProcessState struct {
	mock.Mock
}

// ExitCode provides a mock function with given fields:
func (_m *MockProcessState) ExitCode() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *MockProcessState) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Success provides a mock function with given fields:
func (_m *MockProcessState) Success() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
