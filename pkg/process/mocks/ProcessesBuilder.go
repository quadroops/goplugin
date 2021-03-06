// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	process "github.com/quadroops/goplugin/pkg/process"
	mock "github.com/stretchr/testify/mock"

	rxgo "github.com/reactivex/rxgo/v2"
)

// ProcessesBuilder is an autogenerated mock type for the ProcessesBuilder type
type ProcessesBuilder struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *ProcessesBuilder) Add(_a0 process.Plugin) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(process.Plugin) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0
func (_m *ProcessesBuilder) Get(_a0 string) (process.Plugin, error) {
	ret := _m.Called(_a0)

	var r0 process.Plugin
	if rf, ok := ret.Get(0).(func(string) process.Plugin); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(process.Plugin)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsExist provides a mock function with given fields: name
func (_m *ProcessesBuilder) IsExist(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Listen provides a mock function with given fields:
func (_m *ProcessesBuilder) Listen() (rxgo.Observable, error) {
	ret := _m.Called()

	var r0 rxgo.Observable
	if rf, ok := ret.Get(0).(func() rxgo.Observable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rxgo.Observable)
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

// Remove provides a mock function with given fields: _a0
func (_m *ProcessesBuilder) Remove(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Reset provides a mock function with given fields:
func (_m *ProcessesBuilder) Reset() {
	_m.Called()
}
