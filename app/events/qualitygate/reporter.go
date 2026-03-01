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

package qualitygate

import (
	"context"
	"errors"

	"github.com/harness/gitness/events"
	"github.com/harness/gitness/types"
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

// RuleCreated reports that a quality rule was created.
func (r *Reporter) RuleCreated(ctx context.Context, rule *types.QualityRule) {
	if r == nil || r.innerReporter == nil {
		return
	}
	_ = r.innerReporter.Report(ctx, "rule.created", &RuleCreatedEvent{Rule: rule})
}

// RuleUpdated reports that a quality rule was updated.
func (r *Reporter) RuleUpdated(ctx context.Context, rule *types.QualityRule) {
	if r == nil || r.innerReporter == nil {
		return
	}
	_ = r.innerReporter.Report(ctx, "rule.updated", &RuleUpdatedEvent{Rule: rule})
}

// RuleDeleted reports that a quality rule was deleted.
func (r *Reporter) RuleDeleted(ctx context.Context, rule *types.QualityRule) {
	if r == nil || r.innerReporter == nil {
		return
	}
	_ = r.innerReporter.Report(ctx, "rule.deleted", &RuleDeletedEvent{Rule: rule})
}

// RuleEnabled reports that a quality rule was enabled.
func (r *Reporter) RuleEnabled(ctx context.Context, rule *types.QualityRule) {
	if r == nil || r.innerReporter == nil {
		return
	}
	_ = r.innerReporter.Report(ctx, "rule.enabled", &RuleEnabledEvent{Rule: rule})
}

// RuleDisabled reports that a quality rule was disabled.
func (r *Reporter) RuleDisabled(ctx context.Context, rule *types.QualityRule) {
	if r == nil || r.innerReporter == nil {
		return
	}
	_ = r.innerReporter.Report(ctx, "rule.disabled", &RuleDisabledEvent{Rule: rule})
}

// EvaluationCreated reports that a quality evaluation was created.
func (r *Reporter) EvaluationCreated(ctx context.Context, eval *types.QualityEvaluation) {
	if r == nil || r.innerReporter == nil {
		return
	}
	_ = r.innerReporter.Report(ctx, "evaluation.created", &EvaluationCreatedEvent{Evaluation: eval})
}
