// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/api/controller/qualitygate"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// registerAtomicTools registers all Tier 1 tools (direct controller wrappers).
func registerAtomicTools(s *Server) {
	// --- Auto Pipeline ---
	s.RegisterTool(ToolDefinition{
		Name:        "pipeline_generate",
		Description: "Detect technology stack from file paths and generate a complete CI/CD pipeline YAML. Supports Go, Node, Python, Rust, Java, and more.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference (defaults to SOLODEV_SPACE env)"),
			"files":     ArrayOfStrings("List of file paths in the repository to analyze"),
		}, []string{"files"}),
	}, s.toolPipelineGenerate)

	// --- Security Scan ---
	s.RegisterTool(ToolDefinition{
		Name:        "security_scan",
		Description: "Trigger a SAST + SCA + secret detection security scan on a repository. Returns scan ID for tracking results.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"repo_ref":   StringProp("Repository reference"),
			"scan_type":  StringProp("Type of scan: sast, sca, secret, or full"),
			"commit_sha": StringProp("Commit SHA to scan"),
			"branch":     StringProp("Branch name"),
		}, []string{"repo_ref"}),
	}, s.toolSecurityScan)

	s.RegisterTool(ToolDefinition{
		Name:        "security_findings",
		Description: "List findings from a completed security scan with severity breakdown.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":       StringProp("Space reference"),
			"repo_ref":        StringProp("Repository reference"),
			"scan_identifier": StringProp("Scan identifier to get findings for"),
			"severity":        StringProp("Filter by severity: critical, high, medium, low"),
		}, []string{"repo_ref", "scan_identifier"}),
	}, s.toolSecurityFindings)

	s.RegisterTool(ToolDefinition{
		Name:        "security_fix_finding",
		Description: "Create or reuse an AI remediation task for a specific security finding.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":       StringProp("Space reference"),
			"repo_ref":        StringProp("Repository reference"),
			"scan_identifier": StringProp("Scan identifier"),
			"finding_id":      IntProp("Security finding identifier"),
		}, []string{"space_ref", "repo_ref", "scan_identifier", "finding_id"}),
	}, s.toolSecurityFixFinding)

	// --- Quality Gate ---
	s.RegisterTool(ToolDefinition{
		Name:        "quality_evaluate",
		Description: "Evaluate a commit/branch against all configured quality rules. Returns pass/fail verdict with rule-by-rule results.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"repo_ref":   StringProp("Repository reference for evaluation"),
			"commit_sha": StringProp("Commit SHA to evaluate"),
			"branch":     StringProp("Branch name"),
		}, []string{"space_ref", "repo_ref", "commit_sha"}),
	}, s.toolQualityEvaluate)

	s.RegisterTool(ToolDefinition{
		Name:        "quality_rules_list",
		Description: "List all active quality rules configured in a space.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference"),
		}, []string{"space_ref"}),
	}, s.toolQualityRulesList)

	s.RegisterTool(ToolDefinition{
		Name:        "quality_summary",
		Description: "Get aggregate quality health across all repositories in a space.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference"),
		}, []string{"space_ref"}),
	}, s.toolQualitySummary)

	// --- Error Tracker ---
	s.RegisterTool(ToolDefinition{
		Name:        "error_report",
		Description: "Report a runtime error. Automatically triggers AI remediation via Error Bridge for fatal/error severity.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":     StringProp("Space reference"),
			"identifier":    StringProp("Unique error identifier"),
			"title":         StringProp("Error title/name"),
			"message":       StringProp("Error message"),
			"severity":      StringProp("Severity: fatal, error, warning"),
			"stack_trace":   StringProp("Full stack trace"),
			"file_path":     StringProp("Source file path"),
			"line_number":   IntProp("Line number"),
			"function_name": StringProp("Function/method name"),
			"language":      StringProp("Programming language"),
			"environment":   StringProp("Runtime environment (production, staging, etc.)"),
		}, []string{"space_ref", "title", "message", "stack_trace"}),
	}, s.toolErrorReport)

	s.RegisterTool(ToolDefinition{
		Name:        "error_list",
		Description: "List active error groups with occurrence counts.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference"),
			"status":    StringProp("Filter: open, resolved, ignored, regressed"),
			"severity":  StringProp("Filter: fatal, error, warning"),
		}, []string{"space_ref"}),
	}, s.toolErrorList)

	s.RegisterTool(ToolDefinition{
		Name:        "error_get",
		Description: "Get full error detail including stack trace, occurrences, and assigned user.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"identifier": StringProp("Error group identifier"),
		}, []string{"space_ref", "identifier"}),
	}, s.toolErrorGet)

	// --- AI Remediation ---
	s.RegisterTool(ToolDefinition{
		Name:        "remediation_trigger",
		Description: "Manually create an AI remediation task with full error context. The system will analyze the error and generate a patch.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":      StringProp("Space reference"),
			"title":          StringProp("Remediation title"),
			"description":    StringProp("Description of the issue"),
			"trigger_source": StringProp("Source: manual, error_tracker, pipeline, security_scan, quality_gate"),
			"branch":         StringProp("Branch where the issue exists"),
			"commit_sha":     StringProp("Commit SHA"),
			"error_log":      StringProp("Error log / stack trace"),
			"source_code":    StringProp("Relevant source code snippet"),
			"file_path":      StringProp("File path of the issue"),
		}, []string{"space_ref", "title", "error_log", "branch"}),
	}, s.toolRemediationTrigger)

	s.RegisterTool(ToolDefinition{
		Name:        "remediation_list",
		Description: "Browse remediations by status (pending, processing, completed, failed, applied, dismissed).",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference"),
			"status":    StringProp("Filter by status"),
		}, []string{"space_ref"}),
	}, s.toolRemediationList)

	s.RegisterTool(ToolDefinition{
		Name:        "remediation_get",
		Description: "Get full remediation detail including AI response, patch diff, and fix branch.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"identifier": StringProp("Remediation identifier"),
		}, []string{"space_ref", "identifier"}),
	}, s.toolRemediationGet)

	s.RegisterTool(ToolDefinition{
		Name:        "remediation_apply",
		Description: "Apply a completed remediation diff onto a fix branch and open a draft PR.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"identifier": StringProp("Remediation identifier"),
		}, []string{"space_ref", "identifier"}),
	}, s.toolRemediationApply)

	s.RegisterTool(ToolDefinition{
		Name:        "remediation_update",
		Description: "Push AI-generated patch diff, fix branch, and PR link back to a remediation task.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":   StringProp("Space reference"),
			"identifier":  StringProp("Remediation identifier"),
			"status":      StringProp("New status: processing, completed, failed, applied, dismissed"),
			"patch_diff":  StringProp("Unified diff patch"),
			"ai_response": StringProp("AI analysis response text"),
			"fix_branch":  StringProp("Branch name for the fix"),
			"pr_link":     StringProp("Pull request URL"),
			"confidence":  NumberProp("AI confidence score 0.0-1.0"),
		}, []string{"space_ref", "identifier"}),
	}, s.toolRemediationUpdate)

	// --- Health Check ---
	s.RegisterTool(ToolDefinition{
		Name:        "health_summary",
		Description: "Get endpoint uptime summary across all monitored services in a space.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference"),
		}, []string{"space_ref"}),
	}, s.toolHealthSummary)

	// --- Feature Flag ---
	s.RegisterTool(ToolDefinition{
		Name:        "feature_flag_toggle",
		Description: "Toggle a feature flag on or off live from the AI conversation.",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref":  StringProp("Space reference"),
			"identifier": StringProp("Feature flag identifier"),
			"enabled":    BoolProp("True to enable, false to disable"),
		}, []string{"space_ref", "identifier", "enabled"}),
	}, s.toolFeatureFlagToggle)

	// --- Tech Debt ---
	s.RegisterTool(ToolDefinition{
		Name:        "tech_debt_list",
		Description: "Browse tech debt items by severity (critical, high, medium, low).",
		InputSchema: ObjectSchema(map[string]*Schema{
			"space_ref": StringProp("Space reference"),
			"severity":  StringProp("Filter by severity"),
			"status":    StringProp("Filter by status: open, in_progress, resolved, accepted"),
		}, []string{"space_ref"}),
	}, s.toolTechDebtList)
}

