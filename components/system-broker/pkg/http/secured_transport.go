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

package http

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"sync"
)

type TokenProvider interface {
	GetAuthorizationToken(ctx context.Context) (Token, error)
	Matches(request *http.Request) bool
}

type SecuredTransport struct {
	roundTripper   HTTPRoundTripper
	tokenProviders []TokenProvider
	lock           sync.Mutex

	token Token
}

func (c *SecuredTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := c.refreshToken(request); err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	return c.roundTripper.RoundTrip(request)
}

func (c *SecuredTransport) refreshToken(request *http.Request) error {
	if !c.token.EmptyOrExpired() {
		return nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	for _, tokenProvider := range c.tokenProviders {
		if !tokenProvider.Matches(request) {
			continue
		}

		token, err := tokenProvider.GetAuthorizationToken(request.Context())
		if err != nil {
			return errors.Wrap(err, "error while obtaining token")
		}
		c.token = token
	}

	return nil
}
