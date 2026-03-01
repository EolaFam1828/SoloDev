// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package events

import "context"

// Reader provides a simple interface for publishing events.
// This is the base interface used by SoloDev event sub-packages
// (airemediation, errortracker, etc.).
type Reader interface {
	Publish(ctx context.Context, event interface{}) error
}

// NoOpReader is a Reader that does nothing.
// Used when no event system is configured.
type NoOpReader struct{}

// Publish is a no-op.
func (NoOpReader) Publish(_ context.Context, _ interface{}) error {
	return nil
}
