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

package auth

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
)

type contextKey string

const (
	// CredentialsCtxKey missing godoc
	CredentialsCtxKey contextKey = "CredentialsCtxKey"

	// BasicCredentialType missing godoc
	BasicCredentialType CredentialType = "BasicCredentials"
	// OAuthCredentialType missing godoc
	OAuthCredentialType CredentialType = "OAuthCredentials"
)

// CredentialType specifies a dedicated string type to differentiate every Credentials type
type CredentialType string

// Credentials denotes an authentication credentials type
type Credentials interface {
	Get() interface{}
	Type() CredentialType
}

// BasicCredentials implements the Credentials interface for the Basic Authentication flow
type BasicCredentials struct {
	Username string
	Password string
}

// OAuthCredentials implements the Credentials interface for the OAuth flow
type OAuthCredentials struct {
	ClientID          string
	ClientSecret      string
	TokenURL          string
	Scopes            string
	AdditionalHeaders map[string]string
}

// Get returns the specified Credentials implementation
func (bc *BasicCredentials) Get() interface{} {
	return bc
}

// Type returns the specified Credentials implementation type
func (bc *BasicCredentials) Type() CredentialType {
	return BasicCredentialType
}

// Get returns the specified Credentials implementation
func (oc *OAuthCredentials) Get() interface{} {
	return oc
}

// Type returns the specified Credentials implementation type
func (oc *OAuthCredentials) Type() CredentialType {
	return OAuthCredentialType
}

// LoadFromContext retrieves the credentials from the provided context
func LoadFromContext(ctx context.Context) (Credentials, error) {
	credentials, ok := ctx.Value(CredentialsCtxKey).(Credentials)

	if !ok {
		return nil, apperrors.NewNotFoundErrorWithType("credentials")
	}

	return credentials, nil
}

// SaveToContext saves the given credentials in the specified context
func SaveToContext(ctx context.Context, credentials Credentials) context.Context {
	return context.WithValue(ctx, CredentialsCtxKey, credentials)
}
