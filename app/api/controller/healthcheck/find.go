// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed request parameter extraction and list query parsing.
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

func (c *Controller) Find(ctx context.Context, session *auth.Session, spaceRef string, identifier string) (*types.HealthCheck, error) {
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

	return hc, nil
}

func (c *Controller) List(ctx context.Context, session *auth.Session, spaceRef string, filter types.ListQueryFilter) ([]*types.HealthCheck, error) {
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

	hcs, err := c.healthCheckStore.List(ctx, parentSpace.ID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list health checks: %w", err)
	}

	return hcs, nil
}
