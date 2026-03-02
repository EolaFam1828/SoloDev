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

const ResolvedEvent events.EventType = "resolved"

type ResolvedPayload struct {
	TechDebtID int64  `json:"tech_debt_id"`
	SpaceID    int64  `json:"space_id"`
	Identifier string `json:"identifier"`
	ResolvedBy int64  `json:"resolved_by"`
	ResolvedAt int64  `json:"resolved_at"`
}

func (r *Reporter) Resolved(ctx context.Context, payload *ResolvedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, ResolvedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send tech debt resolved event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported tech debt resolved event with id '%s'", eventID)
}

func (r *Reader) RegisterResolved(fn events.HandlerFunc[*ResolvedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, ResolvedEvent, fn, opts...)
}
