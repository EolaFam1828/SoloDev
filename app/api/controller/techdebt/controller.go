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

package techdebt

import (
	"context"
	"fmt"
	"time"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/errors"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type Controller struct {
	authorizer    authz.Authorizer
	spaceFinder   refcache.SpaceFinder
	techDebtStore store.TechDebtStore
}

func NewController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	techDebtStore store.TechDebtStore,
) *Controller {
	return &Controller{
		authorizer:    authorizer,
		spaceFinder:   spaceFinder,
		techDebtStore: techDebtStore,
	}
}

func (c *Controller) getSpaceCheckAccess(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	permission enum.Permission,
) (*types.SpaceCore, error) {
	if spaceRef == "" {
		return nil, errors.InvalidArgument("A valid space reference must be provided.")
	}

	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err := c.authorizer.Check(ctx, session.Principal, enum.ResourceTypeSpace, space.Identifier, permission); err != nil {
		return nil, fmt.Errorf("failed to verify authorization: %w", err)
	}

	return space, nil
}
