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

import (
	"errors"

	"github.com/harness/gitness/events"
)

// Reader is the event reader for this package.
type Reader struct {
	innerReader *events.GenericReader
}

func NewReader(eventsSystem *events.System) (*Reader, error) {
	innerReader, err := events.NewReader(eventsSystem, category)
	if err != nil {
		return nil, errors.New("failed to create new GenericReader from event system")
	}

	return &Reader{
		innerReader: innerReader,
	}, nil
}

func (r *Reader) RegisterHealthCheckCreated(fn func(event *EventHealthCheckCreated) error) error {
	return r.innerReader.Register(EventHealthCheckCreated, func(data interface{}) error {
		event, ok := data.(*EventHealthCheckCreated)
		if !ok {
			return errors.New("unexpected event type")
		}
		return fn(event)
	})
}

func (r *Reader) RegisterHealthCheckUpdated(fn func(event *EventHealthCheckUpdated) error) error {
	return r.innerReader.Register(EventHealthCheckUpdated, func(data interface{}) error {
		event, ok := data.(*EventHealthCheckUpdated)
		if !ok {
			return errors.New("unexpected event type")
		}
		return fn(event)
	})
}

func (r *Reader) RegisterHealthCheckDeleted(fn func(event *EventHealthCheckDeleted) error) error {
	return r.innerReader.Register(EventHealthCheckDeleted, func(data interface{}) error {
		event, ok := data.(*EventHealthCheckDeleted)
		if !ok {
			return errors.New("unexpected event type")
		}
		return fn(event)
	})
}

func (r *Reader) RegisterHealthCheckStatusChanged(fn func(event *EventHealthCheckStatusChanged) error) error {
	return r.innerReader.Register(EventHealthCheckStatusChanged, func(data interface{}) error {
		event, ok := data.(*EventHealthCheckStatusChanged)
		if !ok {
			return errors.New("unexpected event type")
		}
		return fn(event)
	})
}

func (r *Reader) RegisterHealthCheckResultCreated(fn func(event *EventHealthCheckResultCreated) error) error {
	return r.innerReader.Register(EventHealthCheckResultCreated, func(data interface{}) error {
		event, ok := data.(*EventHealthCheckResultCreated)
		if !ok {
			return errors.New("unexpected event type")
		}
		return fn(event)
	})
}
