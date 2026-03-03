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

package remediationdelivery

import (
	"context"
	stdliberrors "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	controllerpullreq "github.com/harness/gitness/app/api/controller/pullreq"
	controllerrepo "github.com/harness/gitness/app/api/controller/repo"
	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/app/auth"
	airemediationevents "github.com/harness/gitness/app/events/airemediation"
	"github.com/harness/gitness/app/services/remediationnotifier"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/app/url"
	"github.com/harness/gitness/types"
)

type Service struct {
	remStore       store.RemediationStore
	repoStore      store.RepoStore
	principalStore store.PrincipalStore
	repoCtrl       *controllerrepo.Controller
	pullreqCtrl    *controllerpullreq.Controller
	urlProvider    url.Provider
	eventReporter  *airemediationevents.Reporter
	notifier       *remediationnotifier.Service
}

func NewService(
	remStore store.RemediationStore,
	repoStore store.RepoStore,
	principalStore store.PrincipalStore,
	repoCtrl *controllerrepo.Controller,
	pullreqCtrl *controllerpullreq.Controller,
	urlProvider url.Provider,
	eventReporter *airemediationevents.Reporter,
	notifier *remediationnotifier.Service,
) *Service {
	return &Service{
		remStore:       remStore,
		repoStore:      repoStore,
		principalStore: principalStore,
		repoCtrl:       repoCtrl,
		pullreqCtrl:    pullreqCtrl,
		urlProvider:    urlProvider,
		eventReporter:  eventReporter,
		notifier:       notifier,
	}
}

