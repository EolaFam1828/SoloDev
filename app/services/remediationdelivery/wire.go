// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remediationdelivery

import (
	"github.com/harness/gitness/app/api/controller/pullreq"
	"github.com/harness/gitness/app/api/controller/repo"
	airemediationevents "github.com/harness/gitness/app/events/airemediation"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/app/url"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideService,
)

func ProvideService(
	remStore store.RemediationStore,
	repoStore store.RepoStore,
	principalStore store.PrincipalStore,
	repoCtrl *repo.Controller,
	pullreqCtrl *pullreq.Controller,
	urlProvider url.Provider,
	eventReporter *airemediationevents.Reporter,
) *Service {
	return NewService(
		remStore,
		repoStore,
		principalStore,
		repoCtrl,
		pullreqCtrl,
		urlProvider,
		eventReporter,
	)
}
