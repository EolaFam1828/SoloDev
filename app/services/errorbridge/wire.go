// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package errorbridge

import (
	"github.com/EolaFam1828/SoloDev/app/store"

	"github.com/google/wire"
)

// WireSet provides a wire set for the error bridge.
var WireSet = wire.NewSet(
	ProvideBridge,
)

// ProvideBridge wires the default error-to-remediation bridge.
func ProvideBridge(remediationStore store.RemediationStore) *Bridge {
	return NewBridge(remediationStore, true)
}
