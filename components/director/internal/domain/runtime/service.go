package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kyma-incubator/compass/components/director/internal/domain/eventing"
	"regexp"
	"strings"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/kyma-incubator/compass/components/director/internal/domain/label"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/pkg/errors"
)

const IsNormalizedLabel = "isNormalized"

//go:generate mockery --name=RuntimeRepository --output=automock --outpkg=automock --case=underscore
type RuntimeRepository interface {
	Exists(ctx context.Context, tenant, id string) (bool, error)
	GetByID(ctx context.Context, tenant, id string) (*model.Runtime, error)
	GetByFiltersGlobal(ctx context.Context, filter []*labelfilter.LabelFilter) (*model.Runtime, error)
	List(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error)
	Create(ctx context.Context, item *model.Runtime) error
	Update(ctx context.Context, item *model.Runtime) error
	UpdateTenantID(ctx context.Context, runtimeID, newTenantID string) error
	Delete(ctx context.Context, tenant, id string) error
}

//go:generate mockery --name=LabelRepository --output=automock --outpkg=automock --case=underscore
type LabelRepository interface {
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
	ListForObject(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error)
	Delete(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) error
	DeleteAll(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) error
	DeleteByKeyNegationPattern(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, labelKeyPattern string) error
	Upsert(ctx context.Context, label *model.Label) error
}

//go:generate mockery --name=LabelUpsertService --output=automock --outpkg=automock --case=underscore
type LabelUpsertService interface {
	UpsertMultipleLabels(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, labels map[string]interface{}) error
	UpsertLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) error
}

//go:generate mockery --name=ScenariosService --output=automock --outpkg=automock --case=underscore
type ScenariosService interface {
	EnsureScenariosLabelDefinitionExists(ctx context.Context, tenant string) error
	AddDefaultScenarioIfEnabled(ctx context.Context, labels *map[string]interface{})
}

//go:generate mockery --name=ScenarioAssignmentEngine --output=automock --outpkg=automock --case=underscore
type ScenarioAssignmentEngine interface {
	GetScenariosForSelectorLabels(ctx context.Context, inputLabels map[string]string) ([]string, error)
	MergeScenariosFromInputLabelsAndAssignments(ctx context.Context, inputLabels map[string]interface{}) ([]interface{}, error)
	MergeScenarios(baseScenarios, scenariosToDelete, scenariosToAdd []interface{}) []interface{}
}

type ApplicationRepository interface {
	ListAllByLabelFilter(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter) ([]*model.Application, error)
}

//go:generate mockery --name=UIDService --output=automock --outpkg=automock --case=underscore
type UIDService interface {
	Generate() string
}

type service struct {
	repo      RuntimeRepository
	labelRepo LabelRepository

	labelUpsertService       LabelUpsertService
	uidService               UIDService
	scenariosService         ScenariosService
	scenarioAssignmentEngine ScenarioAssignmentEngine
	appRepo                  ApplicationRepository

	protectedLabelPattern string
}

func NewService(repo RuntimeRepository,
	labelRepo LabelRepository,
	scenariosService ScenariosService,
	labelUpsertService LabelUpsertService,
	uidService UIDService,
	scenarioAssignmentEngine ScenarioAssignmentEngine,
	appRepo ApplicationRepository,
	protectedLabelPattern string) *service {
	return &service{
		repo:                     repo,
		labelRepo:                labelRepo,
		scenariosService:         scenariosService,
		labelUpsertService:       labelUpsertService,
		uidService:               uidService,
		scenarioAssignmentEngine: scenarioAssignmentEngine,
		protectedLabelPattern:    protectedLabelPattern,
	}
}

func (s *service) List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error) {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.repo.List(ctx, rtmTenant, filter, pageSize, cursor)
}

func (s *service) Get(ctx context.Context, id string) (*model.Runtime, error) {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	runtime, err := s.repo.GetByID(ctx, rtmTenant, id)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting Runtime with ID %s", id)
	}

	return runtime, nil
}

