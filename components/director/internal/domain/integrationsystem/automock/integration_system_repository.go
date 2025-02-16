// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// IntegrationSystemRepository is an autogenerated mock type for the IntegrationSystemRepository type
type IntegrationSystemRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *IntegrationSystemRepository) Create(ctx context.Context, item model.IntegrationSystem) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.IntegrationSystem) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *IntegrationSystemRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, id
func (_m *IntegrationSystemRepository) Exists(ctx context.Context, id string) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *IntegrationSystemRepository) Get(ctx context.Context, id string) (*model.IntegrationSystem, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.IntegrationSystem
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.IntegrationSystem); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.IntegrationSystem)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, pageSize, cursor
func (_m *IntegrationSystemRepository) List(ctx context.Context, pageSize int, cursor string) (model.IntegrationSystemPage, error) {
	ret := _m.Called(ctx, pageSize, cursor)

	var r0 model.IntegrationSystemPage
	if rf, ok := ret.Get(0).(func(context.Context, int, string) model.IntegrationSystemPage); ok {
		r0 = rf(ctx, pageSize, cursor)
	} else {
		r0 = ret.Get(0).(model.IntegrationSystemPage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *IntegrationSystemRepository) Update(ctx context.Context, _a1 model.IntegrationSystem) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.IntegrationSystem) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
