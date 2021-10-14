// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// ApplicationConverter is an autogenerated mock type for the ApplicationConverter type
type ApplicationConverter struct {
	mock.Mock
}

// CreateInputFromGraphQL provides a mock function with given fields: ctx, in
func (_m *ApplicationConverter) CreateInputFromGraphQL(ctx context.Context, in graphql.ApplicationRegisterInput) (model.ApplicationRegisterInput, error) {
	ret := _m.Called(ctx, in)

	var r0 model.ApplicationRegisterInput
	if rf, ok := ret.Get(0).(func(context.Context, graphql.ApplicationRegisterInput) model.ApplicationRegisterInput); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(model.ApplicationRegisterInput)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, graphql.ApplicationRegisterInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateInputJSONToGQL provides a mock function with given fields: in
func (_m *ApplicationConverter) CreateInputJSONToGQL(in string) (graphql.ApplicationRegisterInput, error) {
	ret := _m.Called(in)

	var r0 graphql.ApplicationRegisterInput
	if rf, ok := ret.Get(0).(func(string) graphql.ApplicationRegisterInput); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(graphql.ApplicationRegisterInput)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQL provides a mock function with given fields: in
func (_m *ApplicationConverter) ToGraphQL(in *model.Application) *graphql.Application {
	ret := _m.Called(in)

	var r0 *graphql.Application
	if rf, ok := ret.Get(0).(func(*model.Application) *graphql.Application); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.Application)
		}
	}

	return r0
}
