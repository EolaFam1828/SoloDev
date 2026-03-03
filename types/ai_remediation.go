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

import (
	"encoding/json"
	"fmt"
)

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

type RemediationDeliveryMode string

const (
	RemediationDeliveryModeManual RemediationDeliveryMode = "manual"
	RemediationDeliveryModeAutoPR RemediationDeliveryMode = "auto_pr"
)

type RemediationDeliveryState string

const (
	RemediationDeliveryStateNotAttempted RemediationDeliveryState = "not_attempted"
	RemediationDeliveryStateBranchReady  RemediationDeliveryState = "branch_ready"
	RemediationDeliveryStateApplied      RemediationDeliveryState = "applied"
	RemediationDeliveryStateFailed       RemediationDeliveryState = "failed"
)

type RemediationDelivery struct {
	Mode        RemediationDeliveryMode  `json:"mode"`
	State       RemediationDeliveryState `json:"state"`
	PRNumber    int64                    `json:"pr_number"`
	LastError   string                   `json:"last_error,omitempty"`
	AttemptedAt int64                    `json:"attempted_at"`
}

func DefaultRemediationDelivery(mode RemediationDeliveryMode) RemediationDelivery {
	if mode == "" {
		mode = RemediationDeliveryModeManual
	}

	return RemediationDelivery{
		Mode:  mode,
		State: RemediationDeliveryStateNotAttempted,
	}
}

func GetRemediationDeliveryMetadata(
	raw json.RawMessage,
	defaultMode RemediationDeliveryMode,
) (RemediationDelivery, error) {
	delivery := DefaultRemediationDelivery(defaultMode)
	if len(raw) == 0 {
		return delivery, nil
	}

	var metadata map[string]json.RawMessage
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return delivery, fmt.Errorf("decode remediation metadata: %w", err)
	}

	rawDelivery, ok := metadata["delivery"]
	if !ok || len(rawDelivery) == 0 {
		return delivery, nil
	}

	if err := json.Unmarshal(rawDelivery, &delivery); err != nil {
		return DefaultRemediationDelivery(defaultMode), fmt.Errorf("decode remediation delivery metadata: %w", err)
	}
	if delivery.Mode == "" {
		delivery.Mode = defaultMode
		if delivery.Mode == "" {
			delivery.Mode = RemediationDeliveryModeManual
		}
	}
	if delivery.State == "" {
		delivery.State = RemediationDeliveryStateNotAttempted
	}

	return delivery, nil
}

func SetRemediationDeliveryMetadata(
	raw json.RawMessage,
	delivery RemediationDelivery,
) (json.RawMessage, error) {
	metadata := map[string]json.RawMessage{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &metadata); err != nil {
			return nil, fmt.Errorf("decode remediation metadata: %w", err)
		}
	}

	encodedDelivery, err := json.Marshal(delivery)
	if err != nil {
		return nil, fmt.Errorf("encode remediation delivery metadata: %w", err)
	}
	metadata["delivery"] = encodedDelivery

	encodedMetadata, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("encode remediation metadata: %w", err)
	}

	return encodedMetadata, nil
}

// RemediationValidationState represents the validation lifecycle.
type RemediationValidationState string

const (
	RemediationValidationNotAttempted RemediationValidationState = "not_attempted"
	RemediationValidationQueued       RemediationValidationState = "queued"
	RemediationValidationRunning      RemediationValidationState = "running"
	RemediationValidationPassed       RemediationValidationState = "passed"
	RemediationValidationFailed       RemediationValidationState = "failed"
	RemediationValidationUnavailable  RemediationValidationState = "unavailable"
)

