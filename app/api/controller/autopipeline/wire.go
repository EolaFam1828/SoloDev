// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package autopipeline

import (
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/refcache"

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
) *Controller {
	return NewController(authorizer, spaceFinder, repoFinder)
}
