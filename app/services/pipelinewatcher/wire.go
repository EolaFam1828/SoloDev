// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pipelinewatcher

import (
	"github.com/harness/gitness/app/services/errorbridge"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideService,
)

// ProvideService wires the pipeline watcher service.
func ProvideService(
	scheduler *job.Scheduler,
	executor *job.Executor,
	executionStore store.ExecutionStore,
	repoStore store.RepoStore,
	remStore store.RemediationStore,
	errorBridge *errorbridge.Bridge,
) *Service {
	return NewService(scheduler, executor, executionStore, repoStore, remStore, errorBridge)
}
