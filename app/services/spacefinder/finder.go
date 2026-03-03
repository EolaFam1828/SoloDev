// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spacefinder

import (
	"context"

	"github.com/harness/gitness/types"
)

// Finder is an interface for finding spaces by reference.
// This wraps the refcache.SpaceFinder to return the full *types.Space
// (as needed by the feature flag module).
type Finder interface {
	FindByRef(ctx context.Context, spaceRef string) (*types.Space, error)
}
