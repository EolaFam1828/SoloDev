// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authn"
)

var (
	ErrMissingToken = errors.New("missing authentication token")
	ErrInvalidToken = errors.New("invalid authentication token")
)

// Authenticator wraps the existing authn.Authenticator to authenticate MCP requests.
type MCPAuthenticator struct {
	authenticator authn.Authenticator
}

// NewMCPAuthenticator creates a new MCP authenticator.
func NewMCPAuthenticator(authenticator authn.Authenticator) *MCPAuthenticator {
	return &MCPAuthenticator{authenticator: authenticator}
}

// AuthenticateHTTP authenticates an HTTP request using the standard Gitness authenticator.
func (a *MCPAuthenticator) AuthenticateHTTP(r *http.Request) (*auth.Session, error) {
	session, err := a.authenticator.Authenticate(r)
	if err != nil {
		if errors.Is(err, authn.ErrNoAuthData) {
			return nil, ErrMissingToken
		}
		return nil, ErrInvalidToken
	}
	return session, nil
}

// AuthenticateToken authenticates a bearer token directly (for stdio transport).
// It builds a synthetic HTTP request with the Authorization header set and delegates.
func (a *MCPAuthenticator) AuthenticateToken(ctx context.Context, token string) (*auth.Session, error) {
	if token == "" {
		return nil, ErrMissingToken
	}

	// Build a fake HTTP request to reuse the Authenticate interface
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+token)
	return a.AuthenticateHTTP(r)
}

// GetSpaceRef resolves the space reference from environment or params.
func GetSpaceRef(params map[string]string) string {
	if ref, ok := params["space_ref"]; ok && ref != "" {
		return ref
	}
	return os.Getenv("SOLODEV_SPACE")
}

// ExtractBearerToken extracts a bearer token from the Authorization header.
func ExtractBearerToken(header string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		return strings.TrimPrefix(header, prefix)
	}
	return ""
}
