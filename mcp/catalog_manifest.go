// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"fmt"
	"strings"

	"github.com/harness/gitness/types"
)

const (
	catalogSurfaceTool     = "tool"
	catalogSurfaceResource = "resource"
	catalogSurfacePrompt   = "prompt"

	catalogStatusActive     = "active"
	catalogStatusComingSoon = "coming_soon"
)

type catalogManifestItem struct {
	Surface     string
	Name        string
	URI         string
	Domain      string
	Description string
	Status      string
	Requires    []string
	Notes       string
}

func catalogManifest() []catalogManifestItem {
	return []catalogManifestItem{
		{Surface: catalogSurfaceTool, Name: "error_report", Domain: "Errors", Description: "Report a runtime error.", Status: catalogStatusActive, Requires: []string{"error_tracker"}},
		{Surface: catalogSurfaceTool, Name: "error_list", Domain: "Errors", Description: "List tracked runtime errors.", Status: catalogStatusActive, Requires: []string{"error_tracker"}},
		{Surface: catalogSurfaceTool, Name: "error_get", Domain: "Errors", Description: "Get error details and context.", Status: catalogStatusActive, Requires: []string{"error_tracker"}},
		{Surface: catalogSurfaceTool, Name: "security_scan", Domain: "Security", Description: "Trigger a security scan.", Status: catalogStatusActive, Requires: []string{"scanner"}},
		{Surface: catalogSurfaceTool, Name: "security_findings", Domain: "Security", Description: "List findings from a security scan.", Status: catalogStatusActive, Requires: []string{"scanner"}},
		{Surface: catalogSurfaceTool, Name: "security_fix_finding", Domain: "Security", Description: "Create or reuse an AI fix task for a security finding.", Status: catalogStatusActive, Requires: []string{"scanner", "ai_remediation"}},
		{Surface: catalogSurfaceTool, Name: "remediation_trigger", Domain: "Remediation", Description: "Create an AI remediation task.", Status: catalogStatusActive, Requires: []string{"ai_remediation"}},
		{Surface: catalogSurfaceTool, Name: "remediation_list", Domain: "Remediation", Description: "List AI remediation tasks.", Status: catalogStatusActive, Requires: []string{"ai_remediation"}},
		{Surface: catalogSurfaceTool, Name: "remediation_get", Domain: "Remediation", Description: "Get remediation details.", Status: catalogStatusActive, Requires: []string{"ai_remediation"}},
		{Surface: catalogSurfaceTool, Name: "remediation_apply", Domain: "Remediation", Description: "Apply a completed remediation into a draft PR.", Status: catalogStatusActive, Requires: []string{"ai_remediation"}},
		{Surface: catalogSurfaceTool, Name: "remediation_update", Domain: "Remediation", Description: "Update remediation output/state.", Status: catalogStatusActive, Requires: []string{"ai_remediation"}},
		{Surface: catalogSurfaceTool, Name: "fix_this", Domain: "Remediation", Description: "Report an error and wait for a patch diff.", Status: catalogStatusActive, Requires: []string{"ai_remediation", "error_tracker"}},
		{Surface: catalogSurfaceTool, Name: "pipeline_generate", Domain: "Pipelines", Description: "Generate CI/CD config.", Status: catalogStatusComingSoon, Notes: "Waiting on the onboarding and pipeline path to be completed end-to-end."},
		{Surface: catalogSurfaceTool, Name: "quality_evaluate", Domain: "Quality", Description: "Evaluate quality rules against a revision.", Status: catalogStatusComingSoon, Notes: "Waiting on fully backed quality rule and evaluation data."},
		{Surface: catalogSurfaceTool, Name: "quality_rules_list", Domain: "Quality", Description: "List quality rules.", Status: catalogStatusComingSoon, Notes: "Waiting on fully backed quality rule and evaluation data."},
		{Surface: catalogSurfaceTool, Name: "quality_summary", Domain: "Quality", Description: "Get quality summary.", Status: catalogStatusComingSoon, Notes: "Waiting on fully backed quality rule and evaluation data."},
		{Surface: catalogSurfaceTool, Name: "health_summary", Domain: "Health", Description: "Get health monitoring summary.", Status: catalogStatusComingSoon, Notes: "Waiting on the health check execution engine."},
		{Surface: catalogSurfaceTool, Name: "feature_flag_toggle", Domain: "Feature Flags", Description: "Toggle a feature flag.", Status: catalogStatusComingSoon, Notes: "Waiting on the feature flag path to be validated end-to-end."},
		{Surface: catalogSurfaceTool, Name: "tech_debt_list", Domain: "Tech Debt", Description: "List tech debt hotspots.", Status: catalogStatusComingSoon, Notes: "Waiting on the tech debt workflow to be completed end-to-end."},
		{Surface: catalogSurfaceTool, Name: "pr_ready", Domain: "Pipelines", Description: "Run the pre-merge readiness workflow.", Status: catalogStatusComingSoon, Notes: "Waiting on multiple backend domains to be fully real."},
		{Surface: catalogSurfaceTool, Name: "pipeline_validate", Domain: "Pipelines", Description: "Validate a generated pipeline.", Status: catalogStatusComingSoon, Notes: "Waiting on the pipeline authoring workflow to be completed end-to-end."},
		{Surface: catalogSurfaceTool, Name: "onboard_repo", Domain: "Pipelines", Description: "Bootstrap repo automation.", Status: catalogStatusComingSoon, Notes: "Waiting on onboarding and quality/health integration."},
		{Surface: catalogSurfaceTool, Name: "incident_triage", Domain: "Health", Description: "Aggregate incident context.", Status: catalogStatusComingSoon, Notes: "Waiting on health and incident context to be fully backed."},
		{Surface: catalogSurfaceResource, URI: "solodev://errors/active", Domain: "Errors", Description: "Open error groups.", Status: catalogStatusActive, Requires: []string{"error_tracker"}},
		{Surface: catalogSurfaceResource, URI: "solodev://security/open-findings", Domain: "Security", Description: "Open findings from the latest completed scans.", Status: catalogStatusActive, Requires: []string{"scanner"}},
		{Surface: catalogSurfaceResource, URI: "solodev://remediations/pending", Domain: "Remediation", Description: "Pending AI remediation tasks.", Status: catalogStatusActive, Requires: []string{"ai_remediation"}},
		{Surface: catalogSurfaceResource, URI: "solodev://quality/rules", Domain: "Quality", Description: "Quality rule catalog.", Status: catalogStatusComingSoon, Notes: "Waiting on fully backed quality rule data."},
		{Surface: catalogSurfaceResource, URI: "solodev://quality/summary", Domain: "Quality", Description: "Quality summary.", Status: catalogStatusComingSoon, Notes: "Waiting on fully backed quality rule data."},
		{Surface: catalogSurfaceResource, URI: "solodev://health/status", Domain: "Health", Description: "Health check status.", Status: catalogStatusComingSoon, Notes: "Waiting on the health check execution engine."},
		{Surface: catalogSurfaceResource, URI: "solodev://tech-debt/hotspots", Domain: "Tech Debt", Description: "Tech debt hotspots.", Status: catalogStatusComingSoon, Notes: "Waiting on the tech debt workflow to be completed end-to-end."},
		{Surface: catalogSurfacePrompt, Name: "solodev_review", Domain: "Prompts", Description: "Repository review workflow.", Status: catalogStatusComingSoon, Notes: "Prompts stay hidden until their backing domains are fully real."},
		{Surface: catalogSurfacePrompt, Name: "solodev_incident", Domain: "Prompts", Description: "Incident response workflow.", Status: catalogStatusComingSoon, Notes: "Prompts stay hidden until their backing domains are fully real."},
		{Surface: catalogSurfacePrompt, Name: "solodev_pipeline_debug", Domain: "Prompts", Description: "Pipeline debugging workflow.", Status: catalogStatusComingSoon, Notes: "Prompts stay hidden until their backing domains are fully real."},
		{Surface: catalogSurfacePrompt, Name: "solodev_security_audit", Domain: "Prompts", Description: "Security review workflow.", Status: catalogStatusComingSoon, Notes: "Prompts stay hidden until their backing domains are fully real."},
		{Surface: catalogSurfacePrompt, Name: "solodev_debt_sprint", Domain: "Prompts", Description: "Tech debt planning workflow.", Status: catalogStatusComingSoon, Notes: "Prompts stay hidden until their backing domains are fully real."},
	}
}

