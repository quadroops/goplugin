// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SourceReader is an autogenerated mock type for the SourceReader type
type SourceReader struct {
	mock.Mock
}

// Read provides a mock function with given fields: sourceAddr
func (_m *SourceReader) Read(sourceAddr string) ([]byte, error) {
	ret := _m.Called(sourceAddr)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(sourceAddr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sourceAddr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
