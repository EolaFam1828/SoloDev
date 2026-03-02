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

package pullreq

import (
	"context"

	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	gitevents "github.com/EolaFam1828/SoloDev/app/events/git"
	pullreqevents "github.com/EolaFam1828/SoloDev/app/events/pullreq"
	"github.com/EolaFam1828/SoloDev/app/services/codecomments"
	"github.com/EolaFam1828/SoloDev/app/services/label"
	"github.com/EolaFam1828/SoloDev/app/services/protection"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/sse"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/events"
	"github.com/EolaFam1828/SoloDev/git"
	"github.com/EolaFam1828/SoloDev/pubsub"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideService,
	ProvideListService,
)

func ProvideService(ctx context.Context,
	config *types.Config,
	gitReaderFactory *events.ReaderFactory[*gitevents.Reader],
	pullReqEvFactory *events.ReaderFactory[*pullreqevents.Reader],
	pullReqEvReporter *pullreqevents.Reporter,
	git git.Interface,
	repoFinder refcache.RepoFinder,
	repoStore store.RepoStore,
	pullreqStore store.PullReqStore,
	activityStore store.PullReqActivityStore,
	principalInfoCache store.PrincipalInfoCache,
	codeCommentView store.CodeCommentView,
	codeCommentMigrator *codecomments.Migrator,
	fileViewStore store.PullReqFileViewStore,
	pubsub pubsub.PubSub,
	urlProvider url.Provider,
	sseStreamer sse.Streamer,
) (*Service, error) {
	return New(ctx,
		config,
		gitReaderFactory,
		pullReqEvFactory,
		pullReqEvReporter,
		git,
		repoFinder,
		repoStore,
		pullreqStore,
		activityStore,
		codeCommentView,
		codeCommentMigrator,
		fileViewStore,
		principalInfoCache,
		pubsub,
		urlProvider,
		sseStreamer,
	)
}

func ProvideListService(
	tx dbtx.Transactor,
	git git.Interface,
	authorizer authz.Authorizer,
	spaceStore store.SpaceStore,
	pullreqStore store.PullReqStore,
	checkStore store.CheckStore,
	repoFinder refcache.RepoFinder,
	labelSvc *label.Service,
	protectionManager *protection.Manager,
) *ListService {
	return NewListService(
		tx,
		git,
		authorizer,
		spaceStore,
		pullreqStore,
		checkStore,
		repoFinder,
		labelSvc,
		protectionManager,
	)
}
