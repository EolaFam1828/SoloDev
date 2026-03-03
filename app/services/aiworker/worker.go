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

	"github.com/harness/gitness/app/services/remediationdelivery"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/git"
	"github.com/harness/gitness/job"
	"github.com/harness/gitness/types"

	"github.com/rs/zerolog/log"
)

const (
	jobTypeRemWorker = "ai-remediation-worker"
	jobTypeRemPoller = "ai-remediation-poller"
	jobCronRemPoller = "*/15 * * * * *" // every 15 seconds
	jobMaxDurPoller  = 30 * time.Second
)

// remWorkerHandler processes a single remediation by calling the LLM.
type remWorkerHandler struct {
	remStore        store.RemediationStore
	repoStore       store.RepoStore
	git             git.Interface
	provider        LLMProvider
	config          Config
	deliveryService *remediationdelivery.Service
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
	if err := h.enrichSourceCode(ctx, rem); err != nil {
		log.Ctx(ctx).Warn().Err(err).Int64("remediation_id", rem.ID).Msg("failed to enrich remediation source code")
	}

	// Build prompts and call LLM.
	userPrompt := BuildUserPrompt(rem)
	rem.AIPrompt = userPrompt
	rem.Updated = time.Now().UnixMilli()
	if err := h.remStore.Update(ctx, rem); err != nil {
		return "", fmt.Errorf("failed to persist remediation prompt: %w", err)
	}

	llmReq := &LLMRequest{
		SystemPrompt: GetSystemPrompt(),
		UserPrompt:   userPrompt,
		MaxTokens:    h.config.MaxTokens,
		Temperature:  h.config.Temperature,
	}

	started := time.Now()
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
	rem.TokensUsed = int64(llmResp.TokensUsed)
	rem.DurationMs = time.Since(started).Milliseconds()
	rem.Updated = time.Now().UnixMilli()

	deliveryMode := types.RemediationDeliveryModeManual
	if h.config.CreateFixBranch {
		deliveryMode = types.RemediationDeliveryModeAutoPR
	}

	if parsed.Diff != "" {
		rem.Status = types.RemediationStatusCompleted
		if err := setDeliveryMetadata(
			rem,
			deliveryMode,
			types.RemediationDeliveryStateNotAttempted,
			"",
			0,
			0,
		); err != nil {
			return "", fmt.Errorf("failed to set remediation delivery metadata: %w", err)
		}
	} else {
		rem.Status = types.RemediationStatusFailed
		if rem.AIResponse == "" {
			rem.AIResponse = "AI returned empty response"
		}
	}

	if err := h.remStore.Update(ctx, rem); err != nil {
		return "", fmt.Errorf("failed to update remediation with result: %w", err)
	}

	if parsed.Diff == "" {
		return fmt.Sprintf("completed: confidence=%.2f diff_len=%d", parsed.Confidence, len(parsed.Diff)), nil
	}

	if h.config.CreateFixBranch {
		if h.deliveryService == nil {
			msg := "auto delivery requested but remediation delivery service is unavailable"
			if err := setDeliveryMetadata(
				rem,
				types.RemediationDeliveryModeAutoPR,
				types.RemediationDeliveryStateFailed,
				msg,
				0,
				types.NowMillis(),
			); err == nil {
				_ = h.remStore.Update(ctx, rem)
			}
			return fmt.Sprintf("completed: confidence=%.2f diff_len=%d auto_apply_failed=%s", parsed.Confidence, len(parsed.Diff), msg), nil
		}

		updatedRem, err := h.deliveryService.ApplyAsRemediationCreator(ctx, rem, types.RemediationDeliveryModeAutoPR)
		if err != nil {
			return fmt.Sprintf("completed: confidence=%.2f diff_len=%d auto_apply_failed=%v", parsed.Confidence, len(parsed.Diff), err), nil
		}
		rem = updatedRem
	}

	return fmt.Sprintf("completed: confidence=%.2f diff_len=%d status=%s", parsed.Confidence, len(parsed.Diff), rem.Status), nil
}
