// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
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

package autopipeline

import (
	"context"
	"fmt"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/controller/space"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/pipeline/autopipeline"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// Controller implements the business logic for auto-pipeline generation.
type Controller struct {
	authorizer  authz.Authorizer
	spaceFinder refcache.SpaceFinder
	repoFinder  refcache.RepoFinder
}

// NewController returns a new Controller.
func NewController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
) *Controller {
	return &Controller{
		authorizer:  authorizer,
		spaceFinder: spaceFinder,
		repoFinder:  repoFinder,
	}
}

// GenerateAutoConfig detects the tech stack and generates a default pipeline.
// In a full implementation, `files` would be fetched from the git tree of the repo.
// For now we accept a list of file paths via the API.
func (c *Controller) GenerateAutoConfig(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	files []string,
) (*types.AutoPipelineConfig, error) {
	_, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	stack := autopipeline.DetectStack(files)
	yaml := autopipeline.GeneratePipelineYAML(stack)

	config := &types.AutoPipelineConfig{
		Identifier: "auto-pipeline",
		YAML:       yaml,
		Stack:      stack,
		Generated:  types.NowMillis(),
	}

	return config, nil
}

// Helper function to get space and check access.
func (c *Controller) getSpaceCheckAccess(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	permission enum.Permission,
) (*types.SpaceCore, error) {
	return space.GetSpaceCheckAuth(ctx, c.spaceFinder, c.authorizer, session, spaceRef, permission)
}

// Ensure apiauth import is used.
var _ = apiauth.CheckSpace

// Ensure fmt import is used.
var _ = fmt.Sprintf
