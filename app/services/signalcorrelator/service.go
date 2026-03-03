// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package signalcorrelator provides cross-domain correlation between
// errors, pipeline failures, health check degradation, and security
// findings. Signals that occur within a configurable time window
// affecting the same entity (repo or space) are grouped into
// correlated incidents.
package signalcorrelator

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"
	"time"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

const (
	defaultWindowMinutes = 30
	defaultMinSignals    = 2
)

// Service correlates signals across modules.
type Service struct {
	errorStore        store.ErrorTrackerStore
	executionStore    store.ExecutionStore
	healthStore       store.HealthCheckStore
	healthResultStore store.HealthCheckResultStore
	scanStore         store.SecurityScanStore
	repoStore         store.RepoStore
}

// NewService creates a new signal correlator service.
func NewService(
	errorStore store.ErrorTrackerStore,
	executionStore store.ExecutionStore,
	healthStore store.HealthCheckStore,
	healthResultStore store.HealthCheckResultStore,
	scanStore store.SecurityScanStore,
	repoStore store.RepoStore,
) *Service {
	return &Service{
		errorStore:        errorStore,
		executionStore:    executionStore,
		healthStore:       healthStore,
		healthResultStore: healthResultStore,
		scanStore:         scanStore,
		repoStore:         repoStore,
	}
}

// Correlate fetches recent signals across all modules for a space and groups
// them into correlated incidents based on time proximity and entity.
func (s *Service) Correlate(
	ctx context.Context,
	spaceID int64,
	filter types.CorrelatedIncidentFilter,
) ([]types.CorrelatedIncident, error) {
	windowMinutes := filter.WindowMinutes
	if windowMinutes <= 0 {
		windowMinutes = defaultWindowMinutes
	}
	minSignals := filter.MinSignals
	if minSignals <= 0 {
		minSignals = defaultMinSignals
	}

	since := time.Now().Add(-time.Duration(windowMinutes) * time.Minute)
	sinceMs := since.UnixMilli()

	// Collect signals from all modules concurrently.
	var signals []types.Signal

	errorSignals, _ := s.collectErrorSignals(ctx, spaceID, sinceMs)
	signals = append(signals, errorSignals...)

	pipelineSignals, _ := s.collectPipelineSignals(ctx, spaceID, sinceMs)
	signals = append(signals, pipelineSignals...)

	healthSignals, _ := s.collectHealthSignals(ctx, spaceID, sinceMs)
	signals = append(signals, healthSignals...)

	securitySignals, _ := s.collectSecuritySignals(ctx, spaceID, sinceMs)
	signals = append(signals, securitySignals...)

	if filter.RepoID != nil {
		var filtered []types.Signal
		for _, sig := range signals {
			if sig.RepoID == *filter.RepoID || sig.RepoID == 0 {
				filtered = append(filtered, sig)
			}
		}
		signals = filtered
	}

	// Group signals into incidents.
	incidents := groupIntoIncidents(signals, spaceID, windowMinutes, minSignals)
	return incidents, nil
}

func (s *Service) collectErrorSignals(ctx context.Context, spaceID, sinceMs int64) ([]types.Signal, error) {
	errors, err := s.errorStore.List(ctx, spaceID, types.ErrorTrackerListOptions{
		ListQueryFilter: types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 50}},
	})
	if err != nil {
		return nil, err
	}

	var signals []types.Signal
	for _, eg := range errors {
		if eg.LastSeen < sinceMs {
			continue
		}
		if eg.Status == types.ErrorGroupStatusResolved || eg.Status == types.ErrorGroupStatusIgnored {
			continue
		}
		severity := "medium"
		switch eg.Severity {
		case types.ErrorSeverityFatal:
			severity = "critical"
		case types.ErrorSeverityError:
			severity = "high"
		case types.ErrorSeverityWarning:
			severity = "low"
		}
		signals = append(signals, types.Signal{
			Type:       types.SignalTypeError,
			SourceID:   eg.Identifier,
			SpaceID:    spaceID,
			RepoID:     eg.RepoID,
			Title:      eg.Title,
			Severity:   severity,
			FilePath:   eg.FilePath,
			OccurredAt: eg.LastSeen,
		})
	}
	return signals, nil
}

