// Code generated by mockery 2.9.0. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// TokenData is an autogenerated mock type for the TokenData type
type TokenData struct {
	mock.Mock
}

// Claims provides a mock function with given fields: v
func (_m *TokenData) Claims(v interface{}) error {
	ret := _m.Called(v)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
