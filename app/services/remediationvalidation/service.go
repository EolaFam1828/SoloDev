// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remediationvalidation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	controllerexecution "github.com/harness/gitness/app/api/controller/execution"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/app/url"
	"github.com/harness/gitness/job"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

const (
	jobTypeValidationPoller = "remediation-validation-poller"
	pollInterval            = 15 * time.Second
	maxPollDuration         = 30 * time.Minute
)

// Service manages remediation validation lifecycle.
type Service struct {
	remStore       store.RemediationStore
	repoStore      store.RepoStore
	pipelineStore  store.PipelineStore
	executionStore store.ExecutionStore
	executionCtrl  *controllerexecution.Controller
	principalStore store.PrincipalStore
	urlProvider    url.Provider
	scheduler      *job.Scheduler
	executor       *job.Executor
}

// NewService creates a new validation service.
func NewService(
	remStore store.RemediationStore,
	repoStore store.RepoStore,
	pipelineStore store.PipelineStore,
	executionStore store.ExecutionStore,
	executionCtrl *controllerexecution.Controller,
	principalStore store.PrincipalStore,
	urlProvider url.Provider,
	scheduler *job.Scheduler,
	executor *job.Executor,
) *Service {
	return &Service{
		remStore:       remStore,
		repoStore:      repoStore,
		pipelineStore:  pipelineStore,
		executionStore: executionStore,
		executionCtrl:  executionCtrl,
		principalStore: principalStore,
		urlProvider:    urlProvider,
		scheduler:      scheduler,
		executor:       executor,
	}
}

// Register registers the validation poller job handler.
func (s *Service) Register(_ context.Context) error {
	return s.executor.Register(jobTypeValidationPoller, &validationPollerHandler{
		remStore:       s.remStore,
		executionStore: s.executionStore,
		pipelineStore:  s.pipelineStore,
	})
}

// Validate triggers a pipeline execution on the remediation's fix branch and starts polling.
func (s *Service) Validate(
	ctx context.Context,
	rem *types.Remediation,
	pipelineIdentifier string,
) (*types.Remediation, error) {
	if rem.Status != types.RemediationStatusApplied {
		return rem, fmt.Errorf("only applied remediations can be validated")
	}
	if rem.FixBranch == "" {
		return rem, fmt.Errorf("remediation has no fix branch")
	}

	repo, err := s.repoStore.Find(ctx, rem.RepoID)
	if err != nil {
		return s.failValidation(ctx, rem, fmt.Sprintf("failed to find repo: %v", err))
	}

	// Resolve pipeline.
	resolvedPipeline, err := s.resolvePipeline(ctx, repo.ID, pipelineIdentifier)
	if err != nil {
		return s.unavailableValidation(ctx, rem, err.Error())
	}

	// Reconstruct session from remediation creator.
	principal, err := s.principalStore.Find(ctx, rem.CreatedBy)
	if err != nil {
		return s.failValidation(ctx, rem, fmt.Sprintf("failed to reconstruct session: %v", err))
	}
	session := &auth.Session{
		Principal: *principal,
		Metadata:  &auth.EmptyMetadata{},
	}

	// Trigger execution on fix branch.
	execution, err := s.executionCtrl.Create(ctx, session, repo.Path, resolvedPipeline.Identifier, rem.FixBranch)
	if err != nil {
		return s.failValidation(ctx, rem, fmt.Sprintf("failed to trigger pipeline: %v", err))
	}

	// Persist queued validation state.
	now := types.NowMillis()
	validation := types.RemediationValidation{
		State:              types.RemediationValidationQueued,
		PipelineIdentifier: resolvedPipeline.Identifier,
		ExecutionNumber:    execution.Number,
		ExecutionStatus:    string(execution.Status),
		StartedAt:          now,
		UpdatedAt:          now,
	}

	rem.Metadata, err = types.SetRemediationValidationMetadata(rem.Metadata, validation)
	if err != nil {
		return nil, fmt.Errorf("set validation metadata: %w", err)
	}
	if err := s.remStore.Update(ctx, rem); err != nil {
		return nil, fmt.Errorf("persist validation state: %w", err)
	}

	// Schedule a background poller job.
	data, _ := json.Marshal(validationPollerInput{
		RemediationID:   rem.ID,
		PipelineID:      resolvedPipeline.ID,
		ExecutionNumber: execution.Number,
	})
	_ = s.scheduler.RunJob(ctx, job.Definition{
		UID:     fmt.Sprintf("rem-validate-%s", rem.Identifier),
		Type:    jobTypeValidationPoller,
		Timeout: maxPollDuration,
		Data:    string(data),
	})

	return rem, nil
}

