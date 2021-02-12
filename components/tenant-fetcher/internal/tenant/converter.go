package tenant

import (
	"github.com/kyma-incubator/compass/components/director/pkg/tenant"
	"github.com/kyma-incubator/compass/components/tenant-fetcher/internal/model"
)

type converter struct{}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToEntity(in model.TenantModel) tenant.Entity {
	return tenant.Entity{
		ID:             "",
		Name:           in.AccountId,
		ExternalTenant: in.AccountId,
		ProviderName:   "",
		Status:         "",
	}
}

func (c *converter) FromEntity(in *tenant.Entity) *model.TenantModel {
	if in == nil {
		return nil
	}
	return &model.TenantModel{
		AccountId: in.ExternalTenant,
	}
}
