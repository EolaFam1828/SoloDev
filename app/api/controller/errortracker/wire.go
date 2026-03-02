// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package errortracker

import (
	errortrackerevents "github.com/EolaFam1828/SoloDev/app/events/errortracker"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"

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
) *Controller {
	return NewController(
		tx,
		authorizer,
		spaceFinder,
		repoFinder,
		errorTrackerStore,
		principalInfoCache,
		eventReporter,
	)
}
