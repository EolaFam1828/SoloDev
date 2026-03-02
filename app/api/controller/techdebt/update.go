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
	"time"

	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/errors"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

func (c *Controller) Update(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
	in *types.TechDebtUpdateInput,
) (*types.TechDebt, error) {
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	if identifier == "" {
		return nil, errors.InvalidArgument("Identifier is required")
	}

	if in == nil {
		return nil, errors.InvalidArgument("Input cannot be nil")
	}

	td, err := c.techDebtStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if in.Title != "" {
		td.Title = in.Title
	}
	if in.Description != "" {
		td.Description = in.Description
	}
	if in.Severity != "" {
		td.Severity = types.TechDebtSeverity(in.Severity)
	}
	if in.Status != "" {
		td.Status = types.TechDebtStatus(in.Status)
	}
	if in.Category != "" {
		td.Category = types.TechDebtCategory(in.Category)
	}
	if in.FilePath != "" {
		td.FilePath = in.FilePath
	}
	if in.LineStart > 0 {
		td.LineStart = in.LineStart
	}
	if in.LineEnd > 0 {
		td.LineEnd = in.LineEnd
	}
	if in.EstimatedEffort != "" {
		td.EstimatedEffort = in.EstimatedEffort
	}
	if len(in.Tags) > 0 {
		td.Tags = in.Tags
	}
	if in.DueDate > 0 {
		td.DueDate = in.DueDate
	}

	td.Updated = time.Now().UnixMilli()
	td.Version++

	if err := c.techDebtStore.Update(ctx, td); err != nil {
		return nil, err
	}

	return td, nil
}
