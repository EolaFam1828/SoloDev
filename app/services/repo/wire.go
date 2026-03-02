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

package repo

import (
	"context"

	repoevents "github.com/EolaFam1828/SoloDev/app/events/repo"
	"github.com/EolaFam1828/SoloDev/app/services/locker"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/events"
	"github.com/EolaFam1828/SoloDev/git"
	"github.com/EolaFam1828/SoloDev/job"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideCalculator,
	ProvideService,
)

func ProvideCalculator(
	config *types.Config,
	git git.Interface,
	repoStore store.RepoStore,
	spaceStore store.SpaceStore,
	scheduler *job.Scheduler,
	executor *job.Executor,
	lfsStore store.LFSObjectStore,
	usageMetricStore store.UsageMetricStore,
) (*SizeCalculator, error) {
	job := &SizeCalculator{
		enabled:          config.RepoSize.Enabled,
		cron:             config.RepoSize.CRON,
		maxDur:           config.RepoSize.MaxDuration,
		numWorkers:       config.RepoSize.NumWorkers,
		git:              git,
		repoStore:        repoStore,
		spaceStore:       spaceStore,
		scheduler:        scheduler,
		lfsStore:         lfsStore,
		usageMetricStore: usageMetricStore,
	}

	err := executor.Register(jobType, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func ProvideService(ctx context.Context,
	config *types.Config,
	repoEvReporter *repoevents.Reporter,
	repoReaderFactory *events.ReaderFactory[*repoevents.Reader],
	repoStore store.RepoStore,
	urlProvider url.Provider,
	git git.Interface,
	locker *locker.Locker,
) (*Service, error) {
	return NewService(ctx, config, repoEvReporter, repoReaderFactory,
		repoStore, urlProvider, git, locker)
}
