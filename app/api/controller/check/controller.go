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

package check

import (
	"context"
	"fmt"

	apiauth "github.com/EolaFam1828/SoloDev/app/api/auth"
	"github.com/EolaFam1828/SoloDev/app/api/controller/space"
	"github.com/EolaFam1828/SoloDev/app/api/usererror"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	checkevents "github.com/EolaFam1828/SoloDev/app/events/check"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/sse"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/git"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

type Controller struct {
	tx            dbtx.Transactor
	authorizer    authz.Authorizer
	spaceStore    store.SpaceStore
	checkStore    store.CheckStore
	spaceFinder   refcache.SpaceFinder
	repoFinder    refcache.RepoFinder
	git           git.Interface
	sanitizers    map[enum.CheckPayloadKind]func(in *ReportInput, s *auth.Session) error
	sseStreamer   sse.Streamer
	eventReporter *checkevents.Reporter
}

func NewController(
	tx dbtx.Transactor,
	authorizer authz.Authorizer,
	spaceStore store.SpaceStore,
	checkStore store.CheckStore,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	git git.Interface,
	sanitizers map[enum.CheckPayloadKind]func(in *ReportInput, s *auth.Session) error,
	sseStreamer sse.Streamer,
	eventReporter *checkevents.Reporter,
) *Controller {
	return &Controller{
		tx:            tx,
		authorizer:    authorizer,
		spaceStore:    spaceStore,
		checkStore:    checkStore,
		spaceFinder:   spaceFinder,
		repoFinder:    repoFinder,
		git:           git,
		sanitizers:    sanitizers,
		sseStreamer:   sseStreamer,
		eventReporter: eventReporter,
	}
}

//nolint:unparam
func (c *Controller) getRepoCheckAccess(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	reqPermission enum.Permission,
	allowedRepoStates ...enum.RepoState,
) (*types.RepositoryCore, error) {
	if repoRef == "" {
		return nil, usererror.BadRequest("A valid repository reference must be provided.")
	}

	repo, err := c.repoFinder.FindByRef(ctx, repoRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find repository: %w", err)
	}

	if err := apiauth.CheckRepoState(ctx, session, repo, reqPermission, allowedRepoStates...); err != nil {
		return nil, err
	}

	if err = apiauth.CheckRepo(ctx, c.authorizer, session, repo, reqPermission); err != nil {
		return nil, fmt.Errorf("access check failed: %w", err)
	}

	return repo, nil
}

func (c *Controller) getSpaceCheckAccess(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	permission enum.Permission,
) (*types.SpaceCore, error) {
	return space.GetSpaceCheckAuth(ctx, c.spaceFinder, c.authorizer, session, spaceRef, permission)
}
