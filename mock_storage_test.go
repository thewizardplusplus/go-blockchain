// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchain

import mock "github.com/stretchr/testify/mock"

// MockStorage is an autogenerated mock type for the Storage type
type MockStorage struct {
	mock.Mock
}

// LoadBlocks provides a mock function with given fields: cursor, count
func (_m *MockStorage) LoadBlocks(cursor interface{}, count int) (BlockGroup, interface{}, error) {
	ret := _m.Called(cursor, count)

	var r0 BlockGroup
	if rf, ok := ret.Get(0).(func(interface{}, int) BlockGroup); ok {
		r0 = rf(cursor, count)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(BlockGroup)
		}
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(interface{}, int) interface{}); ok {
		r1 = rf(cursor, count)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(interface{}, int) error); ok {
		r2 = rf(cursor, count)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// LoadLastBlock provides a mock function with given fields:
func (_m *MockStorage) LoadLastBlock() (Block, error) {
	ret := _m.Called()

	var r0 Block
	if rf, ok := ret.Get(0).(func() Block); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(Block)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreBlock provides a mock function with given fields: block
func (_m *MockStorage) StoreBlock(block Block) error {
	ret := _m.Called(block)

	var r0 error
	if rf, ok := ret.Get(0).(func(Block) error); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
