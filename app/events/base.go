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
