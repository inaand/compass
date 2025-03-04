// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	systemauth "github.com/kyma-incubator/compass/components/director/internal/domain/systemauth"
	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: in
func (_m *Converter) FromEntity(in systemauth.Entity) (model.SystemAuth, error) {
	ret := _m.Called(in)

	var r0 model.SystemAuth
	if rf, ok := ret.Get(0).(func(systemauth.Entity) model.SystemAuth); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(model.SystemAuth)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(systemauth.Entity) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *Converter) ToEntity(in model.SystemAuth) (systemauth.Entity, error) {
	ret := _m.Called(in)

	var r0 systemauth.Entity
	if rf, ok := ret.Get(0).(func(model.SystemAuth) systemauth.Entity); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(systemauth.Entity)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.SystemAuth) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