func (s *service) GetByTokenIssuer(ctx context.Context, issuer string) (*model.Runtime, error) {
	const (
		consoleURLLabelKey = "runtime_consoleUrl"
		dexSubdomain       = "dex"
		consoleSubdomain   = "console"
	)
	consoleURL := strings.Replace(issuer, dexSubdomain, consoleSubdomain, 1)

	filters := []*labelfilter.LabelFilter{
		labelfilter.NewForKeyWithQuery(consoleURLLabelKey, fmt.Sprintf(`"%s"`, consoleURL)),
	}

	runtime, err := s.repo.GetByFiltersGlobal(ctx, filters)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting the Runtime by the console URL label (%s)", consoleURL)
	}

	return runtime, nil
}

func (s *service) GetByFiltersGlobal(ctx context.Context, filters []*labelfilter.LabelFilter) (*model.Runtime, error) {
	runtimes, err := s.repo.GetByFiltersGlobal(ctx, filters)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting runtimes by filters from repo: ")
	}
	return runtimes, nil
}

func (s *service) Exist(ctx context.Context, id string) (bool, error) {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return false, errors.Wrapf(err, "while loading tenant from context")
	}

	exist, err := s.repo.Exists(ctx, rtmTenant, id)
	if err != nil {
		return false, errors.Wrapf(err, "while getting Runtime with ID %s", id)
	}

	return exist, nil
}

func (s *service) Create(ctx context.Context, in model.RuntimeInput) (string, error) {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "while loading tenant from context")
	}
	id := s.uidService.Generate()
	rtm := in.ToRuntime(id, rtmTenant, time.Now(), time.Now())

	err = s.repo.Create(ctx, rtm)
	if err != nil {
		return "", errors.Wrapf(err, "while creating Runtime")
	}

	err = s.scenariosService.EnsureScenariosLabelDefinitionExists(ctx, rtmTenant)
	if err != nil {
		return "", errors.Wrapf(err, "while ensuring Label Definition with key %s exists", model.ScenariosKey)
	}

	scenarios, err := s.scenarioAssignmentEngine.MergeScenariosFromInputLabelsAndAssignments(ctx, in.Labels)
	if err != nil {
		return "", errors.Wrap(err, "while merging scenarios from input and assignments")
	}

	if len(scenarios) > 0 {
		in.Labels[model.ScenariosKey] = scenarios
	} else {
		s.scenariosService.AddDefaultScenarioIfEnabled(ctx, &in.Labels)
	}

	if in.Labels == nil || in.Labels[IsNormalizedLabel] == nil {
		if in.Labels == nil {
			in.Labels = make(map[string]interface{}, 1)
		}
		in.Labels[IsNormalizedLabel] = "true"
	}

	log.C(ctx).Debugf("Removing protected labels. Labels before: %+v", in.Labels)
	in.Labels, err = unsafeExtractUnProtectedLabels(in.Labels, s.protectedLabelPattern)
	if err != nil {
		return "", err
	}
	log.C(ctx).Debugf("Successfully stripped protected labels. Resulting labels after operation are: %+v", in.Labels)

	err = s.labelUpsertService.UpsertMultipleLabels(ctx, rtmTenant, model.RuntimeLabelableObject, id, in.Labels)
	if err != nil {
		return id, errors.Wrapf(err, "while creating multiple labels for Runtime")
	}

	return id, nil
}