func (s *Service) Apply(
	ctx context.Context,
	session *auth.Session,
	rem *types.Remediation,
	mode types.RemediationDeliveryMode,
) (*types.Remediation, error) {
	if rem == nil {
		return nil, usererror.BadRequest("remediation is required")
	}
	if session == nil {
		return s.failDelivery(ctx, rem, mode, "delivery actor session is required", usererror.BadRequest("delivery actor session is required"))
	}
	if s.repoCtrl == nil || s.pullreqCtrl == nil {
		return s.failDelivery(
			ctx,
			rem,
			mode,
			"remediation delivery service is not configured",
			usererror.New(http.StatusServiceUnavailable, "remediation delivery service is not configured"),
		)
	}
	if rem.PRLink != "" {
		return rem, nil
	}

	repo, err := s.repoStore.Find(ctx, rem.RepoID)
	if err != nil {
		return s.failDelivery(ctx, rem, mode, fmt.Sprintf("failed to load remediation repository: %v", err), fmt.Errorf("find remediation repo: %w", err))
	}

	if err := s.validate(rem); err != nil {
		return s.failDelivery(ctx, rem, mode, err.Error(), err)
	}

	delivery, err := s.prepareAttempt(rem, mode)
	if err != nil {
		return nil, err
	}

	if rem.FixBranch == "" {
		branchName := remediationFixBranch(rem.Identifier)
		title := remediationTitle(rem)
		message := remediationCommitMessage(rem)

		out, violations, err := s.repoCtrl.ApplyPatch(ctx, session, repo.Path, &controllerrepo.ApplyPatchOptions{
			Title:     title,
			Message:   message,
			Branch:    rem.Branch,
			NewBranch: branchName,
			Patch:     rem.PatchDiff,
		})
		if err != nil {
			return s.failDelivery(ctx, rem, mode, fmt.Sprintf("failed to apply remediation patch: %v", err), err)
		}
		if out.CommitID.IsEmpty() {
			msg := "remediation delivery was blocked by repository rules"
			if len(violations) > 0 {
				msg = fmt.Sprintf("%s (%d violation(s))", msg, len(violations))
			}
			return s.failDelivery(ctx, rem, mode, msg, usererror.Conflict(msg))
		}

		rem.FixBranch = branchName
		delivery.State = types.RemediationDeliveryStateBranchReady
		delivery.LastError = ""
		delivery.PRNumber = 0
		rem.Metadata, err = types.SetRemediationDeliveryMetadata(rem.Metadata, delivery)
		if err != nil {
			return nil, fmt.Errorf("encode remediation delivery metadata: %w", err)
		}
		if err := s.remStore.Update(ctx, rem); err != nil {
			return nil, fmt.Errorf("persist remediation fix branch: %w", err)
		}
	}

	pr, err := s.pullreqCtrl.Create(ctx, session, repo.Path, &controllerpullreq.CreateInput{
		IsDraft:      true,
		Title:        remediationTitle(rem),
		Description:  remediationPRBody(rem),
		SourceBranch: rem.FixBranch,
		TargetBranch: rem.Branch,
	})
	if err != nil {
		if prNumber, ok := existingPRNumber(err); ok {
			rem.PRLink = s.urlProvider.GenerateUIPRURL(ctx, repo.Path, prNumber)
			rem.Status = types.RemediationStatusApplied
			delivery.State = types.RemediationDeliveryStateApplied
			delivery.PRNumber = prNumber
			delivery.LastError = ""
			rem.Metadata, err = types.SetRemediationDeliveryMetadata(rem.Metadata, delivery)
			if err != nil {
				return nil, fmt.Errorf("encode remediation delivery metadata: %w", err)
			}
			if err := s.remStore.Update(ctx, rem); err != nil {
				return nil, fmt.Errorf("persist existing remediation pull request: %w", err)
			}
			s.reportApplied(ctx, rem)
			return rem, nil
		}
		return s.failWithState(
			ctx,
			rem,
			mode,
			types.RemediationDeliveryStateBranchReady,
			fmt.Sprintf("failed to create remediation pull request: %v", err),
			err,
		)
	}

	rem.PRLink = s.urlProvider.GenerateUIPRURL(ctx, repo.Path, pr.Number)
	rem.Status = types.RemediationStatusApplied
	delivery.State = types.RemediationDeliveryStateApplied
	delivery.PRNumber = pr.Number
	delivery.LastError = ""
	rem.Metadata, err = types.SetRemediationDeliveryMetadata(rem.Metadata, delivery)
	if err != nil {
		return nil, fmt.Errorf("encode remediation delivery metadata: %w", err)
	}
	if err := s.remStore.Update(ctx, rem); err != nil {
		return nil, fmt.Errorf("persist remediation pull request: %w", err)
	}

	s.reportApplied(ctx, rem)
	return rem, nil
}

func (s *Service) ApplyAsRemediationCreator(
	ctx context.Context,
	rem *types.Remediation,
	mode types.RemediationDeliveryMode,
) (*types.Remediation, error) {
	if rem == nil {
		return nil, usererror.BadRequest("remediation is required")
	}

	principal, err := s.principalStore.Find(ctx, rem.CreatedBy)
	if err != nil {
		return s.failDelivery(
			ctx,
			rem,
			mode,
			fmt.Sprintf("failed to reconstruct remediation creator session: %v", err),
			fmt.Errorf("find remediation creator: %w", err),
		)
	}

	return s.Apply(ctx, &auth.Session{
		Principal: *principal,
		Metadata:  &auth.EmptyMetadata{},
	}, rem, mode)
}

func (s *Service) validate(rem *types.Remediation) error {
	switch rem.Status {
	case types.RemediationStatusCompleted, types.RemediationStatusApplied:
	default:
		return usererror.Conflict("only completed or applied remediations can be delivered")
	}
	if rem.RepoID <= 0 {
		return usererror.Conflict("remediation is missing repo_id")
	}
	if strings.TrimSpace(rem.Branch) == "" {
		return usererror.Conflict("remediation is missing target branch")
	}
	if strings.TrimSpace(rem.PatchDiff) == "" {
		return usererror.Conflict("remediation does not contain a patch diff")
	}

	return nil
}

