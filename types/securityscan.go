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

// Package types defines common data structures.
package types

import (
	"github.com/harness/gitness/types/enum"
)

type (
	// ScanResult represents a security scan result for a repository.
	ScanResult struct {
		ID            int64                   `db:"ss_id"             json:"id"`
		SpaceID       int64                   `db:"ss_space_id"       json:"space_id"`
		RepoID        int64                   `db:"ss_repo_id"        json:"repo_id"`
		Identifier    string                  `db:"ss_identifier"     json:"identifier"`
		ScanType      enum.SecurityScanType   `db:"ss_scan_type"      json:"scan_type"`
		Status        enum.SecurityScanStatus `db:"ss_status"         json:"status"`
		CommitSHA     string                  `db:"ss_commit_sha"     json:"commit_sha"`
		Branch        string                  `db:"ss_branch"         json:"branch"`
		TotalIssues   int                     `db:"ss_total_issues"   json:"total_issues"`
		CriticalCount int                     `db:"ss_critical_count" json:"critical_count"`
		HighCount     int                     `db:"ss_high_count"     json:"high_count"`
		MediumCount   int                     `db:"ss_medium_count"   json:"medium_count"`
		LowCount      int                     `db:"ss_low_count"      json:"low_count"`
		Duration      int64                   `db:"ss_duration"       json:"duration"`
		TriggeredBy   enum.SecurityScanTrigger `db:"ss_triggered_by"   json:"triggered_by"`
		CreatedBy     int64                   `db:"ss_created_by"     json:"created_by"`
		Created       int64                   `db:"ss_created"        json:"created"`
		Updated       int64                   `db:"ss_updated"        json:"updated"`
		Version       int64                   `db:"ss_version"        json:"version"`
	}

	// ScanResultInput stores security scan parameters used to create or update a scan.
	ScanResultInput struct {
		ScanType    *enum.SecurityScanType   `json:"scan_type"`
		CommitSHA   *string                  `json:"commit_sha"`
		Branch      *string                  `json:"branch"`
		TriggeredBy *enum.SecurityScanTrigger `json:"triggered_by"`
	}

	// ScanResultFilter stores scan query parameters.
	ScanResultFilter struct {
		Page       int                     `json:"page"`
		Size       int                     `json:"size"`
		Sort       enum.SecurityScanAttr   `json:"sort"`
		Order      enum.Order              `json:"order"`
		Status     enum.SecurityScanStatus `json:"status"`
		ScanType   enum.SecurityScanType   `json:"scan_type"`
		TriggeredBy enum.SecurityScanTrigger `json:"triggered_by"`
	}

	// ScanFinding represents a security finding from a scan.
	ScanFinding struct {
		ID        int64                      `db:"sf_id"            json:"id"`
		ScanID    int64                      `db:"sf_scan_id"       json:"scan_id"`
		Identifier string                    `db:"sf_identifier"    json:"identifier"`
		Severity  enum.SecurityFindingSeverity `db:"sf_severity"      json:"severity"`
		Category  enum.SecurityFindingCategory `db:"sf_category"      json:"category"`
		Title     string                    `db:"sf_title"         json:"title"`
		Description string                  `db:"sf_description"   json:"description"`
		FilePath  string                    `db:"sf_file_path"     json:"file_path"`
		LineStart int                       `db:"sf_line_start"    json:"line_start"`
		LineEnd   int                       `db:"sf_line_end"      json:"line_end"`
		Rule      string                    `db:"sf_rule"          json:"rule"`
		Snippet   string                    `db:"sf_snippet"       json:"snippet"`
		Suggestion string                   `db:"sf_suggestion"    json:"suggestion"`
		Status    enum.SecurityFindingStatus `db:"sf_status"        json:"status"`
		CWE       string                    `db:"sf_cwe"           json:"cwe"`
		Created   int64                     `db:"sf_created"       json:"created"`
		Updated   int64                     `db:"sf_updated"       json:"updated"`
	}

	// ScanFindingInput stores security finding parameters used to create or update a finding.
	ScanFindingInput struct {
		Severity    *enum.SecurityFindingSeverity `json:"severity"`
		Category    *enum.SecurityFindingCategory `json:"category"`
		Title       *string                    `json:"title"`
		Description *string                    `json:"description"`
		FilePath    *string                    `json:"file_path"`
		LineStart   *int                       `json:"line_start"`
		LineEnd     *int                       `json:"line_end"`
		Rule        *string                    `json:"rule"`
		Snippet     *string                    `json:"snippet"`
		Suggestion  *string                    `json:"suggestion"`
		CWE         *string                    `json:"cwe"`
	}

	// ScanFindingFilter stores finding query parameters.
	ScanFindingFilter struct {
		Page     int                          `json:"page"`
		Size     int                          `json:"size"`
		Sort     enum.SecurityFindingAttr     `json:"sort"`
		Order    enum.Order                   `json:"order"`
		Severity enum.SecurityFindingSeverity `json:"severity"`
		Category enum.SecurityFindingCategory `json:"category"`
		Status   enum.SecurityFindingStatus   `json:"status"`
	}

	// SecuritySummary represents the aggregate security posture of a space or repo.
	SecuritySummary struct {
		SpaceID        int64 `json:"space_id"`
		RepoID         int64 `json:"repo_id"`
		LastScanID     int64 `json:"last_scan_id"`
		LastScanTime   int64 `json:"last_scan_time"`
		TotalFindings  int   `json:"total_findings"`
		CriticalIssues int   `json:"critical_issues"`
		HighIssues     int   `json:"high_issues"`
		MediumIssues   int   `json:"medium_issues"`
		LowIssues      int   `json:"low_issues"`
		InfoIssues     int   `json:"info_issues"`
	}

	// ScanFindingStatusUpdate stores the status update for a finding.
	ScanFindingStatusUpdate struct {
		Status *enum.SecurityFindingStatus `json:"status"`
	}
)
