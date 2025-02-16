// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FetchRequestRepository is an autogenerated mock type for the FetchRequestRepository type
type FetchRequestRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, tenant, item
func (_m *FetchRequestRepository) Create(ctx context.Context, tenant string, item *model.FetchRequest) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.FetchRequest) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tenant, id, objectType
func (_m *FetchRequestRepository) Delete(ctx context.Context, tenant string, id string, objectType model.FetchRequestReferenceObjectType) error {
	ret := _m.Called(ctx, tenant, id, objectType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, model.FetchRequestReferenceObjectType) error); ok {
		r0 = rf(ctx, tenant, id, objectType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListByReferenceObjectIDs provides a mock function with given fields: ctx, tenant, objectType, objectIDs
func (_m *FetchRequestRepository) ListByReferenceObjectIDs(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectIDs []string) ([]*model.FetchRequest, error) {
	ret := _m.Called(ctx, tenant, objectType, objectIDs)

	var r0 []*model.FetchRequest
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FetchRequestReferenceObjectType, []string) []*model.FetchRequest); ok {
		r0 = rf(ctx, tenant, objectType, objectIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.FetchRequestReferenceObjectType, []string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
