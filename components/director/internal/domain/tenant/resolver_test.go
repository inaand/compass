package tenant_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"

	tnt "github.com/kyma-incubator/compass/components/director/pkg/tenant"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/stretchr/testify/mock"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence/txtest"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant/automock"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	persistenceautomock "github.com/kyma-incubator/compass/components/director/pkg/persistence/automock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolver_Tenants(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	txGen := txtest.NewTransactionContextGenerator(testError)

	first := 2
	gqlAfter := graphql.PageCursor("test")
	searchTerm := ""
	testFirstParameterMissingError := errors.New("Invalid data [reason=missing required parameter 'first']")

	modelTenants := []*model.BusinessTenantMapping{
		newModelBusinessTenantMapping(testID, testName),
		newModelBusinessTenantMapping("test1", "name1"),
	}

	modelTenantsPage := &model.BusinessTenantMappingPage{
		Data: modelTenants,
		PageInfo: &pagination.Page{
			StartCursor: "",
			EndCursor:   string(gqlAfter),
			HasNextPage: true,
		},
		TotalCount: 3,
	}

	gqlTenants := []*graphql.Tenant{
		newGraphQLTenant(testID, "", testName),
		newGraphQLTenant("test1", "", "name1"),
	}

	gqlTenantsPage := &graphql.TenantPage{
		Data:       gqlTenants,
		TotalCount: modelTenantsPage.TotalCount,
		PageInfo: &graphql.PageInfo{
			StartCursor: graphql.PageCursor(modelTenantsPage.PageInfo.StartCursor),
			EndCursor:   graphql.PageCursor(modelTenantsPage.PageInfo.EndCursor),
			HasNextPage: modelTenantsPage.PageInfo.HasNextPage,
		},
	}

	testCases := []struct {
		Name           string
		TxFn           func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		TenantSvcFn    func() *automock.BusinessTenantMappingService
		TenantConvFn   func() *automock.BusinessTenantMappingConverter
		first          *int
		ExpectedOutput *graphql.TenantPage
		ExpectedError  error
	}{
		{
			Name: "Success",
			TxFn: txGen.ThatSucceeds,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("ListPageBySearchTerm", txtest.CtxWithDBMatcher(), searchTerm, first, string(gqlAfter)).Return(modelTenantsPage, nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				TenantConv := &automock.BusinessTenantMappingConverter{}
				TenantConv.On("MultipleToGraphQL", modelTenants).Return(gqlTenants).Once()
				return TenantConv
			},
			first:          &first,
			ExpectedOutput: gqlTenantsPage,
		},
		{
			Name: "Returns error when getting tenants failed",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("ListPageBySearchTerm", txtest.CtxWithDBMatcher(), searchTerm, first, string(gqlAfter)).Return(nil, testError).Once()
				return TenantSvc
			},
			TenantConvFn:  unusedTenantConverter,
			first:         &first,
			ExpectedError: testError,
		},
		{
			Name:          "Returns error when failing on begin",
			TxFn:          txGen.ThatFailsOnBegin,
			TenantSvcFn:   unusedTenantService,
			TenantConvFn:  unusedTenantConverter,
			first:         &first,
			ExpectedError: testError,
		},
		{
			Name: "Returns error when failing on commit",
			TxFn: txGen.ThatFailsOnCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("ListPageBySearchTerm", txtest.CtxWithDBMatcher(), searchTerm, first, string(gqlAfter)).Return(modelTenantsPage, nil).Once()
				return TenantSvc
			},
			TenantConvFn:  unusedTenantConverter,
			first:         &first,
			ExpectedError: testError,
		},
		{
			Name: "Returns error when 'first' parameter is missing",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.AssertNotCalled(t, "ListPageBySearchTerm")
				return TenantSvc
			},
			TenantConvFn:  unusedTenantConverter,
			first:         nil,
			ExpectedError: testFirstParameterMissingError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tenantSvc := testCase.TenantSvcFn()
			tenantConv := testCase.TenantConvFn()
			persist, transact := testCase.TxFn()
			resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

			// WHEN
			result, err := resolver.Tenants(ctx, testCase.first, &gqlAfter, &searchTerm)

			// THEN
			if testCase.ExpectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.ExpectedOutput, result)

			mock.AssertExpectationsForObjects(t, persist, transact, tenantSvc, tenantConv)
		})
	}
}

