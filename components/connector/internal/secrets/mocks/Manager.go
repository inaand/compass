// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	context "context"

	corev1 "k8s.io/api/core/v1"

	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, name, options
func (_m *Manager) Get(ctx context.Context, name string, options v1.GetOptions) (*corev1.Secret, error) {
	ret := _m.Called(ctx, name, options)

	var r0 *corev1.Secret
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) *corev1.Secret); ok {
		r0 = rf(ctx, name, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Secret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, v1.GetOptions) error); ok {
		r1 = rf(ctx, name, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
