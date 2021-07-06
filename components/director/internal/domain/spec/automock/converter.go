// Code generated by mockery 2.9.0. DO NOT EDIT.

package automock

import (
	spec "github.com/kyma-incubator/compass/components/director/internal/domain/spec"
	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: in
func (_m *Converter) FromEntity(in spec.Entity) (model.Spec, error) {
	ret := _m.Called(in)

	var r0 model.Spec
	if rf, ok := ret.Get(0).(func(spec.Entity) model.Spec); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(model.Spec)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(spec.Entity) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *Converter) ToEntity(in model.Spec) spec.Entity {
	ret := _m.Called(in)

	var r0 spec.Entity
	if rf, ok := ret.Get(0).(func(model.Spec) spec.Entity); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(spec.Entity)
	}

	return r0
}
