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

package qualitygate

import (
	"github.com/harness/gitness/events"
)

// Reader is the event reader for this package.
type Reader struct {
	innerReader *events.GenericReader
}

// NewReader returns a new event reader.
func NewReader(eventsSystem *events.System) (*Reader, error) {
	innerReader, err := events.NewReader(eventsSystem, category)
	if err != nil {
		return nil, err
	}

	return &Reader{
		innerReader: innerReader,
	}, nil
}
