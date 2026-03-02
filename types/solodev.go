// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package types

type SecurityFindingRemediationMode string

const (
	SecurityFindingRemediationModeManual           SecurityFindingRemediationMode = "manual"
	SecurityFindingRemediationModeCriticalHighAuto SecurityFindingRemediationMode = "critical_high_auto"
	SecurityFindingRemediationModeAllAuto          SecurityFindingRemediationMode = "all_auto"
)

type MCPCatalogCounts struct {
	ActiveTools         int `json:"active_tools"`
	ActiveResources     int `json:"active_resources"`
	ActivePrompts       int `json:"active_prompts"`
	BlockedTools        int `json:"blocked_tools"`
	BlockedResources    int `json:"blocked_resources"`
	BlockedPrompts      int `json:"blocked_prompts"`
	ComingSoonTools     int `json:"coming_soon_tools"`
	ComingSoonResources int `json:"coming_soon_resources"`
	ComingSoonPrompts   int `json:"coming_soon_prompts"`
}

type MCPCatalogItem struct {
	Surface     string   `json:"surface"`
	Name        string   `json:"name,omitempty"`
	URI         string   `json:"uri,omitempty"`
	Domain      string   `json:"domain"`
	Description string   `json:"description"`
	Requires    []string `json:"requires,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

type MCPCatalogSection struct {
	Tools     []MCPCatalogItem `json:"tools"`
	Resources []MCPCatalogItem `json:"resources"`
	Prompts   []MCPCatalogItem `json:"prompts"`
}

type MCPCatalog struct {
	ServerName      string            `json:"server_name"`
	ProtocolVersion string            `json:"protocol_version"`
	Counts          MCPCatalogCounts  `json:"counts"`
	Active          MCPCatalogSection `json:"active"`
	Blocked         MCPCatalogSection `json:"blocked"`
	ComingSoon      MCPCatalogSection `json:"coming_soon"`
}

type SoloDevOverviewDomainStatus struct {
	Availability string `json:"availability"`
}

type SoloDevSecurityOverview struct {
	Availability     string `json:"availability"`
	LatestScanStatus string `json:"latest_scan_status"`
	OpenFindings     int    `json:"open_findings"`
	Critical         int    `json:"critical"`
	High             int    `json:"high"`
	Medium           int    `json:"medium"`
	Low              int    `json:"low"`
	LastScanTime     int64  `json:"last_scan_time"`
}

type SoloDevRemediationOverview struct {
	Availability string `json:"availability"`
	Pending      int64  `json:"pending"`
	Processing   int64  `json:"processing"`
	Completed    int64  `json:"completed"`
	Failed       int64  `json:"failed"`
}

type SoloDevErrorsOverview struct {
	Availability string `json:"availability"`
	Open         int64  `json:"open"`
	Fatal        int64  `json:"fatal"`
	Warning      int64  `json:"warning"`
	LastSeen     int64  `json:"last_seen"`
}

type SoloDevMCPOverview struct {
	Tools     int `json:"tools"`
	Resources int `json:"resources"`
	Prompts   int `json:"prompts"`
}

type SoloDevOverview struct {
	SpaceRef        string                     `json:"space_ref"`
	UpdatedAt       int64                      `json:"updated_at"`
	Security        SoloDevSecurityOverview    `json:"security"`
	Remediation     SoloDevRemediationOverview `json:"remediation"`
	Errors          SoloDevErrorsOverview      `json:"errors"`
	MCP             SoloDevMCPOverview         `json:"mcp"`
	DeferredDomains []string                   `json:"deferred_domains"`
}
