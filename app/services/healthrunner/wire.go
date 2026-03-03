// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package healthrunner

import (
	"github.com/harness/gitness/app/services/healthbridge"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"

	"github.com/google/wire"
)

// WireSet provides a wire set for the health runner service.
var WireSet = wire.NewSet(
	ProvideService,
)

// ProvideService creates the health runner service.
func ProvideService(
	scheduler *job.Scheduler,
	executor *job.Executor,
	hcStore store.HealthCheckStore,
	hcResultStore store.HealthCheckResultStore,
	bridge *healthbridge.Bridge,
) *Service {
	return NewService(scheduler, executor, hcStore, hcResultStore, bridge)
}
