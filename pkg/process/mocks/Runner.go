// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	process "github.com/quadroops/goplugin/pkg/process"
	mock "github.com/stretchr/testify/mock"
)

// Runner is an autogenerated mock type for the Runner type
type Runner struct {
	mock.Mock
}

// Run provides a mock function with given fields: toWait, name, execCommand, args
func (_m *Runner) Run(toWait int, name string, execCommand string, args ...string) (<-chan process.Plugin, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, toWait, name, execCommand)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 <-chan process.Plugin
	if rf, ok := ret.Get(0).(func(int, string, string, ...string) <-chan process.Plugin); ok {
		r0 = rf(toWait, name, execCommand, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan process.Plugin)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string, string, ...string) error); ok {
		r1 = rf(toWait, name, execCommand, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}