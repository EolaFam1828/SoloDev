// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package errortracker

import (
	"github.com/harness/gitness/app/auth/authz"
	errortrackerevents "github.com/harness/gitness/app/events/errortracker"
	"github.com/harness/gitness/app/services/errorbridge"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/store/database/dbtx"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
	tx dbtx.Transactor,
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	errorTrackerStore store.ErrorTrackerStore,
	principalInfoCache store.PrincipalInfoCache,
	eventReporter *errortrackerevents.Reporter,
	bridge *errorbridge.Bridge,
) *Controller {
	ctrl := NewController(tx, authorizer, spaceFinder, repoFinder, errorTrackerStore, principalInfoCache, eventReporter)
	ctrl.SetErrorBridge(bridge)
	return ctrl
}
