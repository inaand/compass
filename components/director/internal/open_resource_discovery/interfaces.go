package ord

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

//go:generate mockery --exported --name=labelRepository --output=automock --outpkg=automock --case=underscore
type labelRepository interface {
	ListGlobalByKeyAndObjects(ctx context.Context, objectType model.LabelableObject, objectIDs []string, key string) ([]*model.Label, error)
}

// WebhookService is responsible for the service-layer Webhook operations.
//go:generate mockery --name=WebhookService --output=automock --outpkg=automock --case=underscore
type WebhookService interface {
	ListForApplication(ctx context.Context, applicationID string) ([]*model.Webhook, error)
}

// ApplicationService is responsible for the service-layer Application operations.
//go:generate mockery --name=ApplicationService --output=automock --outpkg=automock --case=underscore
type ApplicationService interface {
	ListGlobal(ctx context.Context, pageSize int, cursor string) (*model.ApplicationPage, error)
}

// BundleService is responsible for the service-layer Bundle operations.
//go:generate mockery --name=BundleService --output=automock --outpkg=automock --case=underscore
type BundleService interface {
	Create(ctx context.Context, applicationID string, in model.BundleCreateInput) (string, error)
	Update(ctx context.Context, id string, in model.BundleUpdateInput) error
	Delete(ctx context.Context, id string) error
	ListByApplicationIDNoPaging(ctx context.Context, appID string) ([]*model.Bundle, error)
}

// BundleReferenceService is responsible for the service-layer BundleReference operations.
//go:generate mockery --name=BundleReferenceService --output=automock --outpkg=automock --case=underscore
type BundleReferenceService interface {
	GetBundleIDsForObject(ctx context.Context, objectType model.BundleReferenceObjectType, objectID *string) ([]string, error)
}

// APIService is responsible for the service-layer API operations.
//go:generate mockery --name=APIService --output=automock --outpkg=automock --case=underscore
type APIService interface {
	Create(ctx context.Context, appID string, bundleID, packageID *string, in model.APIDefinitionInput, spec []*model.SpecInput, targetURLsPerBundle map[string]string, apiHash uint64, defaultBundleID string) (string, error)
	UpdateInManyBundles(ctx context.Context, id string, in model.APIDefinitionInput, specIn *model.SpecInput, defaultTargetURLPerBundle map[string]string, defaultTargetURLPerBundleToBeCreated map[string]string, bundleIDsToBeDeleted []string, apiHash uint64, defaultBundleID string) error
	Delete(ctx context.Context, id string) error
	ListByApplicationID(ctx context.Context, appID string) ([]*model.APIDefinition, error)
}

// EventService is responsible for the service-layer Event operations.
//go:generate mockery --name=EventService --output=automock --outpkg=automock --case=underscore
type EventService interface {
	Create(ctx context.Context, appID string, bundleID, packageID *string, in model.EventDefinitionInput, specs []*model.SpecInput, bundleIDs []string, eventHash uint64, defaultBundleID string) (string, error)
	UpdateInManyBundles(ctx context.Context, id string, in model.EventDefinitionInput, specIn *model.SpecInput, bundleIDsFromBundleReference, bundleIDsForCreation, bundleIDsForDeletion []string, eventHash uint64, defaultBundleID string) error
	Delete(ctx context.Context, id string) error
	ListByApplicationID(ctx context.Context, appID string) ([]*model.EventDefinition, error)
}

// SpecService is responsible for the service-layer Specification operations.
//go:generate mockery --name=SpecService --output=automock --outpkg=automock --case=underscore
type SpecService interface {
	CreateByReferenceObjectID(ctx context.Context, in model.SpecInput, objectType model.SpecReferenceObjectType, objectID string) (string, error)
	DeleteByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) error
	GetFetchRequest(ctx context.Context, specID string, objectType model.SpecReferenceObjectType) (*model.FetchRequest, error)
	ListByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) ([]*model.Spec, error)
	RefetchSpec(ctx context.Context, id string, objectType model.SpecReferenceObjectType) (*model.Spec, error)
}

// PackageService is responsible for the service-layer Package operations.
//go:generate mockery --name=PackageService --output=automock --outpkg=automock --case=underscore
type PackageService interface {
	Create(ctx context.Context, applicationID string, in model.PackageInput, pkgHash uint64) (string, error)
	Update(ctx context.Context, id string, in model.PackageInput, pkgHash uint64) error
	Delete(ctx context.Context, id string) error
	ListByApplicationID(ctx context.Context, appID string) ([]*model.Package, error)
}

// ProductService is responsible for the service-layer Product operations.
//go:generate mockery --name=ProductService --output=automock --outpkg=automock --case=underscore
type ProductService interface {
	Create(ctx context.Context, applicationID string, in model.ProductInput) (string, error)
	Update(ctx context.Context, id string, in model.ProductInput) error
	Delete(ctx context.Context, id string) error
	ListByApplicationID(ctx context.Context, appID string) ([]*model.Product, error)
}

// VendorService is responsible for the service-layer Vendor operations.
//go:generate mockery --name=VendorService --output=automock --outpkg=automock --case=underscore
type VendorService interface {
	Create(ctx context.Context, applicationID string, in model.VendorInput) (string, error)
	Update(ctx context.Context, id string, in model.VendorInput) error
	Delete(ctx context.Context, id string) error
	ListByApplicationID(ctx context.Context, appID string) ([]*model.Vendor, error)
}

// TombstoneService is responsible for the service-layer Tombstone operations.
//go:generate mockery --name=TombstoneService --output=automock --outpkg=automock --case=underscore
type TombstoneService interface {
	Create(ctx context.Context, applicationID string, in model.TombstoneInput) (string, error)
	Update(ctx context.Context, id string, in model.TombstoneInput) error
	Delete(ctx context.Context, id string) error
	ListByApplicationID(ctx context.Context, appID string) ([]*model.Tombstone, error)
}

// TenantService missing godoc
//go:generate mockery --name=TenantService --output=automock --outpkg=automock --case=underscore
type TenantService interface {
	GetLowestOwnerForResource(ctx context.Context, resourceType resource.Type, objectID string) (string, error)
}
