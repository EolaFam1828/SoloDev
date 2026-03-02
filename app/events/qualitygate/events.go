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

import "github.com/EolaFam1828/SoloDev/types"

// RuleCreatedEvent is the event triggered when a quality rule is created.
type RuleCreatedEvent struct {
	Rule *types.QualityRule `json:"rule"`
}

// RuleUpdatedEvent is the event triggered when a quality rule is updated.
type RuleUpdatedEvent struct {
	Rule *types.QualityRule `json:"rule"`
}

// RuleDeletedEvent is the event triggered when a quality rule is deleted.
type RuleDeletedEvent struct {
	Rule *types.QualityRule `json:"rule"`
}

// RuleEnabledEvent is the event triggered when a quality rule is enabled.
type RuleEnabledEvent struct {
	Rule *types.QualityRule `json:"rule"`
}

// RuleDisabledEvent is the event triggered when a quality rule is disabled.
type RuleDisabledEvent struct {
	Rule *types.QualityRule `json:"rule"`
}

// EvaluationCreatedEvent is the event triggered when a quality evaluation is created.
type EvaluationCreatedEvent struct {
	Evaluation *types.QualityEvaluation `json:"evaluation"`
}
