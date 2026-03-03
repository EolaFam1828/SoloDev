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

package remediationdelivery

import (
	"github.com/harness/gitness/app/api/controller/pullreq"
	"github.com/harness/gitness/app/api/controller/repo"
	airemediationevents "github.com/harness/gitness/app/events/airemediation"
	"github.com/harness/gitness/app/services/remediationnotifier"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/app/url"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideService,
)

func ProvideService(
	remStore store.RemediationStore,
	repoStore store.RepoStore,
	principalStore store.PrincipalStore,
	repoCtrl *repo.Controller,
	pullreqCtrl *pullreq.Controller,
	urlProvider url.Provider,
	eventReporter *airemediationevents.Reporter,
	notifier *remediationnotifier.Service,
) *Service {
	return NewService(
		remStore,
		repoStore,
		principalStore,
		repoCtrl,
		pullreqCtrl,
		urlProvider,
		eventReporter,
		notifier,
	)
}
