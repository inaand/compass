package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/domain/eventing"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/normalizer"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	labelpkg "github.com/kyma-incubator/compass/components/director/internal/domain/label"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/timestamp"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/pkg/errors"
)

const (
	intSysKey = "integrationSystemID"
	nameKey   = "name"
)

type repoCreatorFunc func(ctx context.Context, application *model.Application) error

//go:generate mockery --name=ApplicationRepository --output=automock --outpkg=automock --case=underscore
type ApplicationRepository interface {
	Exists(ctx context.Context, tenant, id string) (bool, error)
	GetByID(ctx context.Context, tenant, id string) (*model.Application, error)
	GetGlobalByID(ctx context.Context, id string) (*model.Application, error)
	List(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.ApplicationPage, error)
	ListAll(ctx context.Context, tenant string) ([]*model.Application, error)
	ListGlobal(ctx context.Context, pageSize int, cursor string) (*model.ApplicationPage, error)
	ListByScenarios(ctx context.Context, tenantID uuid.UUID, scenarios []string, pageSize int, cursor string, hidingSelectors map[string][]string) (*model.ApplicationPage, error)
	Create(ctx context.Context, item *model.Application) error
	Update(ctx context.Context, item *model.Application) error
	TechnicalUpdate(ctx context.Context, item *model.Application) error
	Delete(ctx context.Context, tenant, id string) error
	DeleteGlobal(ctx context.Context, id string) error
}

//go:generate mockery --name=LabelRepository --output=automock --outpkg=automock --case=underscore
type LabelRepository interface {
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
	ListForObject(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error)
	Delete(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) error
	DeleteAll(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) error
	Upsert(ctx context.Context, label *model.Label) error
}

//go:generate mockery --name=WebhookRepository --output=automock --outpkg=automock --case=underscore
type WebhookRepository interface {
	CreateMany(ctx context.Context, items []*model.Webhook) error
}

//go:generate mockery --name=RuntimeRepository --output=automock --outpkg=automock --case=underscore
type RuntimeRepository interface {
	Exists(ctx context.Context, tenant, id string) (bool, error)
	ListAll(ctx context.Context, tenantID string, filter []*labelfilter.LabelFilter) ([]*model.Runtime, error)
}

//go:generate mockery --name=IntegrationSystemRepository --output=automock --outpkg=automock --case=underscore
type IntegrationSystemRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
}

//go:generate mockery --name=LabelUpsertService --output=automock --outpkg=automock --case=underscore
type LabelUpsertService interface {
	UpsertMultipleLabels(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, labels map[string]interface{}) error
	UpsertLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) error
}

//go:generate mockery --name=ScenariosDefinitionService --output=automock --outpkg=automock --case=underscore
type ScenariosDefinitionService interface {
	EnsureScenariosLabelDefinitionExists(ctx context.Context, tenant string) error
	AddDefaultScenarioIfEnabled(ctx context.Context, labels *map[string]interface{})
}

//go:generate mockery --name=BundleInstanceAuthService --output=automock --outpkg=automock --case=underscore
type BundleInstanceAuthService interface {
	AssociateBundleInstanceAuthForNewApplicationScenarios(ctx context.Context, existingScenarios, inputScenarios []string, appId string) error
	GetForAppAndAnyMatchingScenarios(ctx context.Context, appId string, scenarios []string) ([]*model.BundleInstanceAuth, error)
}

//go:generate mockery --name=ScenariosService --output=automock --outpkg=automock --case=underscore
type ScenariosService interface {
	GetScenarioNamesForApplication(ctx context.Context, applicationID string) ([]string, error)
}

//go:generate mockery --name=UIDService --output=automock --outpkg=automock --case=underscore
type UIDService interface {
	Generate() string
}

//go:generate mockery --name=ApplicationHideCfgProvider --output=automock --outpkg=automock --case=underscore
type ApplicationHideCfgProvider interface {
	GetApplicationHideSelectors() (map[string][]string, error)
}

type service struct {
	appNameNormalizer  normalizer.Normalizator
	appHideCfgProvider ApplicationHideCfgProvider

	appRepo       ApplicationRepository
	webhookRepo   WebhookRepository
	labelRepo     LabelRepository
	runtimeRepo   RuntimeRepository
	intSystemRepo IntegrationSystemRepository

	labelUpsertService         LabelUpsertService
	scenariosDefinitionService ScenariosDefinitionService
	scenariosService           ScenariosService
	bundleInstanceAuthService  BundleInstanceAuthService
	uidService                 UIDService
	bndlService                BundleService
	timestampGen               timestamp.Generator
}

