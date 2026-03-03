// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package agentorchestrator

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewService,
)
