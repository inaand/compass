package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/tests/pkg/tenantfetcher"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/assertions"
	"github.com/kyma-incubator/compass/tests/pkg/certs"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/kyma-incubator/compass/tests/pkg/ptr"
	"github.com/kyma-incubator/compass/tests/pkg/tenant"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ScenariosLabel          = "scenarios"
	IsNormalizedLabel       = "isNormalized"
	QueryRuntimesCategory   = "query runtimes"
	RegisterRuntimeCategory = "register runtime"
)

func TestRuntimeRegisterUpdateAndUnregister(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	givenInput := graphql.RuntimeInput{
		Name:        "runtime-create-update-delete",
		Description: ptr.String("runtime-1-description"),
		Labels:      graphql.Labels{"ggg": []interface{}{"hhh"}},
	}
	runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(givenInput)
	require.NoError(t, err)
	actualRuntime := graphql.RuntimeExt{}

	// WHEN
	registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
	saveExampleInCustomDir(t, registerReq.Query(), RegisterRuntimeCategory, "register runtime")
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, registerReq, &actualRuntime)
	defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &actualRuntime)

	//THEN
	require.NoError(t, err)
	require.NotEmpty(t, actualRuntime.ID)
	assertions.AssertRuntime(t, givenInput, actualRuntime, conf.DefaultScenarioEnabled, false)

	// add Label
	actualLabel := graphql.Label{}

	// WHEN
	addLabelReq := fixtures.FixSetRuntimeLabelRequest(actualRuntime.ID, "new_label", []string{"bbb"})
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, addLabelReq, &actualLabel)

	//THEN
	require.NoError(t, err)
	assert.Equal(t, "new_label", actualLabel.Key)
	assert.Len(t, actualLabel.Value, 1)
	assert.Contains(t, actualLabel.Value, "bbb")

	// get runtime and validate runtimes
	getRuntimeReq := fixtures.FixGetRuntimeRequest(actualRuntime.ID)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, getRuntimeReq, &actualRuntime)
	require.NoError(t, err)
	if conf.DefaultScenarioEnabled {
		assert.Len(t, actualRuntime.Labels, 4)
	} else {
		assert.Len(t, actualRuntime.Labels, 3)
	}

	// add agent auth
	// GIVEN
	in := fixtures.FixSampleApplicationRegisterInputWithWebhooks("app")

	appInputGQL, err := testctx.Tc.Graphqlizer.ApplicationRegisterInputToGQL(in)
	require.NoError(t, err)
	createAppReq := fixtures.FixRegisterApplicationRequest(appInputGQL)

	//WHEN
	actualApp := graphql.ApplicationExt{}
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, createAppReq, &actualApp)
	defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, tenantId, &actualApp)

	//THEN
	require.NoError(t, err)
	require.NotEmpty(t, actualApp.ID)

	// update runtime, check if only simple values are updated
	//GIVEN
	givenInput.Name = "updated-name"
	givenInput.Description = ptr.String("updated-description")
	givenInput.Labels = graphql.Labels{
		"key": []interface{}{"values", "aabbcc"},
	}
	runtimeStatusCond := graphql.RuntimeStatusConditionConnected
	givenInput.StatusCondition = &runtimeStatusCond

	runtimeInGQL, err = testctx.Tc.Graphqlizer.RuntimeInputToGQL(givenInput)
	require.NoError(t, err)
	updateRuntimeReq := fixtures.FixUpdateRuntimeRequest(actualRuntime.ID, runtimeInGQL)
	saveExample(t, updateRuntimeReq.Query(), "update runtime")
	//WHEN
	actualRuntime = graphql.RuntimeExt{}
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, updateRuntimeReq, &actualRuntime)

	//THEN
	require.NoError(t, err)
	assert.Equal(t, givenInput.Name, actualRuntime.Name)
	assert.Equal(t, *givenInput.Description, *actualRuntime.Description)
	assert.Equal(t, len(actualRuntime.Labels), 2)
	assert.Equal(t, runtimeStatusCond, actualRuntime.Status.Condition)

	// delete runtime

	// WHEN
	delReq := fixtures.FixUnregisterRuntimeRequest(actualRuntime.ID)
	saveExample(t, delReq.Query(), "unregister runtime")
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, delReq, nil)

	//THEN
	require.NoError(t, err)
}

