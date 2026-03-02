// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Updated list return type to ListResponse struct.
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
	"time"

	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

type ListResponse struct {
	Items []types.TechDebt `json:"items"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

func (c *Controller) List(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.TechDebtFilter,
) (*ListResponse, error) {
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &types.TechDebtFilter{}
	}

	// Set defaults for pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	items, err := c.techDebtStore.List(ctx, space.ID, filter)
	if err != nil {
		return nil, err
	}

	total, err := c.techDebtStore.Count(ctx, space.ID, filter)
	if err != nil {
		return nil, err
	}

	if items == nil {
		items = []*types.TechDebt{}
	}

	return &ListResponse{
		Items: convertTechDebtList(items),
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}, nil
}

func (c *Controller) Summary(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.TechDebtFilter,
) (*types.TechDebtSummary, error) {
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &types.TechDebtFilter{}
	}

	return c.techDebtStore.Summary(ctx, space.ID, filter)
}

func (c *Controller) Resolve(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
) (*types.TechDebt, error) {
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	td, err := c.techDebtStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, err
	}

	td.Status = types.TechDebtStatusResolved
	td.ResolvedAt = time.Now().UnixMilli()
	td.ResolvedBy = session.Principal.ID
	td.Version++

	if err := c.techDebtStore.Update(ctx, td); err != nil {
		return nil, err
	}

	return td, nil
}

func convertTechDebtList(items []*types.TechDebt) []types.TechDebt {
	result := make([]types.TechDebt, len(items))
	for i, item := range items {
		result[i] = *item
	}
	return result
}
