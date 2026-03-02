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

	apiauth "github.com/EolaFam1828/SoloDev/app/api/auth"
	"github.com/EolaFam1828/SoloDev/app/api/usererror"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

type UpdateInput struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	URL             *string `json:"url"`
	Method          *string `json:"method"`
	ExpectedStatus  *int    `json:"expected_status"`
	IntervalSeconds *int    `json:"interval_seconds"`
	TimeoutSeconds  *int    `json:"timeout_seconds"`
	Enabled         *bool   `json:"enabled"`
	Headers         *string `json:"headers"`
	Body            *string `json:"body"`
	Tags            *string `json:"tags"`
}

func (c *Controller) Update(ctx context.Context, session *auth.Session, spaceRef string, identifier string, in *UpdateInput) (*types.HealthCheck, error) {
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

	hc, err := c.healthCheckStore.FindByIdentifier(ctx, parentSpace.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find health check: %w", err)
	}

	updated, err := c.healthCheckStore.UpdateOptLock(ctx, hc, func(hc *types.HealthCheck) error {
		if in.Name != nil {
			hc.Name = strings.TrimSpace(*in.Name)
			if len(hc.Name) == 0 {
				return usererror.BadRequest("Name is required")
			}
		}

		if in.Description != nil {
			hc.Description = strings.TrimSpace(*in.Description)
		}

		if in.URL != nil {
			hc.URL = strings.TrimSpace(*in.URL)
			if len(hc.URL) == 0 {
				return errInvalidURL
			}
		}

		if in.Method != nil {
			hc.Method = strings.ToUpper(strings.TrimSpace(*in.Method))
			if hc.Method != "GET" && hc.Method != "POST" && hc.Method != "HEAD" {
				return errInvalidMethod
			}
		}

		if in.ExpectedStatus != nil {
			if *in.ExpectedStatus < 100 || *in.ExpectedStatus > 599 {
				return errInvalidStatus
			}
			hc.ExpectedStatus = *in.ExpectedStatus
		}

		if in.IntervalSeconds != nil {
			if *in.IntervalSeconds < 60 || *in.IntervalSeconds > 86400 {
				return errInvalidInterval
			}
			hc.IntervalSeconds = *in.IntervalSeconds
		}

		if in.TimeoutSeconds != nil {
			if *in.TimeoutSeconds < 1 || *in.TimeoutSeconds > 300 {
				return errInvalidTimeout
			}
			hc.TimeoutSeconds = *in.TimeoutSeconds
		}

		if in.Enabled != nil {
			hc.Enabled = *in.Enabled
		}

		if in.Headers != nil {
			hc.Headers = strings.TrimSpace(*in.Headers)
		}

		if in.Body != nil {
			hc.Body = strings.TrimSpace(*in.Body)
		}

		if in.Tags != nil {
			hc.Tags = strings.TrimSpace(*in.Tags)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update health check: %w", err)
	}

	return updated, nil
}

func (c *Controller) Toggle(ctx context.Context, session *auth.Session, spaceRef string, identifier string) (*types.HealthCheck, error) {
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

	hc, err := c.healthCheckStore.FindByIdentifier(ctx, parentSpace.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find health check: %w", err)
	}

	updated, err := c.healthCheckStore.UpdateOptLock(ctx, hc, func(hc *types.HealthCheck) error {
		hc.Enabled = !hc.Enabled
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to toggle health check: %w", err)
	}

	return updated, nil
}
