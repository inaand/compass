// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *Manager) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

	var r0 watch.Interface
	if rf, ok := ret.Get(0).(func(context.Context, v1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, v1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
