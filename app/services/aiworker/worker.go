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

package aiworker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"
	"github.com/harness/gitness/types"
)

const (
	jobTypeRemWorker = "ai-remediation-worker"
	jobTypeRemPoller = "ai-remediation-poller"
	jobCronRemPoller = "*/15 * * * * *" // every 15 seconds
	jobMaxDurPoller  = 30 * time.Second
)

// remWorkerHandler processes a single remediation by calling the LLM.
type remWorkerHandler struct {
	remStore store.RemediationStore
	provider LLMProvider
	config   Config
}

type remJobInput struct {
	RemediationID int64 `json:"remediation_id"`
}

func (h *remWorkerHandler) Handle(ctx context.Context, data string, _ job.ProgressReporter) (string, error) {
	if h.provider == nil {
		return "", fmt.Errorf("no LLM provider configured")
	}

	var input remJobInput
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("failed to parse remediation job input: %w", err)
	}

	rem, err := h.remStore.Find(ctx, input.RemediationID)
	if err != nil {
		return "", fmt.Errorf("failed to find remediation %d: %w", input.RemediationID, err)
	}

	// Skip if not pending (already picked up by another worker).
	if rem.Status != types.RemediationStatusPending {
		return "skipped: not pending", nil
	}

	// Mark as processing.
	rem.Status = types.RemediationStatusProcessing
	rem.AIModel = fmt.Sprintf("%s/%s", h.provider.Name(), h.config.Model)
	rem.Updated = time.Now().UnixMilli()
	if err := h.remStore.Update(ctx, rem); err != nil {
		return "", fmt.Errorf("failed to mark remediation as processing: %w", err)
	}

	// Build prompts and call LLM.
	llmReq := &LLMRequest{
		SystemPrompt: GetSystemPrompt(),
		UserPrompt:   BuildUserPrompt(rem),
		MaxTokens:    h.config.MaxTokens,
		Temperature:  h.config.Temperature,
	}

	llmResp, err := h.provider.Complete(ctx, llmReq)
	if err != nil {
		rem.Status = types.RemediationStatusFailed
		rem.AIResponse = fmt.Sprintf("LLM call failed: %v", err)
		rem.Updated = time.Now().UnixMilli()
		_ = h.remStore.Update(ctx, rem)
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	// Parse the AI response.
	parsed := ParseAIResponse(llmResp.Content)

	rem.AIResponse = llmResp.Content
	rem.PatchDiff = parsed.Diff
	rem.Confidence = parsed.Confidence
	rem.Updated = time.Now().UnixMilli()

	if parsed.Diff != "" {
		rem.Status = types.RemediationStatusCompleted
	} else {
		rem.Status = types.RemediationStatusFailed
		if rem.AIResponse == "" {
			rem.AIResponse = "AI returned empty response"
		}
	}

	if err := h.remStore.Update(ctx, rem); err != nil {
		return "", fmt.Errorf("failed to update remediation with result: %w", err)
	}

	return fmt.Sprintf("completed: confidence=%.2f diff_len=%d", parsed.Confidence, len(parsed.Diff)), nil
}
