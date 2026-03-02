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

package migrate

import (
	"github.com/EolaFam1828/SoloDev/app/api/controller/limiter"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	repoevents "github.com/EolaFam1828/SoloDev/app/events/repo"
	"github.com/EolaFam1828/SoloDev/app/services/migrate"
	"github.com/EolaFam1828/SoloDev/app/services/publicaccess"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/audit"
	"github.com/EolaFam1828/SoloDev/git"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types/check"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
	authorizer authz.Authorizer,
	publicAccess publicaccess.Service,
	rpcClient git.Interface,
	urlProvider url.Provider,
	pullreqImporter *migrate.PullReq,
	ruleImporter *migrate.Rule,
	webhookImporter *migrate.Webhook,
	labelImporter *migrate.Label,
	resourceLimiter limiter.ResourceLimiter,
	auditService audit.Service,
	identifierCheck check.RepoIdentifier,
	tx dbtx.Transactor,
	spaceStore store.SpaceStore,
	repoStore store.RepoStore,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	eventReporter *repoevents.Reporter,
) *Controller {
	return NewController(
		authorizer,
		publicAccess,
		rpcClient,
		urlProvider,
		pullreqImporter,
		ruleImporter,
		webhookImporter,
		labelImporter,
		resourceLimiter,
		auditService,
		identifierCheck,
		tx,
		spaceStore,
		repoStore,
		spaceFinder,
		repoFinder,
		eventReporter,
	)
}
