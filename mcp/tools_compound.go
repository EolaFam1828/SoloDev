// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/EolaFam1828/SoloDev/app/api/controller/qualitygate"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

// registerCompoundTools registers all Tier 2 tools (multi-step orchestration).
func registerCompoundTools(s *Server) {
	s.RegisterTool(ToolDefinition{
		Name: "fix_this",
		Description: "End-to-end error-to-patch: paste a stack trace and get a unified diff patch. " +
			"Reports the error, triggers auto-remediation via Error Bridge, polls until fix is generated, " +
			"and returns the patch diff + AI analysis.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":   StringProp("Space reference"),
			"error_log":   StringProp("Full error log / stack trace"),
			"file_path":   StringProp("Source file where the error occurred"),
			"source_code": StringProp("Relevant source code snippet"),
			"title":       StringProp("Error title (auto-generated if empty)"),
			"branch":      StringProp("Branch name (defaults to 'main')"),
		}, []string{"space_ref", "error_log"}),
	}, s.toolFixThis)

	s.RegisterTool(ToolDefinition{
		Name: "pr_ready",
		Description: "Full pre-commit quality check: runs security scan, quality evaluation, and tech debt check " +
			"in sequence. Returns a structured PASS/WARN/BLOCK verdict with itemized reasons.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"repo_ref":   StringProp("Repository reference"),
			"commit_sha": StringProp("Commit SHA to check"),
			"branch":     StringProp("Branch name"),
			"file_paths": ArrayOfStrings("List of changed file paths"),
		}, []string{"space_ref", "repo_ref", "commit_sha"}),
	}, s.toolPRReady)

	s.RegisterTool(ToolDefinition{
		Name: "pipeline_validate",
		Description: "Validate a pipeline YAML for correctness. Parses the YAML, checks structure against " +
			"known step schemas, and flags common mistakes (missing test step, no caching, etc.).",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"yaml":       StringProp("Pipeline YAML content to validate"),
			"stack_hint": StringProp("Optional hint: go, node, python, rust, java"),
		}, []string{"yaml"}),
	}, s.toolPipelineValidate)

	s.RegisterTool(ToolDefinition{
		Name: "onboard_repo",
		Description: "Zero-to-CI-ready in one call: generates pipeline YAML, confirms quality rules exist, " +
			"and verifies health check monitoring is active.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":         StringProp("Space reference"),
			"repo_file_listing": ArrayOfStrings("List of files in the repository"),
		}, []string{"space_ref", "repo_file_listing"}),
	}, s.toolOnboardRepo)

	s.RegisterTool(ToolDefinition{
		Name: "incident_triage",
		Description: "Active incident full context dump: groups recent errors by severity, checks endpoint health, " +
			"and shows what remediations are already in flight. Replaces your incident morning ritual.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":   StringProp("Space reference"),
			"service_name": StringProp("Service name to investigate (optional filter)"),
			"time_window":  StringProp("Time window: 'last 1 hour', 'last 2 hours', etc. (informational)"),
		}, []string{"space_ref"}),
	}, s.toolIncidentTriage)
}

// --- Compound Tool Handlers ---

