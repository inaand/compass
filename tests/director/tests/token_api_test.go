//go:build !ignore_external_dependencies
// +build !ignore_external_dependencies

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/tenant"
	"github.com/stretchr/testify/assert"
)

//This test also test runtime/application auths custom resolver
//TODO: Currently we don't save OneTimeToken mutations in examples, because those tests are turn off in gen_examples.sh,
// because we need connector up and running, which requires k8s cluster running.
func TestTokenGeneration(t *testing.T) {
	t.Run("Generate one time token for Runtime", func(t *testing.T) {
		//GIVEN
		ctx := context.Background()

		tenantId := tenant.TestTenants.GetDefaultTenantID()

		input := fixtures.FixRuntimeInput("test")
		runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, dexGraphQLClient, tenantId, &input)
		defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &runtime)
		require.NoError(t, err)
		require.NotEmpty(t, runtime.ID)
		tokenRequestNumber := 3

		//WHEN
		for i := 0; i < tokenRequestNumber; i++ {
			token := fixtures.RequestOneTimeTokenForRuntime(t, ctx, dexGraphQLClient, tenantId, runtime.ID)
			assert.NotEmpty(t, token.Token)
			assert.NotEmpty(t, token.ConnectorURL)
		}
		//THEN
		runtimeExt := fixtures.GetRuntime(t, ctx, dexGraphQLClient, tenantId, runtime.ID)
		assert.Len(t, runtimeExt.Auths, tokenRequestNumber)
	})

	t.Run("Generate one time token for Application", func(t *testing.T) {
		//GIVEN
		ctx := context.TODO()

		tenantId := tenant.TestTenants.GetDefaultTenantID()

		app, err := fixtures.RegisterApplication(t, ctx, dexGraphQLClient, "test", tenantId)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, tenantId, &app)
		require.NoError(t, err)
		require.NotEmpty(t, app.ID)
		tokenRequestNumber := 3

		//WHEN
		for i := 0; i < tokenRequestNumber; i++ {
			token := fixtures.RequestOneTimeTokenForApplication(t, ctx, dexGraphQLClient, app.ID)
			assert.NotEmpty(t, token.Token)
			assert.NotEmpty(t, token.ConnectorURL)
			assert.NotEmpty(t, token.LegacyConnectorURL)
		}

		//THEN
		appExt := fixtures.GetApplication(t, ctx, dexGraphQLClient, tenantId, app.ID)
		assert.Len(t, appExt.Auths, tokenRequestNumber)
	})
}
