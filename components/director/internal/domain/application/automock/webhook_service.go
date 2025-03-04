// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// WebhookService is an autogenerated mock type for the WebhookService type
type WebhookService struct {
	mock.Mock
}

// ListAllApplicationWebhooks provides a mock function with given fields: ctx, applicationTemplateID
func (_m *WebhookService) ListAllApplicationWebhooks(ctx context.Context, applicationTemplateID string) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, applicationTemplateID)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Webhook); ok {
		r0 = rf(ctx, applicationTemplateID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, applicationTemplateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