// --- Tool Handlers ---

func (s *Server) toolPipelineGenerate(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string   `json:"space_ref"`
		Files    []string `json:"files"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.AutoPipeline == nil {
		return ErrorResult("auto-pipeline module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.AutoPipeline.GenerateAutoConfig(ctx, session, spaceRef, p.Files)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolSecurityScan(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef  string `json:"space_ref"`
		RepoRef   string `json:"repo_ref"`
		ScanType  string `json:"scan_type"`
		CommitSHA string `json:"commit_sha"`
		Branch    string `json:"branch"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.SecurityScan == nil {
		return ErrorResult("security scan module not available"), nil
	}
	scanType := enum.SecurityScanType(p.ScanType)
	trigger := enum.SecurityScanTriggerManual
	in := &types.ScanResultInput{
		ScanType:    &scanType,
		CommitSHA:   strPtr(p.CommitSHA),
		Branch:      strPtr(p.Branch),
		TriggeredBy: &trigger,
	}
	result, err := s.controllers.SecurityScan.TriggerScan(ctx, session, p.RepoRef, in)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolSecurityFindings(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef       string `json:"space_ref"`
		RepoRef        string `json:"repo_ref"`
		ScanIdentifier string `json:"scan_identifier"`
		Severity       string `json:"severity"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.SecurityScan == nil {
		return ErrorResult("security scan module not available"), nil
	}
	filter := &types.ScanFindingFilter{
		Page:     1,
		Size:     100,
		Severity: enum.SecurityFindingSeverity(p.Severity),
	}
	findings, count, err := s.controllers.SecurityScan.ListFindings(ctx, session, p.RepoRef, p.ScanIdentifier, filter)
	if err != nil {
		return nil, err
	}
	return SuccessResult(map[string]interface{}{
		"findings":    findings,
		"total_count": count,
	})
}

func (s *Server) toolSecurityFixFinding(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef       string `json:"space_ref"`
		RepoRef        string `json:"repo_ref"`
		ScanIdentifier string `json:"scan_identifier"`
		FindingID      int64  `json:"finding_id"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.Remediation == nil {
		return ErrorResult("remediation module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, _, err := s.controllers.Remediation.TriggerRemediationFromSecurityFinding(
		ctx,
		session,
		spaceRef,
		&types.CreateRemediationFromSecurityFindingInput{
			RepoRef:        p.RepoRef,
			ScanIdentifier: p.ScanIdentifier,
			FindingID:      p.FindingID,
		},
	)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolQualityEvaluate(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef  string `json:"space_ref"`
		RepoRef   string `json:"repo_ref"`
		CommitSHA string `json:"commit_sha"`
		Branch    string `json:"branch"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.QualityGate == nil {
		return ErrorResult("quality gate module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	trigger := enum.QualityTriggerManual
	in := &qualitygate.EvaluateInput{
		RepoRef:   p.RepoRef,
		CommitSHA: p.CommitSHA,
		Branch:    p.Branch,
		Trigger:   trigger,
	}
	result, err := s.controllers.QualityGate.Evaluate(ctx, session, spaceRef, in)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolQualityRulesList(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string `json:"space_ref"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.QualityGate == nil {
		return ErrorResult("quality gate module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	filter := &types.QualityRuleFilter{
		ListQueryFilter: types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 100}},
	}
	result, err := s.controllers.QualityGate.ListRules(ctx, session, spaceRef, filter)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolQualitySummary(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string `json:"space_ref"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.QualityGate == nil {
		return ErrorResult("quality gate module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.QualityGate.GetSummary(ctx, session, spaceRef)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolErrorReport(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef     string          `json:"space_ref"`
		Identifier   string          `json:"identifier"`
		Title        string          `json:"title"`
		Message      string          `json:"message"`
		Severity     string          `json:"severity"`
		StackTrace   string          `json:"stack_trace"`
		FilePath     string          `json:"file_path"`
		LineNumber   int             `json:"line_number"`
		FunctionName string          `json:"function_name"`
		Language     string          `json:"language"`
		Environment  string          `json:"environment"`
		Metadata     json.RawMessage `json:"metadata"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.ErrorTracker == nil {
		return ErrorResult("error tracker module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	severity := types.ErrorSeverityError
	if p.Severity != "" {
		severity = types.ErrorSeverity(p.Severity)
	}
	identifier := p.Identifier
	if identifier == "" {
		identifier = fmt.Sprintf("mcp-%d", types.NowMillis())
	}
	in := &types.ReportErrorInput{
		Identifier:   identifier,
		Title:        p.Title,
		Message:      p.Message,
		Severity:     severity,
		StackTrace:   p.StackTrace,
		FilePath:     p.FilePath,
		LineNumber:   p.LineNumber,
		FunctionName: p.FunctionName,
		Language:     p.Language,
		Environment:  p.Environment,
		Metadata:     p.Metadata,
	}
	result, err := s.controllers.ErrorTracker.ReportError(ctx, session, spaceRef, in)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolErrorList(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string `json:"space_ref"`
		Status   string `json:"status"`
		Severity string `json:"severity"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.ErrorTracker == nil {
		return ErrorResult("error tracker module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	opts := types.ErrorTrackerListOptions{}
	if p.Status != "" {
		status := types.ErrorGroupStatus(p.Status)
		opts.Status = &status
	}
	if p.Severity != "" {
		sev := types.ErrorSeverity(p.Severity)
		opts.Severity = &sev
	}
	result, err := s.controllers.ErrorTracker.ListErrors(ctx, session, spaceRef, opts)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolErrorGet(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef   string `json:"space_ref"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.ErrorTracker == nil {
		return ErrorResult("error tracker module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.ErrorTracker.GetError(ctx, session, spaceRef, p.Identifier)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolRemediationTrigger(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef      string `json:"space_ref"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		TriggerSource string `json:"trigger_source"`
		Branch        string `json:"branch"`
		CommitSHA     string `json:"commit_sha"`
		ErrorLog      string `json:"error_log"`
		SourceCode    string `json:"source_code"`
		FilePath      string `json:"file_path"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.Remediation == nil {
		return ErrorResult("remediation module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	triggerSource := types.RemediationTriggerManual
	if p.TriggerSource != "" {
		triggerSource = types.RemediationTriggerSource(p.TriggerSource)
	}
	in := &types.TriggerRemediationInput{
		Title:         p.Title,
		Description:   p.Description,
		TriggerSource: triggerSource,
		Branch:        p.Branch,
		CommitSHA:     p.CommitSHA,
		ErrorLog:      p.ErrorLog,
		SourceCode:    p.SourceCode,
		FilePath:      p.FilePath,
	}
	result, err := s.controllers.Remediation.TriggerRemediation(ctx, session, spaceRef, in)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolRemediationList(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string `json:"space_ref"`
		Status   string `json:"status"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.Remediation == nil {
		return ErrorResult("remediation module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	filter := &types.RemediationListFilter{}
	if p.Status != "" {
		status := types.RemediationStatus(p.Status)
		filter.Status = &status
	}
	result, err := s.controllers.Remediation.ListRemediations(ctx, session, spaceRef, filter)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolRemediationGet(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef   string `json:"space_ref"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.Remediation == nil {
		return ErrorResult("remediation module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.Remediation.GetRemediation(ctx, session, spaceRef, p.Identifier)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolRemediationApply(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef   string `json:"space_ref"`
		Identifier string `json:"identifier"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.Remediation == nil {
		return ErrorResult("remediation module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.Remediation.ApplyRemediation(ctx, session, spaceRef, p.Identifier)
	if err != nil {
		return nil, err
	}
	return SuccessResult(map[string]any{
		"identifier": result.Identifier,
		"status":     result.Status,
		"fix_branch": result.FixBranch,
		"pr_link":    result.PRLink,
	})
}

func (s *Server) toolRemediationUpdate(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef   string   `json:"space_ref"`
		Identifier string   `json:"identifier"`
		Status     string   `json:"status"`
		PatchDiff  string   `json:"patch_diff"`
		AIResponse string   `json:"ai_response"`
		FixBranch  string   `json:"fix_branch"`
		PRLink     string   `json:"pr_link"`
		Confidence *float64 `json:"confidence"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.Remediation == nil {
		return ErrorResult("remediation module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	in := &types.UpdateRemediationInput{
		PatchDiff:  p.PatchDiff,
		AIResponse: p.AIResponse,
		FixBranch:  p.FixBranch,
		PRLink:     p.PRLink,
		Confidence: p.Confidence,
	}
	if p.Status != "" {
		status := types.RemediationStatus(p.Status)
		in.Status = &status
	}
	result, err := s.controllers.Remediation.UpdateRemediation(ctx, session, spaceRef, p.Identifier, in)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolHealthSummary(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string `json:"space_ref"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.HealthCheck == nil {
		return ErrorResult("health check module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.HealthCheck.GetSummary(ctx, session, spaceRef)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolFeatureFlagToggle(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef   string `json:"space_ref"`
		Identifier string `json:"identifier"`
		Enabled    bool   `json:"enabled"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.FeatureFlag == nil {
		return ErrorResult("feature flag module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	result, err := s.controllers.FeatureFlag.Toggle(ctx, session, spaceRef, p.Identifier, p.Enabled)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

func (s *Server) toolTechDebtList(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error) {
	var p struct {
		SpaceRef string `json:"space_ref"`
		Severity string `json:"severity"`
		Status   string `json:"status"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return ErrorResult("invalid arguments: " + err.Error()), nil
	}
	if s.controllers.TechDebt == nil {
		return ErrorResult("tech debt module not available"), nil
	}
	spaceRef := GetSpaceRef(map[string]string{"space_ref": p.SpaceRef})
	filter := &types.TechDebtFilter{
		Page:  1,
		Limit: 100,
	}
	if p.Severity != "" {
		filter.Severity = []string{p.Severity}
	}
	if p.Status != "" {
		filter.Status = []string{p.Status}
	}
	result, err := s.controllers.TechDebt.List(ctx, session, spaceRef, filter)
	if err != nil {
		return nil, err
	}
	return SuccessResult(result)
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
