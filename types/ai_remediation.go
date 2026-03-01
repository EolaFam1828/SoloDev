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

package types

import "encoding/json"

// RemediationStatus represents the lifecycle status of an AI remediation task.
type RemediationStatus string

const (
	RemediationStatusPending    RemediationStatus = "pending"
	RemediationStatusProcessing RemediationStatus = "processing"
	RemediationStatusCompleted  RemediationStatus = "completed"
	RemediationStatusFailed     RemediationStatus = "failed"
	RemediationStatusApplied    RemediationStatus = "applied"
	RemediationStatusDismissed  RemediationStatus = "dismissed"
)

// RemediationTriggerSource identifies what spawned the remediation.
type RemediationTriggerSource string

const (
	RemediationTriggerPipeline RemediationTriggerSource = "pipeline"
	RemediationTriggerError    RemediationTriggerSource = "error_tracker"
	RemediationTriggerSecurity RemediationTriggerSource = "security_scan"
	RemediationTriggerQuality  RemediationTriggerSource = "quality_gate"
	RemediationTriggerManual   RemediationTriggerSource = "manual"
)

// Remediation represents a single AI-driven code fix task.
type Remediation struct {
	ID          int64             `json:"id"`
	SpaceID     int64             `json:"-"`
	RepoID      int64             `json:"repo_id"`
	Identifier  string            `json:"identifier"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Status      RemediationStatus `json:"status"`

	// Source context — what triggered this remediation.
	TriggerSource RemediationTriggerSource `json:"trigger_source"`
	TriggerRef    string                   `json:"trigger_ref,omitempty"` // e.g. execution ID, error group ID
	Branch        string                   `json:"branch"`
	CommitSHA     string                   `json:"commit_sha,omitempty"`

	// Failure context sent to the LLM.
	ErrorLog   string `json:"error_log"`
	SourceCode string `json:"source_code,omitempty"` // relevant snippet
	FilePath   string `json:"file_path,omitempty"`

	// AI output.
	AIModel      string          `json:"ai_model,omitempty"`
	AIPrompt     string          `json:"ai_prompt,omitempty"`
	AIResponse   string          `json:"ai_response,omitempty"`
	PatchDiff    string          `json:"patch_diff,omitempty"`
	FixBranch    string          `json:"fix_branch,omitempty"`
	PRLink       string          `json:"pr_link,omitempty"`
	Confidence   float64         `json:"confidence,omitempty"`  // 0.0 – 1.0
	TokensUsed   int64           `json:"tokens_used,omitempty"`
	DurationMs   int64           `json:"duration_ms,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`

	CreatedBy int64 `json:"-"`
	Created   int64 `json:"created"`
	Updated   int64 `json:"updated"`
	Version   int64 `json:"version"`
}

// RemediationListFilter holds query parameters for listing remediations.
type RemediationListFilter struct {
	ListQueryFilter
	Status        *RemediationStatus        `json:"status,omitempty"`
	TriggerSource *RemediationTriggerSource `json:"trigger_source,omitempty"`
}

// TriggerRemediationInput is the request body for triggering a new AI remediation.
type TriggerRemediationInput struct {
	Title         string                   `json:"title" binding:"required"`
	Description   string                   `json:"description,omitempty"`
	TriggerSource RemediationTriggerSource `json:"trigger_source" binding:"required"`
	TriggerRef    string                   `json:"trigger_ref,omitempty"`
	Branch        string                   `json:"branch" binding:"required"`
	CommitSHA     string                   `json:"commit_sha,omitempty"`
	ErrorLog      string                   `json:"error_log" binding:"required"`
	SourceCode    string                   `json:"source_code,omitempty"`
	FilePath      string                   `json:"file_path,omitempty"`
	AIModel       string                   `json:"ai_model,omitempty"`
	Metadata      json.RawMessage          `json:"metadata,omitempty"`
}

// UpdateRemediationInput is the request body for updating a remediation status.
type UpdateRemediationInput struct {
	Status     *RemediationStatus `json:"status,omitempty"`
	PatchDiff  string             `json:"patch_diff,omitempty"`
	AIResponse string             `json:"ai_response,omitempty"`
	FixBranch  string             `json:"fix_branch,omitempty"`
	PRLink     string             `json:"pr_link,omitempty"`
	Confidence *float64           `json:"confidence,omitempty"`
}

// RemediationSummary summarises remediation activity for a space.
type RemediationSummary struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	Processing int64 `json:"processing"`
	Completed  int64 `json:"completed"`
	Applied    int64 `json:"applied"`
	Failed     int64 `json:"failed"`
	Dismissed  int64 `json:"dismissed"`
}
