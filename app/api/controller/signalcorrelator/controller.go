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

package signalcorrelator

import (
	"context"
	"fmt"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/services/signalcorrelator"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// Controller exposes the signal correlator to the REST API.
type Controller struct {
	authorizer  authz.Authorizer
	spaceFinder refcache.SpaceFinder
	service     *signalcorrelator.Service
}

// ProvideController creates a new signal correlator controller.
func ProvideController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	service *signalcorrelator.Service,
) *Controller {
	return &Controller{
		authorizer:  authorizer,
		spaceFinder: spaceFinder,
		service:     service,
	}
}

// Correlate runs cross-domain signal correlation for a space.
func (c *Controller) Correlate(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter types.CorrelatedIncidentFilter,
) ([]types.CorrelatedIncident, error) {
	parentSpace, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find parent space: %w", err)
	}

	err = apiauth.CheckSpace(ctx, c.authorizer, session, parentSpace, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	incidents, err := c.service.Correlate(ctx, parentSpace.ID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to correlate signals: %w", err)
	}

	return incidents, nil
}
