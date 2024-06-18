// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	entity "yandex-gophkeeper-client/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// LocalStorager is an autogenerated mock type for the LocalStorager type
type LocalStorager struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *LocalStorager) GetAll() (*[]entity.ItemDto, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 *[]entity.ItemDto
	var r1 error
	if rf, ok := ret.Get(0).(func() (*[]entity.ItemDto, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *[]entity.ItemDto); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]entity.ItemDto)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: item
func (_m *LocalStorager) Save(item *entity.ItemDto) error {
	ret := _m.Called(item)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.ItemDto) error); ok {
		r0 = rf(item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLocalStorager creates a new instance of LocalStorager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLocalStorager(t interface {
	mock.TestingT
	Cleanup(func())
}) *LocalStorager {
	mock := &LocalStorager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}