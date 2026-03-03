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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
)

// registerPrompts registers all Tier 4 expert prompts.
func registerPrompts(s *Server) {
	s.RegisterPrompt(PromptDefinition{
		Name:        "solodev_review",
		Description: "Full code review against quality rules and security findings. Structured as: summary, rule violations, security issues, suggestions.",
		Arguments: []PromptArgument{
			{Name: "space_ref", Description: "Space reference", Required: true},
			{Name: "repo_ref", Description: "Repository reference", Required: false},
			{Name: "commit_sha", Description: "Commit to review", Required: false},
			{Name: "file_paths", Description: "Comma-separated file paths to review", Required: false},
		},
	}, s.promptReview)

	s.RegisterPrompt(PromptDefinition{
		Name:        "solodev_incident",
		Description: "Incident response runbook. Pulls active errors + remediation queue + health status, walks through mitigation steps.",
		Arguments: []PromptArgument{
			{Name: "space_ref", Description: "Space reference", Required: true},
			{Name: "service_name", Description: "Affected service name", Required: false},
		},
	}, s.promptIncident)

	s.RegisterPrompt(PromptDefinition{
		Name:        "solodev_pipeline_debug",
		Description: "CI failure post-mortem. Takes a failed pipeline log, identifies the failing step, suggests exactly what to fix.",
		Arguments: []PromptArgument{
			{Name: "pipeline_log", Description: "Failed pipeline log output", Required: true},
			{Name: "stack", Description: "Technology stack (go, node, python, etc.)", Required: false},
		},
	}, s.promptPipelineDebug)

	s.RegisterPrompt(PromptDefinition{
		Name:        "solodev_security_audit",
		Description: "Security deep-dive. Reviews code for OWASP Top 10 + secret exposure + dependency vulnerabilities.",
		Arguments: []PromptArgument{
			{Name: "space_ref", Description: "Space reference", Required: true},
			{Name: "repo_ref", Description: "Repository reference", Required: false},
		},
	}, s.promptSecurityAudit)

	s.RegisterPrompt(PromptDefinition{
		Name:        "solodev_debt_sprint",
		Description: "Tech debt prioritization. Ranks open debt items by blast radius and suggests a sprint-sized chunk to tackle.",
		Arguments: []PromptArgument{
			{Name: "space_ref", Description: "Space reference", Required: true},
			{Name: "sprint_days", Description: "Sprint length in days (default: 14)", Required: false},
		},
	}, s.promptDebtSprint)
}

// --- Prompt Handlers ---

