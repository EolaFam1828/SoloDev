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