func (s *service) Update(ctx context.Context, id string, in model.RuntimeInput) error {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	rtm, err := s.repo.GetByID(ctx, rtmTenant, id)
	if err != nil {
		return errors.Wrapf(err, "while getting Runtime with id %s", id)
	}

	rtm = in.ToRuntime(id, rtm.Tenant, rtm.CreationTimestamp, time.Now())

	err = s.repo.Update(ctx, rtm)
	if err != nil {
		return errors.Wrap(err, "while updating Runtime")
	}

	if in.Labels == nil || in.Labels[IsNormalizedLabel] == nil {
		if in.Labels == nil {
			in.Labels = make(map[string]interface{}, 1)
		}
		in.Labels[IsNormalizedLabel] = "true"
	}

	log.C(ctx).Debugf("Removing protected labels. Labels before: %+v", in.Labels)
	in.Labels, err = unsafeExtractUnProtectedLabels(in.Labels, s.protectedLabelPattern)
	if err != nil {
		return err
	}
	log.C(ctx).Debugf("Successfully stripped protected labels. Resulting labels after operation are: %+v", in.Labels)

	// NOTE: The db layer does not support OR currently so multiple label patterns can't be implemented easily
	err = s.labelRepo.DeleteByKeyNegationPattern(ctx, rtmTenant, model.RuntimeLabelableObject, id, s.protectedLabelPattern)
	if err != nil {
		return errors.Wrapf(err, "while deleting all labels for Runtime")
	}

	if in.Labels == nil {
		return nil
	}

	scenarios, err := s.scenarioAssignmentEngine.MergeScenariosFromInputLabelsAndAssignments(ctx, in.Labels)
	if err != nil {
		return errors.Wrap(err, "while merging scenarios from input and assignments")
	}

	if len(scenarios) > 0 {
		in.Labels[model.ScenariosKey] = scenarios
	}

	err = s.labelUpsertService.UpsertMultipleLabels(ctx, rtmTenant, model.RuntimeLabelableObject, id, in.Labels)
	if err != nil {
		return errors.Wrapf(err, "while creating multiple labels for Runtime")
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	err = s.repo.Delete(ctx, rtmTenant, id)
	if err != nil {
		return errors.Wrapf(err, "while deleting Runtime")
	}

	// All labels are deleted (cascade delete)

	return nil
}

func (s *service) SetLabel(ctx context.Context, labelInput *model.LabelInput) error {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	err = s.ensureRuntimeExists(ctx, rtmTenant, labelInput.ObjectID)
	if err != nil {
		return err
	}
	currentRuntimeLabels, err := s.getCurrentLabelsForRuntime(ctx, rtmTenant, labelInput.ObjectID)
	if err != nil {
		return err
	}

	inputScenarios, exist := getScenarioLabelsFromInput(labelInput)
	if exist {
		existingScenarios, err := s.GetScenarioNamesForRuntime(ctx, labelInput.ObjectID)
		if err != nil {
			//TODO: handle properly
			return err
		}

		scToRemove := getScenariosToRemove(existingScenarios, inputScenarios)
		if s.isAnyBundleInstanceAuthForScenariosExist(ctx, scToRemove, labelInput.ObjectID) {
			return errors.New("Unable to delete label .....Bundle Instance Auths should be deleted first")
		}

		scToAdd := getScenariosToAdd(existingScenarios, inputScenarios)
		scenariosToKeep := GetScenariosToKeep(existingScenarios, inputScenarios)
		commonApplications := s.getCommonApplications(ctx, rtmTenant, scenariosToKeep, scToAdd)

		for appId, scenario := range commonApplications {
			bundleInstanceAuthsLabels := getBundleInstanceAuthsLabels(ctx, appId, labelInput.ObjectID)
			//TODO: Reuse common logic in engine.go -> upsertScenario(...)
			for _, currentLabel := range bundleInstanceAuthsLabels {
				if currentLabel.Key == model.ScenariosKey {
					var scenariosValue []string
					err := json.Unmarshal([]byte(currentLabel.Value), &scenariosValue)
					if err != nil {
						//todo handleErr
						continue
					}
					scenariosValue = append(scenariosValue, scenario)
					updatedValue, _ := json.Marshal(scenariosValue)
					currentLabel.Value = string(updatedValue)

					labelConverter := label.NewConverter()
					labelInput, err := labelConverter.FromEntity(currentLabel)
					if err != nil {
						// todo handle
					}
					labelInput.ObjectID = currentLabel.BundleInstanceAuthId.String // TODO find better way to set ObjectID
					if err := s.labelRepo.Upsert(ctx, &labelInput); err != nil {
						return err
					}
				}
			}
		}
	}

	newRuntimeLabels := make(map[string]interface{})
	for k, v := range currentRuntimeLabels {
		newRuntimeLabels[k] = v
	}

	newRuntimeLabels[labelInput.Key] = labelInput.Value

	err = s.upsertScenariosLabelIfShould(ctx, labelInput.ObjectID, labelInput.Key, currentRuntimeLabels, newRuntimeLabels)
	if err != nil {
		return err
	}

	protected, err := isProtected(labelInput.Key, s.protectedLabelPattern)
	if err != nil {
		return err
	}
	if protected {
		return apperrors.NewInvalidDataError("could not set protected label key %s", labelInput.Key)
	}
	if labelInput.Key != model.ScenariosKey {
		err = s.labelUpsertService.UpsertLabel(ctx, rtmTenant, labelInput)
		if err != nil {
			return errors.Wrapf(err, "while creating label for Runtime")
		}
	}

	return nil
}

func (s *service) GetLabel(ctx context.Context, runtimeID string, key string) (*model.Label, error) {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	rtmExists, err := s.repo.Exists(ctx, rtmTenant, runtimeID)
	if err != nil {
		return nil, errors.Wrap(err, "while checking Runtime existence")
	}
	if !rtmExists {
		return nil, fmt.Errorf("Runtime with ID %s doesn't exist", runtimeID)
	}

	label, err := s.labelRepo.GetByKey(ctx, rtmTenant, model.RuntimeLabelableObject, runtimeID, key)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Runtime")
	}

	return label, nil
}