func TestResolver_Tenant(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	txGen := txtest.NewTransactionContextGenerator(testError)

	tenantParent := ""
	tenantInternalID := "internal"

	expectedTenantsModel := []*model.BusinessTenantMapping{
		{
			ID:             testExternal,
			Name:           testName,
			ExternalTenant: testExternal,
			Parent:         tenantParent,
			Type:           tnt.Account,
			Provider:       testProvider,
			Status:         tnt.Active,
			Initialized:    nil,
		},
	}

	expectedTenantsGQL := []*graphql.Tenant{
		{
			ID:          testExternal,
			InternalID:  tenantInternalID,
			Name:        str.Ptr(testName),
			Type:        string(tnt.Account),
			ParentID:    tenantParent,
			Initialized: nil,
			Labels:      nil,
		},
	}

	testCases := []struct {
		Name           string
		TxFn           func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		TenantSvcFn    func() *automock.BusinessTenantMappingService
		TenantConvFn   func() *automock.BusinessTenantMappingConverter
		TenantInput    graphql.BusinessTenantMappingInput
		IDInput        string
		ExpectedError  error
		ExpectedResult *graphql.Tenant
	}{
		{
			Name: "Success",
			TxFn: txGen.ThatSucceeds,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), testExternal).Return(expectedTenantsModel[0], nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				conv := &automock.BusinessTenantMappingConverter{}
				conv.On("MultipleToGraphQL", expectedTenantsModel).Return(expectedTenantsGQL)
				return conv
			},
			IDInput:        testExternal,
			ExpectedError:  nil,
			ExpectedResult: expectedTenantsGQL[0],
		},
		{
			Name:           "That returns error when can not start transaction",
			TxFn:           txGen.ThatFailsOnBegin,
			TenantSvcFn:    unusedTenantService,
			TenantConvFn:   unusedTenantConverter,
			IDInput:        testExternal,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
		{
			Name: "That returns error when can not get tenant by external ID",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), testExternal).Return(nil, testError).Once()
				return TenantSvc
			},
			TenantConvFn:   unusedTenantConverter,
			IDInput:        testExternal,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
		{
			Name: "That returns error when cannot commit",
			TxFn: txGen.ThatFailsOnCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), testExternal).Return(expectedTenantsModel[0], nil).Once()
				return TenantSvc
			},
			TenantConvFn:   unusedTenantConverter,
			IDInput:        testExternal,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tenantSvc := testCase.TenantSvcFn()
			tenantConv := testCase.TenantConvFn()
			persist, transact := testCase.TxFn()
			resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

			// WHEN
			result, err := resolver.Tenant(ctx, testCase.IDInput)

			// THEN
			if testCase.ExpectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, result)
			}

			mock.AssertExpectationsForObjects(t, persist, transact, tenantSvc, tenantConv)
		})
	}
}

