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

package director

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	graphqlbroker "github.com/kyma-incubator/compass/components/system-broker/pkg/graphql"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/types"
	gcli "github.com/machinebox/graphql"
)

// client implements the DirectorClient interface
type client struct {
	types.ApplicationLister
	httpClient        *http.Client
	directorURL       string
	operationEndpoint string
	gqlClient         *gcli.Client
}

type Request struct {
	OperationType graphql.OperationType `json:"operation_type"`
	ResourceType  resource.Type         `json:"resource_type"`
	ResourceID    string                `json:"resource_id"`
	Error         string                `json:"error,omitempty"`
}

// NewClient constructs a default implementation of the Client interface
func NewClient(operationEndpoint string, cfg *graphqlbroker.Config, httpClient *http.Client) (*client, error) {
	graphqlClient, err := graphqlbroker.PrepareGqlClientWithHttpClient(cfg, httpClient)
	if err != nil {
		return nil, err
	}

	return &client{
		ApplicationLister: graphqlClient,
		gqlClient:         gcli.NewClient(cfg.GraphqlEndpoint, gcli.WithHTTPClient(httpClient)),
		httpClient:        httpClient,
		operationEndpoint: operationEndpoint,
	}, nil
}

// UpdateOperation makes an http request to the Director to notify about any operation state changes
func (c *client) UpdateOperation(ctx context.Context, request *Request) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.operationEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code when notifying director for application state: %d", resp.StatusCode)
	}

	return nil
}
// TODO add other types
// SetBundleInstanceAuth makes a graphql request to the Director to set instance auths
func (c *client) SetBundleInstanceAuth(ctx context.Context, instanceAuthID string, in *graphql.APIKeyCredentialDataInput) error {
	mutation := fmt.Sprintf(`mutation {
		result: setBundleInstanceAuth(authID: "%s", in: {
			auth:{
				credential:  {
					apiKey: {
						apiKey: "%s",
						tokenServerURL: "%s",
					},
				}
			}
		}) {
			id
		}}`, instanceAuthID, in.APIKey, in.TokenServerURL)

	gqlRequest := gcli.NewRequest(mutation)
	resp := struct{ id string }{}
	return c.gqlClient.Run(ctx, gqlRequest, &resp)
}
