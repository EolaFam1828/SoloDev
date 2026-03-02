// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package airemediation

import (
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	airemediationevents "github.com/EolaFam1828/SoloDev/app/events/airemediation"
	"github.com/EolaFam1828/SoloDev/app/services/aiworker"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/services/securityremediation"
	"github.com/EolaFam1828/SoloDev/app/store"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	remediationStore store.RemediationStore,
	scanResultStore store.SecurityScanStore,
	scanFindingStore store.ScanFindingStore,
	eventReporter *airemediationevents.Reporter,
	securityRemediation *securityremediation.Service,
	aiWorker *aiworker.Service,
) *Controller {
	return NewController(
		authorizer,
		spaceFinder,
		repoFinder,
		remediationStore,
		scanResultStore,
		scanFindingStore,
		eventReporter,
		securityRemediation,
		aiWorker != nil && aiWorker.Available(),
	)
}