func TestResolver_Labels(t *testing.T) {
	// GIVEN
	ctx := context.TODO()

	txGen := txtest.NewTransactionContextGenerator(testError)

	tenantID := "2af44425-d02d-4aed-9086-b0fc3122b508"
	testTenant := &graphql.Tenant{ID: "externalID", InternalID: tenantID}

	testLabelKey := "my-key"
	testLabels := map[string]*model.Label{
		testLabelKey: {
			ID:         "5d0ec128-47da-418a-99f5-8409105ce82d",
			Tenant:     str.Ptr(tenantID),
			Key:        testLabelKey,
			Value:      "value",
			ObjectID:   tenantID,
			ObjectType: model.TenantLabelableObject,
		},
	}

	t.Run("Succeeds", func(t *testing.T) {
		tenantSvc := unusedTenantService()
		tenantSvc.On("ListLabels", txtest.CtxWithDBMatcher(), testTenant.InternalID).Return(testLabels, nil)
		tenantConv := unusedTenantConverter()
		persist, transact := txGen.ThatSucceeds()

		defer mock.AssertExpectationsForObjects(t, tenantSvc, tenantConv, persist, transact)

		resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

		result, err := resolver.Labels(ctx, testTenant, nil)
		assert.NoError(t, err)

		assert.NotNil(t, result)
		assert.Len(t, result, len(testLabels))
		assert.Equal(t, testLabels[testLabelKey].Value, result[testLabelKey])
	})
	t.Run("Succeeds when labels do not exist", func(t *testing.T) {
		tenantSvc := unusedTenantService()
		tenantSvc.On("ListLabels", txtest.CtxWithDBMatcher(), testTenant.InternalID).Return(nil, apperrors.NewNotFoundError(resource.Tenant, testTenant.InternalID))
		tenantConv := unusedTenantConverter()
		persist, transact := txGen.ThatSucceeds()

		defer mock.AssertExpectationsForObjects(t, tenantSvc, tenantConv, persist, transact)

		resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

		labels, err := resolver.Labels(ctx, testTenant, nil)
		assert.NoError(t, err)
		assert.Nil(t, labels)
	})
	t.Run("Returns error when the provided tenant is nil", func(t *testing.T) {
		tenantSvc := unusedTenantService()
		tenantConv := unusedTenantConverter()
		persist, transact := txGen.ThatDoesntStartTransaction()

		defer mock.AssertExpectationsForObjects(t, tenantSvc, tenantConv, persist, transact)

		resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

		_, err := resolver.Labels(ctx, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Tenant cannot be empty")
	})
	t.Run("Returns error when starting transaction fails", func(t *testing.T) {
		tenantSvc := unusedTenantService()
		tenantConv := unusedTenantConverter()
		persist, transact := txGen.ThatFailsOnBegin()

		defer mock.AssertExpectationsForObjects(t, tenantSvc, tenantConv, persist, transact)

		resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

		result, err := resolver.Labels(ctx, testTenant, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Returns error when it fails to list labels", func(t *testing.T) {
		tenantSvc := unusedTenantService()
		tenantSvc.On("ListLabels", txtest.CtxWithDBMatcher(), testTenant.InternalID).Return(nil, testError)
		tenantConv := unusedTenantConverter()
		persist, transact := txGen.ThatDoesntExpectCommit()

		defer mock.AssertExpectationsForObjects(t, tenantSvc, tenantConv, persist, transact)

		resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

		_, err := resolver.Labels(ctx, testTenant, nil)
		assert.Error(t, err)
		assert.Equal(t, testError, err)
	})
	t.Run("Returns error when commit fails", func(t *testing.T) {
		tenantSvc := unusedTenantService()
		tenantSvc.On("ListLabels", txtest.CtxWithDBMatcher(), testTenant.InternalID).Return(testLabels, nil)
		tenantConv := unusedTenantConverter()
		persist, transact := txGen.ThatFailsOnCommit()

		defer mock.AssertExpectationsForObjects(t, tenantSvc, tenantConv, persist, transact)

		resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

		_, err := resolver.Labels(ctx, testTenant, nil)
		assert.Error(t, err)
		assert.Equal(t, testError, err)
	})
}

func TestResolver_Write(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	txGen := txtest.NewTransactionContextGenerator(testError)

	tenantNames := []string{"name1", "name2"}
	tenantExternalTenants := []string{"external1", "external2"}
	tenantParent := ""
	tenantSubdomain := "subdomain"
	tenantRegion := "region"
	tenantProvider := "test"

	tenantsToUpsertGQL := []*graphql.BusinessTenantMappingInput{
		{
			Name:           tenantNames[0],
			ExternalTenant: tenantExternalTenants[0],
			Parent:         str.Ptr(tenantParent),
			Subdomain:      str.Ptr(tenantSubdomain),
			Region:         str.Ptr(tenantRegion),
			Type:           string(tnt.Account),
			Provider:       tenantProvider,
		},
		{
			Name:           tenantNames[1],
			ExternalTenant: tenantExternalTenants[1],
			Parent:         str.Ptr(tenantParent),
			Subdomain:      str.Ptr(tenantSubdomain),
			Region:         str.Ptr(tenantRegion),
			Type:           string(tnt.Account),
			Provider:       tenantProvider,
		},
	}
	tenantsToUpsertModel := []model.BusinessTenantMappingInput{
		{
			Name:           tenantNames[0],
			ExternalTenant: tenantExternalTenants[0],
			Parent:         tenantParent,
			Subdomain:      tenantSubdomain,
			Region:         tenantRegion,
			Type:           string(tnt.Account),
			Provider:       tenantProvider,
		},
		{
			Name:           tenantNames[1],
			ExternalTenant: tenantExternalTenants[1],
			Parent:         tenantParent,
			Subdomain:      tenantSubdomain,
			Region:         tenantRegion,
			Type:           string(tnt.Account),
			Provider:       tenantProvider,
		},
	}

	testCases := []struct {
		Name           string
		TxFn           func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		TenantSvcFn    func() *automock.BusinessTenantMappingService
		TenantConvFn   func() *automock.BusinessTenantMappingConverter
		TenantsInput   []*graphql.BusinessTenantMappingInput
		ExpectedError  error
		ExpectedResult int
	}{
		{
			Name: "Success",
			TxFn: txGen.ThatSucceeds,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := unusedTenantService()
				TenantSvc.On("UpsertMany", txtest.CtxWithDBMatcher(), tenantsToUpsertModel[0], tenantsToUpsertModel[1]).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				TenantConv := &automock.BusinessTenantMappingConverter{}
				TenantConv.On("MultipleInputFromGraphQL", tenantsToUpsertGQL).Return(tenantsToUpsertModel).Once()
				return TenantConv
			},
			TenantsInput:   tenantsToUpsertGQL,
			ExpectedError:  nil,
			ExpectedResult: 2,
		},
		{
			Name:           "Returns error when can not start transaction",
			TxFn:           txGen.ThatFailsOnBegin,
			TenantSvcFn:    unusedTenantService,
			TenantConvFn:   unusedTenantConverter,
			TenantsInput:   tenantsToUpsertGQL,
			ExpectedError:  testError,
			ExpectedResult: -1,
		},
		{
			Name: "Returns error when can not create the tenants",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("UpsertMany", txtest.CtxWithDBMatcher(), tenantsToUpsertModel[0], tenantsToUpsertModel[1]).Return(testError).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				TenantConv := &automock.BusinessTenantMappingConverter{}
				TenantConv.On("MultipleInputFromGraphQL", tenantsToUpsertGQL).Return(tenantsToUpsertModel).Once()
				return TenantConv
			},
			TenantsInput:   tenantsToUpsertGQL,
			ExpectedError:  testError,
			ExpectedResult: -1,
		},
		{
			Name: "Returns error when can not commit",
			TxFn: txGen.ThatFailsOnCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("UpsertMany", txtest.CtxWithDBMatcher(), tenantsToUpsertModel[0], tenantsToUpsertModel[1]).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				TenantConv := &automock.BusinessTenantMappingConverter{}
				TenantConv.On("MultipleInputFromGraphQL", tenantsToUpsertGQL).Return(tenantsToUpsertModel).Once()
				return TenantConv
			},
			TenantsInput:   tenantsToUpsertGQL,
			ExpectedError:  testError,
			ExpectedResult: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tenantSvc := testCase.TenantSvcFn()
			tenantConv := testCase.TenantConvFn()
			persist, transact := testCase.TxFn()
			resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

			// WHEN
			result, err := resolver.Write(ctx, testCase.TenantsInput)

			// THEN
			if testCase.ExpectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, result)
			}

			mock.AssertExpectationsForObjects(t, persist, transact, tenantSvc, tenantConv)
		})
	}
}