func BuildCatalog(controllers *Controllers) *types.MCPCatalog {
	catalog := &types.MCPCatalog{
		ServerName:      ServerName,
		ProtocolVersion: ProtocolVersion,
	}

	for _, item := range catalogManifest() {
		apiItem := toCatalogItem(item)

		switch item.Status {
		case catalogStatusComingSoon:
			appendCatalogItem(&catalog.ComingSoon, apiItem)
		case catalogStatusActive:
			missing := missingRequirements(item.Requires, controllers)
			if len(missing) == 0 {
				appendCatalogItem(&catalog.Active, apiItem)
			} else {
				if apiItem.Notes == "" {
					apiItem.Notes = fmt.Sprintf("Missing runtime requirements: %s.", strings.Join(missing, ", "))
				}
				appendCatalogItem(&catalog.Blocked, apiItem)
			}
		}
	}

	catalog.Counts = types.MCPCatalogCounts{
		ActiveTools:         len(catalog.Active.Tools),
		ActiveResources:     len(catalog.Active.Resources),
		ActivePrompts:       len(catalog.Active.Prompts),
		BlockedTools:        len(catalog.Blocked.Tools),
		BlockedResources:    len(catalog.Blocked.Resources),
		BlockedPrompts:      len(catalog.Blocked.Prompts),
		ComingSoonTools:     len(catalog.ComingSoon.Tools),
		ComingSoonResources: len(catalog.ComingSoon.Resources),
		ComingSoonPrompts:   len(catalog.ComingSoon.Prompts),
	}

	return catalog
}