func NewService(appNameNormalizer normalizer.Normalizator, appHideCfgProvider ApplicationHideCfgProvider, app ApplicationRepository, webhook WebhookRepository, runtimeRepo RuntimeRepository, labelRepo LabelRepository, intSystemRepo IntegrationSystemRepository, labelUpsertService LabelUpsertService, scenariosDefinitionService ScenariosDefinitionService, scenariosService ScenariosService, bndlService BundleService, uidService UIDService, bundleInstanceAuthService BundleInstanceAuthService) *service {
	return &service{
		appNameNormalizer:          appNameNormalizer,
		appHideCfgProvider:         appHideCfgProvider,
		appRepo:                    app,
		webhookRepo:                webhook,
		runtimeRepo:                runtimeRepo,
		labelRepo:                  labelRepo,
		intSystemRepo:              intSystemRepo,
		labelUpsertService:         labelUpsertService,
		scenariosDefinitionService: scenariosDefinitionService,
		scenariosService:           scenariosService,
		bndlService:                bndlService,
		uidService:                 uidService,
		bundleInstanceAuthService:  bundleInstanceAuthService,
		timestampGen:               timestamp.DefaultGenerator(),
	}
}

func (s *service) List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.ApplicationPage, error) {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.appRepo.List(ctx, appTenant, filter, pageSize, cursor)
}

func (s *service) ListGlobal(ctx context.Context, pageSize int, cursor string) (*model.ApplicationPage, error) {
	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.appRepo.ListGlobal(ctx, pageSize, cursor)
}

func (s *service) ListByRuntimeID(ctx context.Context, runtimeID uuid.UUID, pageSize int, cursor string) (*model.ApplicationPage, error) {
	tenantID, err := tenant.LoadFromContext(ctx)

	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, apperrors.NewInvalidDataError("tenantID is not UUID")
	}

	exist, err := s.runtimeRepo.Exists(ctx, tenantID, runtimeID.String())
	if err != nil {
		return nil, errors.Wrap(err, "while checking if runtime exits")
	}

	if !exist {
		return nil, apperrors.NewInvalidDataError("runtime does not exist")
	}

	scenariosLabel, err := s.labelRepo.GetByKey(ctx, tenantID, model.RuntimeLabelableObject, runtimeID.String(), model.ScenariosKey)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return &model.ApplicationPage{
				Data:       []*model.Application{},
				PageInfo:   &pagination.Page{},
				TotalCount: 0,
			}, nil
		}
		return nil, errors.Wrap(err, "while getting scenarios for runtime")
	}

	scenarios, err := labelpkg.ValueToStringsSlice(scenariosLabel.Value)
	if err != nil {
		return nil, errors.Wrap(err, "while converting scenarios labels")
	}
	if len(scenarios) == 0 {
		return &model.ApplicationPage{
			Data:       []*model.Application{},
			TotalCount: 0,
			PageInfo: &pagination.Page{
				StartCursor: "",
				EndCursor:   "",
				HasNextPage: false,
			},
		}, nil
	}

	hidingSelectors, err := s.appHideCfgProvider.GetApplicationHideSelectors()
	if err != nil {
		return nil, errors.Wrap(err, "while getting application hide selectors from config")
	}

	return s.appRepo.ListByScenarios(ctx, tenantUUID, scenarios, pageSize, cursor, hidingSelectors)
}

func (s *service) Get(ctx context.Context, id string) (*model.Application, error) {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	app, err := s.appRepo.GetByID(ctx, appTenant, id)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting Application with id %s", id)
	}

	return app, nil
}

func (s *service) Exist(ctx context.Context, id string) (bool, error) {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return false, errors.Wrapf(err, "while loading tenant from context")
	}

	exist, err := s.appRepo.Exists(ctx, appTenant, id)
	if err != nil {
		return false, errors.Wrapf(err, "while getting Application with ID %s", id)
	}

	return exist, nil
}

func (s *service) Create(ctx context.Context, in model.ApplicationRegisterInput) (string, error) {
	creator := func(ctx context.Context, application *model.Application) (err error) {
		err = s.appRepo.Create(ctx, application)
		if err != nil {
			return errors.Wrapf(err, "while creating Application with name %s", application.Name)
		}
		return
	}

	return s.genericCreate(ctx, in, creator)
}

