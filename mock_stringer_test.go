// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchain

import mock "github.com/stretchr/testify/mock"

// MockStringer is an autogenerated mock type for the Stringer type
type MockStringer struct {
	mock.Mock
}

// String provides a mock function with given fields:
func (_m *MockStringer) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
