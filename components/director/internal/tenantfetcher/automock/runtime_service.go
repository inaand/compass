// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	labelfilter "github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// RuntimeService is an autogenerated mock type for the RuntimeService type
type RuntimeService struct {
	mock.Mock
}

// ListByFilters provides a mock function with given fields: ctx, filters
func (_m *RuntimeService) ListByFilters(ctx context.Context, filters []*labelfilter.LabelFilter) ([]*model.Runtime, error) {
	ret := _m.Called(ctx, filters)

	var r0 []*model.Runtime
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter) []*model.Runtime); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