func (s *Service) prepareAttempt(
	rem *types.Remediation,
	mode types.RemediationDeliveryMode,
) (types.RemediationDelivery, error) {
	delivery, err := types.GetRemediationDeliveryMetadata(rem.Metadata, mode)
	if err != nil {
		return types.RemediationDelivery{}, fmt.Errorf("decode remediation delivery metadata: %w", err)
	}
	if mode == "" {
		mode = types.RemediationDeliveryModeManual
	}
	delivery.Mode = mode
	delivery.AttemptedAt = types.NowMillis()
	if delivery.State == "" {
		delivery.State = types.RemediationDeliveryStateNotAttempted
	}
	return delivery, nil
}

func (s *Service) failDelivery(
	ctx context.Context,
	rem *types.Remediation,
	mode types.RemediationDeliveryMode,
	message string,
	cause error,
) (*types.Remediation, error) {
	return s.failWithState(ctx, rem, mode, types.RemediationDeliveryStateFailed, message, cause)
}

func (s *Service) failWithState(
	ctx context.Context,
	rem *types.Remediation,
	mode types.RemediationDeliveryMode,
	state types.RemediationDeliveryState,
	message string,
	cause error,
) (*types.Remediation, error) {
	if rem == nil {
		return nil, cause
	}

	delivery, err := s.prepareAttempt(rem, mode)
	if err == nil {
		delivery.State = state
		delivery.LastError = message
		updatedMetadata, setErr := types.SetRemediationDeliveryMetadata(rem.Metadata, delivery)
		if setErr == nil {
			rem.Metadata = updatedMetadata
			if rem.Status == types.RemediationStatusApplied && rem.PRLink == "" {
				rem.Status = types.RemediationStatusCompleted
			}
			if storeErr := s.remStore.Update(ctx, rem); storeErr != nil {
				return rem, fmt.Errorf("%w: %v", cause, storeErr)
			}
		}
	}

	return rem, cause
}

func (s *Service) reportApplied(ctx context.Context, rem *types.Remediation) {
	if s.eventReporter != nil {
		s.eventReporter.RemediationApplied(ctx, rem)
	}
	if s.notifier != nil {
		s.notifier.NotifyApplied(ctx, rem)
	}
}

func remediationFixBranch(identifier string) string {
	return "solodev/rem-" + strings.TrimSpace(identifier)
}

func remediationTitle(rem *types.Remediation) string {
	return "[SoloDev] " + strings.TrimSpace(rem.Title)
}

func remediationCommitMessage(rem *types.Remediation) string {
	return fmt.Sprintf(
		"Remediation: %s\nSource: %s\nRef: %s\nConfidence: %.2f",
		rem.Identifier,
		rem.TriggerSource,
		emptyFallback(rem.TriggerRef, "n/a"),
		rem.Confidence,
	)
}

func remediationPRBody(rem *types.Remediation) string {
	return fmt.Sprintf(
		"Remediation: %s\nTrigger source: %s\nTrigger ref: %s\nConfidence: %.2f\n\nAI-generated draft PR, human review required.",
		rem.Identifier,
		rem.TriggerSource,
		emptyFallback(rem.TriggerRef, "n/a"),
		rem.Confidence,
	)
}

func emptyFallback(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func existingPRNumber(err error) (int64, bool) {
	var userErr *usererror.Error
	if !stdliberrors.As(err, &userErr) {
		return 0, false
	}
	if userErr.Status != http.StatusConflict || userErr.Values == nil {
		return 0, false
	}
	if kind, _ := userErr.Values["type"].(string); kind != "pr already exists" {
		return 0, false
	}

	switch value := userErr.Values["number"].(type) {
	case int64:
		return value, true
	case int:
		return int64(value), true
	case float64:
		return int64(value), true
	case string:
		number, convErr := strconv.ParseInt(value, 10, 64)
		if convErr == nil {
			return number, true
		}
	}

	return 0, false
}
