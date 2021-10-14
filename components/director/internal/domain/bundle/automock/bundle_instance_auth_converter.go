// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// BundleInstanceAuthConverter is an autogenerated mock type for the BundleInstanceAuthConverter type
type BundleInstanceAuthConverter struct {
	mock.Mock
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *BundleInstanceAuthConverter) MultipleToGraphQL(in []*model.BundleInstanceAuth) ([]*graphql.BundleInstanceAuth, error) {
	ret := _m.Called(in)

	var r0 []*graphql.BundleInstanceAuth
	if rf, ok := ret.Get(0).(func([]*model.BundleInstanceAuth) []*graphql.BundleInstanceAuth); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.BundleInstanceAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*model.BundleInstanceAuth) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQL provides a mock function with given fields: in
func (_m *BundleInstanceAuthConverter) ToGraphQL(in *model.BundleInstanceAuth) (*graphql.BundleInstanceAuth, error) {
	ret := _m.Called(in)

	var r0 *graphql.BundleInstanceAuth
	if rf, ok := ret.Get(0).(func(*model.BundleInstanceAuth) *graphql.BundleInstanceAuth); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.BundleInstanceAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.BundleInstanceAuth) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