func TestRuntimeUnregisterDeletesScenarioAssignments(t *testing.T) {
	const (
		testScenario = "test-scenario"
	)
	// GIVEN
	ctx := context.Background()
	subaccount := tenant.TestTenants.GetIDByName(t, tenant.TestProviderSubaccount)

	givenInput := graphql.RuntimeInput{
		Name:        "runtime-with-scenario-assignments",
		Description: ptr.String("runtime-1-description"),
		Labels:      graphql.Labels{"global_subaccount_id": []interface{}{subaccount}},
	}
	runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(givenInput)
	require.NoError(t, err)
	actualRuntime := graphql.RuntimeExt{}

	// WHEN
	registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, registerReq, &actualRuntime)
	defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenant.TestTenants.GetDefaultTenantID(), &actualRuntime)

	//THEN
	require.NoError(t, err)
	require.NotEmpty(t, actualRuntime.ID)
	assertions.AssertRuntime(t, givenInput, actualRuntime, conf.DefaultScenarioEnabled, true)

	// update label definition
	_ = fixtures.UpdateScenariosLabelDefinitionWithinTenant(t, ctx, dexGraphQLClient, tenant.TestTenants.GetDefaultTenantID(), []string{testScenario, "DEFAULT"})
	defer fixtures.UpdateScenariosLabelDefinitionWithinTenant(t, ctx, dexGraphQLClient, tenant.TestTenants.GetDefaultTenantID(), []string{"DEFAULT"})

	// register automatic sccenario assignment
	givenScenarioAssignment := graphql.AutomaticScenarioAssignmentSetInput{
		ScenarioName: testScenario,
		Selector: &graphql.LabelSelectorInput{
			Key:   "global_subaccount_id",
			Value: subaccount,
		},
	}

	scenarioAssignmentInGQL, err := testctx.Tc.Graphqlizer.AutomaticScenarioAssignmentSetInputToGQL(givenScenarioAssignment)
	require.NoError(t, err)
	actualScenarioAssignment := graphql.AutomaticScenarioAssignment{}

	// WHEN
	createAutomaticScenarioAssignmentReq := fixtures.FixCreateAutomaticScenarioAssignmentRequest(scenarioAssignmentInGQL)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, createAutomaticScenarioAssignmentReq, &actualScenarioAssignment)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, givenScenarioAssignment.ScenarioName, actualScenarioAssignment.ScenarioName)
	assert.Equal(t, givenScenarioAssignment.Selector.Key, actualScenarioAssignment.Selector.Key)
	assert.Equal(t, givenScenarioAssignment.Selector.Value, actualScenarioAssignment.Selector.Value)

	// get runtime - verify it is in scenario
	getRuntimeReq := fixtures.FixGetRuntimeRequest(actualRuntime.ID)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, getRuntimeReq, &actualRuntime)

	require.NoError(t, err)
	scenarios, hasScenarios := actualRuntime.Labels["scenarios"]
	assert.True(t, hasScenarios)
	assert.Len(t, scenarios, 1)
	assert.Contains(t, scenarios, testScenario)

	// delete runtime

	// WHEN
	delReq := fixtures.FixUnregisterRuntimeRequest(actualRuntime.ID)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, delReq, nil)

	//THEN
	require.NoError(t, err)

	// get automatic scenario assignment - see that it's deleted
	actualScenarioAssignments := graphql.AutomaticScenarioAssignmentPage{}
	getScenarioAssignmentsReq := fixtures.FixAutomaticScenarioAssignmentsRequest()
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, getScenarioAssignmentsReq, &actualScenarioAssignments)
	require.NoError(t, err)
	assert.Equal(t, actualScenarioAssignments.TotalCount, 0)
}

