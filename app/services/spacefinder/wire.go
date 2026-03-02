// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package spacefinder

import (
	"github.com/harness/gitness/app/store"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideFinder,
)

func ProvideFinder(spaceStore store.SpaceStore) Finder {
	return spaceStore
}