func TestResolver_Delete(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	txGen := txtest.NewTransactionContextGenerator(testError)

	tenantExternalTenants := []string{"external1", "external2"}

	testCases := []struct {
		Name           string
		TxFn           func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		TenantSvcFn    func() *automock.BusinessTenantMappingService
		TenantConvFn   func() *automock.BusinessTenantMappingConverter
		TenantsInput   []string
		ExpectedError  error
		ExpectedResult int
	}{
		{
			Name: "Success",
			TxFn: txGen.ThatSucceeds,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("DeleteMany", txtest.CtxWithDBMatcher(), tenantExternalTenants).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn:   unusedTenantConverter,
			TenantsInput:   tenantExternalTenants,
			ExpectedError:  nil,
			ExpectedResult: 2,
		},
		{
			Name:           "Returns error when can not start transaction",
			TxFn:           txGen.ThatFailsOnBegin,
			TenantSvcFn:    unusedTenantService,
			TenantConvFn:   unusedTenantConverter,
			TenantsInput:   tenantExternalTenants,
			ExpectedError:  testError,
			ExpectedResult: -1,
		},
		{
			Name: "Returns error when can not create the tenants",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("DeleteMany", txtest.CtxWithDBMatcher(), tenantExternalTenants).Return(testError).Once()
				return TenantSvc
			},
			TenantConvFn:   unusedTenantConverter,
			TenantsInput:   tenantExternalTenants,
			ExpectedError:  testError,
			ExpectedResult: -1,
		},
		{
			Name: "Returns error when can not commit",
			TxFn: txGen.ThatFailsOnCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("DeleteMany", txtest.CtxWithDBMatcher(), tenantExternalTenants).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn:   unusedTenantConverter,
			TenantsInput:   tenantExternalTenants,
			ExpectedError:  testError,
			ExpectedResult: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tenantSvc := testCase.TenantSvcFn()
			tenantConv := testCase.TenantConvFn()
			persist, transact := testCase.TxFn()
			resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

			// WHEN
			result, err := resolver.Delete(ctx, testCase.TenantsInput)

			// THEN
			if testCase.ExpectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, result)
			}

			mock.AssertExpectationsForObjects(t, persist, transact, tenantSvc, tenantConv)
		})
	}
}