func (s *service) ListLabels(ctx context.Context, runtimeID string) (map[string]*model.Label, error) {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	rtmExists, err := s.repo.Exists(ctx, rtmTenant, runtimeID)
	if err != nil {
		return nil, errors.Wrap(err, "while checking Runtime existence")
	}

	if !rtmExists {
		return nil, fmt.Errorf("Runtime with ID %s doesn't exist", runtimeID)
	}

	labels, err := s.labelRepo.ListForObject(ctx, rtmTenant, model.RuntimeLabelableObject, runtimeID)
	if err != nil {
		return nil, errors.Wrap(err, "while getting label for Runtime")
	}

	return extractUnProtectedLabels(labels, s.protectedLabelPattern)
}

func (s *service) UpdateTenantID(ctx context.Context, runtimeID, newTenantID string) error {
	if err := s.repo.UpdateTenantID(ctx, runtimeID, newTenantID); err != nil {
		return errors.Wrapf(err, "while updating tenant_id for runtime with ID %s", runtimeID)
	}
	return nil
}

func (s *service) DeleteLabel(ctx context.Context, runtimeID string, key string) error {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	err = s.ensureRuntimeExists(ctx, rtmTenant, runtimeID)
	if err != nil {
		return err
	}

	currentRuntimeLabels, err := s.getCurrentLabelsForRuntime(ctx, rtmTenant, runtimeID)
	if err != nil {
		return err
	}

	if key == model.ScenariosKey {
		scenarios, err := s.GetScenarioNamesForRuntime(ctx, runtimeID)
		if err != nil {
			// todo handle
		}
		if s.isAnyBundleInstanceAuthForScenariosExist(ctx, scenarios, runtimeID) {
			return errors.New("Unable to delete label .....Bundle Instance Auths should be deleted first")
		}
	}

	newRuntimeLabels := make(map[string]interface{})
	for k, v := range currentRuntimeLabels {
		newRuntimeLabels[k] = v
	}

	delete(newRuntimeLabels, key)

	err = s.upsertScenariosLabelIfShould(ctx, runtimeID, key, currentRuntimeLabels, newRuntimeLabels)
	if err != nil {
		return err
	}

	protected, err := isProtected(key, s.protectedLabelPattern)
	if err != nil {
		return err
	}
	if protected {
		return apperrors.NewInvalidDataError("could not delete protected label key %s", key)
	}
	if key != model.ScenariosKey {
		err = s.labelRepo.Delete(ctx, rtmTenant, model.RuntimeLabelableObject, runtimeID, key)
		if err != nil {
			return errors.Wrapf(err, "while deleting Runtime label")
		}
	}

	return nil
}

