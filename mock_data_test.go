// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchain

import mock "github.com/stretchr/testify/mock"

// MockData is an autogenerated mock type for the Data type
type MockData struct {
	mock.Mock
}

// Equal provides a mock function with given fields: data
func (_m *MockData) Equal(data Data) bool {
	ret := _m.Called(data)

	var r0 bool
	if rf, ok := ret.Get(0).(func(Data) bool); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *MockData) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}