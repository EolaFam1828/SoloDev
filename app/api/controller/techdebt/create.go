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

func (c *Controller) Create(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *types.TechDebtCreateInput,
) (*types.TechDebt, error) {
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	if in == nil {
		return nil, errors.InvalidArgument("Input cannot be nil")
	}

	if in.Identifier == "" {
		return nil, errors.InvalidArgument("Identifier is required")
	}

	if in.Title == "" {
		return nil, errors.InvalidArgument("Title is required")
	}

	// Check if identifier already exists
	_, err = c.techDebtStore.FindByIdentifier(ctx, space.ID, in.Identifier)
	if err == nil {
		return nil, errors.Conflict("Technical debt item with this identifier already exists")
	}
	if !errors.IsNotFound(err) {
		return nil, err
	}

	now := time.Now().UnixMilli()
	td := &types.TechDebt{
		SpaceID:         space.ID,
		RepoID:          in.RepoID,
		Identifier:      in.Identifier,
		Title:           in.Title,
		Description:     in.Description,
		Severity:        types.TechDebtSeverity(in.Severity),
		Status:          types.TechDebtStatusOpen,
		Category:        types.TechDebtCategory(in.Category),
		FilePath:        in.FilePath,
		LineStart:       in.LineStart,
		LineEnd:         in.LineEnd,
		EstimatedEffort: in.EstimatedEffort,
		Tags:            in.Tags,
		DueDate:         in.DueDate,
		CreatedBy:       session.Principal.ID,
		Created:         now,
		Updated:         now,
		Version:         1,
	}

	if in.Status != "" {
		td.Status = types.TechDebtStatus(in.Status)
	}

	if err := c.techDebtStore.Create(ctx, td); err != nil {
		return nil, err
	}

	return td, nil
}
