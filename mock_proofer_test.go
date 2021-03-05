// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchain

import mock "github.com/stretchr/testify/mock"

// MockProofer is an autogenerated mock type for the Proofer type
type MockProofer struct {
	mock.Mock
}

// Hash provides a mock function with given fields: block
func (_m *MockProofer) Hash(block Block) string {
	ret := _m.Called(block)

	var r0 string
	if rf, ok := ret.Get(0).(func(Block) string); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Validate provides a mock function with given fields: block
func (_m *MockProofer) Validate(block Block) bool {
	ret := _m.Called(block)

	var r0 bool
	if rf, ok := ret.Get(0).(func(Block) bool); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}