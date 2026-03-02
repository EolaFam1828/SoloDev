// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package spacefinder

import (
	"github.com/EolaFam1828/SoloDev/app/store"

	"github.com/google/wire"
)

// WireSet provides a wire set for the Space finder adapter.
var WireSet = wire.NewSet(
	ProvideFinder,
)

// ProvideFinder adapts the full space store to the Finder interface expected
// by SoloDev feature-flag controllers.
func ProvideFinder(spaceStore store.SpaceStore) Finder {
	return spaceStore
}
