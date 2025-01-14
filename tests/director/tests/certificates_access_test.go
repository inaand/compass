package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/certs"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/kyma-incubator/compass/tests/pkg/tenant"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSystemAccess(t *testing.T) {
	testCases := []struct {
		name           string
		tenant         string
		resourceSuffix string
		expectErr      bool
	}{
		{
			name:           "Integration System can manage account tenant entities",
			tenant:         tenant.TestTenants.GetIDByName(t, tenant.TestIntegrationSystemManagedAccount),
			resourceSuffix: "account-owned",
		},
		{
			name:           "Integration System can manage subaccount tenant entities",
			tenant:         tenant.TestTenants.GetIDByName(t, tenant.TestIntegrationSystemManagedSubaccount),
			resourceSuffix: "subaccount-owned",
		},
		{
			name:           "Integration System cannot manage customer tenant entities",
			tenant:         tenant.TestTenants.GetIDByName(t, tenant.TestDefaultCustomerTenant),
			resourceSuffix: "customer-owned",
			expectErr:      true,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			// Build graphql director client configured with certificate
			clientKey, rawCertChain := certs.ClientCertPair(t, conf.ExternalCA.Certificate, conf.ExternalCA.Key)
			directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, clientKey, rawCertChain, conf.SkipSSLValidation)

			t.Log(fmt.Sprintf("Trying to create application in account tenant %s", test.tenant))
			app, err := fixtures.RegisterApplication(t, ctx, directorCertSecuredClient, fmt.Sprintf("app-%s", test.resourceSuffix), test.tenant)
			defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, test.tenant, &app)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, app.ID)
			}

			t.Log(fmt.Sprintf("Trying to list applications in account tenant %s", test.tenant))
			getAppReq := fixtures.FixGetApplicationsRequestWithPagination()
			apps := graphql.ApplicationPage{}
			err = testctx.Tc.RunOperationWithCustomTenant(ctx, directorCertSecuredClient, test.tenant, getAppReq, &apps)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, apps.Data)
			}

			t.Log(fmt.Sprintf("Trying to register runtime in account tenant %s", test.tenant))
			rt, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, directorCertSecuredClient, test.tenant, &graphql.RuntimeInput{Name: fmt.Sprintf("runtime-%s", test.resourceSuffix)})
			defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, test.tenant, &rt)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, rt.ID)
			}

			t.Log(fmt.Sprintf("Trying to create application template in account tenant %s via client certificate", test.tenant))
			at, err := fixtures.CreateApplicationTemplate(t, ctx, directorCertSecuredClient, test.tenant, fmt.Sprintf("app-template-%s", test.resourceSuffix))
			defer fixtures.CleanupApplicationTemplate(t, ctx, dexGraphQLClient, test.tenant, &at)
			if test.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, at.ID)
			}
		})
	}
}
