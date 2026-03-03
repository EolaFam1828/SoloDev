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

package airemediation

import (
	"github.com/harness/gitness/app/auth/authz"
	airemediationevents "github.com/harness/gitness/app/events/airemediation"
	"github.com/harness/gitness/app/services/aiworker"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/services/remediationdelivery"
	"github.com/harness/gitness/app/services/remediationvalidation"
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
	validationService *remediationvalidation.Service,
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
		validationService,
		securityRemediation,
		aiWorker != nil && aiWorker.Available(),
	)
}