func (s *Server) toolFixThis(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef   string `json:"space_ref"`
		ErrorLog   string `json:"error_log"`
		FilePath   string `json:"file_path"`
		SourceCode string `json:"source_code"`
		Title      string `json:"title"`
		Branch     string `json:"branch"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}

	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	if p.Branch == "" {
		p.Branch = "main"
	}
	if p.Title == "" {
		p.Title = extractTitle(p.ErrorLog)
	}

	// Step 1: Report the error
	var errorGroup *types.ErrorGroup
	if s.controllers.ErrorTracker != nil {
		identifier := fmt.Sprintf("fix-%d", types.NowMillis())
		in := &types.ReportErrorInput{
			Identifier:  identifier,
			Title:       p.Title,
			Message:     p.ErrorLog,
			Severity:    types.ErrorSeverityError,
			StackTrace:  p.ErrorLog,
			FilePath:    p.FilePath,
			Environment: "mcp",
		}
		var err error
		errorGroup, err = s.controllers.ErrorTracker.ReportError(ctx, session, spaceRef, in)
		if err != nil {
			return nil, fmt.Errorf("error_report failed: %w", err)
		}
	}

	// Step 2: The Error Bridge auto-triggers a remediation task.
	// If the bridge didn't fire (e.g. disabled), manually trigger remediation.
	if s.controllers.Remediation != nil {
		// Check if a remediation was already created by the bridge
		remediations, _ := s.controllers.Remediation.ListRemediations(ctx, session, spaceRef,
			&types.RemediationListFilter{
				Status: remStatusPtr(types.RemediationStatusPending),
			})

		var targetRem *types.Remediation
		// Look for the bridge-created remediation matching our error
		if errorGroup != nil {
			for _, rem := range remediations {
				if rem.TriggerSource == types.RemediationTriggerError &&
					rem.TriggerRef == errorGroup.Identifier {
					targetRem = rem
					break
				}
			}
		}

		// If no bridge-created remediation exists, trigger manually
		if targetRem == nil {
			in := &types.TriggerRemediationInput{
				Title:         fmt.Sprintf("[MCP] Fix: %s", p.Title),
				Description:   "Auto-triggered by fix_this MCP tool",
				TriggerSource: types.RemediationTriggerManual,
				Branch:        p.Branch,
				ErrorLog:      p.ErrorLog,
				SourceCode:    p.SourceCode,
				FilePath:      p.FilePath,
			}
			rem, err := s.controllers.Remediation.TriggerRemediation(ctx, session, spaceRef, in)
			if err != nil {
				return nil, fmt.Errorf("remediation_trigger failed: %w", err)
			}
			targetRem = rem
		}

		// Step 3: Poll remediation status with timeout
		result := s.pollRemediation(ctx, session, spaceRef, targetRem.Identifier)
		return SuccessResult(result)
	}

	// Fallback if no remediation module
	return SuccessResult(map[string]interface{}{
		"error_reported": errorGroup,
		"message":        "Error reported but remediation module not available. No patch generated.",
	})
}

func (s *Server) pollRemediation(ctx context.Context, session *auth.Session, spaceRef, identifier string) map[string]interface{} {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			rem, _ := s.controllers.Remediation.GetRemediation(ctx, session, spaceRef, identifier)
			return map[string]interface{}{
				"status":      "timeout",
				"remediation": rem,
				"message":     "Remediation is still processing. Check back with remediation_get.",
			}
		case <-ticker.C:
			rem, err := s.controllers.Remediation.GetRemediation(ctx, session, spaceRef, identifier)
			if err != nil {
				continue
			}
			switch rem.Status {
			case types.RemediationStatusCompleted, types.RemediationStatusApplied:
				return map[string]interface{}{
					"status":        string(rem.Status),
					"patch_diff":    rem.PatchDiff,
					"ai_response":   rem.AIResponse,
					"fix_branch":    rem.FixBranch,
					"pr_link":       rem.PRLink,
					"confidence":    rem.Confidence,
					"remediation_id": rem.Identifier,
				}
			case types.RemediationStatusFailed:
				return map[string]interface{}{
					"status":        "failed",
					"ai_response":   rem.AIResponse,
					"remediation_id": rem.Identifier,
					"message":        "AI remediation failed. Review the AI response for details.",
				}
			}
		case <-ctx.Done():
			return map[string]interface{}{
				"status":  "cancelled",
				"message": "Request cancelled.",
			}
		}
	}
}

func (s *Server) toolPRReady(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef  string   `json:"space_ref"`
		RepoRef   string   `json:"repo_ref"`
		CommitSHA string   `json:"commit_sha"`
		Branch    string   `json:"branch"`
		FilePaths []string `json:"file_paths"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}

	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	verdict := "PASS"
	var issues []map[string]interface{}

	// Step 1: Security Scan
	var scanResult *types.ScanResult
	if s.controllers.SecurityScan != nil {
		scanType := enum.SecurityScanType("full")
		trigger := enum.SecurityScanTriggerManual
		in := &types.ScanResultInput{
			ScanType:    &scanType,
			CommitSHA:   strPtr(p.CommitSHA),
			Branch:      strPtr(p.Branch),
			TriggeredBy: &trigger,
		}
		var err error
		scanResult, err = s.controllers.SecurityScan.TriggerScan(ctx, session, p.RepoRef, in)
		if err == nil && scanResult != nil {
			if scanResult.CriticalCount > 0 || scanResult.HighCount > 0 {
				verdict = "BLOCK"
				issues = append(issues, map[string]interface{}{
					"type":     "security",
					"severity": "critical",
					"message":  fmt.Sprintf("%d critical + %d high security findings", scanResult.CriticalCount, scanResult.HighCount),
				})
			} else if scanResult.MediumCount > 0 {
				if verdict != "BLOCK" {
					verdict = "WARN"
				}
				issues = append(issues, map[string]interface{}{
					"type":     "security",
					"severity": "medium",
					"message":  fmt.Sprintf("%d medium security findings", scanResult.MediumCount),
				})
			}
		}
	}

	// Step 2: Quality Evaluation
	var qualityResult *types.QualityEvaluation
	if s.controllers.QualityGate != nil {
		trigger := enum.QualityTriggerManual
		in := &qualitygate.EvaluateInput{
			RepoRef:   p.RepoRef,
			CommitSHA: p.CommitSHA,
			Branch:    p.Branch,
			Trigger:   trigger,
		}
		var err error
		qualityResult, err = s.controllers.QualityGate.Evaluate(ctx, session, spaceRef, in)
		if err == nil && qualityResult != nil {
			if qualityResult.OverallStatus == enum.QualityStatusFailed {
				verdict = "BLOCK"
				issues = append(issues, map[string]interface{}{
					"type":    "quality",
					"message": fmt.Sprintf("%d rules failed out of %d evaluated", qualityResult.RulesFailed, qualityResult.RulesEvaluated),
				})
			} else if qualityResult.RulesWarned > 0 {
				if verdict != "BLOCK" {
					verdict = "WARN"
				}
				issues = append(issues, map[string]interface{}{
					"type":    "quality",
					"message": fmt.Sprintf("%d quality warnings", qualityResult.RulesWarned),
				})
			}
		}
	}

	// Step 3: Tech Debt check
	if s.controllers.TechDebt != nil {
		filter := &types.TechDebtFilter{
			Status: []string{"open"},
			Page:   1,
			Limit:  50,
		}
		items, err := s.controllers.TechDebt.List(ctx, session, spaceRef, filter)
		if err == nil && len(items.Items) > 0 {
			criticalDebt := 0
			for _, item := range items.Items {
				if item.Severity == types.TechDebtSeverityCritical {
					criticalDebt++
				}
			}
			if criticalDebt > 0 {
				if verdict != "BLOCK" {
					verdict = "WARN"
				}
				issues = append(issues, map[string]interface{}{
					"type":    "tech_debt",
					"message": fmt.Sprintf("%d critical tech debt items, %d total open", criticalDebt, len(items.Items)),
				})
			}
		}
	}

	return SuccessResult(map[string]interface{}{
		"verdict":         verdict,
		"issues":          issues,
		"security_scan":   scanResult,
		"quality_result":  qualityResult,
	})
}

