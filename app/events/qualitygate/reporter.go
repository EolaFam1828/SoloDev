// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Simplified to log-based event reporter for MCP integration.
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

	"github.com/EolaFam1828/SoloDev/types"

	"github.com/rs/zerolog/log"
)

// Reporter is the event reporter for this package.
type Reporter struct{}

func NewReporter() *Reporter {
	return &Reporter{}
}

// RuleCreated reports that a quality rule was created.
func (r *Reporter) RuleCreated(ctx context.Context, rule *types.QualityRule) {
	if r == nil {
		return
	}
	log.Ctx(ctx).Info().
		Str("event", "rule.created").
		Str("identifier", rule.Identifier).
		Msg("quality rule created")
}

// RuleUpdated reports that a quality rule was updated.
func (r *Reporter) RuleUpdated(ctx context.Context, rule *types.QualityRule) {
	if r == nil {
		return
	}
	log.Ctx(ctx).Info().
		Str("event", "rule.updated").
		Str("identifier", rule.Identifier).
		Msg("quality rule updated")
}

// RuleDeleted reports that a quality rule was deleted.
func (r *Reporter) RuleDeleted(ctx context.Context, rule *types.QualityRule) {
	if r == nil {
		return
	}
	log.Ctx(ctx).Info().
		Str("event", "rule.deleted").
		Str("identifier", rule.Identifier).
		Msg("quality rule deleted")
}

// RuleEnabled reports that a quality rule was enabled.
func (r *Reporter) RuleEnabled(ctx context.Context, rule *types.QualityRule) {
	if r == nil {
		return
	}
	log.Ctx(ctx).Info().
		Str("event", "rule.enabled").
		Str("identifier", rule.Identifier).
		Msg("quality rule enabled")
}

// RuleDisabled reports that a quality rule was disabled.
func (r *Reporter) RuleDisabled(ctx context.Context, rule *types.QualityRule) {
	if r == nil {
		return
	}
	log.Ctx(ctx).Info().
		Str("event", "rule.disabled").
		Str("identifier", rule.Identifier).
		Msg("quality rule disabled")
}

// EvaluationCreated reports that a quality evaluation was created.
func (r *Reporter) EvaluationCreated(ctx context.Context, eval *types.QualityEvaluation) {
	if r == nil {
		return
	}
	log.Ctx(ctx).Info().
		Str("event", "evaluation.created").
		Str("identifier", eval.Identifier).
		Msg("quality evaluation created")
}
