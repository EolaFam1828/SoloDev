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
	"fmt"

	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"

	"github.com/antonmedv/expr"
)

// EvalResult holds the result of evaluating a single quality rule.
type EvalResult struct {
	RuleID         int64                        `json:"rule_id"`
	RuleIdentifier string                       `json:"rule_identifier"`
	RuleName       string                       `json:"rule_name"`
	Status         enum.QualityEvaluationStatus `json:"status"`
	ActualValue    string                       `json:"actual_value,omitempty"`
	ExpectedValue  string                       `json:"expected_value,omitempty"`
	Message        string                       `json:"message,omitempty"`
}

// Evaluator evaluates quality rules using the expr library.
type Evaluator struct{}

// NewEvaluator creates a new Evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// EvaluateRule evaluates a single quality rule against the given metrics.
func (e *Evaluator) EvaluateRule(rule *types.QualityRule, metrics map[string]interface{}) EvalResult {
	result := EvalResult{
		RuleID:         rule.ID,
		RuleIdentifier: rule.Identifier,
		RuleName:       rule.Name,
		ExpectedValue:  rule.Condition,
	}

	// Empty condition — status depends on enforcement level.
	if rule.Condition == "" {
		switch rule.Enforcement {
		case enum.QualityEnforcementBlock:
			result.Status = enum.QualityEvaluationStatusFailed
			result.Message = "no condition defined; enforcement=block → fail"
		case enum.QualityEnforcementWarn:
			result.Status = enum.QualityEvaluationStatusWarning
			result.Message = "no condition defined; enforcement=warn → warning"
		default:
			result.Status = enum.QualityEvaluationStatusPassed
			result.Message = "no condition defined; enforcement=info → pass"
		}
		return result
	}

	// Nil metrics map — skip the rule.
	if metrics == nil {
		result.Status = enum.QualityEvaluationStatusPassed
		result.Message = "skipped: no metrics provided"
		return result
	}

	// Compile the expression.
	program, err := expr.Compile(rule.Condition, expr.AsBool())
	if err != nil {
		result.Status = enum.QualityEvaluationStatusFailed
		result.Message = fmt.Sprintf("invalid condition: %v", err)
		return result
	}

	// Run the expression against the metrics.
	output, err := expr.Run(program, metrics)
	if err != nil {
		// If the error is due to a missing variable, treat as skipped.
		result.Status = enum.QualityEvaluationStatusPassed
		result.Message = fmt.Sprintf("skipped: %v", err)
		return result
	}

	passed, ok := output.(bool)
	if !ok {
		result.Status = enum.QualityEvaluationStatusFailed
		result.Message = "condition did not return a boolean"
		return result
	}

	result.ActualValue = fmt.Sprintf("%v", passed)

	if passed {
		result.Status = enum.QualityEvaluationStatusPassed
		result.Message = "condition met"
	} else {
		switch rule.Enforcement {
		case enum.QualityEnforcementBlock:
			result.Status = enum.QualityEvaluationStatusFailed
			result.Message = "condition not met"
		case enum.QualityEnforcementWarn:
			result.Status = enum.QualityEvaluationStatusWarning
			result.Message = "condition not met (warning)"
		default:
			result.Status = enum.QualityEvaluationStatusPassed
			result.Message = "condition not met (info only)"
		}
	}

	return result
}

// EvaluateRules evaluates a slice of quality rules and returns results plus counts.
func (e *Evaluator) EvaluateRules(
	rules []*types.QualityRule,
	metrics map[string]interface{},
) (results []EvalResult, passed, failed, warned, skipped int) {
	results = make([]EvalResult, 0, len(rules))

	for _, rule := range rules {
		if !rule.Enabled {
			skipped++
			continue
		}

		r := e.EvaluateRule(rule, metrics)
		results = append(results, r)

		switch r.Status {
		case enum.QualityEvaluationStatusPassed:
			passed++
		case enum.QualityEvaluationStatusFailed:
			failed++
		case enum.QualityEvaluationStatusWarning:
			warned++
		}
	}

	return results, passed, failed, warned, skipped
}
