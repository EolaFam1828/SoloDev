// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
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

package airemediation

import (
	"context"

	"github.com/harness/gitness/types"
)

// Reporter reports AI remediation events.
type Reporter struct {
	reader *Reader
}

// NewReporter returns a new Reporter.
func NewReporter(reader *Reader) *Reporter {
	return &Reporter{
		reader: reader,
	}
}

// RemediationTriggered reports that a new remediation was triggered.
func (r *Reporter) RemediationTriggered(ctx context.Context, rem *types.Remediation) {
	if r.reader == nil {
		return
	}

	event := RemediationTriggered{
		Base: Base{
			RemediationID: rem.ID,
			Identifier:    rem.Identifier,
			SpaceID:       rem.SpaceID,
		},
		Remediation: rem,
	}

	_ = r.reader.Publish(ctx, &event)
}

// RemediationCompleted reports that a remediation completed.
func (r *Reporter) RemediationCompleted(
	ctx context.Context,
	rem *types.Remediation,
	status types.RemediationStatus,
) {
	if r.reader == nil {
		return
	}

	event := RemediationCompleted{
		Base: Base{
			RemediationID: rem.ID,
			Identifier:    rem.Identifier,
			SpaceID:       rem.SpaceID,
		},
		Status: status,
	}

	_ = r.reader.Publish(ctx, &event)
}

// RemediationApplied reports that a remediation fix was applied.
func (r *Reporter) RemediationApplied(
	ctx context.Context,
	rem *types.Remediation,
) {
	if r.reader == nil {
		return
	}

	event := RemediationApplied{
		Base: Base{
			RemediationID: rem.ID,
			Identifier:    rem.Identifier,
			SpaceID:       rem.SpaceID,
		},
		FixBranch: rem.FixBranch,
		PRLink:    rem.PRLink,
	}

	_ = r.reader.Publish(ctx, &event)
}
