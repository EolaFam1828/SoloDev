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

package errortracker

import (
	"context"

	"github.com/EolaFam1828/SoloDev/app/events"
)

// Reader provides read access to error tracker events.
type Reader struct {
	innerReader events.Reader
}

// NewReader returns a new Reader.
func NewReader(innerReader events.Reader) *Reader {
	return &Reader{
		innerReader: innerReader,
	}
}

// Publish publishes an error tracker event.
func (r *Reader) Publish(ctx context.Context, event interface{}) error {
	return r.innerReader.Publish(ctx, event)
}
