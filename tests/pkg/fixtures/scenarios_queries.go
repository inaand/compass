/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fixtures

import (
	"context"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	gcli "github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
)

func CreateAutomaticScenarioAssignmentInTenant(t require.TestingT, ctx context.Context, gqlClient *gcli.Client, in graphql.AutomaticScenarioAssignmentSetInput, tenantID string) *graphql.AutomaticScenarioAssignment {
	assignmentInput, err := testctx.Tc.Graphqlizer.AutomaticScenarioAssignmentSetInputToGQL(in)
	require.NoError(t, err)

	createRequest := FixCreateAutomaticScenarioAssignmentRequest(assignmentInput)

	assignment := graphql.AutomaticScenarioAssignment{}

	require.NoError(t, testctx.Tc.RunOperationWithCustomTenant(ctx, gqlClient, tenantID, createRequest, &assignment))
	require.NotEmpty(t, assignment.ScenarioName)
	return &assignment
}

func ListAutomaticScenarioAssignmentsWithinTenant(t require.TestingT, ctx context.Context, gqlClient *gcli.Client, tenantID string) graphql.AutomaticScenarioAssignmentPage {
	assignmentsPage := graphql.AutomaticScenarioAssignmentPage{}
	req := FixAutomaticScenarioAssignmentsRequest()
	err := testctx.Tc.RunOperationWithCustomTenant(ctx, gqlClient, tenantID, req, &assignmentsPage)
	require.NoError(t, err)
	return assignmentsPage
}

func DeleteAutomaticScenarioAssignmentForScenarioWithinTenant(t require.TestingT, ctx context.Context, gqlClient *gcli.Client, tenantID, scenarioName string) graphql.AutomaticScenarioAssignment {
	assignment := graphql.AutomaticScenarioAssignment{}
	req := FixDeleteAutomaticScenarioAssignmentForScenarioRequest(scenarioName)
	err := testctx.Tc.RunOperationWithCustomTenant(ctx, gqlClient, tenantID, req, &assignment)
	require.NoError(t, err)
	return assignment
}

func CleanUpAutomaticScenarioAssignmentForScenarioWithinTenant(t *testing.T, ctx context.Context, gqlClient *gcli.Client, tenantID, scenarioName string) {
	assignment := graphql.AutomaticScenarioAssignment{}
	req := FixDeleteAutomaticScenarioAssignmentForScenarioRequest(scenarioName)
	err := testctx.Tc.RunOperationWithCustomTenant(ctx, gqlClient, tenantID, req, &assignment)
	if err != nil {
		t.Logf("Error while cleaning up ASA: %s", err.Error())
	}
}

func DeleteAutomaticScenarioAssigmentForSelector(t require.TestingT, ctx context.Context, gqlClient *gcli.Client, tenantID string, selector graphql.LabelSelectorInput) []graphql.AutomaticScenarioAssignment {
	paylaod, err := testctx.Tc.Graphqlizer.LabelSelectorInputToGQL(selector)
	require.NoError(t, err)
	req := FixDeleteAutomaticScenarioAssignmentsForSelectorRequest(paylaod)

	assignment := []graphql.AutomaticScenarioAssignment{}
	err = testctx.Tc.RunOperationWithCustomTenant(ctx, gqlClient, tenantID, req, &assignment)
	require.NoError(t, err)
	return assignment
}
