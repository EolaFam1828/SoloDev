// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed request parameter extraction for API compliance.
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

package healthcheck

import (
	"context"
	"fmt"
	"strings"
	"time"

	apiauth "github.com/EolaFam1828/SoloDev/app/api/auth"
	"github.com/EolaFam1828/SoloDev/app/api/usererror"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/check"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

var (
	errHealthCheckRequiresParent = usererror.BadRequest(
		"Parent space required - standalone health checks are not supported.")
	errInvalidURL = usererror.BadRequest(
		"Invalid URL provided.")
	errInvalidMethod = usererror.BadRequest(
		"Invalid HTTP method. Must be GET, POST, or HEAD.")
	errInvalidInterval = usererror.BadRequest(
		"Invalid interval. Must be between 60 and 86400 seconds.")
	errInvalidTimeout = usererror.BadRequest(
		"Invalid timeout. Must be between 1 and 300 seconds.")
	errInvalidStatus = usererror.BadRequest(
		"Invalid expected status code. Must be between 100 and 599.")
)

type CreateInput struct {
	Identifier      string `json:"identifier"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	Method          string `json:"method"`
	ExpectedStatus  int    `json:"expected_status"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	Enabled         bool   `json:"enabled"`
	Headers         string `json:"headers"`
	Body            string `json:"body"`
	Tags            string `json:"tags"`
}

func (c *Controller) Create(ctx context.Context, session *auth.Session, spaceRef string, in *CreateInput) (*types.HealthCheck, error) {
	if err := c.sanitizeCreateInput(in); err != nil {
		return nil, fmt.Errorf("failed to sanitize input: %w", err)
	}

	parentSpace, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find parent space: %w", err)
	}

	err = apiauth.CheckSpace(
		ctx,
		c.authorizer,
		session,
		parentSpace,
		enum.PermissionSpaceView,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	hc := &types.HealthCheck{
		SpaceID:         parentSpace.ID,
		Identifier:      in.Identifier,
		Name:            in.Name,
		Description:     in.Description,
		URL:             in.URL,
		Method:          in.Method,
		ExpectedStatus:  in.ExpectedStatus,
		IntervalSeconds: in.IntervalSeconds,
		TimeoutSeconds:  in.TimeoutSeconds,
		Enabled:         in.Enabled,
		Headers:         in.Headers,
		Body:            in.Body,
		Tags:            in.Tags,
		LastStatus:      string(types.HealthCheckStatusUnknown),
		LastCheckedAt:   0,
		LastResponseTime: 0,
		ConsecutiveFailures: 0,
		CreatedBy:       session.Principal.ID,
		Created:         now,
		Updated:         now,
		Version:         0,
	}

	err = c.healthCheckStore.Create(ctx, hc)
	if err != nil {
		return nil, fmt.Errorf("health check creation failed: %w", err)
	}

	return hc, nil
}

func (c *Controller) sanitizeCreateInput(in *CreateInput) error {
	if err := check.Identifier(in.Identifier); err != nil {
		return err
	}

	in.Name = strings.TrimSpace(in.Name)
	if len(in.Name) == 0 {
		return usererror.BadRequest("Name is required")
	}

	in.Description = strings.TrimSpace(in.Description)

	in.URL = strings.TrimSpace(in.URL)
	if len(in.URL) == 0 {
		return errInvalidURL
	}

	in.Method = strings.ToUpper(strings.TrimSpace(in.Method))
	if in.Method == "" {
		in.Method = "GET"
	}
	if in.Method != "GET" && in.Method != "POST" && in.Method != "HEAD" {
		return errInvalidMethod
	}

	if in.ExpectedStatus < 100 || in.ExpectedStatus > 599 {
		return errInvalidStatus
	}

	if in.IntervalSeconds < 60 || in.IntervalSeconds > 86400 {
		return errInvalidInterval
	}

	if in.TimeoutSeconds < 1 || in.TimeoutSeconds > 300 {
		return errInvalidTimeout
	}

	in.Headers = strings.TrimSpace(in.Headers)
	in.Body = strings.TrimSpace(in.Body)
	in.Tags = strings.TrimSpace(in.Tags)

	return nil
}
