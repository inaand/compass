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

// GraphQLToModel provides a mock function with given fields: obj, tenantID
func (_m *ApplicationConverter) GraphQLToModel(obj *graphql.Application, tenantID string) *model.Application {
	ret := _m.Called(obj, tenantID)

	var r0 *model.Application
	if rf, ok := ret.Get(0).(func(*graphql.Application, string) *model.Application); ok {
		r0 = rf(obj, tenantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	return r0
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *ApplicationConverter) MultipleToGraphQL(in []*model.Application) []*graphql.Application {
	ret := _m.Called(in)

	var r0 []*graphql.Application
	if rf, ok := ret.Get(0).(func([]*model.Application) []*graphql.Application); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.Application)
		}
	}

	return r0
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

// UpdateInputFromGraphQL provides a mock function with given fields: in
func (_m *ApplicationConverter) UpdateInputFromGraphQL(in graphql.ApplicationUpdateInput) model.ApplicationUpdateInput {
	ret := _m.Called(in)

	var r0 model.ApplicationUpdateInput
	if rf, ok := ret.Get(0).(func(graphql.ApplicationUpdateInput) model.ApplicationUpdateInput); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(model.ApplicationUpdateInput)
	}

	return r0
}
