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

const FailedEvent events.EventType = "failed"

type FailedPayload struct {
	ScanID    int64  `json:"scan_id"`
	RepoID    int64  `json:"repo_id"`
	SpaceID   int64  `json:"space_id"`
	ErrorMsg  string `json:"error_msg"`
	Duration  int64  `json:"duration"`
}

func (r *Reporter) Failed(ctx context.Context, payload *FailedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, FailedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send security scan failed event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported security scan failed event with id '%s'", eventID)
}

func (r *Reader) RegisterFailed(fn events.HandlerFunc[*FailedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, FailedEvent, fn, opts...)
}
