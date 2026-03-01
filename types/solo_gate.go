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

package types

// EnforcementMode controls how quality/security gates behave.
type EnforcementMode string

const (
	// EnforcementModeStrict blocks the pipeline on any gate failure.
	EnforcementModeStrict EnforcementMode = "strict"

	// EnforcementModePrototype allows everything through but logs tech debt for later.
	EnforcementModePrototype EnforcementMode = "prototype"

	// EnforcementModeBalanced blocks only critical/high, warns on medium/low.
	EnforcementModeBalanced EnforcementMode = "balanced"
)

// SoloGateConfig stores the solopreneur-mode gate configuration for a space.
type SoloGateConfig struct {
	ID              int64           `json:"id"`
	SpaceID         int64           `json:"-"`
	EnforcementMode EnforcementMode `json:"enforcement_mode"`
	AutoIgnoreLow   bool            `json:"auto_ignore_low"`   // auto-dismiss low severity findings
	AutoTriageKnown bool            `json:"auto_triage_known"` // auto-dismiss repeat false positives
	AIAutoFix       bool            `json:"ai_auto_fix"`       // automatically trigger AI remediation on failures
	LogTechDebt     bool            `json:"log_tech_debt"`     // log skipped issues as tech debt items
	Created         int64           `json:"created"`
	Updated         int64           `json:"updated"`
}

// UpdateSoloGateConfigInput is the request body for updating gate configuration.
type UpdateSoloGateConfigInput struct {
	EnforcementMode *EnforcementMode `json:"enforcement_mode,omitempty"`
	AutoIgnoreLow   *bool            `json:"auto_ignore_low,omitempty"`
	AutoTriageKnown *bool            `json:"auto_triage_known,omitempty"`
	AIAutoFix       *bool            `json:"ai_auto_fix,omitempty"`
	LogTechDebt     *bool            `json:"log_tech_debt,omitempty"`
}
