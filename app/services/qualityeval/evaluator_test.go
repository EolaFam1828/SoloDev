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

package qualityeval

import (
	"testing"

	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

func TestEvaluateRule(t *testing.T) {
	e := NewEvaluator()

	tests := []struct {
		name       string
		rule       *types.QualityRule
		metrics    map[string]interface{}
		wantStatus enum.QualityEvaluationStatus
	}{
		{
			name: "threshold_pass",
			rule: &types.QualityRule{
				ID: 1, Identifier: "cov", Name: "Coverage",
				Condition: "coverage >= 80", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{"coverage": 90.0},
			wantStatus: enum.QualityEvaluationStatusPassed,
		},
		{
			name: "threshold_fail",
			rule: &types.QualityRule{
				ID: 2, Identifier: "cov", Name: "Coverage",
				Condition: "coverage >= 80", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{"coverage": 75.0},
			wantStatus: enum.QualityEvaluationStatusFailed,
		},
		{
			name: "threshold_fail_warn_enforcement",
			rule: &types.QualityRule{
				ID: 3, Identifier: "cov", Name: "Coverage",
				Condition: "coverage >= 80", Enforcement: enum.QualityEnforcementWarn, Enabled: true,
			},
			metrics:    map[string]interface{}{"coverage": 75.0},
			wantStatus: enum.QualityEvaluationStatusWarning,
		},
		{
			name: "compound_condition_pass",
			rule: &types.QualityRule{
				ID: 4, Identifier: "comp", Name: "Compound",
				Condition: "coverage >= 80 && bugs == 0", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{"coverage": 85.0, "bugs": 0},
			wantStatus: enum.QualityEvaluationStatusPassed,
		},
		{
			name: "missing_metric_skipped",
			rule: &types.QualityRule{
				ID: 6, Identifier: "cov", Name: "Coverage",
				Condition: "coverage >= 80", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{"bugs": 0},
			wantStatus: enum.QualityEvaluationStatusPassed,
		},
		{
			name: "empty_condition_block",
			rule: &types.QualityRule{
				ID: 8, Identifier: "empty", Name: "Empty Block",
				Condition: "", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{},
			wantStatus: enum.QualityEvaluationStatusFailed,
		},
		{
			name: "empty_condition_warn",
			rule: &types.QualityRule{
				ID: 9, Identifier: "empty", Name: "Empty Warn",
				Condition: "", Enforcement: enum.QualityEnforcementWarn, Enabled: true,
			},
			metrics:    map[string]interface{}{},
			wantStatus: enum.QualityEvaluationStatusWarning,
		},
		{
			name: "empty_condition_info",
			rule: &types.QualityRule{
				ID: 10, Identifier: "empty", Name: "Empty Info",
				Condition: "", Enforcement: enum.QualityEnforcementInfo, Enabled: true,
			},
			metrics:    map[string]interface{}{},
			wantStatus: enum.QualityEvaluationStatusPassed,
		},
		{
			name: "nil_metrics",
			rule: &types.QualityRule{
				ID: 11, Identifier: "cov", Name: "Coverage",
				Condition: "coverage >= 80", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    nil,
			wantStatus: enum.QualityEvaluationStatusPassed,
		},
		{
			name: "integer_metric_pass",
			rule: &types.QualityRule{
				ID: 12, Identifier: "bugs", Name: "Zero Bugs",
				Condition: "bugs == 0", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{"bugs": 0},
			wantStatus: enum.QualityEvaluationStatusPassed,
		},
		{
			name: "integer_metric_fail",
			rule: &types.QualityRule{
				ID: 13, Identifier: "bugs", Name: "Zero Bugs",
				Condition: "bugs == 0", Enforcement: enum.QualityEnforcementBlock, Enabled: true,
			},
			metrics:    map[string]interface{}{"bugs": 5},
			wantStatus: enum.QualityEvaluationStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := e.EvaluateRule(tt.rule, tt.metrics)
			if result.Status != tt.wantStatus {
				t.Errorf("got status %q, want %q (message: %s)", result.Status, tt.wantStatus, result.Message)
			}
		})
	}
}

func TestEvaluateRules_DisabledSkipped(t *testing.T) {
	e := NewEvaluator()

	rules := []*types.QualityRule{
		{ID: 1, Identifier: "a", Name: "Enabled", Condition: "x > 0", Enforcement: enum.QualityEnforcementBlock, Enabled: true},
		{ID: 2, Identifier: "b", Name: "Disabled", Condition: "x > 0", Enforcement: enum.QualityEnforcementBlock, Enabled: false},
	}

	metrics := map[string]interface{}{"x": 1}
	results, passed, failed, warned, skipped := e.EvaluateRules(rules, metrics)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if passed != 1 || failed != 0 || warned != 0 || skipped != 1 {
		t.Errorf("counts: passed=%d failed=%d warned=%d skipped=%d", passed, failed, warned, skipped)
	}
}
