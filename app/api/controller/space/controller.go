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

package space

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	apiauth "github.com/EolaFam1828/SoloDev/app/api/auth"
	"github.com/EolaFam1828/SoloDev/app/api/controller/limiter"
	"github.com/EolaFam1828/SoloDev/app/api/controller/repo"
	"github.com/EolaFam1828/SoloDev/app/api/usererror"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/autolink"
	"github.com/EolaFam1828/SoloDev/app/services/exporter"
	"github.com/EolaFam1828/SoloDev/app/services/gitspace"
	"github.com/EolaFam1828/SoloDev/app/services/importer"
	"github.com/EolaFam1828/SoloDev/app/services/infraprovider"
	"github.com/EolaFam1828/SoloDev/app/services/instrument"
	"github.com/EolaFam1828/SoloDev/app/services/label"
	"github.com/EolaFam1828/SoloDev/app/services/publicaccess"
	"github.com/EolaFam1828/SoloDev/app/services/pullreq"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/services/rules"
	"github.com/EolaFam1828/SoloDev/app/services/space"
	"github.com/EolaFam1828/SoloDev/app/sse"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/audit"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/check"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

var (
	// TODO (Nested Spaces): Remove once full support is added
	errNestedSpacesNotSupported    = usererror.BadRequestf("Nested spaces are not supported.")
	errPublicSpaceCreationDisabled = usererror.BadRequestf("Public space creation is disabled.")
)

//nolint:revive
type SpaceOutput struct {
	types.Space
	IsPublic bool `json:"is_public" yaml:"is_public"`
}

// TODO [CODE-1363]: remove after identifier migration.
func (s SpaceOutput) MarshalJSON() ([]byte, error) {
	// alias allows us to embed the original object while avoiding an infinite loop of marshaling.
	type alias SpaceOutput
	return json.Marshal(&struct {
		alias
		UID string `json:"uid"`
	}{
		alias: (alias)(s),
		UID:   s.Identifier,
	})
}

type Controller struct {
	nestedSpacesEnabled bool

	tx                  dbtx.Transactor
	urlProvider         url.Provider
	sseStreamer         sse.Streamer
	identifierCheck     check.SpaceIdentifier
	authorizer          authz.Authorizer
	spacePathStore      store.SpacePathStore
	pipelineStore       store.PipelineStore
	secretStore         store.SecretStore
	connectorStore      store.ConnectorStore
	templateStore       store.TemplateStore
	spaceStore          store.SpaceStore
	repoStore           store.RepoStore
	principalStore      store.PrincipalStore
	repoCtrl            *repo.Controller
	membershipStore     store.MembershipStore
	prListService       *pullreq.ListService
	spaceFinder         refcache.SpaceFinder
	repoFinder          refcache.RepoFinder
	importer            *importer.JobRepository
	exporter            *exporter.Repository
	resourceLimiter     limiter.ResourceLimiter
	publicAccess        publicaccess.Service
	auditService        audit.Service
	gitspaceSvc         *gitspace.Service
	labelSvc            *label.Service
	instrumentation     instrument.Service
	executionStore      store.ExecutionStore
	rulesSvc            *rules.Service
	usageMetricStore    store.UsageMetricStore
	repoIdentifierCheck check.RepoIdentifier
	infraProviderSvc    *infraprovider.Service
	favoriteStore       store.FavoriteStore
	autolinkSvc         *autolink.Service
	spaceSvc            *space.Service
}

