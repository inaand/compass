package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	"github.com/kyma-incubator/compass/tests/pkg/token"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewerQuery(t *testing.T) {
	ctx := context.Background()

	t.Run("Test viewer as Integration System", func(t *testing.T) {
		t.Log("Register Integration System with Dex id token")
		intSys, err := fixtures.RegisterIntegrationSystem(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, "integration-system")
		defer fixtures.CleanupIntegrationSystem(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, intSys)

		require.NoError(t, err)
		require.NotEmpty(t, intSys.ID)
		t.Logf("Registered Integration System with [id=%s]", intSys.ID)

		t.Log("Request Client Credentials for Integration System")
		intSystemAuth := fixtures.RequestClientCredentialsForIntegrationSystem(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, intSys.ID)

		intSysOauthCredentialData, ok := intSystemAuth.Auth.Credential.(*graphql.OAuthCredentialData)
		require.True(t, ok)

		t.Log("Issue a Hydra token with Client Credentials")
		accessToken := token.GetAccessToken(t, intSysOauthCredentialData, "")
		oauthGraphQLClient := gql.NewAuthorizedGraphQLClientWithCustomURL(accessToken, testConfig.DirectorURL)

		t.Log("Requesting Viewer as Integration System")
		viewer := graphql.Viewer{}
		req := fixtures.FixGetViewerRequest()

		err = testctx.Tc.RunOperationWithCustomTenant(ctx, oauthGraphQLClient, testConfig.DefaultTestTenant, req, &viewer)
		require.NoError(t, err)
		assert.Equal(t, intSys.ID, viewer.ID)
		assert.Equal(t, graphql.ViewerTypeIntegrationSystem, viewer.Type)
	})

	t.Run("Test viewer as Application", func(t *testing.T) {
		appInput := graphql.ApplicationRegisterInput{
			Name: "test-app",
			Labels: graphql.Labels{
				"scenarios": []interface{}{"DEFAULT"},
			},
		}

		t.Log("Register Application with Dex id token")
		app, err := fixtures.RegisterApplicationFromInput(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, appInput)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &app)
		require.NoError(t, err)
		t.Logf("Registered Application with [id=%s]", app.ID)

		t.Log("Request Client Credentials for Application")
		appAuth := fixtures.RequestClientCredentialsForApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, app.ID)
		appOauthCredentialData, ok := appAuth.Auth.Credential.(*graphql.OAuthCredentialData)
		require.True(t, ok)
		require.NotEmpty(t, appOauthCredentialData.ClientSecret)
		require.NotEmpty(t, appOauthCredentialData.ClientID)

		t.Log("Issue a Hydra token with Client Credentials")
		accessToken := token.GetAccessToken(t, appOauthCredentialData, "")
		oauthGraphQLClient := gql.NewAuthorizedGraphQLClientWithCustomURL(accessToken, fmt.Sprintf("https://compass-gateway-auth-oauth.%s/director/graphql", testConfig.Domain))

		t.Log("Requesting Viewer as Application")
		viewer := graphql.Viewer{}
		req := fixtures.FixGetViewerRequest()

		err = testctx.Tc.RunOperationWithCustomTenant(ctx, oauthGraphQLClient, testConfig.DefaultTestTenant, req, &viewer)
		require.NoError(t, err)
		assert.Equal(t, app.ID, viewer.ID)
		assert.Equal(t, graphql.ViewerTypeApplication, viewer.Type)
	})

	t.Run("Test viewer as Runtime", func(t *testing.T) {
		runtimeInput := graphql.RuntimeInput{
			Name: "test-runtime",
			Labels: graphql.Labels{
				"scenarios": []interface{}{"DEFAULT"},
			},
		}

		t.Log("Register Runtime with Dex id token")
		runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &runtimeInput)
		defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &runtime)
		require.NoError(t, err)
		require.NotEmpty(t, runtime.ID)

		t.Logf("Registered Runtime with [id=%s]", runtime.ID)

		t.Log("Request Client Credentials for Runtime")
		rtmAuth := fixtures.RequestClientCredentialsForRuntime(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, runtime.ID)
		rtmOauthCredentialData, ok := rtmAuth.Auth.Credential.(*graphql.OAuthCredentialData)
		require.True(t, ok)
		require.NotEmpty(t, rtmOauthCredentialData.ClientSecret)
		require.NotEmpty(t, rtmOauthCredentialData.ClientID)

		t.Log("Issue a Hydra token with Client Credentials")
		accessToken := token.GetAccessToken(t, rtmOauthCredentialData, "")
		oauthGraphQLClient := gql.NewAuthorizedGraphQLClientWithCustomURL(accessToken, fmt.Sprintf("https://compass-gateway-auth-oauth.%s/director/graphql", testConfig.Domain))

		t.Log("Requesting Viewer as Runtime")
		viewer := graphql.Viewer{}
		req := fixtures.FixGetViewerRequest()

		err = testctx.Tc.RunOperationWithCustomTenant(ctx, oauthGraphQLClient, testConfig.DefaultTestTenant, req, &viewer)
		require.NoError(t, err)
		assert.Equal(t, runtime.ID, viewer.ID)
		assert.Equal(t, graphql.ViewerTypeRuntime, viewer.Type)
	})

}
