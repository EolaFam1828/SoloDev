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

package gitspace

import (
	gitspaceevents "github.com/EolaFam1828/SoloDev/app/events/gitspace"
	gitspacedeleteevents "github.com/EolaFam1828/SoloDev/app/events/gitspacedelete"
	"github.com/EolaFam1828/SoloDev/app/gitspace/orchestrator"
	"github.com/EolaFam1828/SoloDev/app/gitspace/orchestrator/ide"
	"github.com/EolaFam1828/SoloDev/app/gitspace/scm"
	"github.com/EolaFam1828/SoloDev/app/services/infraprovider"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/services/tokengenerator"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideGitspace,
)

func ProvideGitspace(
	tx dbtx.Transactor,
	gitspaceStore store.GitspaceConfigStore,
	gitspaceInstanceStore store.GitspaceInstanceStore,
	eventReporter *gitspaceevents.Reporter,
	gitspaceEventStore store.GitspaceEventStore,
	spaceFinder refcache.SpaceFinder,
	infraProviderSvc *infraprovider.Service,
	orchestrator orchestrator.Orchestrator,
	scm *scm.SCM,
	config *types.Config,
	gitspaceDeleteEventReporter *gitspacedeleteevents.Reporter,
	ideFactory ide.Factory,
	spaceStore store.SpaceStore,
	tokenGenerator tokengenerator.TokenGenerator,
) *Service {
	return NewService(tx, gitspaceStore, gitspaceInstanceStore, eventReporter,
		gitspaceEventStore, spaceFinder, infraProviderSvc, orchestrator, scm, config,
		gitspaceDeleteEventReporter, ideFactory, spaceStore, tokenGenerator,
	)
}