func TestQueryRuntimes(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	idsToRemove := make([]string, 0)
	defer func() {
		for _, id := range idsToRemove {
			if id != "" {
				fixtures.UnregisterRuntime(t, ctx, dexGraphQLClient, tenantId, id)
			}
		}
	}()

	inputRuntimes := []*graphql.Runtime{
		{Name: "runtime-query-1", Description: ptr.String("test description")},
		{Name: "runtime-query-2", Description: ptr.String("another description")},
		{Name: "runtime-query-3"},
	}

	for _, rtm := range inputRuntimes {
		givenInput := graphql.RuntimeInput{
			Name:        rtm.Name,
			Description: rtm.Description,
		}
		runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(givenInput)
		require.NoError(t, err)
		createReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
		actualRuntime := graphql.Runtime{}
		err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, createReq, &actualRuntime)
		require.NoError(t, err)
		require.NotEmpty(t, actualRuntime.ID)
		rtm.ID = actualRuntime.ID
		idsToRemove = append(idsToRemove, actualRuntime.ID)
	}
	actualPage := graphql.RuntimePage{}

	// WHEN
	queryReq := fixtures.FixGetRuntimesRequestWithPagination()
	err := testctx.Tc.RunOperation(ctx, dexGraphQLClient, queryReq, &actualPage)
	saveExampleInCustomDir(t, queryReq.Query(), QueryRuntimesCategory, "query runtimes")

	//THEN
	require.NoError(t, err)
	assert.Len(t, actualPage.Data, len(inputRuntimes))
	assert.Equal(t, len(inputRuntimes), actualPage.TotalCount)

	for _, inputRtm := range inputRuntimes {
		found := false
		for _, actualRtm := range actualPage.Data {
			if inputRtm.ID == actualRtm.ID {
				found = true
				assert.Equal(t, inputRtm.Name, actualRtm.Name)
				assert.Equal(t, inputRtm.Description, actualRtm.Description)
				break
			}
		}
		assert.True(t, found)
	}
}

func TestQuerySpecificRuntime(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	givenInput := graphql.RuntimeInput{
		Name: "runtime-specific-runtime",
	}
	runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(givenInput)
	require.NoError(t, err)
	registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
	createdRuntime := graphql.RuntimeExt{}
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, registerReq, &createdRuntime)
	defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &createdRuntime)

	require.NoError(t, err)
	require.NotEmpty(t, createdRuntime.ID)

	// WHEN
	queriedRuntime := graphql.Runtime{}
	queryReq := fixtures.FixGetRuntimeRequest(createdRuntime.ID)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, queryReq, &queriedRuntime)
	saveExample(t, queryReq.Query(), "query runtime")

	//THEN
	require.NoError(t, err)
	assert.Equal(t, createdRuntime.ID, queriedRuntime.ID)
	assert.Equal(t, createdRuntime.Name, queriedRuntime.Name)
	assert.Equal(t, createdRuntime.Description, queriedRuntime.Description)
}

func TestQueryRuntimesWithPagination(t *testing.T) {
	//GIVEN
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	runtimes := make(map[string]*graphql.Runtime)
	runtimesAmount := 10
	for i := 0; i < runtimesAmount; i++ {
		runtimeInput := graphql.RuntimeInput{
			Name: fmt.Sprintf("runtime-%d", i),
		}
		runtimeInputGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(runtimeInput)
		require.NoError(t, err)

		registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInputGQL)

		runtime := graphql.Runtime{}
		err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, registerReq, &runtime)
		defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &graphql.RuntimeExt{Runtime: runtime})

		require.NoError(t, err)
		require.NotEmpty(t, runtime.ID)
		runtimes[runtime.ID] = &runtime
	}

	after := 3
	cursor := ""
	queriesForFullPage := int(runtimesAmount / after)

	for i := 0; i < queriesForFullPage; i++ {
		runtimesRequest := fixtures.FixRuntimeRequestWithPaginationRequest(after, cursor)

		//WHEN
		runtimePage := graphql.RuntimePage{}
		err := testctx.Tc.RunOperation(ctx, dexGraphQLClient, runtimesRequest, &runtimePage)
		require.NoError(t, err)

		//THEN
		assert.Equal(t, cursor, string(runtimePage.PageInfo.StartCursor))
		assert.True(t, runtimePage.PageInfo.HasNextPage)
		assert.Len(t, runtimePage.Data, after)
		assert.Equal(t, runtimesAmount, runtimePage.TotalCount)
		for _, runtime := range runtimePage.Data {
			assert.Equal(t, runtime, runtimes[runtime.ID])
			delete(runtimes, runtime.ID)
		}
		cursor = string(runtimePage.PageInfo.EndCursor)
	}

	//WHEN get last page with last runtime
	runtimesRequest := fixtures.FixRuntimeRequestWithPaginationRequest(after, cursor)
	lastRuntimePage := graphql.RuntimePage{}
	err := testctx.Tc.RunOperation(ctx, dexGraphQLClient, runtimesRequest, &lastRuntimePage)
	require.NoError(t, err)
	saveExampleInCustomDir(t, runtimesRequest.Query(), QueryRuntimesCategory, "query runtimes with pagination")

	//THEN
	assert.False(t, lastRuntimePage.PageInfo.HasNextPage)
	assert.Empty(t, lastRuntimePage.PageInfo.EndCursor)
	require.Len(t, lastRuntimePage.Data, 1)
	assert.Equal(t, lastRuntimePage.Data[0], runtimes[lastRuntimePage.Data[0].ID])
	delete(runtimes, lastRuntimePage.Data[0].ID)
	assert.Len(t, runtimes, 0)
}

