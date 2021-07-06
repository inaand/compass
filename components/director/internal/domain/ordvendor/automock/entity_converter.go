// Code generated by mockery 2.9.0. DO NOT EDIT.

package automock

import (
	ordvendor "github.com/kyma-incubator/compass/components/director/internal/domain/ordvendor"
	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *ordvendor.Entity) (*model.Vendor, error) {
	ret := _m.Called(entity)

	var r0 *model.Vendor
	if rf, ok := ret.Get(0).(func(*ordvendor.Entity) *model.Vendor); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Vendor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ordvendor.Entity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.Vendor) *ordvendor.Entity {
	ret := _m.Called(in)

	var r0 *ordvendor.Entity
	if rf, ok := ret.Get(0).(func(*model.Vendor) *ordvendor.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ordvendor.Entity)
		}
	}

	return r0
}
