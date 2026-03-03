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
	"github.com/harness/gitness/events"
)

// Reader is the event reader for this package.
type Reader struct {
	innerReader *events.GenericReader
}

func (r *Reader) Configure(opts ...events.ReaderOption) {
	r.innerReader.Configure(opts...)
}

func NewReaderFactory(eventsSystem *events.System) (*events.ReaderFactory[*Reader], error) {
	return events.NewReaderFactory(eventsSystem, category, func(innerReader *events.GenericReader) (*Reader, error) {
		return &Reader{innerReader: innerReader}, nil
	})
}

func (r *Reader) RegisterHealthCheckCreated(fn events.HandlerFunc[*HealthCheckCreatedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, HealthCheckCreatedEvent, fn, opts...)
}

func (r *Reader) RegisterHealthCheckUpdated(fn events.HandlerFunc[*HealthCheckUpdatedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, HealthCheckUpdatedEvent, fn, opts...)
}

func (r *Reader) RegisterHealthCheckDeleted(fn events.HandlerFunc[*HealthCheckDeletedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, HealthCheckDeletedEvent, fn, opts...)
}

func (r *Reader) RegisterHealthCheckStatusChanged(fn events.HandlerFunc[*HealthCheckStatusChangedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, HealthCheckStatusChangedEvent, fn, opts...)
}

func (r *Reader) RegisterHealthCheckResultCreated(fn events.HandlerFunc[*HealthCheckResultCreatedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, HealthCheckResultCreatedEvent, fn, opts...)
}
