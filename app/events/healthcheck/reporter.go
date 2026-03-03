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
	"context"
	"errors"

	"github.com/harness/gitness/events"

	"github.com/rs/zerolog/log"
)

// Reporter is the event reporter for this package.
type Reporter struct {
	innerReporter *events.GenericReporter
}

func NewReporter(eventsSystem *events.System) (*Reporter, error) {
	innerReporter, err := events.NewReporter(eventsSystem, category)
	if err != nil {
		return nil, errors.New("failed to create new GenericReporter from event system")
	}

	return &Reporter{
		innerReporter: innerReporter,
	}, nil
}

func (r *Reporter) HealthCheckCreated(ctx context.Context, payload *HealthCheckCreatedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, HealthCheckCreatedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to send health check created event")
		return
	}
	log.Ctx(ctx).Debug().Msgf("reported health check created event with id '%s'", eventID)
}

func (r *Reporter) HealthCheckUpdated(ctx context.Context, payload *HealthCheckUpdatedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, HealthCheckUpdatedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to send health check updated event")
		return
	}
	log.Ctx(ctx).Debug().Msgf("reported health check updated event with id '%s'", eventID)
}

func (r *Reporter) HealthCheckDeleted(ctx context.Context, payload *HealthCheckDeletedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, HealthCheckDeletedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to send health check deleted event")
		return
	}
	log.Ctx(ctx).Debug().Msgf("reported health check deleted event with id '%s'", eventID)
}

func (r *Reporter) HealthCheckStatusChanged(ctx context.Context, payload *HealthCheckStatusChangedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, HealthCheckStatusChangedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to send health check status changed event")
		return
	}
	log.Ctx(ctx).Debug().Msgf("reported health check status changed event with id '%s'", eventID)
}

func (r *Reporter) HealthCheckResultCreated(ctx context.Context, payload *HealthCheckResultCreatedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, HealthCheckResultCreatedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to send health check result created event")
		return
	}
	log.Ctx(ctx).Debug().Msgf("reported health check result created event with id '%s'", eventID)
}
