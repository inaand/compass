// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	tenantfetchersvc "github.com/kyma-incubator/compass/components/director/internal/tenantfetchersvc"
	mock "github.com/stretchr/testify/mock"
)

// TenantSubscriber is an autogenerated mock type for the TenantSubscriber type
type TenantSubscriber struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: ctx, tenantSubscriptionRequest, region
func (_m *TenantSubscriber) Subscribe(ctx context.Context, tenantSubscriptionRequest *tenantfetchersvc.TenantSubscriptionRequest, region string) error {
	ret := _m.Called(ctx, tenantSubscriptionRequest, region)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *tenantfetchersvc.TenantSubscriptionRequest, string) error); ok {
		r0 = rf(ctx, tenantSubscriptionRequest, region)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsubscribe provides a mock function with given fields: ctx, tenantSubscriptionRequest, region
func (_m *TenantSubscriber) Unsubscribe(ctx context.Context, tenantSubscriptionRequest *tenantfetchersvc.TenantSubscriptionRequest, region string) error {
	ret := _m.Called(ctx, tenantSubscriptionRequest, region)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *tenantfetchersvc.TenantSubscriptionRequest, string) error); ok {
		r0 = rf(ctx, tenantSubscriptionRequest, region)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
