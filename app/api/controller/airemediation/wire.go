// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package airemediation

import (
	"github.com/harness/gitness/app/auth/authz"
	airemediationevents "github.com/harness/gitness/app/events/airemediation"
	"github.com/harness/gitness/app/services/aiworker"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/services/remediationdelivery"
	"github.com/harness/gitness/app/services/securityremediation"
	"github.com/harness/gitness/app/store"

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
	deliveryService *remediationdelivery.Service,
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
		deliveryService,
		securityRemediation,
		aiWorker != nil && aiWorker.Available(),
	)
}
