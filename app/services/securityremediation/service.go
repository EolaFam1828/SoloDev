// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package securityremediation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/services/aiworker"
	"github.com/harness/gitness/app/services/settings"
	"github.com/harness/gitness/app/store"
	basestore "github.com/harness/gitness/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type Service struct {
	settingsService  *settings.Service
	remediationStore store.RemediationStore
	enabled          bool
}

func NewService(
	settingsService *settings.Service,
	remediationStore store.RemediationStore,
	aiWorker *aiworker.Service,
) *Service {
	enabled := aiWorker != nil && aiWorker.Available()
	if !enabled {
		return nil
	}

	return &Service{
		settingsService:  settingsService,
		remediationStore: remediationStore,
		enabled:          enabled,
	}
}

func (s *Service) Enabled() bool {
	return s != nil && s.enabled
}

func (s *Service) RepoMode(ctx context.Context, repoID int64) (types.SecurityFindingRemediationMode, error) {
	if s == nil {
		return types.SecurityFindingRemediationModeManual, nil
	}

	mode := types.SecurityFindingRemediationMode(settings.DefaultFindingRemediationMode)
	found, err := s.settingsService.RepoGet(ctx, repoID, settings.KeyFindingRemediationMode, &mode)
	if err != nil {
		return "", fmt.Errorf("failed to load remediation mode: %w", err)
	}
	if !found {
		return types.SecurityFindingRemediationModeManual, nil
	}

	switch mode {
	case types.SecurityFindingRemediationModeManual,
		types.SecurityFindingRemediationModeCriticalHighAuto,
		types.SecurityFindingRemediationModeAllAuto:
		return mode, nil
	default:
		return types.SecurityFindingRemediationModeManual, nil
	}
}

func (s *Service) AutoCreateForFindings(
	ctx context.Context,
	scan *types.ScanResult,
	findings []*types.ScanFinding,
	createdBy int64,
) error {
	if !s.Enabled() || scan == nil {
		return nil
	}

	mode, err := s.RepoMode(ctx, scan.RepoID)
	if err != nil {
		return err
	}
	if mode == types.SecurityFindingRemediationModeManual {
		return nil
	}

	for _, finding := range findings {
		if finding == nil || finding.Status != enum.SecurityFindingStatusOpen {
			continue
		}
		if !shouldAutoCreate(mode, finding.Severity) {
			continue
		}
		if _, _, err := s.CreateFromFinding(ctx, scan, finding, createdBy, true); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) CreateFromFinding(
	ctx context.Context,
	scan *types.ScanResult,
	finding *types.ScanFinding,
	createdBy int64,
	auto bool,
) (*types.Remediation, bool, error) {
	if !s.Enabled() {
		return nil, false, fmt.Errorf("AI remediation is not configured")
	}
	if scan == nil || finding == nil {
		return nil, false, fmt.Errorf("scan and finding are required")
	}

	triggerRef := BuildTriggerRef(scan.RepoID, finding)
	existing, err := s.remediationStore.FindActiveByTriggerRef(ctx, scan.RepoID, triggerRef)
	if err == nil {
		return existing, false, nil
	}
	if err != nil && err != basestore.ErrResourceNotFound {
		return nil, false, fmt.Errorf("failed to check existing remediation: %w", err)
	}

	metadata, err := buildMetadata(scan, finding, auto)
	if err != nil {
		return nil, false, fmt.Errorf("failed to build remediation metadata: %w", err)
	}

	remediation := &types.Remediation{
		SpaceID:       scan.SpaceID,
		RepoID:        scan.RepoID,
		Identifier:    generateIdentifier(),
		Title:         fmt.Sprintf("Fix security finding: %s", finding.Title),
		Description:   finding.Description,
		Status:        types.RemediationStatusPending,
		TriggerSource: types.RemediationTriggerSecurity,
		TriggerRef:    triggerRef,
		Branch:        scan.Branch,
		CommitSHA:     scan.CommitSHA,
		ErrorLog:      finding.Description,
		SourceCode:    finding.Snippet,
		FilePath:      finding.FilePath,
		Metadata:      metadata,
		CreatedBy:     createdBy,
		Created:       types.NowMillis(),
		Updated:       types.NowMillis(),
		Version:       1,
	}

	if remediation.Description == "" {
		remediation.Description = finding.Title
	}

	if err := s.remediationStore.Create(ctx, remediation); err != nil {
		return nil, false, fmt.Errorf("failed to create remediation: %w", err)
	}

	return remediation, true, nil
}

func BuildTriggerRef(repoID int64, finding *types.ScanFinding) string {
	return fmt.Sprintf(
		"security:%d:%s:%s:%d:%d",
		repoID,
		finding.Identifier,
		finding.FilePath,
		finding.LineStart,
		finding.LineEnd,
	)
}

func shouldAutoCreate(mode types.SecurityFindingRemediationMode, severity enum.SecurityFindingSeverity) bool {
	switch mode {
	case types.SecurityFindingRemediationModeAllAuto:
		return true
	case types.SecurityFindingRemediationModeCriticalHighAuto:
		return severity == enum.SecurityFindingSeverityCritical || severity == enum.SecurityFindingSeverityHigh
	default:
		return false
	}
}

func buildMetadata(scan *types.ScanResult, finding *types.ScanFinding, auto bool) (json.RawMessage, error) {
	payload := map[string]any{
		"scan_id":            scan.ID,
		"scan_identifier":    scan.Identifier,
		"finding_id":         finding.ID,
		"finding_identifier": finding.Identifier,
		"severity":           finding.Severity,
		"category":           finding.Category,
		"rule":               finding.Rule,
		"cwe":                finding.CWE,
		"auto_triggered":     auto,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func generateIdentifier() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return "rem-" + hex.EncodeToString(buf)
}
