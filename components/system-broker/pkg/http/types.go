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
	"net/http"
	"time"
)

type HTTPRoundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}

type Token struct {
	AccessToken string `json:"access_token"`
	Expiration  int64  `json:"expires_in"`
}

//TODO check if we need to add buffer time?
func (token Token) EmptyOrExpired() bool {
	if token.AccessToken == "" {
		return true
	}

	if token.Expiration == 0 {
		return false
	}

	expiration := time.Unix(token.Expiration, 0)
	return time.Now().After(expiration)
}
