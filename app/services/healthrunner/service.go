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

// Package healthrunner provides a background service that executes HTTP
// health checks on their configured schedule and records results.
package healthrunner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/harness/gitness/app/services/healthbridge"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"
	"github.com/harness/gitness/types"

	"github.com/rs/zerolog/log"
)

const (
	jobTypeHealthPoller = "health-check-poller"
	jobCronPoller       = "*/30 * * * * *" // every 30 seconds
	jobMaxDurPoller     = 60 * time.Second
)

// Service manages the health check runner background jobs.
type Service struct {
	scheduler     *job.Scheduler
	executor      *job.Executor
	hcStore       store.HealthCheckStore
	hcResultStore store.HealthCheckResultStore
	bridge        *healthbridge.Bridge
}

// NewService creates a new health runner service.
func NewService(
	scheduler *job.Scheduler,
	executor *job.Executor,
	hcStore store.HealthCheckStore,
	hcResultStore store.HealthCheckResultStore,
	bridge *healthbridge.Bridge,
) *Service {
	return &Service{
		scheduler:     scheduler,
		executor:      executor,
		hcStore:       hcStore,
		hcResultStore: hcResultStore,
		bridge:        bridge,
	}
}

// Register registers job handlers and schedules the recurring poller.
func (s *Service) Register(ctx context.Context) error {
	if err := s.executor.Register(jobTypeHealthPoller, &healthPollerHandler{
		hcStore:       s.hcStore,
		hcResultStore: s.hcResultStore,
		bridge:        s.bridge,
	}); err != nil {
		return fmt.Errorf("failed to register health poller handler: %w", err)
	}

	if err := s.scheduler.AddRecurring(ctx, jobTypeHealthPoller, jobTypeHealthPoller, jobCronPoller, jobMaxDurPoller); err != nil {
		return fmt.Errorf("failed to schedule health poller: %w", err)
	}

	return nil
}

// healthPollerHandler finds due health checks and runs them.
type healthPollerHandler struct {
	hcStore       store.HealthCheckStore
	hcResultStore store.HealthCheckResultStore
	bridge        *healthbridge.Bridge
}

func (h *healthPollerHandler) Handle(ctx context.Context, _ string, _ job.ProgressReporter) (string, error) {
	due, err := h.hcStore.ListDue(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list due health checks: %w", err)
	}

	checked := 0
	for _, hc := range due {
		h.runCheck(ctx, hc)
		checked++
	}

	return fmt.Sprintf("checked %d health endpoints", checked), nil
}

func (h *healthPollerHandler) runCheck(ctx context.Context, hc *types.HealthCheck) {
	result := executeHTTPCheck(hc)

	if err := h.hcResultStore.Create(ctx, result); err != nil {
		log.Error().Err(err).Str("hc", hc.Identifier).Msg("healthrunner: failed to store result")
		return
	}

	_, err := h.hcStore.UpdateOptLock(ctx, hc, func(hc *types.HealthCheck) error {
		hc.LastCheckedAt = result.CreatedAt
		hc.LastResponseTime = result.ResponseTime
		hc.LastStatus = result.Status

		if result.Status == string(types.HealthCheckStatusDown) {
			hc.ConsecutiveFailures++
		} else {
			hc.ConsecutiveFailures = 0
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("hc", hc.Identifier).Msg("healthrunner: failed to update health check")
		return
	}

	// Notify bridge on failure
	if result.Status == string(types.HealthCheckStatusDown) && h.bridge != nil {
		// Re-read the health check to get the updated consecutive_failures count
		updated, findErr := h.hcStore.Find(ctx, hc.ID)
		if findErr == nil {
			h.bridge.OnHealthCheckFailed(ctx, updated, result)
		}
	}
}

func executeHTTPCheck(hc *types.HealthCheck) *types.HealthCheckResult {
	result := &types.HealthCheckResult{
		HealthCheckID: hc.ID,
		CreatedAt:     time.Now().UnixMilli(),
	}

	client := &http.Client{
		Timeout: time.Duration(hc.TimeoutSeconds) * time.Second,
	}

	method := hc.Method
	if method == "" {
		method = http.MethodGet
	}

	req, err := http.NewRequest(method, hc.URL, nil)
	if err != nil {
		result.Status = string(types.HealthCheckStatusDown)
		result.ErrorMessage = fmt.Sprintf("invalid request: %s", err)
		return result
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	result.ResponseTime = elapsed.Milliseconds()

	if err != nil {
		result.Status = string(types.HealthCheckStatusDown)
		result.ErrorMessage = err.Error()
		return result
	}
	defer func() {
		_ = resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()

	result.StatusCode = resp.StatusCode

	expectedStatus := hc.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	if resp.StatusCode == expectedStatus {
		result.Status = string(types.HealthCheckStatusUp)
	} else {
		result.Status = string(types.HealthCheckStatusDown)
		result.ErrorMessage = fmt.Sprintf("expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	return result
}
