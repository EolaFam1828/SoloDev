// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package healthbridge

import (
	"github.com/harness/gitness/app/services/aiworker"
	"github.com/harness/gitness/app/store"

	"github.com/google/wire"
)

// WireSet provides a wire set for the health bridge.
var WireSet = wire.NewSet(
	ProvideBridge,
)

// ProvideBridge wires the health-to-remediation bridge.
func ProvideBridge(remediationStore store.RemediationStore, aiWorker *aiworker.Service) *Bridge {
	return NewBridge(remediationStore, aiWorker != nil && aiWorker.Available())
}
