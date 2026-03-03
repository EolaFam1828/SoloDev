// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package types

// SignalType identifies the source module of a signal.
type SignalType string

const (
	SignalTypeError           SignalType = "error"
	SignalTypePipelineFailure SignalType = "pipeline_failure"
	SignalTypeHealthCheck     SignalType = "health_check"
	SignalTypeSecurity        SignalType = "security_finding"
)

// Signal represents a normalized cross-module event.
type Signal struct {
	Type       SignalType `json:"type"`
	SourceID   string     `json:"source_id"` // module-specific identifier
	SpaceID    int64      `json:"space_id"`
	RepoID     int64      `json:"repo_id,omitempty"`
	Title      string     `json:"title"`
	Severity   string     `json:"severity"` // critical, high, medium, low
	FilePath   string     `json:"file_path,omitempty"`
	Branch     string     `json:"branch,omitempty"`
	OccurredAt int64      `json:"occurred_at"` // millis
}

// IncidentSeverity scores how severe a correlated incident is.
type IncidentSeverity string

const (
	IncidentSeverityCritical IncidentSeverity = "critical"
	IncidentSeverityHigh     IncidentSeverity = "high"
	IncidentSeverityMedium   IncidentSeverity = "medium"
	IncidentSeverityLow      IncidentSeverity = "low"
)

// CorrelatedIncident groups related signals that occurred in proximity.
type CorrelatedIncident struct {
	ID          string           `json:"id"`
	SpaceID     int64            `json:"space_id"`
	RepoID      int64            `json:"repo_id,omitempty"`
	Severity    IncidentSeverity `json:"severity"`
	Title       string           `json:"title"`
	Summary     string           `json:"summary"`
	Signals     []Signal         `json:"signals"`
	SignalCount int              `json:"signal_count"`
	FirstSeen   int64            `json:"first_seen"`
	LastSeen    int64            `json:"last_seen"`
}

// CorrelatedIncidentFilter holds query parameters for listing correlated incidents.
type CorrelatedIncidentFilter struct {
	WindowMinutes int    `json:"window_minutes"` // correlation window (default 30)
	MinSignals    int    `json:"min_signals"`    // minimum signals to form an incident (default 2)
	RepoID        *int64 `json:"repo_id,omitempty"`
}