func TestResolver_Update(t *testing.T) {
	// GIVEN
	ctx := context.TODO()
	txGen := txtest.NewTransactionContextGenerator(testError)

	tenantParent := ""
	tenantInternalID := "internal"

	tenantsToUpdateGQL := []*graphql.BusinessTenantMappingInput{
		{
			Name:           testName,
			ExternalTenant: testExternal,
			Parent:         str.Ptr(tenantParent),
			Subdomain:      str.Ptr(testSubdomain),
			Region:         str.Ptr(testRegion),
			Type:           string(tnt.Account),
			Provider:       testProvider,
		},
	}

	tenantsToUpdateModel := []model.BusinessTenantMappingInput{
		{
			Name:           testName,
			ExternalTenant: testExternal,
			Parent:         tenantParent,
			Subdomain:      testSubdomain,
			Region:         testRegion,
			Type:           string(tnt.Account),
			Provider:       testProvider,
		},
	}

	expectedTenantModel := &model.BusinessTenantMapping{
		ID:             testExternal,
		Name:           testName,
		ExternalTenant: testExternal,
		Parent:         tenantParent,
		Type:           tnt.Account,
		Provider:       testProvider,
		Status:         tnt.Active,
		Initialized:    nil,
	}

	expectedTenantGQL := &graphql.Tenant{
		ID:          testExternal,
		InternalID:  tenantInternalID,
		Name:        str.Ptr(testName),
		Type:        string(tnt.Account),
		ParentID:    tenantParent,
		Initialized: nil,
		Labels:      nil,
	}

	testCases := []struct {
		Name           string
		TxFn           func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		TenantSvcFn    func() *automock.BusinessTenantMappingService
		TenantConvFn   func() *automock.BusinessTenantMappingConverter
		TenantInput    graphql.BusinessTenantMappingInput
		IDInput        string
		ExpectedError  error
		ExpectedResult *graphql.Tenant
	}{
		{
			Name: "Success",
			TxFn: txGen.ThatSucceeds,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), tenantsToUpdateGQL[0].ExternalTenant).Return(expectedTenantModel, nil).Once()
				TenantSvc.On("Update", txtest.CtxWithDBMatcher(), tenantInternalID, tenantsToUpdateModel[0]).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				conv := &automock.BusinessTenantMappingConverter{}
				conv.On("MultipleInputFromGraphQL", tenantsToUpdateGQL).Return(tenantsToUpdateModel)
				conv.On("ToGraphQL", expectedTenantModel).Return(expectedTenantGQL)
				return conv
			},
			TenantInput:    *tenantsToUpdateGQL[0],
			IDInput:        tenantInternalID,
			ExpectedError:  nil,
			ExpectedResult: expectedTenantGQL,
		},
		{
			Name:           "Returns error when can not start transaction",
			TxFn:           txGen.ThatFailsOnBegin,
			TenantSvcFn:    unusedTenantService,
			TenantConvFn:   unusedTenantConverter,
			TenantInput:    *tenantsToUpdateGQL[0],
			IDInput:        tenantInternalID,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
		{
			Name: "Returns error when updating tenant fails",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("Update", txtest.CtxWithDBMatcher(), tenantInternalID, tenantsToUpdateModel[0]).Return(testError).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				conv := &automock.BusinessTenantMappingConverter{}
				conv.On("MultipleInputFromGraphQL", tenantsToUpdateGQL).Return(tenantsToUpdateModel)
				return conv
			},
			TenantInput:    *tenantsToUpdateGQL[0],
			IDInput:        tenantInternalID,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
		{
			Name: "Returns error when can not get tenant by external ID",
			TxFn: txGen.ThatDoesntExpectCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), tenantsToUpdateGQL[0].ExternalTenant).Return(nil, testError).Once()
				TenantSvc.On("Update", txtest.CtxWithDBMatcher(), tenantInternalID, tenantsToUpdateModel[0]).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				conv := &automock.BusinessTenantMappingConverter{}
				conv.On("MultipleInputFromGraphQL", tenantsToUpdateGQL).Return(tenantsToUpdateModel)
				return conv
			},
			TenantInput:    *tenantsToUpdateGQL[0],
			IDInput:        tenantInternalID,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
		{
			Name: "Returns error when can not commit",
			TxFn: txGen.ThatFailsOnCommit,
			TenantSvcFn: func() *automock.BusinessTenantMappingService {
				TenantSvc := &automock.BusinessTenantMappingService{}
				TenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), tenantsToUpdateGQL[0].ExternalTenant).Return(expectedTenantModel, nil).Once()
				TenantSvc.On("Update", txtest.CtxWithDBMatcher(), tenantInternalID, tenantsToUpdateModel[0]).Return(nil).Once()
				return TenantSvc
			},
			TenantConvFn: func() *automock.BusinessTenantMappingConverter {
				conv := &automock.BusinessTenantMappingConverter{}
				conv.On("MultipleInputFromGraphQL", tenantsToUpdateGQL).Return(tenantsToUpdateModel)
				return conv
			},
			TenantInput:    *tenantsToUpdateGQL[0],
			IDInput:        tenantInternalID,
			ExpectedError:  testError,
			ExpectedResult: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tenantSvc := testCase.TenantSvcFn()
			tenantConv := testCase.TenantConvFn()
			persist, transact := testCase.TxFn()
			resolver := tenant.NewResolver(transact, tenantSvc, tenantConv)

			// WHEN
			result, err := resolver.Update(ctx, testCase.IDInput, testCase.TenantInput)

			// THEN
			if testCase.ExpectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, result)
			}

			mock.AssertExpectationsForObjects(t, persist, transact, tenantSvc, tenantConv)
		})
	}
}

func unusedTenantConverter() *automock.BusinessTenantMappingConverter {
	return &automock.BusinessTenantMappingConverter{}
}

func unusedTenantService() *automock.BusinessTenantMappingService {
	return &automock.BusinessTenantMappingService{}
}