func catalogHasActiveTool(catalog *types.MCPCatalog, name string) bool {
	for _, item := range catalog.Active.Tools {
		if item.Name == name {
			return true
		}
	}
	return false
}

func catalogHasActiveResource(catalog *types.MCPCatalog, uri string) bool {
	for _, item := range catalog.Active.Resources {
		if item.URI == uri {
			return true
		}
	}
	return false
}

func catalogHasActivePrompt(catalog *types.MCPCatalog, name string) bool {
	for _, item := range catalog.Active.Prompts {
		if item.Name == name {
			return true
		}
	}
	return false
}

func appendCatalogItem(section *types.MCPCatalogSection, item types.MCPCatalogItem) {
	switch item.Surface {
	case catalogSurfaceTool:
		section.Tools = append(section.Tools, item)
	case catalogSurfaceResource:
		section.Resources = append(section.Resources, item)
	case catalogSurfacePrompt:
		section.Prompts = append(section.Prompts, item)
	}
}

func toCatalogItem(item catalogManifestItem) types.MCPCatalogItem {
	return types.MCPCatalogItem{
		Surface:     item.Surface,
		Name:        item.Name,
		URI:         item.URI,
		Domain:      item.Domain,
		Description: item.Description,
		Requires:    item.Requires,
		Notes:       item.Notes,
	}
}

func missingRequirements(requirements []string, controllers *Controllers) []string {
	missing := make([]string, 0)
	for _, requirement := range requirements {
		if !hasRequirement(requirement, controllers) {
			missing = append(missing, requirement)
		}
	}
	return missing
}

func hasRequirement(requirement string, controllers *Controllers) bool {
	switch requirement {
	case "scanner":
		return controllers != nil &&
			controllers.SecurityScan != nil &&
			controllers.SecurityScan.ScannerStatus().Ready
	case "ai_remediation":
		return controllers != nil &&
			controllers.Remediation != nil &&
			controllers.Remediation.AIAvailable()
	case "error_tracker":
		return controllers != nil && controllers.ErrorTracker != nil
	default:
		return false
	}
}