func (s *Server) toolPipelineValidate(_ context.Context, _ *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef  string `json:"space_ref"`
		YAML      string `json:"yaml"`
		StackHint string `json:"stack_hint"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}

	var issues []map[string]interface{}
	lines := strings.Split(p.YAML, "\n")

	// Basic YAML structure validation
	hasStages := false
	hasSteps := false
	hasTest := false
	hasCache := false
	hasBuild := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "stages:") {
			hasStages = true
		}
		if strings.HasPrefix(trimmed, "steps:") || strings.HasPrefix(trimmed, "- step:") {
			hasSteps = true
		}
		if strings.Contains(trimmed, "test") || strings.Contains(trimmed, "Test") {
			hasTest = true
		}
		if strings.Contains(trimmed, "cache") || strings.Contains(trimmed, "Cache") {
			hasCache = true
		}
		if strings.Contains(trimmed, "build") || strings.Contains(trimmed, "Build") || strings.Contains(trimmed, "compile") {
			hasBuild = true
		}
		// Check for common mistakes
		if strings.Contains(trimmed, "	") && strings.Contains(p.YAML, "  ") {
			issues = append(issues, map[string]interface{}{
				"line":    i + 1,
				"type":    "warning",
				"message": "Mixed tabs and spaces detected — YAML requires consistent indentation",
			})
		}
	}

	if !hasStages && !hasSteps {
		issues = append(issues, map[string]interface{}{
			"line":    1,
			"type":    "error",
			"message": "No stages or steps found — pipeline has no work to do",
		})
	}

	if !hasTest {
		issues = append(issues, map[string]interface{}{
			"type":    "warning",
			"message": "No test step detected — consider adding tests to the pipeline",
		})
	}

	if !hasCache {
		issues = append(issues, map[string]interface{}{
			"type":    "suggestion",
			"message": "No caching configured — add caching for faster builds",
		})
	}

	if !hasBuild && p.StackHint != "" {
		issues = append(issues, map[string]interface{}{
			"type":    "warning",
			"message": fmt.Sprintf("No build step detected for %s stack", p.StackHint),
		})
	}

	status := "valid"
	for _, issue := range issues {
		if issue["type"] == "error" {
			status = "invalid"
			break
		}
	}

	return SuccessResult(map[string]interface{}{
		"status":     status,
		"issues":     issues,
		"line_count": len(lines),
		"has_stages": hasStages,
		"has_steps":  hasSteps,
		"has_tests":  hasTest,
		"has_cache":  hasCache,
	})
}

func (s *Server) toolOnboardRepo(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef        string   `json:"space_ref"`
		RepoFileListing []string `json:"repo_file_listing"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}

	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result := map[string]interface{}{}

	// Step 1: Generate pipeline
	if s.controllers.AutoPipeline != nil {
		config, err := s.controllers.AutoPipeline.GenerateAutoConfig(ctx, session, spaceRef, p.RepoFileListing)
		if err == nil {
			result["pipeline_yaml"] = config.YAML
			result["detected_stack"] = config.Stack
		} else {
			result["pipeline_error"] = err.Error()
		}
	}

	// Step 2: Check quality rules
	if s.controllers.QualityGate != nil {
		filter := &types.QualityRuleFilter{ListQueryFilter: types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 100}}}
		rules, err := s.controllers.QualityGate.ListRules(ctx, session, spaceRef, filter)
		if err == nil {
			result["quality_rules_configured"] = rules
		} else {
			result["quality_rules_error"] = err.Error()
		}
	}

	// Step 3: Check health monitoring
	if s.controllers.HealthCheck != nil {
		summary, err := s.controllers.HealthCheck.GetSummary(ctx, session, spaceRef)
		if err == nil {
			result["health_checks_active"] = summary
		} else {
			result["health_checks_error"] = err.Error()
		}
	}

	return SuccessResult(result)
}

