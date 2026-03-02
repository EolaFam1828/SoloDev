// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Updated constructor and ListResponse type for MCP integration.
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

	apiauth "github.com/EolaFam1828/SoloDev/app/api/auth"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/errors"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
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

	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, permission); err != nil {
		return nil, fmt.Errorf("failed to verify authorization: %w", err)
	}

	return space, nil
}
