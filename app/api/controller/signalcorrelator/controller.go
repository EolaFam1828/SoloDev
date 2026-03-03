// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