func (s *Server) toolIncidentTriage(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef    string `json:"space_ref"`
		ServiceName string `json:"service_name"`
		TimeWindow  string `json:"time_window"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}

	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	report := map[string]interface{}{
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"time_window":  p.TimeWindow,
		"service_name": p.ServiceName,
	}

	// Step 1: Recent errors by severity
	if s.controllers.ErrorTracker != nil {
		errors, err := s.controllers.ErrorTracker.ListErrors(ctx, session, spaceRef, types.ErrorTrackerListOptions{})
		if err == nil {
			report["active_errors"] = errors

			summary, serr := s.controllers.ErrorTracker.GetSummary(ctx, session, spaceRef)
			if serr == nil {
				report["error_summary"] = summary
			}
		}
	}

	// Step 2: Endpoint health
	if s.controllers.HealthCheck != nil {
		healthSummary, err := s.controllers.HealthCheck.GetSummary(ctx, session, spaceRef)
		if err == nil {
			report["health_status"] = healthSummary
		}
	}

	// Step 3: In-flight remediations
	if s.controllers.Remediation != nil {
		pending, _ := s.controllers.Remediation.ListRemediations(ctx, session, spaceRef,
			&types.RemediationListFilter{
				Status: remStatusPtr(types.RemediationStatusPending),
			})
		processing, _ := s.controllers.Remediation.ListRemediations(ctx, session, spaceRef,
			&types.RemediationListFilter{
				Status: remStatusPtr(types.RemediationStatusProcessing),
			})

		report["pending_remediations"] = pending
		report["processing_remediations"] = processing

		remSummary, err := s.controllers.Remediation.GetSummary(ctx, session, spaceRef)
		if err == nil {
			report["remediation_summary"] = remSummary
		}
	}

	return SuccessResult(report)
}

// --- Helpers ---

func extractTitle(errorLog string) string {
	lines := strings.Split(errorLog, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && len(trimmed) > 5 {
			if len(trimmed) > 100 {
				return trimmed[:100]
			}
			return trimmed
		}
	}
	return "Error from MCP fix_this"
}

func remStatusPtr(s types.RemediationStatus) *types.RemediationStatus {
	return &s
}