// RemediationValidation tracks the validation state for a remediation.
type RemediationValidation struct {
	State              RemediationValidationState `json:"state"`
	PipelineIdentifier string                     `json:"pipeline_identifier,omitempty"`
	ExecutionNumber    int64                      `json:"execution_number,omitempty"`
	ExecutionStatus    string                     `json:"execution_status,omitempty"`
	URL                string                     `json:"url,omitempty"`
	LastError          string                     `json:"last_error,omitempty"`
	StartedAt          int64                      `json:"started_at,omitempty"`
	UpdatedAt          int64                      `json:"updated_at,omitempty"`
	CompletedAt        int64                      `json:"completed_at,omitempty"`
}

// DefaultRemediationValidation returns a not_attempted validation.
func DefaultRemediationValidation() RemediationValidation {
	return RemediationValidation{
		State: RemediationValidationNotAttempted,
	}
}

// GetRemediationValidationMetadata extracts validation from metadata JSON.
func GetRemediationValidationMetadata(raw json.RawMessage) (RemediationValidation, error) {
	v := DefaultRemediationValidation()
	if len(raw) == 0 {
		return v, nil
	}

	var metadata map[string]json.RawMessage
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return v, fmt.Errorf("decode remediation metadata: %w", err)
	}

	rawVal, ok := metadata["validation"]
	if !ok || len(rawVal) == 0 {
		return v, nil
	}

	if err := json.Unmarshal(rawVal, &v); err != nil {
		return DefaultRemediationValidation(), fmt.Errorf("decode remediation validation metadata: %w", err)
	}
	if v.State == "" {
		v.State = RemediationValidationNotAttempted
	}

	return v, nil
}

// SetRemediationValidationMetadata persists validation into metadata JSON.
func SetRemediationValidationMetadata(raw json.RawMessage, validation RemediationValidation) (json.RawMessage, error) {
	metadata := map[string]json.RawMessage{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &metadata); err != nil {
			return nil, fmt.Errorf("decode remediation metadata: %w", err)
		}
	}

	encoded, err := json.Marshal(validation)
	if err != nil {
		return nil, fmt.Errorf("encode remediation validation metadata: %w", err)
	}
	metadata["validation"] = encoded

	encodedMetadata, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("encode remediation metadata: %w", err)
	}

	return encodedMetadata, nil
}

// IsValidationTerminal returns true if the validation state is terminal.
func (v RemediationValidation) IsTerminal() bool {
	switch v.State {
	case RemediationValidationPassed, RemediationValidationFailed, RemediationValidationUnavailable:
		return true
	default:
		return false
	}
}

// RemediationContextFragment represents a single context piece with provenance.
type RemediationContextFragment struct {
	Label        string `json:"label"`
	Source       string `json:"source"`
	FilePath     string `json:"file_path,omitempty"`
	TrimmedBytes int64  `json:"trimmed_bytes,omitempty"`
	CharCount    int    `json:"char_count"`
}

// RemediationContext is the provenance summary of what context was sent to the AI.
type RemediationContext struct {
	Fragments     []RemediationContextFragment `json:"fragments"`
	TotalCharsEst int                          `json:"total_chars_est"`
	TriggerSource string                       `json:"trigger_source"`
}

// GetRemediationContextMetadata extracts context provenance from metadata JSON.
func GetRemediationContextMetadata(raw json.RawMessage) *RemediationContext {
	if len(raw) == 0 {
		return nil
	}
	var metadata map[string]json.RawMessage
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return nil
	}
	rawCtx, ok := metadata["context_bundle"]
	if !ok || len(rawCtx) == 0 {
		return nil
	}

	// Parse the full bundle to extract provenance summary.
	var bundle struct {
		Fragments []struct {
			Label        string `json:"label"`
			Content      string `json:"content"`
			Source       string `json:"source"`
			FilePath     string `json:"file_path"`
			TrimmedBytes int64  `json:"trimmed_bytes"`
		} `json:"fragments"`
		TotalCharsEst int    `json:"total_chars_est"`
		TriggerSource string `json:"trigger_source"`
	}
	if err := json.Unmarshal(rawCtx, &bundle); err != nil {
		return nil
	}

	ctx := &RemediationContext{
		TotalCharsEst: bundle.TotalCharsEst,
		TriggerSource: bundle.TriggerSource,
	}
	for _, f := range bundle.Fragments {
		ctx.Fragments = append(ctx.Fragments, RemediationContextFragment{
			Label:        f.Label,
			Source:       f.Source,
			FilePath:     f.FilePath,
			TrimmedBytes: f.TrimmedBytes,
			CharCount:    len(f.Content),
		})
	}
	return ctx
}

