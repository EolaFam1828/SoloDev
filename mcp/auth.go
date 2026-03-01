// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
