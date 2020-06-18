// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package xexec

import (
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// mockOsProcessCtrl is an autogenerated mock type for the osProcessCtrl type
type mockOsProcessCtrl struct {
	mock.Mock
}

// Kill provides a mock function with given fields:
func (_m *mockOsProcessCtrl) Kill() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Release provides a mock function with given fields:
func (_m *mockOsProcessCtrl) Release() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Signal provides a mock function with given fields: sig
func (_m *mockOsProcessCtrl) Signal(sig os.Signal) error {
	ret := _m.Called(sig)

	var r0 error
	if rf, ok := ret.Get(0).(func(os.Signal) error); ok {
		r0 = rf(sig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Wait provides a mock function with given fields:
func (_m *mockOsProcessCtrl) Wait() (*os.ProcessState, error) {
	ret := _m.Called()

	var r0 *os.ProcessState
	if rf, ok := ret.Get(0).(func() *os.ProcessState); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*os.ProcessState)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}