func (s *service) ensureRuntimeExists(ctx context.Context, tnt string, runtimeID string) error {
	rtmExists, err := s.repo.Exists(ctx, tnt, runtimeID)
	if err != nil {
		return errors.Wrap(err, "while checking Runtime existence")
	}
	if !rtmExists {
		return fmt.Errorf("Runtime with ID %s doesn't exist", runtimeID)
	}

	return nil
}

func (s *service) upsertScenariosLabelIfShould(ctx context.Context, runtimeID string, modifiedLabelKey string, currentRuntimeLabels, newRuntimeLabels map[string]interface{}) error {
	rtmTenant, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	finalScenarios := make([]interface{}, 0)

	if modifiedLabelKey == model.ScenariosKey {
		scenarios, err := s.scenarioAssignmentEngine.MergeScenariosFromInputLabelsAndAssignments(ctx, newRuntimeLabels)
		if err != nil {
			return errors.Wrap(err, "while merging scenarios from input and assignments")
		}

		for _, scenario := range scenarios {
			finalScenarios = append(finalScenarios, scenario)
		}
	} else {
		oldScenariosLabel, err := getScenariosLabel(currentRuntimeLabels)
		if err != nil {
			return err
		}

		previousScenariosFromAssignments, err := s.getScenariosFromAssignments(ctx, currentRuntimeLabels)
		if err != nil {
			return errors.Wrap(err, "while getting old scenarios label and scenarios from assignments")
		}

		newScenariosFromAssignments, err := s.getScenariosFromAssignments(ctx, newRuntimeLabels)
		if err != nil {
			return errors.Wrap(err, "while getting new scenarios from assignments")
		}

		finalScenarios = s.scenarioAssignmentEngine.MergeScenarios(oldScenariosLabel, previousScenariosFromAssignments, newScenariosFromAssignments)
	}

	//TODO compare finalScenarios and oldScenariosLabel to determine when to delete scenarios label
	if len(finalScenarios) == 0 {
		err := s.labelRepo.Delete(ctx, rtmTenant, model.RuntimeLabelableObject, runtimeID, model.ScenariosKey)
		if err != nil {
			return errors.Wrapf(err, "while deleting scenarios label from runtime with id [%s]", runtimeID)
		}
		return nil
	}

	scenariosLabelInput := &model.LabelInput{
		Key:        model.ScenariosKey,
		Value:      finalScenarios,
		ObjectID:   runtimeID,
		ObjectType: model.RuntimeLabelableObject,
	}

	err = s.labelUpsertService.UpsertLabel(ctx, rtmTenant, scenariosLabelInput)
	if err != nil {
		return errors.Wrapf(err, "while creating scenarios label for Runtime with id [%s]", runtimeID)
	}

	return nil
}

func (s *service) getCurrentLabelsForRuntime(ctx context.Context, tenantID, runtimeID string) (map[string]interface{}, error) {
	labels, err := s.labelRepo.ListForObject(ctx, tenantID, model.RuntimeLabelableObject, runtimeID)
	if err != nil {
		return nil, err
	}

	currentLabels := make(map[string]interface{})
	for _, v := range labels {
		currentLabels[v.Key] = v.Value
	}
	return currentLabels, nil
}

func extractUnProtectedLabels(labels map[string]*model.Label, protectedLabelsKeyPattern string) (map[string]*model.Label, error) {
	result := make(map[string]*model.Label)
	for labelKey, label := range labels {
		protected, err := isProtected(labelKey, protectedLabelsKeyPattern)
		if err != nil {
			return nil, err
		}
		if !protected {
			result[labelKey] = label
		}
	}
	return result, nil
}