func (s *service) CreateFromTemplate(ctx context.Context, in model.ApplicationRegisterInput, appTemplateId *string) (string, error) {
	creator := func(ctx context.Context, application *model.Application) (err error) {
		application.ApplicationTemplateID = appTemplateId
		err = s.appRepo.Create(ctx, application)
		if err != nil {
			return errors.Wrapf(err, "while creating Application with name %s from template", application.Name)
		}
		return
	}

	return s.genericCreate(ctx, in, creator)
}

func (s *service) Update(ctx context.Context, id string, in model.ApplicationUpdateInput) error {
	exists, err := s.ensureIntSysExists(ctx, in.IntegrationSystemID)
	if err != nil {
		return errors.Wrap(err, "while validating Integration System ID")
	}

	if !exists {
		return apperrors.NewNotFoundError(resource.IntegrationSystem, *in.IntegrationSystemID)
	}

	app, err := s.Get(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "while getting Application with id %s", id)
	}

	app.SetFromUpdateInput(in, s.timestampGen())

	err = s.appRepo.Update(ctx, app)
	if err != nil {
		return errors.Wrapf(err, "while updating Application with id %s", id)
	}

	if in.IntegrationSystemID != nil {
		intSysLabel := createLabel(intSysKey, *in.IntegrationSystemID, id)
		err = s.SetLabel(ctx, intSysLabel)
		if err != nil {
			return errors.Wrapf(err, "while setting the integration system label for %s with id %s", intSysLabel.ObjectType, intSysLabel.ObjectID)
		}
		log.C(ctx).Debugf("Successfully set Label for %s with id %s", intSysLabel.ObjectType, intSysLabel.ObjectID)
	}

	label := createLabel(nameKey, s.appNameNormalizer.Normalize(app.Name), app.ID)
	err = s.SetLabel(ctx, label)
	if err != nil {
		return errors.Wrap(err, "while setting application name label")
	}
	log.C(ctx).Debugf("Successfully set Label for Application with id %s", app.ID)
	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	appExists, err := s.appRepo.Exists(ctx, appTenant, id)
	if err != nil {
		return errors.Wrap(err, "while checking Application existence")
	}
	if !appExists {
		return fmt.Errorf("application with ID %s doesn't exist", id)
	}

	scenarios, err := s.scenariosService.GetScenarioNamesForApplication(ctx, id)
	if err != nil {
		return errors.Wrap(err, "while fetching scenario names for application")
	}

	validScenarios := removeDefaultScenario(scenarios)
	if len(validScenarios) > 0 {
		runtimes, err := s.getRuntimeNamesForScenarios(ctx, appTenant, validScenarios)
		if err != nil {
			return err
		}

		if len(runtimes) > 0 {
			application, err := s.appRepo.GetByID(ctx, appTenant, id)
			if err != nil {
				return errors.Wrapf(err, "while getting application with id %s", id)
			}
			msg := fmt.Sprintf("System %s is still used and cannot be deleted. Unassign the system from the following formations first: %s. Then, unassign the system from the following runtimes, too: %s", application.Name, strings.Join(validScenarios, ", "), strings.Join(runtimes, ", "))
			return apperrors.NewInvalidOperationError(msg)
		}
	}

	err = s.appRepo.Delete(ctx, appTenant, id)
	if err != nil {
		return errors.Wrapf(err, "while deleting Application with id %s", id)
	}

	return nil
}

func (s *service) SetLabel(ctx context.Context, labelInput *model.LabelInput) error {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	appExists, err := s.appRepo.Exists(ctx, appTenant, labelInput.ObjectID)
	if err != nil {
		return errors.Wrap(err, "while checking Application existence")
	}

	if !appExists {
		return apperrors.NewNotFoundError(resource.Application, labelInput.ObjectID)
	}

	if labelInput.Key == model.ScenariosKey {
		if err = s.associateRelatedBundleInstanceAuthWithNewScenarios(ctx, labelInput); err != nil {
			return err
		}
	}

	err = s.labelUpsertService.UpsertLabel(ctx, appTenant, labelInput)
	if err != nil {
		return errors.Wrapf(err, "while creating label for Application")
	}

	return nil
}

