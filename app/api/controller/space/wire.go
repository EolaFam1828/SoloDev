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
	"github.com/EolaFam1828/SoloDev/app/api/controller/limiter"
	"github.com/EolaFam1828/SoloDev/app/api/controller/repo"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/autolink"
	"github.com/EolaFam1828/SoloDev/app/services/exporter"
	"github.com/EolaFam1828/SoloDev/app/services/gitspace"
	"github.com/EolaFam1828/SoloDev/app/services/importer"
	infraprovider2 "github.com/EolaFam1828/SoloDev/app/services/infraprovider"
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

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(config *types.Config, tx dbtx.Transactor, urlProvider url.Provider, sseStreamer sse.Streamer,
	identifierCheck check.SpaceIdentifier, authorizer authz.Authorizer, spacePathStore store.SpacePathStore,
	pipelineStore store.PipelineStore, secretStore store.SecretStore,
	connectorStore store.ConnectorStore, templateStore store.TemplateStore,
	spaceStore store.SpaceStore, repoStore store.RepoStore, principalStore store.PrincipalStore,
	repoCtrl *repo.Controller, membershipStore store.MembershipStore, prListService *pullreq.ListService,
	spaceFinder refcache.SpaceFinder, repoFinder refcache.RepoFinder,
	importer *importer.JobRepository, exporter *exporter.Repository,
	limiter limiter.ResourceLimiter, publicAccess publicaccess.Service,
	auditService audit.Service, gitspaceService *gitspace.Service,
	labelSvc *label.Service, instrumentation instrument.Service, executionStore store.ExecutionStore,
	rulesSvc *rules.Service, usageMetricStore store.UsageMetricStore, repoIdentifierCheck check.RepoIdentifier,
	infraProviderSvc *infraprovider2.Service, favoriteStore store.FavoriteStore, autolinkSvc *autolink.Service,
	spaceSvc *space.Service,
) *Controller {
	return NewController(config, tx, urlProvider,
		sseStreamer, identifierCheck, authorizer,
		spacePathStore, pipelineStore, secretStore,
		connectorStore, templateStore,
		spaceStore, repoStore, principalStore,
		repoCtrl, membershipStore, prListService,
		spaceFinder, repoFinder,
		importer, exporter, limiter, publicAccess,
		auditService, gitspaceService,
		labelSvc, instrumentation, executionStore,
		rulesSvc, usageMetricStore, repoIdentifierCheck,
		infraProviderSvc, favoriteStore, autolinkSvc, spaceSvc,
	)
}