func unsafeExtractUnProtectedLabels(labels map[string]interface{}, protectedLabelsKeyPattern string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for labelKey, label := range labels {
		protected, err := isProtected(labelKey, protectedLabelsKeyPattern)
		if err != nil {
			return nil, err
		}
		if !protected {
			result[labelKey] = label
		}
	}
	return result, nil
}

func isProtected(labelKey string, labelKeyPattern string) (bool, error) {
	matched, err := regexp.MatchString(labelKeyPattern, labelKey)
	if err != nil {
		return false, err
	}
	return matched, nil
}

func getScenariosLabel(currentRuntimeLabels map[string]interface{}) ([]interface{}, error) {
	oldScenariosLabel, ok := currentRuntimeLabels[model.ScenariosKey]

	var oldScenariosLabelInterfaceSlice []interface{}
	if ok {
		oldScenariosLabelInterfaceSlice, ok = oldScenariosLabel.([]interface{})
		if !ok {
			return nil, apperrors.NewInternalError("value for scenarios label must be []interface{}")
		}
	}
	return oldScenariosLabelInterfaceSlice, nil
}

func (s *service) getScenariosFromAssignments(ctx context.Context, currentRuntimeLabels map[string]interface{}) ([]interface{}, error) {
	selectors := s.convertMapStringInterfaceToMapStringString(currentRuntimeLabels)

	ScenariosFromAssignments, err := s.scenarioAssignmentEngine.GetScenariosForSelectorLabels(ctx, selectors)
	if err != nil {
		return nil, errors.Wrap(err, "while getting scenarios for selector labels")
	}

	newScenariosInterfaceSlice := s.convertStringSliceToInterfaceSlice(ScenariosFromAssignments)

	return newScenariosInterfaceSlice, nil
}

func (s *service) convertMapStringInterfaceToMapStringString(in map[string]interface{}) map[string]string {
	out := make(map[string]string)

	for k, v := range in {
		val, ok := v.(string)
		if ok {
			out[k] = val
		}
	}

	return out
}

func (s *service) convertStringSliceToInterfaceSlice(in []string) []interface{} {
	out := make([]interface{}, 0)
	for _, v := range in {
		out = append(out, v)
	}

	return out
}

func (s *service) GetScenarioNamesForRuntime(ctx context.Context, runtimeID string) ([]string, error) {
	log.C(ctx).Infof("Getting scenarios for runtime with id %s", runtimeID)

	runtimeLabels, err := s.GetLabel(ctx, runtimeID, model.ScenariosKey)
	if err != nil {
		if apperrors.ErrorCode(err) == apperrors.NotFound {
			log.C(ctx).Infof("No scenarios found for runtime")
			return nil, nil
		}
		return nil, err
	}

	scenarios, err := label.ValueToStringsSlice(runtimeLabels.Value)
	if err != nil {
		return nil, errors.Wrapf(err, "while parsing runtime label values")
	}

	return scenarios, nil
}

func (s *service) isAnyBundleInstanceAuthForScenariosExist(ctx context.Context, scenarios []string, runtimeId string) bool {
	for _, scenario := range scenarios {
		if s.isBundleInstanceAuthForScenarioExist(ctx, scenario, runtimeId) {
			return true
		}
	}
	return false
}

func (s *service) isBundleInstanceAuthForScenarioExist(ctx context.Context, scenario, runtimeId string) bool {
	persist, _ := persistence.FromCtx(ctx)

	var count int
	query := "SELECT 1 FROM labels INNER JOIN bundle_instance_auths ON labels.bundle_instance_auth_id = bundle_instance_auths.id WHERE json_build_array($1::text)::jsonb <@ labels.value AND bundle_instance_auths.runtime_id=$2 AND bundle_instance_auths.status_condition='SUCCEEDED'"
	err := persist.Get(&count, query, scenario, runtimeId)
	if err != nil {
		return false
	}

	return count != 0
}