// PopulateDelivery extracts delivery state from Metadata and sets the top-level Delivery field.
func (r *Remediation) PopulateDelivery() {
	d, _ := GetRemediationDeliveryMetadata(r.Metadata, RemediationDeliveryModeManual)
	r.Delivery = &d
}

// PopulateValidation extracts validation state from Metadata and sets the top-level Validation field.
func (r *Remediation) PopulateValidation() {
	v, _ := GetRemediationValidationMetadata(r.Metadata)
	r.Validation = &v
}

// PopulateContext extracts context provenance from Metadata and sets the top-level Context field.
func (r *Remediation) PopulateContext() {
	r.Context = GetRemediationContextMetadata(r.Metadata)
}

// PopulateAPIFields extracts delivery, validation, and context from metadata.
func (r *Remediation) PopulateAPIFields() {
	r.PopulateDelivery()
	r.PopulateValidation()
	r.PopulateContext()
}

// PopulateAPIFieldsSlice calls PopulateAPIFields on each remediation in the slice.
func PopulateAPIFieldsSlice(rems []*Remediation) {
	for _, r := range rems {
		r.PopulateAPIFields()
	}
}

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
	AIModel    string          `json:"ai_model,omitempty"`
	AIPrompt   string          `json:"ai_prompt,omitempty"`
	AIResponse string          `json:"ai_response,omitempty"`
	PatchDiff  string          `json:"patch_diff,omitempty"`
	FixBranch  string          `json:"fix_branch,omitempty"`
	PRLink     string          `json:"pr_link,omitempty"`
	Confidence float64         `json:"confidence,omitempty"` // 0.0 – 1.0
	TokensUsed int64           `json:"tokens_used,omitempty"`
	DurationMs int64           `json:"duration_ms,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`

	// Delivery is a top-level projection of metadata.delivery for API consumers.
	// Not stored as a DB column — populated before returning API responses.
	Delivery *RemediationDelivery `json:"delivery,omitempty"`

	// Validation is a top-level projection of metadata.validation for API consumers.
	// Not stored as a DB column — populated before returning API responses.
	Validation *RemediationValidation `json:"validation,omitempty"`

	// Context is a provenance summary of what context fragments were sent to the AI.
	// Not stored as a DB column — populated before returning API responses.
	Context *RemediationContext `json:"context,omitempty"`

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

type CreateRemediationFromSecurityFindingInput struct {
	RepoRef        string `json:"repo_ref"`
	ScanIdentifier string `json:"scan_identifier"`
	FindingID      int64  `json:"finding_id"`
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

// RemediationMetrics contains time-windowed remediation metrics.
type RemediationMetrics struct {
	WindowDays        int              `json:"window_days"`
	Total             int64            `json:"total"`
	Completed         int64            `json:"completed"`
	Applied           int64            `json:"applied"`
	Failed            int64            `json:"failed"`
	AvgConfidence     float64          `json:"avg_confidence"`
	AvgDurationMs     int64            `json:"avg_duration_ms"`
	ValidationsPassed int64            `json:"validations_passed"`
	ValidationsFailed int64            `json:"validations_failed"`
	ByTrigger         map[string]int64 `json:"by_trigger"`
	SuccessRate       float64          `json:"success_rate"`
	MeanTimeToFixMs   int64            `json:"mean_time_to_fix_ms"`
}
