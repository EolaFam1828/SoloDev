// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
)

// registerResources registers all Tier 3 live resources.
func registerResources(s *Server) {
	spaceRef := GetSpaceRef(nil)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://errors/active",
		Name:        "Active Errors",
		Description: "All open error groups with counts — AI knows what's currently broken before you even ask.",
		MimeType:    "application/json",
	}, s.resourceActiveErrors)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://remediations/pending",
		Name:        "Pending Remediations",
		Description: "All pending + in-progress remediations — AI knows what fixes are already in flight.",
		MimeType:    "application/json",
	}, s.resourcePendingRemediations)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://quality/rules",
		Name:        "Quality Rules",
		Description: "All active quality rules + thresholds — AI enforces your standards while writing code.",
		MimeType:    "application/json",
	}, s.resourceQualityRules)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://quality/summary",
		Name:        "Quality Summary",
		Description: "Space-wide quality health snapshot — AI knows the pass/fail state of every repo.",
		MimeType:    "application/json",
	}, s.resourceQualitySummary)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://security/open-findings",
		Name:        "Open Security Findings",
		Description: "Unresolved security findings + severity — AI never suggests code patterns that match open CVEs.",
		MimeType:    "application/json",
	}, s.resourceSecurityFindings)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://health/status",
		Name:        "Health Status",
		Description: "Current endpoint health across all monitors — AI knows if the backend is up.",
		MimeType:    "application/json",
	}, s.resourceHealthStatus)

	s.RegisterResource(ResourceDefinition{
		URI:         "solodev://tech-debt/hotspots",
		Name:        "Tech Debt Hotspots",
		Description: "Highest-debt files/modules — AI steers you away from hotspot files.",
		MimeType:    "application/json",
	}, s.resourceTechDebtHotspots)

	_ = spaceRef // spaceRef resolved at read time, not registration
}

func (s *Server) resourceActiveErrors(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.ErrorTracker == nil {
		return emptyResource("solodev://errors/active", "error tracker not available")
	}

	openStatus := types.ErrorGroupStatusOpen
	errors, err := s.controllers.ErrorTracker.ListErrors(ctx, session, spaceRef,
		types.ErrorTrackerListOptions{Status: &openStatus})
	if err != nil {
		return nil, fmt.Errorf("list active errors: %w", err)
	}

	summary, _ := s.controllers.ErrorTracker.GetSummary(ctx, session, spaceRef)

	return jsonResource("solodev://errors/active", map[string]interface{}{
		"errors":  errors,
		"summary": summary,
	})
}

func (s *Server) resourcePendingRemediations(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.Remediation == nil {
		return emptyResource("solodev://remediations/pending", "remediation module not available")
	}

	pending, _ := s.controllers.Remediation.ListRemediations(ctx, session, spaceRef,
		&types.RemediationListFilter{Status: remStatusPtr(types.RemediationStatusPending)})
	processing, _ := s.controllers.Remediation.ListRemediations(ctx, session, spaceRef,
		&types.RemediationListFilter{Status: remStatusPtr(types.RemediationStatusProcessing)})

	summary, _ := s.controllers.Remediation.GetSummary(ctx, session, spaceRef)

	return jsonResource("solodev://remediations/pending", map[string]interface{}{
		"pending":    pending,
		"processing": processing,
		"summary":    summary,
	})
}

func (s *Server) resourceQualityRules(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.QualityGate == nil {
		return emptyResource("solodev://quality/rules", "quality gate module not available")
	}

	filter := &types.QualityRuleFilter{ListQueryFilter: types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 200}}}
	rules, err := s.controllers.QualityGate.ListRules(ctx, session, spaceRef, filter)
	if err != nil {
		return nil, fmt.Errorf("list quality rules: %w", err)
	}

	return jsonResource("solodev://quality/rules", rules)
}

func (s *Server) resourceQualitySummary(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.QualityGate == nil {
		return emptyResource("solodev://quality/summary", "quality gate module not available")
	}

	summary, err := s.controllers.QualityGate.GetSummary(ctx, session, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("get quality summary: %w", err)
	}

	return jsonResource("solodev://quality/summary", summary)
}

func (s *Server) resourceSecurityFindings(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.SecurityScan == nil {
		return emptyResource("solodev://security/open-findings", "security scan module not available")
	}

	summary, err := s.controllers.SecurityScan.GetSecuritySummary(ctx, session, spaceRef, nil)
	if err != nil {
		return nil, fmt.Errorf("get security summary: %w", err)
	}

	return jsonResource("solodev://security/open-findings", summary)
}

func (s *Server) resourceHealthStatus(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.HealthCheck == nil {
		return emptyResource("solodev://health/status", "health check module not available")
	}

	summary, err := s.controllers.HealthCheck.GetSummary(ctx, session, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("get health summary: %w", err)
	}

	return jsonResource("solodev://health/status", summary)
}

func (s *Server) resourceTechDebtHotspots(ctx context.Context, session *auth.Session) (*ResourceReadResult, error) {
	spaceRef := GetSpaceRef(nil)
	if s.controllers.TechDebt == nil {
		return emptyResource("solodev://tech-debt/hotspots", "tech debt module not available")
	}

	filter := &types.TechDebtFilter{
		Severity: []string{"critical", "high"},
		Status:   []string{"open"},
		Page:     1,
		Limit:    50,
		Sort:     "severity",
	}
	items, err := s.controllers.TechDebt.List(ctx, session, spaceRef, filter)
	if err != nil {
		return nil, fmt.Errorf("list tech debt: %w", err)
	}

	summary, _ := s.controllers.TechDebt.Summary(ctx, session, spaceRef, &types.TechDebtFilter{})

	return jsonResource("solodev://tech-debt/hotspots", map[string]interface{}{
		"hotspots": items,
		"summary":  summary,
	})
}

// --- Resource helpers ---

func jsonResource(uri string, data interface{}) (*ResourceReadResult, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal resource: %w", err)
	}
	return &ResourceReadResult{
		Contents: []ResourceContent{{
			URI:      uri,
			MimeType: "application/json",
			Text:     string(b),
		}},
	}, nil
}

func emptyResource(uri, message string) (*ResourceReadResult, error) {
	return &ResourceReadResult{
		Contents: []ResourceContent{{
			URI:      uri,
			MimeType: "application/json",
			Text:     fmt.Sprintf(`{"message": "%s"}`, message),
		}},
	}, nil
}
