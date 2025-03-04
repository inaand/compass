// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	labelfilter "github.com/kyma-incubator/compass/components/director/internal/labelfilter"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// RuntimeRepository is an autogenerated mock type for the RuntimeRepository type
type RuntimeRepository struct {
	mock.Mock
}

// GetByFiltersAndID provides a mock function with given fields: ctx, tenant, id, filter
func (_m *RuntimeRepository) GetByFiltersAndID(ctx context.Context, tenant string, id string, filter []*labelfilter.LabelFilter) (*model.Runtime, error) {
	ret := _m.Called(ctx, tenant, id, filter)

	var r0 *model.Runtime
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []*labelfilter.LabelFilter) *model.Runtime); ok {
		r0 = rf(ctx, tenant, id, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, tenant, id, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOldestForFilters provides a mock function with given fields: ctx, tenant, filter
func (_m *RuntimeRepository) GetOldestForFilters(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter) (*model.Runtime, error) {
	ret := _m.Called(ctx, tenant, filter)

	var r0 *model.Runtime
	if rf, ok := ret.Get(0).(func(context.Context, string, []*labelfilter.LabelFilter) *model.Runtime); ok {
		r0 = rf(ctx, tenant, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, tenant, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, tenant, filter, pageSize, cursor
func (_m *RuntimeRepository) List(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error) {
	ret := _m.Called(ctx, tenant, filter, pageSize, cursor)

	var r0 *model.RuntimePage
	if rf, ok := ret.Get(0).(func(context.Context, string, []*labelfilter.LabelFilter, int, string) *model.RuntimePage); ok {
		r0 = rf(ctx, tenant, filter, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimePage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []*labelfilter.LabelFilter, int, string) error); ok {
		r1 = rf(ctx, tenant, filter, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