func TestRegisterUpdateRuntimeWithoutLabels(t *testing.T) {
	//GIVEN
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	name := "test-create-runtime-without-labels"
	runtimeInput := graphql.RuntimeInput{Name: name}

	runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, dexGraphQLClient, tenantId, &runtimeInput)
	defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)

	//WHEN
	fetchedRuntime := fixtures.GetRuntime(t, ctx, dexGraphQLClient, tenantId, runtime.ID)

	//THEN
	require.Equal(t, runtime.ID, fetchedRuntime.ID)
	assertions.AssertRuntime(t, runtimeInput, fetchedRuntime, conf.DefaultScenarioEnabled, false)

	//GIVEN
	secondRuntime := graphql.RuntimeExt{}
	secondInput := graphql.RuntimeInput{
		Name:        name,
		Description: ptr.String("runtime-1-description"),
		Labels:      graphql.Labels{ScenariosLabel: []interface{}{"DEFAULT"}},
	}
	runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(secondInput)
	require.NoError(t, err)
	updateReq := fixtures.FixUpdateRuntimeRequest(fetchedRuntime.ID, runtimeInGQL)

	// WHEN
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, updateReq, &secondRuntime)

	//THEN
	require.NoError(t, err)
	assertions.AssertRuntime(t, secondInput, secondRuntime, conf.DefaultScenarioEnabled, false)
}

func TestRegisterUpdateRuntimeWithIsNormalizedLabel(t *testing.T) {
	//GIVEN
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	name := "test-create-runtime-without-labels"
	runtimeInput := graphql.RuntimeInput{
		Name:   name,
		Labels: graphql.Labels{IsNormalizedLabel: "false"},
	}

	runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, dexGraphQLClient, tenantId, &runtimeInput)
	defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)

	//WHEN
	fetchedRuntime := fixtures.GetRuntime(t, ctx, dexGraphQLClient, tenantId, runtime.ID)

	//THEN
	require.Equal(t, runtime.ID, fetchedRuntime.ID)
	assertions.AssertRuntime(t, runtimeInput, fetchedRuntime, conf.DefaultScenarioEnabled, false)

	//GIVEN
	secondRuntime := graphql.RuntimeExt{}
	secondInput := graphql.RuntimeInput{
		Name:        name,
		Description: ptr.String("runtime-1-description"),
		Labels:      graphql.Labels{IsNormalizedLabel: "true", ScenariosLabel: []interface{}{"DEFAULT"}},
	}
	runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(secondInput)
	require.NoError(t, err)
	updateReq := fixtures.FixUpdateRuntimeRequest(fetchedRuntime.ID, runtimeInGQL)

	// WHEN
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, updateReq, &secondRuntime)

	//THEN
	require.NoError(t, err)
	assertions.AssertRuntime(t, secondInput, secondRuntime, conf.DefaultScenarioEnabled, false)
}