func (s *Service) resolvePipeline(ctx context.Context, repoID int64, identifier string) (*types.Pipeline, error) {
	if identifier != "" {
		p, err := s.pipelineStore.FindByIdentifier(ctx, repoID, identifier)
		if err != nil {
			return nil, fmt.Errorf("pipeline %q not found: %w", identifier, err)
		}
		if p.Disabled {
			return nil, fmt.Errorf("pipeline %q is disabled", identifier)
		}
		return p, nil
	}

	pipelines, err := s.pipelineStore.List(ctx, repoID, &types.ListPipelinesFilter{Latest: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipelines: %w", err)
	}

	var enabled []*types.Pipeline
	for _, p := range pipelines {
		if !p.Disabled {
			enabled = append(enabled, p)
		}
	}

	if len(enabled) == 0 {
		return nil, fmt.Errorf("no enabled pipelines in repo")
	}
	if len(enabled) == 1 {
		return enabled[0], nil
	}

	// Look for "default" pipeline.
	for _, p := range enabled {
		if p.Identifier == "default" {
			return p, nil
		}
	}

	return nil, fmt.Errorf("multiple pipelines exist; specify pipeline_identifier")
}

func (s *Service) failValidation(ctx context.Context, rem *types.Remediation, msg string) (*types.Remediation, error) {
	return s.setValidationState(ctx, rem, types.RemediationValidationFailed, msg)
}

func (s *Service) unavailableValidation(ctx context.Context, rem *types.Remediation, msg string) (*types.Remediation, error) {
	return s.setValidationState(ctx, rem, types.RemediationValidationUnavailable, msg)
}

func (s *Service) setValidationState(
	ctx context.Context,
	rem *types.Remediation,
	state types.RemediationValidationState,
	msg string,
) (*types.Remediation, error) {
	now := types.NowMillis()
	v := types.RemediationValidation{
		State:       state,
		LastError:   msg,
		UpdatedAt:   now,
		CompletedAt: now,
	}
	var err error
	rem.Metadata, err = types.SetRemediationValidationMetadata(rem.Metadata, v)
	if err != nil {
		return rem, fmt.Errorf("set validation metadata: %w", err)
	}
	if storeErr := s.remStore.Update(ctx, rem); storeErr != nil {
		return rem, fmt.Errorf("persist validation state: %w", storeErr)
	}
	return rem, fmt.Errorf("%s", msg)
}

// --- Background poller ---

type validationPollerInput struct {
	RemediationID   int64 `json:"remediation_id"`
	PipelineID      int64 `json:"pipeline_id"`
	ExecutionNumber int64 `json:"execution_number"`
}

type validationPollerHandler struct {
	remStore       store.RemediationStore
	executionStore store.ExecutionStore
	pipelineStore  store.PipelineStore
}

func (h *validationPollerHandler) Handle(ctx context.Context, data string, _ job.ProgressReporter) (string, error) {
	var input validationPollerInput
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("decode poller input: %w", err)
	}

	rem, err := h.remStore.Find(ctx, input.RemediationID)
	if err != nil {
		return "", fmt.Errorf("find remediation: %w", err)
	}

	deadline := time.Now().Add(maxPollDuration)
	for time.Now().Before(deadline) {
		execution, err := h.executionStore.FindByNumber(ctx, input.PipelineID, input.ExecutionNumber)
		if err != nil {
			return "", fmt.Errorf("find execution: %w", err)
		}

		state := mapCIStatusToValidation(execution.Status)
		now := types.NowMillis()

		validation, _ := types.GetRemediationValidationMetadata(rem.Metadata)
		validation.State = state
		validation.ExecutionStatus = string(execution.Status)
		validation.UpdatedAt = now

		if state == types.RemediationValidationPassed || state == types.RemediationValidationFailed {
			validation.CompletedAt = now
		}

		rem.Metadata, err = types.SetRemediationValidationMetadata(rem.Metadata, validation)
		if err != nil {
			return "", fmt.Errorf("set validation metadata: %w", err)
		}
		if err := h.remStore.Update(ctx, rem); err != nil {
			// Re-read on version conflict.
			rem, _ = h.remStore.Find(ctx, input.RemediationID)
			continue
		}

		if validation.IsTerminal() {
			return fmt.Sprintf("validation %s for %s", state, rem.Identifier), nil
		}

		select {
		case <-ctx.Done():
			return "context canceled", ctx.Err()
		case <-time.After(pollInterval):
		}

		// Re-read remediation for version freshness.
		rem, err = h.remStore.Find(ctx, input.RemediationID)
		if err != nil {
			return "", fmt.Errorf("re-read remediation: %w", err)
		}
	}

	return "validation poll timeout", nil
}

func mapCIStatusToValidation(status enum.CIStatus) types.RemediationValidationState {
	switch status {
	case enum.CIStatusPending:
		return types.RemediationValidationQueued
	case enum.CIStatusRunning:
		return types.RemediationValidationRunning
	case enum.CIStatusSuccess:
		return types.RemediationValidationPassed
	case enum.CIStatusFailure, enum.CIStatusKilled, enum.CIStatusError,
		enum.CIStatusSkipped, enum.CIStatusDeclined:
		return types.RemediationValidationFailed
	default:
		return types.RemediationValidationQueued
	}
}