func (s *service) associateRelatedBundleInstanceAuthWithNewScenarios(ctx context.Context, labelInput *model.LabelInput) error {
	inputScenarios, err := labelpkg.GetScenariosFromValueAsStringSlice(labelInput.Value)
	if err != nil {
		return errors.Wrap(err, "while parsing label value")
	}

	if len(inputScenarios) < 1 {
		return nil
	}

	existingScenarios, err := s.scenariosService.GetScenarioNamesForApplication(ctx, labelInput.ObjectID)
	if err != nil {
		return err
	}

	scenariosToRemove := str.SubstractSlice(existingScenarios, inputScenarios)
	if err = s.validateNoBundleInstanceAuthsExist(ctx, labelInput.ObjectID, scenariosToRemove); err != nil {
		return err
	}

	if err := s.bundleInstanceAuthService.AssociateBundleInstanceAuthForNewApplicationScenarios(ctx, existingScenarios, inputScenarios, labelInput.ObjectID); err != nil {
		return errors.Wrap(err, "while associating existing bundle instance auths with the new scenarios")
	}
	return nil
}

func (s *service) GetLabel(ctx context.Context, applicationID string, key string) (*model.Label, error) {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	appExists, err := s.appRepo.Exists(ctx, appTenant, applicationID)
	if err != nil {
		return nil, errors.Wrap(err, "while checking Application existence")
	}
	if !appExists {
		return nil, fmt.Errorf("application with ID %s doesn't exist", applicationID)
	}

	label, err := s.labelRepo.GetByKey(ctx, appTenant, model.ApplicationLabelableObject, applicationID, key)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Application")
	}

	return label, nil
}

func (s *service) ListLabels(ctx context.Context, applicationID string) (map[string]*model.Label, error) {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	appExists, err := s.appRepo.Exists(ctx, appTenant, applicationID)
	if err != nil {
		return nil, errors.Wrap(err, "while checking Application existence")
	}

	if !appExists {
		return nil, fmt.Errorf("application with ID %s doesn't exist", applicationID)
	}

	labels, err := s.labelRepo.ListForObject(ctx, appTenant, model.ApplicationLabelableObject, applicationID)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Application")
	}

	return labels, nil
}

func (s *service) DeleteLabel(ctx context.Context, applicationID string, key string) error {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	if key == model.ScenariosKey {
		scenarios, err := s.scenariosService.GetScenarioNamesForApplication(ctx, applicationID)
		if err != nil {
			return err
		}

		if err = s.validateNoBundleInstanceAuthsExist(ctx, applicationID, scenarios); err != nil {
			return err
		}
	}

	err = s.labelRepo.Delete(ctx, appTenant, model.ApplicationLabelableObject, applicationID, key)
	if err != nil {
		return errors.Wrapf(err, "while deleting Application label")
	}

	return nil
}

func (s *service) createRelatedResources(ctx context.Context, in model.ApplicationRegisterInput, tenant string, applicationID string) error {
	var err error
	var webhooks []*model.Webhook
	for _, item := range in.Webhooks {
		webhooks = append(webhooks, item.ToApplicationWebhook(s.uidService.Generate(), &tenant, applicationID))
	}
	err = s.webhookRepo.CreateMany(ctx, webhooks)
	if err != nil {
		return errors.Wrapf(err, "while creating Webhooks for application")
	}

	return nil
}

func (s *service) validateNoBundleInstanceAuthsExist(ctx context.Context, appId string, scenariosToRemove []string) error {
	bndlAuths, err := s.bundleInstanceAuthService.GetForAppAndAnyMatchingScenarios(ctx, appId, scenariosToRemove)
	if err != nil {
		return errors.Wrapf(err, "while getting existing bundle instance auths for old scenarios")
	}

	if len(bndlAuths) == 0 {
		return nil
	}

	app, err := s.Get(ctx, appId)
	if err != nil {
		return errors.Wrapf(err, "while getting application")
	}

	authCtx := make([]*string, 0, len(bndlAuths))
	for _, auth := range bndlAuths {
		authCtx = append(authCtx, auth.Context)
	}

	return apperrors.NewScenarioUnassignWhenCredentialsExistsError(resource.Application, app.Name, authCtx)
}