func TestRuntimeRegisterUpdateAndUnregisterWithCertificate(t *testing.T) {
	t.Run("Test runtime operations(CUD) with externally issued certificate", func(t *testing.T) {
		// GIVEN
		ctx := context.Background()
		subscriptionProviderSubaccountID := tenant.TestTenants.GetIDByName(t, tenant.TestProviderSubaccount)

		// Build graphql director client configured with certificate
		clientKey, rawCertChain := certs.ClientCertPair(t, conf.ExternalCA.Certificate, conf.ExternalCA.Key)
		directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, clientKey, rawCertChain, conf.SkipSSLValidation)

		protectedConsumerSubaccountIdsLabel := "consumer_subaccount_ids"

		runtimeInput := graphql.RuntimeInput{
			Name:        "register-runtime-with-protected-labels",
			Description: ptr.String("register-runtime-with-protected-labels-description"),
			Labels:      graphql.Labels{protectedConsumerSubaccountIdsLabel: []string{"subaccountID-1", "subaccountID-2"}},
		}
		actualRtm := graphql.RuntimeExt{}

		t.Log("Successfully register runtime using certificate with protected labels and validate that they are excluded")
		// GIVEN
		runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(runtimeInput)
		require.NoError(t, err)

		// WHEN
		registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, registerReq, &actualRtm)
		defer fixtures.CleanupRuntime(t, ctx, directorCertSecuredClient, subscriptionProviderSubaccountID, &actualRtm)

		//THEN
		require.NoError(t, err)
		require.NotEmpty(t, actualRtm.ID)
		require.Equal(t, runtimeInput.Name, actualRtm.Name)
		require.Equal(t, runtimeInput.Description, actualRtm.Description)
		require.Empty(t, actualRtm.Labels[protectedConsumerSubaccountIdsLabel])

		t.Log("Successfully register runtime with certificate")
		// GIVEN
		runtimeInput = graphql.RuntimeInput{
			Name:        "runtime-create-update-delete",
			Description: ptr.String("runtime-create-update-delete-description"),
			Labels:      graphql.Labels{conf.SelfRegisterDistinguishLabelKey: []interface{}{"distinguish-label-value"}, tenantfetcher.RegionKey: tenantfetcher.RegionPathParamValue},
		}
		actualRuntime := graphql.RuntimeExt{}

		runtimeInGQL, err = testctx.Tc.Graphqlizer.RuntimeInputToGQL(runtimeInput)
		require.NoError(t, err)

		// WHEN
		registerReq = fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, registerReq, &actualRuntime)
		defer fixtures.CleanupRuntime(t, ctx, directorCertSecuredClient, subscriptionProviderSubaccountID, &actualRuntime)

		//THEN
		require.NoError(t, err)
		require.NotEmpty(t, actualRuntime.ID)
		assertions.AssertRuntime(t, runtimeInput, actualRuntime, conf.DefaultScenarioEnabled, true)

		t.Log("Successfully set regular runtime label using certificate")
		// GIVEN
		actualLabel := graphql.Label{}

		// WHEN
		addLabelReq := fixtures.FixSetRuntimeLabelRequest(actualRuntime.ID, "regular_label", []string{"labelValue"})
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, addLabelReq, &actualLabel)

		//THEN
		require.NoError(t, err)
		require.Equal(t, "regular_label", actualLabel.Key)
		require.Len(t, actualLabel.Value, 1)
		require.Contains(t, actualLabel.Value, "labelValue")

		t.Log("Fail setting protected label on runtime")
		// GIVEN
		protectedLabel := graphql.Label{}

		// WHEN
		pLabelReq := fixtures.FixSetRuntimeLabelRequest(actualRuntime.ID, protectedConsumerSubaccountIdsLabel, []string{"subaccountID-1", "subaccountID-2"})
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, pLabelReq, &protectedLabel)

		//THEN
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not set unmodifiable label with key consumer_subaccount_ids")
		require.Empty(t, protectedLabel)

		t.Log("Successfully get runtime")
		getRuntimeReq := fixtures.FixGetRuntimeRequest(actualRuntime.ID)
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, getRuntimeReq, &actualRuntime)
		require.NoError(t, err)
		require.NotEmpty(t, actualRuntime.ID)
		assert.Len(t, actualRuntime.Labels, 5) // three labels from the different runtime inputs plus two additional during runtime registration - isNormalized and "self register" label

		t.Log("Successfully update runtime and validate the protected labels are excluded")
		//GIVEN
		runtimeInput.Name = "updated-runtime"
		runtimeInput.Description = ptr.String("updated-runtime-description")
		runtimeInput.Labels = graphql.Labels{
			conf.SelfRegisterDistinguishLabelKey: []interface{}{"distinguish-label-value"}, tenantfetcher.RegionKey: tenantfetcher.RegionPathParamValue, protectedConsumerSubaccountIdsLabel: []interface{}{"subaccountID-1", "subaccountID-2"},
		}
		runtimeStatusCond := graphql.RuntimeStatusConditionConnected
		runtimeInput.StatusCondition = &runtimeStatusCond

		runtimeInGQL, err = testctx.Tc.Graphqlizer.RuntimeInputToGQL(runtimeInput)
		require.NoError(t, err)
		updateRuntimeReq := fixtures.FixUpdateRuntimeRequest(actualRuntime.ID, runtimeInGQL)

		//WHEN
		actualRuntime = graphql.RuntimeExt{}
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, updateRuntimeReq, &actualRuntime)

		//THEN
		require.NoError(t, err)
		require.Equal(t, runtimeInput.Name, actualRuntime.Name)
		require.Equal(t, *runtimeInput.Description, *actualRuntime.Description)
		require.Equal(t, runtimeStatusCond, actualRuntime.Status.Condition)
		require.Equal(t, len(actualRuntime.Labels), 3) // two labels from the runtime input plus one additional label, added during runtime update(isNormalized)
		labelValues, ok := actualRuntime.Labels[protectedConsumerSubaccountIdsLabel]
		require.False(t, ok)
		require.Empty(t, labelValues)

		t.Log("Successfully delete runtime using certificate")
		// WHEN
		delReq := fixtures.FixUnregisterRuntimeRequest(actualRuntime.ID)
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, delReq, nil)

		//THEN
		require.NoError(t, err)
	})
}

