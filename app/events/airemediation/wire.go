// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package airemediation

import (
	appevents "github.com/harness/gitness/app/events"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideReader,
	NewReporter,
)

func ProvideReader() *Reader {
	return NewReader(appevents.NoOpReader{})
}
