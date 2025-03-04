// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// LabelDefinitionRepository is an autogenerated mock type for the LabelDefinitionRepository type
type LabelDefinitionRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, def
func (_m *LabelDefinitionRepository) Create(ctx context.Context, def model.LabelDefinition) error {
	ret := _m.Called(ctx, def)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelDefinition) error); ok {
		r0 = rf(ctx, def)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, tenant, key
func (_m *LabelDefinitionRepository) Exists(ctx context.Context, tenant string, key string) (bool, error) {
	ret := _m.Called(ctx, tenant, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, tenant, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByKey provides a mock function with given fields: ctx, tenant, key
func (_m *LabelDefinitionRepository) GetByKey(ctx context.Context, tenant string, key string) (*model.LabelDefinition, error) {
	ret := _m.Called(ctx, tenant, key)

	var r0 *model.LabelDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.LabelDefinition); ok {
		r0 = rf(ctx, tenant, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.LabelDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