func TestQueryRuntimesWithCertificate(t *testing.T) {
	t.Run("Query runtime with externally issued certificate", func(t *testing.T) {
		// GIVEN
		ctx := context.Background()
		subscriptionProviderSubaccountID := tenant.TestTenants.GetIDByName(t, tenant.TestProviderSubaccount)

		// Build graphql director client configured with certificate
		clientKey, rawCertChain := certs.ClientCertPair(t, conf.ExternalCA.Certificate, conf.ExternalCA.Key)
		directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, clientKey, rawCertChain, conf.SkipSSLValidation)

		idsToRemove := make([]string, 0)
		defer func() {
			for _, id := range idsToRemove {
				if id != "" {
					fixtures.UnregisterRuntime(t, ctx, directorCertSecuredClient, subscriptionProviderSubaccountID, id)
				}
			}
		}()

		inputRuntimes := []*graphql.Runtime{
			{Name: "runtime-query-1", Description: ptr.String("test description")},
			{Name: "runtime-query-2", Description: ptr.String("another description")},
			{Name: "runtime-query-3"},
		}

		for _, rtm := range inputRuntimes {
			givenInput := graphql.RuntimeInput{
				Name:        rtm.Name,
				Description: rtm.Description,
			}
			runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(givenInput)
			require.NoError(t, err)
			createReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
			actualRuntime := graphql.Runtime{}
			err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, createReq, &actualRuntime)
			require.NoError(t, err)
			require.NotEmpty(t, actualRuntime.ID)
			rtm.ID = actualRuntime.ID
			idsToRemove = append(idsToRemove, actualRuntime.ID)
		}
		actualPage := graphql.RuntimePage{}

		// WHEN
		queryReq := fixtures.FixGetRuntimesRequestWithPagination()
		err := testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, queryReq, &actualPage)

		//THEN
		require.NoError(t, err)
		assert.Len(t, actualPage.Data, len(inputRuntimes))
		assert.Equal(t, len(inputRuntimes), actualPage.TotalCount)

		for _, inputRtm := range inputRuntimes {
			found := false
			for _, actualRtm := range actualPage.Data {
				if inputRtm.ID == actualRtm.ID {
					found = true
					assert.Equal(t, inputRtm.Name, actualRtm.Name)
					assert.Equal(t, inputRtm.Description, actualRtm.Description)
					break
				}
			}
			assert.True(t, found)
		}
	})
}

func TestQuerySpecificRuntimeWithCertificate(t *testing.T) {
	t.Run("Query specific runtime with externally issued certificate", func(t *testing.T) {
		// GIVEN
		ctx := context.Background()
		subscriptionProviderSubaccountID := tenant.TestTenants.GetIDByName(t, tenant.TestProviderSubaccount)

		// Build graphql director client configured with certificate
		clientKey, rawCertChain := certs.ClientCertPair(t, conf.ExternalCA.Certificate, conf.ExternalCA.Key)
		directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, clientKey, rawCertChain, conf.SkipSSLValidation)

		runtimeInput := graphql.RuntimeInput{
			Name: "runtime-specific-runtime",
		}
		runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(runtimeInput)
		require.NoError(t, err)
		registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
		createdRuntime := graphql.RuntimeExt{}
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, registerReq, &createdRuntime)
		defer fixtures.CleanupRuntime(t, ctx, directorCertSecuredClient, subscriptionProviderSubaccountID, &createdRuntime)

		require.NoError(t, err)
		require.NotEmpty(t, createdRuntime.ID)

		// WHEN
		queriedRuntime := graphql.Runtime{}
		queryReq := fixtures.FixGetRuntimeRequest(createdRuntime.ID)
		err = testctx.Tc.RunOperationWithoutTenant(ctx, directorCertSecuredClient, queryReq, &queriedRuntime)

		//THEN
		require.NoError(t, err)
		assert.Equal(t, createdRuntime.ID, queriedRuntime.ID)
		assert.Equal(t, createdRuntime.Name, queriedRuntime.Name)
		assert.Equal(t, createdRuntime.Description, queriedRuntime.Description)
	})
}