func NewController(config *types.Config, tx dbtx.Transactor, urlProvider url.Provider,
	sseStreamer sse.Streamer, identifierCheck check.SpaceIdentifier, authorizer authz.Authorizer,
	spacePathStore store.SpacePathStore, pipelineStore store.PipelineStore, secretStore store.SecretStore,
	connectorStore store.ConnectorStore, templateStore store.TemplateStore, spaceStore store.SpaceStore,
	repoStore store.RepoStore, principalStore store.PrincipalStore, repoCtrl *repo.Controller,
	membershipStore store.MembershipStore, prListService *pullreq.ListService,
	spaceFinder refcache.SpaceFinder, repoFinder refcache.RepoFinder,
	importer *importer.JobRepository, exporter *exporter.Repository,
	limiter limiter.ResourceLimiter, publicAccess publicaccess.Service, auditService audit.Service,
	gitspaceSvc *gitspace.Service, labelSvc *label.Service,
	instrumentation instrument.Service, executionStore store.ExecutionStore,
	rulesSvc *rules.Service, usageMetricStore store.UsageMetricStore, repoIdentifierCheck check.RepoIdentifier,
	infraProviderSvc *infraprovider.Service, favoriteStore store.FavoriteStore, autolinkSvc *autolink.Service,
	spaceSvc *space.Service,
) *Controller {
	return &Controller{
		nestedSpacesEnabled: config.NestedSpacesEnabled,
		tx:                  tx,
		urlProvider:         urlProvider,
		sseStreamer:         sseStreamer,
		identifierCheck:     identifierCheck,
		authorizer:          authorizer,
		spacePathStore:      spacePathStore,
		pipelineStore:       pipelineStore,
		secretStore:         secretStore,
		connectorStore:      connectorStore,
		templateStore:       templateStore,
		spaceStore:          spaceStore,
		repoStore:           repoStore,
		principalStore:      principalStore,
		repoCtrl:            repoCtrl,
		membershipStore:     membershipStore,
		prListService:       prListService,
		spaceFinder:         spaceFinder,
		repoFinder:          repoFinder,
		importer:            importer,
		exporter:            exporter,
		resourceLimiter:     limiter,
		publicAccess:        publicAccess,
		auditService:        auditService,
		gitspaceSvc:         gitspaceSvc,
		labelSvc:            labelSvc,
		instrumentation:     instrumentation,
		executionStore:      executionStore,
		rulesSvc:            rulesSvc,
		usageMetricStore:    usageMetricStore,
		repoIdentifierCheck: repoIdentifierCheck,
		infraProviderSvc:    infraProviderSvc,
		favoriteStore:       favoriteStore,
		autolinkSvc:         autolinkSvc,
		spaceSvc:            spaceSvc,
	}
}

// getSpaceCheckAuth checks whether the user has the requested permission on the provided space and returns the space.
func (c *Controller) getSpaceCheckAuth(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	permission enum.Permission,
) (*types.SpaceCore, error) {
	return GetSpaceCheckAuth(ctx, c.spaceFinder, c.authorizer, session, spaceRef, permission)
}

func (c *Controller) getSpaceCheckAuthRepoCreation(
	ctx context.Context,
	session *auth.Session,
	parentRef string,
) (*types.SpaceCore, error) {
	return repo.GetSpaceCheckAuthRepoCreation(ctx, c.spaceFinder, c.authorizer, session, parentRef)
}

func (c *Controller) getSpaceCheckAuthSpaceCreation(
	ctx context.Context,
	session *auth.Session,
	parentRef string,
) (*types.SpaceCore, error) {
	parentRefAsID, err := strconv.ParseInt(parentRef, 10, 64)
	if (parentRefAsID <= 0 && err == nil) || (len(strings.TrimSpace(parentRef)) == 0) {
		// TODO: Restrict top level space creation - should be move to authorizer?
		if auth.IsAnonymousSession(session) {
			return nil, fmt.Errorf("anonymous user not allowed to create top level spaces: %w", usererror.ErrUnauthorized)
		}

		return &types.SpaceCore{}, nil
	}

	parentSpace, err := c.spaceFinder.FindByRef(ctx, parentRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent space: %w", err)
	}

	if err = apiauth.CheckSpaceScope(
		ctx,
		c.authorizer,
		session,
		parentSpace,
		enum.ResourceTypeSpace,
		enum.PermissionSpaceEdit,
	); err != nil {
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	return parentSpace, nil
}
