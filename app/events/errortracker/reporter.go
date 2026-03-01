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

	"github.com/harness/gitness/app/events"
	"github.com/harness/gitness/types"
)

// Reporter reports error tracker events.
type Reporter struct {
	reader *Reader
}

// NewReporter returns a new Reporter.
func NewReporter(reader *Reader) *Reporter {
	return &Reporter{
		reader: reader,
	}
}

// ErrorReported reports that an error was reported.
func (r *Reporter) ErrorReported(ctx context.Context, errorGroup *types.ErrorGroup) {
	if r.reader == nil {
		return
	}

	event := ErrorReported{
		Base: Base{
			ErrorGroupID: errorGroup.ID,
			Identifier:   errorGroup.Identifier,
			SpaceID:      errorGroup.SpaceID,
		},
		ErrorGroup: errorGroup,
	}

	_ = r.reader.Publish(ctx, &event)
}

// ErrorStatusChanged reports that an error group status changed.
func (r *Reporter) ErrorStatusChanged(
	ctx context.Context,
	errorGroup *types.ErrorGroup,
	newStatus types.ErrorGroupStatus,
) {
	if r.reader == nil {
		return
	}

	event := ErrorStatusChanged{
		Base: Base{
			ErrorGroupID: errorGroup.ID,
			Identifier:   errorGroup.Identifier,
			SpaceID:      errorGroup.SpaceID,
		},
		OldStatus: errorGroup.Status,
		NewStatus: newStatus,
	}

	_ = r.reader.Publish(ctx, &event)
}

// ErrorAssigned reports that an error group was assigned.
func (r *Reporter) ErrorAssigned(
	ctx context.Context,
	errorGroup *types.ErrorGroup,
	assignedTo *int64,
) {
	if r.reader == nil {
		return
	}

	event := ErrorAssigned{
		Base: Base{
			ErrorGroupID: errorGroup.ID,
			Identifier:   errorGroup.Identifier,
			SpaceID:      errorGroup.SpaceID,
		},
		AssignedTo: assignedTo,
	}

	_ = r.reader.Publish(ctx, &event)
}