func (s *Server) promptReview(ctx context.Context, session *auth.Session, args map[string]string) (*PromptGetResult, error) {
	spaceRef := GetSpaceRef(args)

	var contextParts []string
	contextParts = append(contextParts, "# SoloDev Code Review Context\n")

	// Fetch quality rules
	if s.controllers.QualityGate != nil {
		filter := &types.QualityRuleFilter{ListQueryFilter: types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 100}}}
		rules, err := s.controllers.QualityGate.ListRules(ctx, session, spaceRef, filter)
		if err == nil {
			b, _ := json.MarshalIndent(rules, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Active Quality Rules\n```json\n%s\n```\n", string(b)))
		}
	}

	// Fetch security findings
	if s.controllers.SecurityScan != nil {
		summary, err := s.controllers.SecurityScan.GetSecuritySummary(ctx, session, spaceRef, nil)
		if err == nil {
			b, _ := json.MarshalIndent(summary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Security Posture\n```json\n%s\n```\n", string(b)))
		}
	}

	// Fetch tech debt
	if s.controllers.TechDebt != nil {
		items, err := s.controllers.TechDebt.List(ctx, session, spaceRef, &types.TechDebtFilter{
			Severity: []string{"critical", "high"},
			Status:   []string{"open"},
			Page:     1,
			Limit:    20,
		})
		if err == nil && len(items.Items) > 0 {
			b, _ := json.MarshalIndent(items.Items, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Critical/High Tech Debt\n```json\n%s\n```\n", string(b)))
		}
	}

	systemContext := strings.Join(contextParts, "\n")

	return &PromptGetResult{
		Description: "SoloDev Code Review — quality rules, security findings, and tech debt context loaded",
		Messages: []PromptMessage{
			{
				Role: "user",
				Content: TextContent(systemContext + `

## Instructions

You are performing a code review with the context of the SoloDev quality rules and security posture above.

Review the code changes and provide:

1. **Summary** — One-paragraph overview of the changes
2. **Quality Rule Violations** — Check each active rule against the changes. List any violations.
3. **Security Issues** — Cross-reference the changes against the open security findings and OWASP Top 10
4. **Tech Debt Impact** — Note if changes touch any high-debt files or create new debt
5. **Suggestions** — Specific, actionable improvements with code examples where helpful

Be direct and concise. Focus on what matters most.`),
			},
		},
	}, nil
}

func (s *Server) promptIncident(ctx context.Context, session *auth.Session, args map[string]string) (*PromptGetResult, error) {
	spaceRef := GetSpaceRef(args)
	serviceName := args["service_name"]

	var contextParts []string
	contextParts = append(contextParts, "# SoloDev Incident Response Context\n")

	if serviceName != "" {
		contextParts = append(contextParts, fmt.Sprintf("**Affected Service:** %s\n", serviceName))
	}

	// Active errors
	if s.controllers.ErrorTracker != nil {
		errors, err := s.controllers.ErrorTracker.ListErrors(ctx, session, spaceRef, types.ErrorTrackerListOptions{})
		if err == nil {
			b, _ := json.MarshalIndent(errors, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Active Errors\n```json\n%s\n```\n", string(b)))
		}
		summary, err := s.controllers.ErrorTracker.GetSummary(ctx, session, spaceRef)
		if err == nil {
			b, _ := json.MarshalIndent(summary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Error Summary\n```json\n%s\n```\n", string(b)))
		}
	}

	// Health status
	if s.controllers.HealthCheck != nil {
		healthSummary, err := s.controllers.HealthCheck.GetSummary(ctx, session, spaceRef)
		if err == nil {
			b, _ := json.MarshalIndent(healthSummary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Health Status\n```json\n%s\n```\n", string(b)))
		}
	}

	// Remediation queue
	if s.controllers.Remediation != nil {
		remSummary, err := s.controllers.Remediation.GetSummary(ctx, session, spaceRef)
		if err == nil {
			b, _ := json.MarshalIndent(remSummary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Remediation Queue\n```json\n%s\n```\n", string(b)))
		}
	}

	systemContext := strings.Join(contextParts, "\n")

	return &PromptGetResult{
		Description: "SoloDev Incident Response — active errors, health status, and remediation context loaded",
		Messages: []PromptMessage{
			{
				Role: "user",
				Content: TextContent(systemContext + `

## Incident Response Runbook

Using the context above, walk through these steps:

1. **Impact Assessment** — What's broken? How many users/services affected? What severity?
2. **Root Cause Hypothesis** — Based on error patterns and health status, what's the likely cause?
3. **Active Remediations** — What fixes are already in progress? Do they address the issue?
4. **Immediate Mitigation** — What can be done right now to reduce impact?
5. **Fix Plan** — Step-by-step plan to resolve the root cause
6. **Prevention** — What monitoring, tests, or quality rules would prevent recurrence?

Prioritize urgency. Be specific about commands, file changes, or tool calls to make.`),
			},
		},
	}, nil
}

func (s *Server) promptPipelineDebug(_ context.Context, _ *auth.Session, args map[string]string) (*PromptGetResult, error) {
	pipelineLog := args["pipeline_log"]
	stack := args["stack"]

	var contextParts []string
	contextParts = append(contextParts, "# CI Pipeline Failure Post-Mortem\n")

	if stack != "" {
		contextParts = append(contextParts, fmt.Sprintf("**Stack:** %s\n", stack))
	}

	if pipelineLog != "" {
		contextParts = append(contextParts, fmt.Sprintf("## Failed Pipeline Log\n```\n%s\n```\n", pipelineLog))
	}

	systemContext := strings.Join(contextParts, "\n")

	return &PromptGetResult{
		Description: "SoloDev Pipeline Debug — CI failure analysis",
		Messages: []PromptMessage{
			{
				Role: "user",
				Content: TextContent(systemContext + `

## Pipeline Debug Instructions

Analyze the failed pipeline log above and provide:

1. **Failing Step** — Which exact step failed? (build, test, lint, deploy, etc.)
2. **Error Identification** — The specific error message and what it means
3. **Root Cause** — Why did it fail? (dependency issue, test failure, config error, etc.)
4. **Fix** — Exact changes needed to fix the pipeline:
   - File path + line changes
   - Commands to run locally to verify
   - Config changes if applicable
5. **Prevention** — How to prevent this type of failure in the future

Be precise. Include the exact file changes as a diff when possible.`),
			},
		},
	}, nil
}

func (s *Server) promptSecurityAudit(ctx context.Context, session *auth.Session, args map[string]string) (*PromptGetResult, error) {
	spaceRef := GetSpaceRef(args)

	var contextParts []string
	contextParts = append(contextParts, "# SoloDev Security Audit Context\n")

	// Fetch security summary
	if s.controllers.SecurityScan != nil {
		summary, err := s.controllers.SecurityScan.GetSecuritySummary(ctx, session, spaceRef, nil)
		if err == nil {
			b, _ := json.MarshalIndent(summary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Current Security Posture\n```json\n%s\n```\n", string(b)))
		}
	}

	// Fetch quality rules related to security
	if s.controllers.QualityGate != nil {
		filter := &types.QualityRuleFilter{ListQueryFilter: types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 100}}}
		rules, err := s.controllers.QualityGate.ListRules(ctx, session, spaceRef, filter)
		if err == nil {
			b, _ := json.MarshalIndent(rules, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Quality Rules (Security-Related)\n```json\n%s\n```\n", string(b)))
		}
	}

	systemContext := strings.Join(contextParts, "\n")

	return &PromptGetResult{
		Description: "SoloDev Security Audit — OWASP Top 10 + secrets + dependency analysis",
		Messages: []PromptMessage{
			{
				Role: "user",
				Content: TextContent(systemContext + `

## Security Audit Instructions

Perform a comprehensive security review:

1. **OWASP Top 10 Assessment** — Check for injection, broken auth, XSS, SSRF, etc.
2. **Secret Exposure** — Identify any hardcoded secrets, API keys, tokens, or credentials
3. **Dependency Vulnerabilities** — Flag known CVEs in dependencies
4. **Configuration Security** — Review security headers, CORS, CSP, TLS settings
5. **Data Protection** — Check encryption at rest/in transit, PII handling
6. **Access Control** — Review authentication and authorization patterns
7. **Recommendations** — Prioritized list of security improvements with severity ratings

Cross-reference all findings against the existing security scan results above.
Flag any NEW issues not already captured by the scanner.`),
			},
		},
	}, nil
}

func (s *Server) promptDebtSprint(ctx context.Context, session *auth.Session, args map[string]string) (*PromptGetResult, error) {
	spaceRef := GetSpaceRef(args)
	sprintDays := args["sprint_days"]
	if sprintDays == "" {
		sprintDays = "14"
	}

	var contextParts []string
	contextParts = append(contextParts, "# SoloDev Tech Debt Sprint Planning\n")
	contextParts = append(contextParts, fmt.Sprintf("**Sprint Length:** %s days\n", sprintDays))

	// Fetch all tech debt
	if s.controllers.TechDebt != nil {
		items, err := s.controllers.TechDebt.List(ctx, session, spaceRef, &types.TechDebtFilter{
			Status: []string{"open"},
			Page:   1,
			Limit:  100,
		})
		if err == nil {
			b, _ := json.MarshalIndent(items, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Open Tech Debt Items\n```json\n%s\n```\n", string(b)))
		}

		summary, err := s.controllers.TechDebt.Summary(ctx, session, spaceRef, &types.TechDebtFilter{})
		if err == nil {
			b, _ := json.MarshalIndent(summary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Tech Debt Summary\n```json\n%s\n```\n", string(b)))
		}
	}

	// Fetch error hotspots for cross-reference
	if s.controllers.ErrorTracker != nil {
		summary, err := s.controllers.ErrorTracker.GetSummary(ctx, session, spaceRef)
		if err == nil {
			b, _ := json.MarshalIndent(summary, "", "  ")
			contextParts = append(contextParts, fmt.Sprintf("## Error Tracker Summary (for blast radius)\n```json\n%s\n```\n", string(b)))
		}
	}

	systemContext := strings.Join(contextParts, "\n")

	return &PromptGetResult{
		Description: "SoloDev Debt Sprint — prioritized tech debt backlog for sprint planning",
		Messages: []PromptMessage{
			{
				Role: "user",
				Content: TextContent(systemContext + fmt.Sprintf(`

## Tech Debt Sprint Planning Instructions

Using the tech debt backlog above, plan a %s-day sprint:

1. **Blast Radius Ranking** — Rank each debt item by:
   - Frequency of related errors (cross-reference with error data)
   - Number of dependent modules
   - Security implications
   - Developer friction (how much it slows development)

2. **Sprint Scope** — Select a realistic subset that fits in %s days:
   - Mix of quick wins (< 2 hours) and deeper fixes
   - Group related items that share code paths
   - Ensure no single item takes more than 3 days

3. **Execution Order** — Sequence the work:
   - Dependencies between items
   - Quick wins first for momentum
   - Pair risky changes with tests

4. **Success Metrics** — How to measure improvement:
   - Error count reduction targets
   - Code quality score improvements
   - Performance benchmarks

Be specific about effort estimates and provide a day-by-day breakdown.`, sprintDays, sprintDays)),
			},
		},
	}, nil
}
