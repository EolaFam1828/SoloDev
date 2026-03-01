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

package types

import (
	"encoding/json"
	"time"
)

// ErrorGroup represents a group of similar errors grouped by fingerprint.
type ErrorGroup struct {
	ID           int64             `json:"id"`
	SpaceID      int64             `json:"-"`
	RepoID       int64             `json:"-"`
	Identifier   string            `json:"identifier"`
	Title        string            `json:"title"`
	Message      string            `json:"message"`
	Fingerprint  string            `json:"fingerprint"`
	Status       ErrorGroupStatus  `json:"status"`
	Severity     ErrorSeverity     `json:"severity"`
	FirstSeen    int64             `json:"first_seen"`
	LastSeen     int64             `json:"last_seen"`
	OccurrenceCount int64           `json:"occurrence_count"`
	FilePath     string            `json:"file_path,omitempty"`
	LineNumber   int               `json:"line_number,omitempty"`
	FunctionName string            `json:"function_name,omitempty"`
	Language     string            `json:"language,omitempty"`
	Tags         json.RawMessage   `json:"tags,omitempty"`
	AssignedTo   *int64            `json:"assigned_to,omitempty"`
	ResolvedAt   *int64            `json:"resolved_at,omitempty"`
	ResolvedBy   *int64            `json:"resolved_by,omitempty"`
	CreatedBy    int64             `json:"-"`
	Created      int64             `json:"created"`
	Updated      int64             `json:"updated"`
	Version      int64             `json:"version"`
}

// ErrorOccurrence represents a single error instance.
type ErrorOccurrence struct {
	ID          int64           `json:"id"`
	ErrorGroupID int64          `json:"error_group_id"`
	StackTrace  string          `json:"stack_trace"`
	Environment string          `json:"environment"`
	Runtime     string          `json:"runtime,omitempty"`
	OS          string          `json:"os,omitempty"`
	Arch        string          `json:"arch,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   int64           `json:"created_at"`
}

// ErrorGroupStatus represents the status of an error group.
type ErrorGroupStatus string

const (
	ErrorGroupStatusOpen       ErrorGroupStatus = "open"
	ErrorGroupStatusResolved   ErrorGroupStatus = "resolved"
	ErrorGroupStatusIgnored    ErrorGroupStatus = "ignored"
	ErrorGroupStatusRegressed  ErrorGroupStatus = "regressed"
)

// ErrorSeverity represents the severity level of an error.
type ErrorSeverity string

const (
	ErrorSeverityFatal   ErrorSeverity = "fatal"
	ErrorSeverityError   ErrorSeverity = "error"
	ErrorSeverityWarning ErrorSeverity = "warning"
)

// ErrorTrackerListOptions holds list error groups query parameters.
type ErrorTrackerListOptions struct {
	ListQueryFilter
	Status   *ErrorGroupStatus
	Severity *ErrorSeverity
	Language *string
}

// ErrorTrackerSummary represents summary statistics for error groups.
type ErrorTrackerSummary struct {
	TotalErrors      int64 `json:"total_errors"`
	OpenErrors       int64 `json:"open_errors"`
	ResolvedErrors   int64 `json:"resolved_errors"`
	IgnoredErrors    int64 `json:"ignored_errors"`
	RegressionErrors int64 `json:"regression_errors"`
	FatalCount       int64 `json:"fatal_count"`
	ErrorCount       int64 `json:"error_count"`
	WarningCount     int64 `json:"warning_count"`
	LastUpdated      int64 `json:"last_updated"`
}

// ReportErrorInput is the request body for reporting an error.
type ReportErrorInput struct {
	Identifier   string            `json:"identifier" binding:"required"`
	Title        string            `json:"title" binding:"required"`
	Message      string            `json:"message" binding:"required"`
	Severity     ErrorSeverity     `json:"severity"`
	FilePath     string            `json:"file_path,omitempty"`
	LineNumber   int               `json:"line_number,omitempty"`
	FunctionName string            `json:"function_name,omitempty"`
	Language     string            `json:"language,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	StackTrace   string            `json:"stack_trace" binding:"required"`
	Environment  string            `json:"environment"`
	Runtime      string            `json:"runtime,omitempty"`
	OS           string            `json:"os,omitempty"`
	Arch         string            `json:"arch,omitempty"`
	Metadata     json.RawMessage   `json:"metadata,omitempty"`
}

// UpdateErrorGroupInput is the request body for updating an error group.
type UpdateErrorGroupInput struct {
	Status     *ErrorGroupStatus `json:"status,omitempty"`
	AssignedTo *int64            `json:"assigned_to,omitempty"`
	Title      string            `json:"title,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
}

// ErrorGroupDetail includes ErrorGroup with additional computed fields.
type ErrorGroupDetail struct {
	*ErrorGroup
	OccurrencesSample []ErrorOccurrence `json:"occurrences_sample,omitempty"`
	AssignedUser      *PrincipalInfo    `json:"assigned_user,omitempty"`
	ResolvedByUser    *PrincipalInfo    `json:"resolved_by_user,omitempty"`
	CreatedByUser     *PrincipalInfo    `json:"created_by_user,omitempty"`
}

// Timestamp helpers for unix milliseconds conversion.
func NowMillis() int64 {
	return time.Now().UnixMilli()
}
