// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchain

import mock "github.com/stretchr/testify/mock"

// MockGroupStorage is an autogenerated mock type for the GroupStorage type
type MockGroupStorage struct {
	mock.Mock
}

// DeleteBlock provides a mock function with given fields: block
func (_m *MockGroupStorage) DeleteBlock(block Block) error {
	ret := _m.Called(block)

	var r0 error
	if rf, ok := ret.Get(0).(func(Block) error); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteBlockGroup provides a mock function with given fields: blocks
func (_m *MockGroupStorage) DeleteBlockGroup(blocks BlockGroup) error {
	ret := _m.Called(blocks)

	var r0 error
	if rf, ok := ret.Get(0).(func(BlockGroup) error); ok {
		r0 = rf(blocks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadBlocks provides a mock function with given fields: cursor, count
func (_m *MockGroupStorage) LoadBlocks(cursor interface{}, count int) (BlockGroup, interface{}, error) {
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
func (_m *MockGroupStorage) LoadLastBlock() (Block, error) {
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
func (_m *MockGroupStorage) StoreBlock(block Block) error {
	ret := _m.Called(block)

	var r0 error
	if rf, ok := ret.Get(0).(func(Block) error); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreBlockGroup provides a mock function with given fields: blocks
func (_m *MockGroupStorage) StoreBlockGroup(blocks BlockGroup) error {
	ret := _m.Called(blocks)

	var r0 error
	if rf, ok := ret.Get(0).(func(BlockGroup) error); ok {
		r0 = rf(blocks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