func (s *Service) collectPipelineSignals(ctx context.Context, spaceID, sinceMs int64) ([]types.Signal, error) {
	// Pipeline failures are repo-scoped; find repos in this space.
	failures, err := s.executionStore.ListRecentFailed(ctx, sinceMs, 50)
	if err != nil {
		return nil, err
	}

	var signals []types.Signal
	for _, exec := range failures {
		// Verify the execution belongs to a repo in this space.
		repo, err := s.repoStore.Find(ctx, exec.RepoID)
		if err != nil || repo.ParentID != spaceID {
			continue
		}

		title := fmt.Sprintf("Pipeline #%d failed", exec.Number)
		if exec.Error != "" {
			title = exec.Error
			if len(title) > 120 {
				title = title[:120]
			}
		}

		signals = append(signals, types.Signal{
			Type:       types.SignalTypePipelineFailure,
			SourceID:   fmt.Sprintf("exec-%d-%d", exec.RepoID, exec.Number),
			SpaceID:    spaceID,
			RepoID:     exec.RepoID,
			Title:      title,
			Severity:   "high",
			Branch:     exec.Target,
			OccurredAt: exec.Finished,
		})
	}
	return signals, nil
}

func (s *Service) collectHealthSignals(ctx context.Context, spaceID, sinceMs int64) ([]types.Signal, error) {
	checks, err := s.healthStore.List(ctx, spaceID, types.ListQueryFilter{Pagination: types.Pagination{Page: 1, Size: 50}})
	if err != nil {
		return nil, err
	}

	var signals []types.Signal
	for _, hc := range checks {
		if hc.ConsecutiveFailures == 0 {
			continue
		}
		if hc.LastCheckedAt < sinceMs {
			continue
		}

		severity := "medium"
		if hc.ConsecutiveFailures >= 5 {
			severity = "critical"
		} else if hc.ConsecutiveFailures >= 3 {
			severity = "high"
		}

		signals = append(signals, types.Signal{
			Type:       types.SignalTypeHealthCheck,
			SourceID:   hc.Identifier,
			SpaceID:    spaceID,
			Title:      fmt.Sprintf("Health check %q failing (%d consecutive)", hc.Name, hc.ConsecutiveFailures),
			Severity:   severity,
			OccurredAt: hc.LastCheckedAt,
		})
	}
	return signals, nil
}

func (s *Service) collectSecuritySignals(ctx context.Context, spaceID, sinceMs int64) ([]types.Signal, error) {
	// ListByStatus returns scans across all repos; we filter by space membership.
	scans, err := s.scanStore.ListByStatus(ctx, enum.SecurityScanStatusCompleted, 50)
	if err != nil {
		return nil, err
	}

	var signals []types.Signal
	for _, scan := range scans {
		if scan.Created < sinceMs {
			continue
		}
		if scan.SpaceID != spaceID {
			continue
		}
		if scan.CriticalCount == 0 && scan.HighCount == 0 {
			continue
		}

		severity := "medium"
		if scan.CriticalCount > 0 {
			severity = "critical"
		} else if scan.HighCount > 0 {
			severity = "high"
		}

		signals = append(signals, types.Signal{
			Type:       types.SignalTypeSecurity,
			SourceID:   scan.Identifier,
			SpaceID:    spaceID,
			RepoID:     scan.RepoID,
			Title:      fmt.Sprintf("Security scan: %d critical, %d high findings", scan.CriticalCount, scan.HighCount),
			Severity:   severity,
			Branch:     scan.Branch,
			OccurredAt: scan.Created,
		})
	}
	return signals, nil
}

