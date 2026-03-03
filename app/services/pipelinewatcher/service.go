// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package pipelinewatcher monitors pipeline executions for failures
// and triggers remediation via the error bridge, forming the entry
// point of the self-healing pipeline loop.
package pipelinewatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/harness/gitness/app/services/errorbridge"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"

	"github.com/rs/zerolog/log"
)

const (
	jobTypePipelineWatcher = "pipeline-failure-watcher"
	jobCronWatcher         = "0 * * * * *" // every 60 seconds
	jobMaxDurWatcher       = 45 * time.Second
	// lookbackWindow is how far back we scan for failures (2 minutes to overlap with cron interval).
	lookbackWindow = 2 * time.Minute
	// maxFailures is the maximum number of failures to process per poll cycle.
	maxFailures = 20
)

// Service watches for pipeline failures and creates remediations.
type Service struct {
	scheduler      *job.Scheduler
	executor       *job.Executor
	executionStore store.ExecutionStore
	repoStore      store.RepoStore
	remStore       store.RemediationStore
	errorBridge    *errorbridge.Bridge
}

// NewService creates a new pipeline watcher service.
func NewService(
	scheduler *job.Scheduler,
	executor *job.Executor,
	executionStore store.ExecutionStore,
	repoStore store.RepoStore,
	remStore store.RemediationStore,
	errorBridge *errorbridge.Bridge,
) *Service {
	if errorBridge == nil {
		return nil
	}
	return &Service{
		scheduler:      scheduler,
		executor:       executor,
		executionStore: executionStore,
		repoStore:      repoStore,
		remStore:       remStore,
		errorBridge:    errorBridge,
	}
}

// Register registers the watcher job handler and schedules recurring execution.
func (s *Service) Register(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if err := s.executor.Register(jobTypePipelineWatcher, &watcherHandler{
		executionStore: s.executionStore,
		repoStore:      s.repoStore,
		remStore:       s.remStore,
		errorBridge:    s.errorBridge,
	}); err != nil {
		return fmt.Errorf("register pipeline watcher handler: %w", err)
	}

	return s.scheduler.AddRecurring(ctx, jobTypePipelineWatcher, jobTypePipelineWatcher, jobCronWatcher, jobMaxDurWatcher)
}

// watcherHandler is the job handler that scans for recent pipeline failures.
type watcherHandler struct {
	executionStore store.ExecutionStore
	repoStore      store.RepoStore
	remStore       store.RemediationStore
	errorBridge    *errorbridge.Bridge
}

func (h *watcherHandler) Handle(ctx context.Context, _ string, _ job.ProgressReporter) (string, error) {
	sinceMs := time.Now().Add(-lookbackWindow).UnixMilli()

	failures, err := h.executionStore.ListRecentFailed(ctx, sinceMs, maxFailures)
	if err != nil {
		return "", fmt.Errorf("list recent failed executions: %w", err)
	}

	if len(failures) == 0 {
		return "no recent failures", nil
	}

	created := 0
	for _, exec := range failures {
		triggerRef := fmt.Sprintf("%d", exec.Number)

		// Check if a remediation already exists for this execution.
		existing, _ := h.remStore.FindActiveByTriggerRef(ctx, exec.RepoID, triggerRef)
		if existing != nil {
			continue
		}

		// Resolve the repo to get the space ID.
		repo, err := h.repoStore.Find(ctx, exec.RepoID)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Int64("repo_id", exec.RepoID).Msg("pipeline watcher: repo not found, skipping")
			continue
		}

		// Build error context from execution fields.
		errorLog := exec.Error
		if errorLog == "" {
			errorLog = fmt.Sprintf("Pipeline execution #%d failed with status: %s", exec.Number, exec.Status)
		}

		// Extract branch from ref (refs/heads/main → main).
		branch := extractBranch(exec.Ref)
		if branch == "" {
			branch = exec.Target
		}

		h.errorBridge.OnPipelineFailed(
			ctx,
			repo.ParentID, // space ID
			exec.RepoID,
			exec.Number,
			branch,
			exec.After, // commit SHA
			errorLog,
			exec.CreatedBy,
		)
		created++

		log.Ctx(ctx).Info().
			Int64("execution_number", exec.Number).
			Int64("repo_id", exec.RepoID).
			Str("branch", branch).
			Msg("pipeline watcher: auto-created remediation for pipeline failure")
	}

	result := map[string]int{
		"failures_found": len(failures),
		"remediations":   created,
	}
	out, _ := json.Marshal(result)
	return string(out), nil
}

func extractBranch(ref string) string {
	const prefix = "refs/heads/"
	if len(ref) > len(prefix) && ref[:len(prefix)] == prefix {
		return ref[len(prefix):]
	}
	return ref
}
