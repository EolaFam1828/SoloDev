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

const CompletedEvent events.EventType = "completed"

type CompletedPayload struct {
	ScanID        int64 `json:"scan_id"`
	RepoID        int64 `json:"repo_id"`
	SpaceID       int64 `json:"space_id"`
	TotalIssues   int   `json:"total_issues"`
	CriticalCount int   `json:"critical_count"`
	HighCount     int   `json:"high_count"`
	MediumCount   int   `json:"medium_count"`
	LowCount      int   `json:"low_count"`
	Duration      int64 `json:"duration"`
}

func (r *Reporter) Completed(ctx context.Context, payload *CompletedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, CompletedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send security scan completed event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported security scan completed event with id '%s'", eventID)
}

func (r *Reader) RegisterCompleted(fn events.HandlerFunc[*CompletedPayload],
	opts ...events.HandlerOption) error {
	return events.ReaderRegisterEvent(r.innerReader, CompletedEvent, fn, opts...)
}
