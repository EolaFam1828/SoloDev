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

	"github.com/EolaFam1828/SoloDev/events"

	"github.com/rs/zerolog/log"
)

const UpdatedEvent events.EventType = "updated"

type UpdatedPayload struct {
	TechDebtID int64  `json:"tech_debt_id"`
	SpaceID    int64  `json:"space_id"`
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	Status     string `json:"status"`
}

func (r *Reporter) Updated(ctx context.Context, payload *UpdatedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, UpdatedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send tech debt updated event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported tech debt updated event with id '%s'", eventID)
}

func (r *Reader) RegisterUpdated(fn events.HandlerFunc[*UpdatedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, UpdatedEvent, fn, opts...)
}
