// Code generated by mockery 2.9.0. DO NOT EDIT.

package automock

import (
	context "context"

	jwt "github.com/form3tech-oss/jwt-go"
	mock "github.com/stretchr/testify/mock"
)

// TokenVerifier is an autogenerated mock type for the TokenVerifier type
type TokenVerifier struct {
	mock.Mock
}

// Verify provides a mock function with given fields: ctx, token
func (_m *TokenVerifier) Verify(ctx context.Context, token string) (*jwt.MapClaims, error) {
	ret := _m.Called(ctx, token)

	var r0 *jwt.MapClaims
	if rf, ok := ret.Get(0).(func(context.Context, string) *jwt.MapClaims); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.MapClaims)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
