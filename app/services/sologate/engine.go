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

// Package sologate implements the solopreneur-mode quality and security gate engine.
// It sits between pipeline execution and the quality/security modules, deciding
// whether to block, warn, or auto-remediate based on the user's configured mode.
package sologate

import (
	"context"
	"fmt"
	"log"

	"github.com/harness/gitness/app/services/errorbridge"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
)

// Engine evaluates quality/security gate results through the solo-mode lens.
type Engine struct {
	remediationStore store.RemediationStore
	bridge           *errorbridge.Bridge
}

// NewEngine creates a new gate engine.
func NewEngine(
	remediationStore store.RemediationStore,
	bridge *errorbridge.Bridge,
) *Engine {
	return &Engine{
		remediationStore: remediationStore,
		bridge:           bridge,
	}
}

// EvaluateResult represents the outcome of gate evaluation.
type EvaluateResult struct {
	Action         string   `json:"action"` // "block", "warn", "pass"
	Reasons        []string `json:"reasons"`
	AutoRemediate  bool     `json:"auto_remediate"`
	TechDebtLogged bool     `json:"tech_debt_logged"`
}

// Evaluate takes gate findings and applies the solo-mode enforcement rules.
func (e *Engine) Evaluate(
	ctx context.Context,
	config *types.SoloGateConfig,
	findings []Finding,
) *EvaluateResult {
	result := &EvaluateResult{
		Action: "pass",
	}

	if config == nil {
		// No config = strict mode by default
		config = &types.SoloGateConfig{
			EnforcementMode: types.EnforcementModeStrict,
		}
	}

	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, f := range findings {
		switch f.Severity {
		case "critical":
			criticalCount++
		case "high":
			highCount++
		case "medium":
			mediumCount++
		case "low", "info":
			lowCount++
		}
	}

	switch config.EnforcementMode {
	case types.EnforcementModeStrict:
		if criticalCount+highCount+mediumCount+lowCount > 0 {
			result.Action = "block"
			result.Reasons = append(result.Reasons, fmt.Sprintf(
				"%d critical, %d high, %d medium, %d low findings",
				criticalCount, highCount, mediumCount, lowCount,
			))
		}

	case types.EnforcementModeBalanced:
		if criticalCount+highCount > 0 {
			result.Action = "block"
			result.Reasons = append(result.Reasons, fmt.Sprintf(
				"%d critical, %d high findings (blocking)", criticalCount, highCount,
			))
		} else if mediumCount > 0 {
			result.Action = "warn"
			result.Reasons = append(result.Reasons, fmt.Sprintf(
				"%d medium findings (warning only)", mediumCount,
			))
		}
		// Low findings are auto-ignored in balanced mode

	case types.EnforcementModePrototype:
		// Never block in prototype mode
		if criticalCount+highCount+mediumCount > 0 {
			result.Action = "warn"
			result.Reasons = append(result.Reasons, fmt.Sprintf(
				"prototype mode: %d findings logged as tech debt",
				criticalCount+highCount+mediumCount+lowCount,
			))
			result.TechDebtLogged = true
		}
	}

	// Auto-ignore low severity if configured
	if config.AutoIgnoreLow && lowCount > 0 {
		log.Printf("[sologate] auto-ignoring %d low severity findings", lowCount)
	}

	// Auto-trigger AI remediation if configured and there are actionable findings
	if config.AIAutoFix && (criticalCount+highCount) > 0 {
		result.AutoRemediate = true
	}

	return result
}

// Finding is a simplified gate finding for evaluation.
type Finding struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FilePath    string `json:"file_path,omitempty"`
}
