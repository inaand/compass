// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// SpecService is an autogenerated mock type for the SpecService type
type SpecService struct {
	mock.Mock
}

// CreateByReferenceObjectID provides a mock function with given fields: ctx, in, objectType, objectID
func (_m *SpecService) CreateByReferenceObjectID(ctx context.Context, in model.SpecInput, objectType model.SpecReferenceObjectType, objectID string) (string, error) {
	ret := _m.Called(ctx, in, objectType, objectID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.SpecInput, model.SpecReferenceObjectType, string) string); ok {
		r0 = rf(ctx, in, objectType, objectID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.SpecInput, model.SpecReferenceObjectType, string) error); ok {
		r1 = rf(ctx, in, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByReferenceObjectID provides a mock function with given fields: ctx, objectType, objectID
func (_m *SpecService) DeleteByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) error {
	ret := _m.Called(ctx, objectType, objectID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.SpecReferenceObjectType, string) error); ok {
		r0 = rf(ctx, objectType, objectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFetchRequest provides a mock function with given fields: ctx, specID, objectType
func (_m *SpecService) GetFetchRequest(ctx context.Context, specID string, objectType model.SpecReferenceObjectType) (*model.FetchRequest, error) {
	ret := _m.Called(ctx, specID, objectType)

	var r0 *model.FetchRequest
	if rf, ok := ret.Get(0).(func(context.Context, string, model.SpecReferenceObjectType) *model.FetchRequest); ok {
		r0 = rf(ctx, specID, objectType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.SpecReferenceObjectType) error); ok {
		r1 = rf(ctx, specID, objectType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByReferenceObjectID provides a mock function with given fields: ctx, objectType, objectID
func (_m *SpecService) ListByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) ([]*model.Spec, error) {
	ret := _m.Called(ctx, objectType, objectID)

	var r0 []*model.Spec
	if rf, ok := ret.Get(0).(func(context.Context, model.SpecReferenceObjectType, string) []*model.Spec); ok {
		r0 = rf(ctx, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Spec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.SpecReferenceObjectType, string) error); ok {
		r1 = rf(ctx, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefetchSpec provides a mock function with given fields: ctx, id, objectType
func (_m *SpecService) RefetchSpec(ctx context.Context, id string, objectType model.SpecReferenceObjectType) (*model.Spec, error) {
	ret := _m.Called(ctx, id, objectType)

	var r0 *model.Spec
	if rf, ok := ret.Get(0).(func(context.Context, string, model.SpecReferenceObjectType) *model.Spec); ok {
		r0 = rf(ctx, id, objectType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Spec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.SpecReferenceObjectType) error); ok {
		r1 = rf(ctx, id, objectType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
