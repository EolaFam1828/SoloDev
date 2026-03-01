// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Added ListResponse type and updated TechDebt struct.
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

package types

import (
	"database/sql/driver"
	"encoding/json"
)

type TechDebt struct {
	ID               int64            `json:"id"`
	SpaceID          int64            `json:"space_id"`
	RepoID           int64            `json:"repo_id"`
	Identifier       string           `json:"identifier"`
	Title            string           `json:"title"`
	Description      string           `json:"description,omitempty"`
	Severity         TechDebtSeverity `json:"severity"`
	Status           TechDebtStatus   `json:"status"`
	Category         TechDebtCategory `json:"category"`
	FilePath         string           `json:"file_path,omitempty"`
	LineStart        int              `json:"line_start,omitempty"`
	LineEnd          int              `json:"line_end,omitempty"`
	EstimatedEffort  string           `json:"estimated_effort"`
	Tags             []string         `json:"tags,omitempty"`
	DueDate          int64            `json:"due_date,omitempty"`
	ResolvedAt       int64            `json:"resolved_at,omitempty"`
	ResolvedBy       int64            `json:"resolved_by,omitempty"`
	CreatedBy        int64            `json:"created_by"`
	Created          int64            `json:"created"`
	Updated          int64            `json:"updated"`
	Version          int64            `json:"version"`
}

type TechDebtSeverity string

const (
	TechDebtSeverityCritical TechDebtSeverity = "critical"
	TechDebtSeverityHigh     TechDebtSeverity = "high"
	TechDebtSeverityMedium   TechDebtSeverity = "medium"
	TechDebtSeverityLow      TechDebtSeverity = "low"
)

type TechDebtStatus string

const (
	TechDebtStatusOpen       TechDebtStatus = "open"
	TechDebtStatusInProgress TechDebtStatus = "in_progress"
	TechDebtStatusResolved   TechDebtStatus = "resolved"
	TechDebtStatusAccepted   TechDebtStatus = "accepted"
)

type TechDebtCategory string

const (
	TechDebtCategoryCodeSmell      TechDebtCategory = "code_smell"
	TechDebtCategoryBugRisk        TechDebtCategory = "bug_risk"
	TechDebtCategoryPerformance    TechDebtCategory = "performance"
	TechDebtCategorySecurity       TechDebtCategory = "security"
	TechDebtCategoryDocumentation  TechDebtCategory = "documentation"
	TechDebtCategoryTestCoverage   TechDebtCategory = "test_coverage"
	TechDebtCategoryDependency     TechDebtCategory = "dependency"
	TechDebtCategoryArchitecture   TechDebtCategory = "architecture"
)

type TechDebtCreateInput struct {
	Identifier      string   `json:"identifier"`
	Title           string   `json:"title"`
	Description     string   `json:"description,omitempty"`
	Severity        string   `json:"severity"`
	Status          string   `json:"status,omitempty"`
	Category        string   `json:"category"`
	FilePath        string   `json:"file_path,omitempty"`
	LineStart       int      `json:"line_start,omitempty"`
	LineEnd         int      `json:"line_end,omitempty"`
	EstimatedEffort string   `json:"estimated_effort"`
	Tags            []string `json:"tags,omitempty"`
	DueDate         int64    `json:"due_date,omitempty"`
	RepoID          int64    `json:"repo_id,omitempty"`
}

type TechDebtUpdateInput struct {
	Title           string   `json:"title,omitempty"`
	Description     string   `json:"description,omitempty"`
	Severity        string   `json:"severity,omitempty"`
	Status          string   `json:"status,omitempty"`
	Category        string   `json:"category,omitempty"`
	FilePath        string   `json:"file_path,omitempty"`
	LineStart       int      `json:"line_start,omitempty"`
	LineEnd         int      `json:"line_end,omitempty"`
	EstimatedEffort string   `json:"estimated_effort,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	DueDate         int64    `json:"due_date,omitempty"`
}

type TechDebtFilter struct {
	Severity   []string
	Status     []string
	Category   []string
	RepoID     int64
	Page       int
	Limit      int
	Sort       string
}

type TechDebtSummary struct {
	BySeverity map[string]int `json:"by_severity"`
	ByStatus   map[string]int `json:"by_status"`
	ByCategory map[string]int `json:"by_category"`
	Total      int            `json:"total"`
}

// StringSlice is a named type for []string to allow implementing driver.Valuer/sql.Scanner.
type StringSlice []string

// Value implements the driver.Valuer interface for storing StringSlice as JSON.
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for reading JSON into StringSlice.
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, s)
}