// groupIntoIncidents clusters signals by repo and time window into incidents.
func groupIntoIncidents(signals []types.Signal, spaceID int64, windowMinutes, minSignals int) []types.CorrelatedIncident {
	if len(signals) < minSignals {
		return nil
	}

	// Sort by time.
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].OccurredAt < signals[j].OccurredAt
	})

	windowMs := int64(windowMinutes) * 60 * 1000

	// Group by repo (0 = space-level).
	byRepo := map[int64][]types.Signal{}
	for _, sig := range signals {
		byRepo[sig.RepoID] = append(byRepo[sig.RepoID], sig)
	}

	var incidents []types.CorrelatedIncident

	for repoID, repoSignals := range byRepo {
		// Sliding window clustering: start a new cluster when a signal
		// falls outside the window of the cluster's first signal.
		var clusters [][]types.Signal
		var current []types.Signal

		for _, sig := range repoSignals {
			if len(current) == 0 {
				current = append(current, sig)
				continue
			}
			if sig.OccurredAt-current[0].OccurredAt <= windowMs {
				current = append(current, sig)
			} else {
				clusters = append(clusters, current)
				current = []types.Signal{sig}
			}
		}
		if len(current) > 0 {
			clusters = append(clusters, current)
		}

		for _, cluster := range clusters {
			if len(cluster) < minSignals {
				continue
			}

			incident := buildIncident(cluster, spaceID, repoID)
			incidents = append(incidents, incident)
		}
	}

	// Sort incidents by severity then recency.
	sort.Slice(incidents, func(i, j int) bool {
		si := severityRank(incidents[i].Severity)
		sj := severityRank(incidents[j].Severity)
		if si != sj {
			return si > sj
		}
		return incidents[i].LastSeen > incidents[j].LastSeen
	})

	return incidents
}

func buildIncident(signals []types.Signal, spaceID, repoID int64) types.CorrelatedIncident {
	first := signals[0].OccurredAt
	last := signals[len(signals)-1].OccurredAt

	// Compute severity as the max severity across signals.
	maxSev := types.IncidentSeverityLow
	for _, sig := range signals {
		sev := mapSeverity(sig.Severity)
		if severityRank(sev) > severityRank(maxSev) {
			maxSev = sev
		}
	}

	// Boost severity when multiple signal types are involved.
	typeSet := map[types.SignalType]bool{}
	for _, sig := range signals {
		typeSet[sig.Type] = true
	}
	if len(typeSet) >= 3 && maxSev != types.IncidentSeverityCritical {
		maxSev = types.IncidentSeverityCritical
	} else if len(typeSet) >= 2 && severityRank(maxSev) < severityRank(types.IncidentSeverityHigh) {
		maxSev = types.IncidentSeverityHigh
	}

	title := buildTitle(signals, typeSet)
	summary := buildSummary(signals, typeSet)

	// Deterministic ID from signal sources.
	idInput := fmt.Sprintf("%d-%d-%d-%d", spaceID, repoID, first, len(signals))
	hash := sha256.Sum256([]byte(idInput))
	id := fmt.Sprintf("inc-%x", hash[:8])

	return types.CorrelatedIncident{
		ID:          id,
		SpaceID:     spaceID,
		RepoID:      repoID,
		Severity:    maxSev,
		Title:       title,
		Summary:     summary,
		Signals:     signals,
		SignalCount: len(signals),
		FirstSeen:   first,
		LastSeen:    last,
	}
}

func buildTitle(signals []types.Signal, typeSet map[types.SignalType]bool) string {
	if len(typeSet) == 1 {
		for t := range typeSet {
			return fmt.Sprintf("%d correlated %s signals", len(signals), t)
		}
	}
	var typeNames []string
	for t := range typeSet {
		typeNames = append(typeNames, string(t))
	}
	sort.Strings(typeNames)
	return fmt.Sprintf("Correlated incident: %d signals across %v", len(signals), typeNames)
}

func buildSummary(signals []types.Signal, typeSet map[types.SignalType]bool) string {
	counts := map[types.SignalType]int{}
	for _, sig := range signals {
		counts[sig.Type]++
	}
	summary := ""
	for t, c := range counts {
		if summary != "" {
			summary += ", "
		}
		summary += fmt.Sprintf("%d %s", c, t)
	}
	return summary
}

func mapSeverity(s string) types.IncidentSeverity {
	switch s {
	case "critical":
		return types.IncidentSeverityCritical
	case "high":
		return types.IncidentSeverityHigh
	case "medium":
		return types.IncidentSeverityMedium
	default:
		return types.IncidentSeverityLow
	}
}

func severityRank(s types.IncidentSeverity) int {
	switch s {
	case types.IncidentSeverityCritical:
		return 4
	case types.IncidentSeverityHigh:
		return 3
	case types.IncidentSeverityMedium:
		return 2
	case types.IncidentSeverityLow:
		return 1
	default:
		return 0
	}
}
