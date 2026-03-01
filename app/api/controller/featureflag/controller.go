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

package featureflag

import (
	"context"
	"fmt"
	"time"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/spacefinder"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type Controller struct {
	featureFlagStore store.FeatureFlagStore
	spaceFinder      spacefinder.Finder
	authorizer       authz.Authorizer
}

func NewController(
	featureFlagStore store.FeatureFlagStore,
	spaceFinder spacefinder.Finder,
	authorizer authz.Authorizer,
) *Controller {
	return &Controller{
		featureFlagStore: featureFlagStore,
		spaceFinder:      spaceFinder,
		authorizer:       authorizer,
	}
}

// getSpaceCheckFeatureFlagAccess fetches a space and checks if the current user has permission to access feature flags.
func (c *Controller) getSpaceCheckFeatureFlagAccess(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	reqPermission enum.Permission,
) (*types.Space, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space by ref: %w", err)
	}

	if err = c.authorizer.Check(ctx, session, &authz.ResourceSpace{Space: space}, reqPermission); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	return space, nil
}

// Create creates a new feature flag in the specified space.
func (c *Controller) Create(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *CreateInput,
) (*types.FeatureFlag, error) {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagCreate)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()

	featureFlag := &types.FeatureFlag{
		SpaceID:               space.ID,
		Identifier:            in.Identifier,
		Name:                  in.Name,
		Description:           in.Description,
		Kind:                  in.Kind,
		DefaultOnVariation:    in.DefaultOnVariation,
		DefaultOffVariation:   in.DefaultOffVariation,
		Enabled:               in.Enabled,
		Variations:            in.Variations,
		Tags:                  in.Tags,
		Permanent:             in.Permanent,
		CreatedBy:             session.Principal.ID,
		Created:               now,
		Updated:               now,
		Version:               0,
	}

	if err = c.featureFlagStore.Create(ctx, featureFlag); err != nil {
		return nil, fmt.Errorf("failed to create feature flag: %w", err)
	}

	return featureFlag, nil
}

// Find retrieves a feature flag by identifier.
func (c *Controller) Find(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
) (*types.FeatureFlag, error) {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagView)
	if err != nil {
		return nil, err
	}

	featureFlag, err := c.featureFlagStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find feature flag: %w", err)
	}

	return featureFlag, nil
}

// Update updates an existing feature flag.
func (c *Controller) Update(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
	in *UpdateInput,
) (*types.FeatureFlag, error) {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagEdit)
	if err != nil {
		return nil, err
	}

	featureFlag, err := c.featureFlagStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find feature flag: %w", err)
	}

	// Update fields if provided
	if in.Name != nil {
		featureFlag.Name = *in.Name
	}
	if in.Description != nil {
		featureFlag.Description = *in.Description
	}
	if in.Kind != nil {
		featureFlag.Kind = *in.Kind
	}
	if in.DefaultOnVariation != nil {
		featureFlag.DefaultOnVariation = *in.DefaultOnVariation
	}
	if in.DefaultOffVariation != nil {
		featureFlag.DefaultOffVariation = *in.DefaultOffVariation
	}
	if in.Enabled != nil {
		featureFlag.Enabled = *in.Enabled
	}
	if in.Variations != nil {
		featureFlag.Variations = *in.Variations
	}
	if in.Tags != nil {
		featureFlag.Tags = *in.Tags
	}
	if in.Permanent != nil {
		featureFlag.Permanent = *in.Permanent
	}

	if err = c.featureFlagStore.Update(ctx, featureFlag); err != nil {
		return nil, fmt.Errorf("failed to update feature flag: %w", err)
	}

	return featureFlag, nil
}

// Toggle enables or disables a feature flag.
func (c *Controller) Toggle(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
	enabled bool,
) (*types.FeatureFlag, error) {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagEdit)
	if err != nil {
		return nil, err
	}

	featureFlag, err := c.featureFlagStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find feature flag: %w", err)
	}

	featureFlag.Enabled = enabled

	if err = c.featureFlagStore.Update(ctx, featureFlag); err != nil {
		return nil, fmt.Errorf("failed to update feature flag: %w", err)
	}

	return featureFlag, nil
}

// Delete deletes a feature flag.
func (c *Controller) Delete(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
) error {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagDelete)
	if err != nil {
		return err
	}

	featureFlag, err := c.featureFlagStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return fmt.Errorf("failed to find feature flag: %w", err)
	}

	if err = c.featureFlagStore.Delete(ctx, featureFlag.ID); err != nil {
		return fmt.Errorf("failed to delete feature flag: %w", err)
	}

	return nil
}

// List lists feature flags in a space.
func (c *Controller) List(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.FeatureFlagFilter,
) ([]*types.FeatureFlag, error) {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagView)
	if err != nil {
		return nil, err
	}

	featureFlags, err := c.featureFlagStore.List(ctx, space.ID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list feature flags: %w", err)
	}

	return featureFlags, nil
}

// Count counts feature flags in a space.
func (c *Controller) Count(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.FeatureFlagFilter,
) (int64, error) {
	space, err := c.getSpaceCheckFeatureFlagAccess(ctx, session, spaceRef, enum.PermissionFeatureFlagView)
	if err != nil {
		return 0, err
	}

	count, err := c.featureFlagStore.Count(ctx, space.ID, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count feature flags: %w", err)
	}

	return count, nil
}
