// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remediationnotifier

import (
	"context"
	"fmt"

	"github.com/harness/gitness/app/services/webhook"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"

	"github.com/rs/zerolog/log"
)

// RemediationPayload is the webhook payload for remediation events.
type RemediationPayload struct {
	Trigger       enum.WebhookTrigger            `json:"trigger"`
	Remediation   RemediationInfo                `json:"remediation"`
	TriggerSource types.RemediationTriggerSource `json:"trigger_source"`
	Delivery      *types.RemediationDelivery     `json:"delivery,omitempty"`
	Validation    *types.RemediationValidation   `json:"validation,omitempty"`
}

// RemediationInfo describes remediation data in a webhook payload.
type RemediationInfo struct {
	ID         int64                   `json:"id"`
	Identifier string                  `json:"identifier"`
	Title      string                  `json:"title"`
	Status     types.RemediationStatus `json:"status"`
	RepoID     int64                   `json:"repo_id"`
	Branch     string                  `json:"branch"`
	AIModel    string                  `json:"ai_model,omitempty"`
	Confidence float64                 `json:"confidence,omitempty"`
	FixBranch  string                  `json:"fix_branch,omitempty"`
	PRLink     string                  `json:"pr_link,omitempty"`
	DurationMs int64                   `json:"duration_ms,omitempty"`
	Created    int64                   `json:"created"`
	Updated    int64                   `json:"updated"`
}

// Service delivers webhook notifications for remediation events.
type Service struct {
	webhookExecutor *webhook.WebhookExecutor
	spaceStore      store.SpaceStore
}

// NewService creates a new remediation notifier service.
func NewService(
	webhookExecutor *webhook.WebhookExecutor,
	spaceStore store.SpaceStore,
) *Service {
	return &Service{
		webhookExecutor: webhookExecutor,
		spaceStore:      spaceStore,
	}
}

// NotifyCompleted sends webhooks for a completed remediation.
func (s *Service) NotifyCompleted(ctx context.Context, rem *types.Remediation) {
	if s.webhookExecutor == nil {
		return
	}
	s.notify(ctx, rem, enum.WebhookTriggerRemediationCompleted)
}

// NotifyApplied sends webhooks for an applied remediation.
func (s *Service) NotifyApplied(ctx context.Context, rem *types.Remediation) {
	if s.webhookExecutor == nil {
		return
	}
	s.notify(ctx, rem, enum.WebhookTriggerRemediationApplied)
}

func (s *Service) notify(ctx context.Context, rem *types.Remediation, trigger enum.WebhookTrigger) {
	parents, err := s.getSpaceParents(ctx, rem.SpaceID)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).
			Int64("remediation_id", rem.ID).
			Str("trigger", string(trigger)).
			Msg("failed to resolve webhook parents for remediation notification")
		return
	}

	rem.PopulateAPIFields()

	payload := RemediationPayload{
		Trigger: trigger,
		Remediation: RemediationInfo{
			ID:         rem.ID,
			Identifier: rem.Identifier,
			Title:      rem.Title,
			Status:     rem.Status,
			RepoID:     rem.RepoID,
			Branch:     rem.Branch,
			AIModel:    rem.AIModel,
			Confidence: rem.Confidence,
			FixBranch:  rem.FixBranch,
			PRLink:     rem.PRLink,
			DurationMs: rem.DurationMs,
			Created:    rem.Created,
			Updated:    rem.Updated,
		},
		TriggerSource: rem.TriggerSource,
		Delivery:      rem.Delivery,
		Validation:    rem.Validation,
	}

	eventID := fmt.Sprintf("remediation-%d-%s", rem.ID, trigger)
	err = s.webhookExecutor.TriggerForEvent(ctx, eventID, parents, trigger, payload)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).
			Int64("remediation_id", rem.ID).
			Str("trigger", string(trigger)).
			Msg("failed to trigger webhooks for remediation event")
	}
}

func (s *Service) getSpaceParents(ctx context.Context, spaceID int64) ([]types.WebhookParentInfo, error) {
	ids, err := s.spaceStore.GetAncestorIDs(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("get ancestor space IDs: %w", err)
	}

	parents := make([]types.WebhookParentInfo, len(ids))
	for i, id := range ids {
		parents[i] = types.WebhookParentInfo{
			Type: enum.WebhookParentSpace,
			ID:   id,
		}
	}
	return parents, nil
}
