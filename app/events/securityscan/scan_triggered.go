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

	"github.com/harness/gitness/events"
	"github.com/rs/zerolog/log"
)

const TriggeredEvent events.EventType = "triggered"

type TriggeredPayload struct {
	ScanID      int64  `json:"scan_id"`
	ScanType    string `json:"scan_type"`
	RepoID      int64  `json:"repo_id"`
	SpaceID     int64  `json:"space_id"`
	CommitSHA   string `json:"commit_sha"`
	Branch      string `json:"branch"`
	TriggeredBy string `json:"triggered_by"`
}

func (r *Reporter) Triggered(ctx context.Context, payload *TriggeredPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, TriggeredEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send security scan triggered event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported security scan triggered event with id '%s'", eventID)
}

func (r *Reader) RegisterTriggered(fn events.HandlerFunc[*TriggeredPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, TriggeredEvent, fn, opts...)
}