func (s *service) genericCreate(ctx context.Context, in model.ApplicationRegisterInput, repoCreatorFunc repoCreatorFunc) (string, error) {
	appTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return "", err
	}
	log.C(ctx).Debugf("Loaded Application Tenant %s from context", appTenant)

	applications, err := s.appRepo.ListAll(ctx, appTenant)
	if err != nil {
		return "", err
	}

	normalizedName := s.appNameNormalizer.Normalize(in.Name)
	for _, app := range applications {
		if normalizedName == s.appNameNormalizer.Normalize(app.Name) {
			return "", apperrors.NewNotUniqueNameError(resource.Application)
		}
	}

	exists, err := s.ensureIntSysExists(ctx, in.IntegrationSystemID)
	if err != nil {
		return "", errors.Wrap(err, "while ensuring integration system exists")
	}

	if !exists {
		return "", apperrors.NewNotFoundError(resource.IntegrationSystem, *in.IntegrationSystemID)
	}

	id := s.uidService.Generate()
	log.C(ctx).Debugf("ID %s generated for Application with name %s", id, in.Name)

	app := in.ToApplication(s.timestampGen(), id, appTenant)

	err = repoCreatorFunc(ctx, app)
	if err != nil {
		return "", err
	}

	log.C(ctx).Debugf("Ensuring Scenarios label definition exists for Tenant %s", appTenant)
	err = s.scenariosDefinitionService.EnsureScenariosLabelDefinitionExists(ctx, appTenant)
	if err != nil {
		return "", err
	}

	s.scenariosDefinitionService.AddDefaultScenarioIfEnabled(ctx, &in.Labels)

	if in.Labels == nil {
		in.Labels = map[string]interface{}{}
	}
	in.Labels[intSysKey] = ""
	if in.IntegrationSystemID != nil {
		in.Labels[intSysKey] = *in.IntegrationSystemID
	}
	in.Labels[nameKey] = normalizedName

	err = s.labelUpsertService.UpsertMultipleLabels(ctx, appTenant, model.ApplicationLabelableObject, id, in.Labels)
	if err != nil {
		return id, errors.Wrapf(err, "while creating multiple labels for Application with id %s", id)
	}

	err = s.createRelatedResources(ctx, in, app.Tenant, app.ID)
	if err != nil {
		return "", errors.Wrapf(err, "while creating related resources for Application with id %s", id)
	}

	if in.Bundles != nil {
		err = s.bndlService.CreateMultiple(ctx, id, in.Bundles)
		if err != nil {
			return "", errors.Wrapf(err, "while creating related Bundle resources for Application with id %s", id)
		}
	}

	return id, nil
}

func createLabel(key string, value string, objectID string) *model.LabelInput {
	return &model.LabelInput{
		Key:        key,
		Value:      value,
		ObjectID:   objectID,
		ObjectType: model.ApplicationLabelableObject,
	}
}

func (s *service) ensureIntSysExists(ctx context.Context, id *string) (bool, error) {
	if id == nil {
		return true, nil
	}

	log.C(ctx).Infof("Ensuring Integration System with id %s exists", *id)
	exists, err := s.intSystemRepo.Exists(ctx, *id)
	if err != nil {
		return false, err
	}

	if !exists {
		log.C(ctx).Infof("Integration System with id %s does not exist", *id)
		return false, nil
	}
	log.C(ctx).Infof("Integration System with id %s exists", *id)
	return true, nil
}

func (s *service) getRuntimeNamesForScenarios(ctx context.Context, tenant string, scenarios []string) ([]string, error) {
	scenariosQuery := eventing.BuildQueryForScenarios(scenarios)
	runtimeScenariosFilter := []*labelfilter.LabelFilter{labelfilter.NewForKeyWithQuery(model.ScenariosKey, scenariosQuery)}

	log.C(ctx).Debugf("Listing runtimes matching the query %s", scenariosQuery)
	runtimes, err := s.runtimeRepo.ListAll(ctx, tenant, runtimeScenariosFilter)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting runtimes")
	}

	var runtimesNames []string
	for _, r := range runtimes {
		runtimesNames = append(runtimesNames, r.Name)
	}

	return runtimesNames, nil
}

func removeDefaultScenario(scenarios []string) []string {
	defaultScenarioIndex := -1
	for idx, scenario := range scenarios {
		if scenario == model.ScenariosDefaultValue[0] {
			defaultScenarioIndex = idx
			break
		}
	}

	if defaultScenarioIndex >= 0 {
		return append(scenarios[:defaultScenarioIndex], scenarios[defaultScenarioIndex+1:]...)
	}

	return scenarios
}
