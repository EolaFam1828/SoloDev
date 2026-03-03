// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package healthbridge connects the Health Monitor to the AI Remediation system.
// When a health check transitions to "down" with consecutive failures exceeding
// a threshold, this bridge automatically creates a remediation task.
package healthbridge

import (
	"context"
	"fmt"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"

	"github.com/rs/zerolog/log"
)

const defaultConsecutiveFailureThreshold = 3

// Bridge watches for health check failures and creates remediation tasks.
type Bridge struct {
	remediationStore store.RemediationStore
	enabled          bool
	threshold        int // consecutive failures before triggering remediation
}

// NewBridge creates a new health-to-remediation bridge.
func NewBridge(
	remediationStore store.RemediationStore,
	enabled bool,
) *Bridge {
	return &Bridge{
		remediationStore: remediationStore,
		enabled:          enabled,
		threshold:        defaultConsecutiveFailureThreshold,
	}
}

// OnHealthCheckFailed is called when a health check records a failure result.
// It creates a remediation task once consecutive failures exceed the threshold.
func (b *Bridge) OnHealthCheckFailed(ctx context.Context, hc *types.HealthCheck, result *types.HealthCheckResult) {
	if !b.enabled {
		return
	}

	if hc.ConsecutiveFailures < b.threshold {
		return
	}

	// Only trigger once at the threshold boundary to avoid duplicate remediations
	if hc.ConsecutiveFailures != b.threshold {
		return
	}

	rem := &types.Remediation{
		SpaceID:       hc.SpaceID,
		Identifier:    fmt.Sprintf("hc-%s-%d", hc.Identifier, types.NowMillis()),
		Title:         fmt.Sprintf("[Auto] Health check failing: %s", hc.Name),
		Description:   fmt.Sprintf("Health check %q (%s %s) has failed %d consecutive times.", hc.Name, hc.Method, hc.URL, hc.ConsecutiveFailures),
		Status:        types.RemediationStatusPending,
		TriggerSource: types.RemediationTriggerHealthCheck,
		TriggerRef:    hc.Identifier,
		ErrorLog:      result.ErrorMessage,
		CreatedBy:     hc.CreatedBy,
		Created:       types.NowMillis(),
		Updated:       types.NowMillis(),
		Version:       1,
	}

	if err := b.remediationStore.Create(ctx, rem); err != nil {
		log.Error().Err(err).
			Str("health_check", hc.Identifier).
			Msg("healthbridge: failed to auto-create remediation")
		return
	}

	log.Info().
		Str("remediation_id", rem.Identifier).
		Str("health_check", hc.Identifier).
		Int("consecutive_failures", hc.ConsecutiveFailures).
		Msg("healthbridge: auto-created remediation for health check failure")
}
