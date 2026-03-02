// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package airemediation

import (
	airemediationevents "github.com/EolaFam1828/SoloDev/app/events/airemediation"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
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
	remediationStore store.RemediationStore,
	eventReporter *airemediationevents.Reporter,
) *Controller {
	return NewController(authorizer, spaceFinder, remediationStore, eventReporter)
}
