// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	apperrors "github.com/kyma-incubator/compass/components/connector/internal/apperrors"
	mock "github.com/stretchr/testify/mock"

	types "k8s.io/apimachinery/pkg/types"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Get provides a mock function with given fields: name
func (_m *Repository) Get(name types.NamespacedName) (map[string][]byte, apperrors.AppError) {
	ret := _m.Called(name)

	var r0 map[string][]byte
	if rf, ok := ret.Get(0).(func(types.NamespacedName) map[string][]byte); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]byte)
		}
	}

	var r1 apperrors.AppError
	if rf, ok := ret.Get(1).(func(types.NamespacedName) apperrors.AppError); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(apperrors.AppError)
		}
	}

	return r0, r1
}
