// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// TenantConverter is an autogenerated mock type for the TenantConverter type
type TenantConverter struct {
	mock.Mock
}

// MultipleInputToGraphQLInput provides a mock function with given fields: _a0
func (_m *TenantConverter) MultipleInputToGraphQLInput(_a0 []model.BusinessTenantMappingInput) []graphql.BusinessTenantMappingInput {
	ret := _m.Called(_a0)

	var r0 []graphql.BusinessTenantMappingInput
	if rf, ok := ret.Get(0).(func([]model.BusinessTenantMappingInput) []graphql.BusinessTenantMappingInput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]graphql.BusinessTenantMappingInput)
		}
	}

	return r0
}

// ToGraphQLInput provides a mock function with given fields: _a0
func (_m *TenantConverter) ToGraphQLInput(_a0 model.BusinessTenantMappingInput) graphql.BusinessTenantMappingInput {
	ret := _m.Called(_a0)

	var r0 graphql.BusinessTenantMappingInput
	if rf, ok := ret.Get(0).(func(model.BusinessTenantMappingInput) graphql.BusinessTenantMappingInput); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(graphql.BusinessTenantMappingInput)
	}

	return r0
}