func getScenarioLabelsFromInput(label *model.LabelInput) ([]string, bool) {
	if model.ScenariosKey != label.Key {
		return nil, false
	}

	var result []string
	switch val := label.Value.(type) {
	case string:
		return []string{val}, true
	case []interface{}:
		for _, elem := range val {
			valAsString, ok := elem.(string)
			if !ok {
				// todo think about handling non string types
				continue
			} else {
				result = append(result, valAsString)
			}
		}
		return result, true
	default:
		return nil, false
	}
}

func getScenariosToRemove(existing, new []string) []string {
	newScenariosMap := make(map[string]bool, 0)
	for _, scenario := range new {
		newScenariosMap[scenario] = true
	}

	result := make([]string, 0)
	for _, scenario := range existing {
		if _, ok := newScenariosMap[scenario]; !ok {
			result = append(result, scenario)
		}
	}
	return result
}

func GetScenariosToKeep(existing []string, input []string) []string {
	existingScenarioMap := make(map[string]bool, 0)
	for _, scenario := range existing {
		existingScenarioMap[scenario] = true
	}

	result := make([]string, 0)
	for _, scenario := range input {
		if _, ok := existingScenarioMap[scenario]; ok {
			result = append(result, scenario)
		}
	}
	return result
}

func (s *service) getCommonApplications(ctx context.Context, tenant string, scenariosToKeep []string, scenariosToAdd []string) map[string]string {
	appIdsForScenariosToKeep, err := s.getApplicationIdsForScenario(ctx, tenant, scenariosToKeep)
	if err != nil {
		// TODO HANDLE ERROR
	}

	commonApplicationScenarios := make(map[string]string)
	for _, scenario := range scenariosToAdd {
		appIds, err := s.getApplicationIdsForScenario(ctx, tenant, []string{scenario})
		if err != nil {
			// todo handle
		}
		for _, appId := range appIds {
			if contains(appIdsForScenariosToKeep, appId) {
				commonApplicationScenarios[appId] = scenario
			}
		}
	}
	return commonApplicationScenarios
}

func (s *service) getApplicationIdsForScenario(ctx context.Context, tenant string, scenarios []string) ([]string, error) {
	scenariosQuery := eventing.BuildQueryForScenarios(scenarios)
	appScenariosFilter := []*labelfilter.LabelFilter{labelfilter.NewForKeyWithQuery(model.ScenariosKey, scenariosQuery)}

	log.C(ctx).Debugf("Listing runtimes matching the query %s", scenariosQuery)
	applications, err := s.appRepo.ListAllByLabelFilter(ctx, tenant, appScenariosFilter)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting runtimes")
	}

	var appIds []string
	for _, r := range applications {
		appIds = append(appIds, r.ID)
	}

	return appIds, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getScenariosToAdd(existing, new []string) []string {
	existingScenarioMap := make(map[string]bool, 0)
	for _, scenario := range existing {
		existingScenarioMap[scenario] = true
	}

	result := make([]string, 0)
	for _, scenario := range new {
		if _, ok := existingScenarioMap[scenario]; !ok {
			result = append(result, scenario)
		}
	}
	return result
}

func getBundleInstanceAuthsLabels(ctx context.Context, appId, runtimeId string) []label.Entity {
	persist, _ := persistence.FromCtx(ctx)

	var entities []label.Entity
	query := "SELECT labels.* FROM bundle_instance_auths INNER JOIN bundles ON bundle_instance_auths.bundle_id=bundles.id INNER JOIN labels on bundle_instance_auths.id=labels.bundle_instance_auth_id WHERE bundles.app_id=$1 AND bundle_instance_auths.runtime_id=$2"
	err := persist.Select(&entities, query, appId, runtimeId)
	if err != nil {
		return []label.Entity{}
	}

	return entities
}
