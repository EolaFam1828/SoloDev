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

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

func (c *Controller) GetResults(ctx context.Context, session *auth.Session, spaceRef string, identifier string, limit int) ([]*types.HealthCheckResult, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
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

	hc, err := c.healthCheckStore.FindByIdentifier(ctx, parentSpace.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find health check: %w", err)
	}

	results, err := c.healthCheckResultStore.ListByHealthCheckID(ctx, hc.ID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get results: %w", err)
	}

	return results, nil
}

func (c *Controller) GetSummary(ctx context.Context, session *auth.Session, spaceRef string) ([]*types.HealthCheckSummary, error) {
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

	hcs, err := c.healthCheckStore.ListAll(ctx, parentSpace.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list health checks: %w", err)
	}

	summaries := make([]*types.HealthCheckSummary, 0, len(hcs))

	for _, hc := range hcs {
		total, err := c.healthCheckResultStore.CountTotal(ctx, hc.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to count total results: %w", err)
		}

		successCount, err := c.healthCheckResultStore.CountByStatus(ctx, hc.ID, string(types.HealthCheckStatusUp))
		if err != nil {
			return nil, fmt.Errorf("failed to count successful results: %w", err)
		}

		failedCount, err := c.healthCheckResultStore.CountByStatus(ctx, hc.ID, string(types.HealthCheckStatusDown))
		if err != nil {
			return nil, fmt.Errorf("failed to count failed results: %w", err)
		}

		avgResponseTime, err := c.healthCheckResultStore.GetAverageResponseTime(ctx, hc.ID, 24)
		if err != nil {
			avgResponseTime = 0
		}

		uptimePercentage := 0.0
		if total > 0 {
			uptimePercentage = (float64(successCount) / float64(total)) * 100
		}

		summary := &types.HealthCheckSummary{
			HealthCheckID:       hc.ID,
			Identifier:          hc.Identifier,
			Name:                hc.Name,
			CurrentStatus:       hc.LastStatus,
			UptimePercentage:    uptimePercentage,
			TotalChecks:         total,
			SuccessfulChecks:    successCount,
			FailedChecks:        failedCount,
			AverageResponseTime: avgResponseTime,
			LastCheckedAt:       hc.LastCheckedAt,
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
