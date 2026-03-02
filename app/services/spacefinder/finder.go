// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package spacefinder

import (
	"context"

	"github.com/EolaFam1828/SoloDev/types"
)

// Finder is an interface for finding spaces by reference.
// This wraps the refcache.SpaceFinder to return the full *types.Space
// (as needed by the feature flag module).
type Finder interface {
	FindByRef(ctx context.Context, spaceRef string) (*types.Space, error)
}
