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

package anomalydetector

import (
	"context"
	"fmt"
	"math"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
)

const (
	// Number of recent results to analyze per health check.
	analysisWindow = 100
	// Z-score threshold for spike detection (2.5 standard deviations).
	spikeZThreshold = 2.5
	// Minimum results needed for statistical analysis.
	minResults = 10
	// Window size for drift detection (compare first half vs second half).
	driftRatioThreshold = 1.5
	// Minimum status transitions per window for flapping detection.
	flappingThreshold = 5
	// Error ratio threshold for burst detection.
	errorBurstThreshold = 0.3
)

// Service detects anomalies in health check response time trends.
type Service struct {
	healthCheckStore       store.HealthCheckStore
	healthCheckResultStore store.HealthCheckResultStore
}

// NewService creates a new anomaly detection service.
func NewService(
	healthCheckStore store.HealthCheckStore,
	healthCheckResultStore store.HealthCheckResultStore,
) *Service {
	return &Service{
		healthCheckStore:       healthCheckStore,
		healthCheckResultStore: healthCheckResultStore,
	}
}

// Analyze scans all enabled health checks in a space and returns detected anomalies.
func (s *Service) Analyze(ctx context.Context, spaceID int64) (*types.AnomalyReport, error) {
	checks, err := s.healthCheckStore.List(ctx, spaceID, types.ListQueryFilter{
		Pagination: types.Pagination{Page: 1, Size: 200},
	})
	if err != nil {
		return nil, fmt.Errorf("list health checks: %w", err)
	}

	report := &types.AnomalyReport{
		SpaceID:    spaceID,
		AnalyzedAt: types.NowMillis(),
	}

	for _, check := range checks {
		if !check.Enabled {
			continue
		}
		report.ChecksScanned++

		results, err := s.healthCheckResultStore.ListByHealthCheckID(ctx, check.ID, analysisWindow)
		if err != nil || len(results) < minResults {
			continue
		}

		// Detect various anomaly types.
		if anomaly := detectLatencySpike(check, results); anomaly != nil {
			report.Anomalies = append(report.Anomalies, *anomaly)
		}
		if anomaly := detectLatencyDrift(check, results); anomaly != nil {
			report.Anomalies = append(report.Anomalies, *anomaly)
		}
		if anomaly := detectErrorBurst(check, results); anomaly != nil {
			report.Anomalies = append(report.Anomalies, *anomaly)
		}
		if anomaly := detectFlapping(check, results); anomaly != nil {
			report.Anomalies = append(report.Anomalies, *anomaly)
		}
	}

	return report, nil
}

// detectLatencySpike identifies recent response times that are significantly above the baseline.
func detectLatencySpike(check *types.HealthCheck, results []*types.HealthCheckResult) *types.Anomaly {
	times := extractResponseTimes(results)
	if len(times) < minResults {
		return nil
	}

	mean, stddev := meanStdDev(times)
	if stddev == 0 {
		return nil
	}

	// Check the most recent result.
	latest := times[0]
	zScore := (latest - mean) / stddev

	if zScore < spikeZThreshold {
		return nil
	}

	severity := types.AnomalySeverityLow
	if zScore >= 4.0 {
		severity = types.AnomalySeverityCritical
	} else if zScore >= 3.5 {
		severity = types.AnomalySeverityHigh
	} else if zScore >= 3.0 {
		severity = types.AnomalySeverityMedium
	}

	return &types.Anomaly{
		HealthCheckID:   check.ID,
		Identifier:      check.Identifier,
		Name:            check.Name,
		Type:            types.AnomalyTypeLatencySpike,
		Severity:        severity,
		Description:     fmt.Sprintf("Response time %.0fms is %.1f standard deviations above mean %.0fms", latest, zScore, mean),
		CurrentValue:    latest,
		BaselineValue:   mean,
		DeviationFactor: zScore,
		DetectedAt:      results[0].CreatedAt,
	}
}

// detectLatencyDrift identifies a gradual upward trend in response times.
func detectLatencyDrift(check *types.HealthCheck, results []*types.HealthCheckResult) *types.Anomaly {
	times := extractResponseTimes(results)
	if len(times) < minResults*2 {
		return nil
	}

	mid := len(times) / 2
	// Results are ordered most-recent-first, so first half = recent, second half = older.
	recentMean, _ := meanStdDev(times[:mid])
	olderMean, _ := meanStdDev(times[mid:])

	if olderMean == 0 {
		return nil
	}

	ratio := recentMean / olderMean
	if ratio < driftRatioThreshold {
		return nil
	}

	severity := types.AnomalySeverityLow
	if ratio >= 3.0 {
		severity = types.AnomalySeverityHigh
	} else if ratio >= 2.0 {
		severity = types.AnomalySeverityMedium
	}

	return &types.Anomaly{
		HealthCheckID:   check.ID,
		Identifier:      check.Identifier,
		Name:            check.Name,
		Type:            types.AnomalyTypeLatencyDrift,
		Severity:        severity,
		Description:     fmt.Sprintf("Average response time drifted from %.0fms to %.0fms (%.1fx increase)", olderMean, recentMean, ratio),
		CurrentValue:    recentMean,
		BaselineValue:   olderMean,
		DeviationFactor: ratio,
		DetectedAt:      results[0].CreatedAt,
	}
}

// detectErrorBurst identifies a cluster of errors in the recent window.
func detectErrorBurst(check *types.HealthCheck, results []*types.HealthCheckResult) *types.Anomaly {
	// Look at the 20 most recent results.
	window := results
	if len(window) > 20 {
		window = window[:20]
	}

	errors := 0
	for _, r := range window {
		if r.Status != string(types.HealthCheckStatusUp) {
			errors++
		}
	}

	ratio := float64(errors) / float64(len(window))
	if ratio < errorBurstThreshold {
		return nil
	}

	severity := types.AnomalySeverityMedium
	if ratio >= 0.8 {
		severity = types.AnomalySeverityCritical
	} else if ratio >= 0.5 {
		severity = types.AnomalySeverityHigh
	}

	return &types.Anomaly{
		HealthCheckID:   check.ID,
		Identifier:      check.Identifier,
		Name:            check.Name,
		Type:            types.AnomalyTypeErrorBurst,
		Severity:        severity,
		Description:     fmt.Sprintf("%d of %d recent checks failed (%.0f%% error rate)", errors, len(window), ratio*100),
		CurrentValue:    ratio,
		BaselineValue:   0,
		DeviationFactor: ratio,
		DetectedAt:      results[0].CreatedAt,
	}
}

// detectFlapping identifies rapid up/down oscillation.
func detectFlapping(check *types.HealthCheck, results []*types.HealthCheckResult) *types.Anomaly {
	if len(results) < minResults {
		return nil
	}

	// Count status transitions in recent results.
	window := results
	if len(window) > 30 {
		window = window[:30]
	}

	transitions := 0
	for i := 1; i < len(window); i++ {
		if window[i].Status != window[i-1].Status {
			transitions++
		}
	}

	if transitions < flappingThreshold {
		return nil
	}

	severity := types.AnomalySeverityMedium
	if transitions >= 10 {
		severity = types.AnomalySeverityHigh
	}

	return &types.Anomaly{
		HealthCheckID:   check.ID,
		Identifier:      check.Identifier,
		Name:            check.Name,
		Type:            types.AnomalyTypeFlapping,
		Severity:        severity,
		Description:     fmt.Sprintf("%d status transitions in %d checks (flapping detected)", transitions, len(window)),
		CurrentValue:    float64(transitions),
		BaselineValue:   0,
		DeviationFactor: float64(transitions) / float64(len(window)),
		DetectedAt:      results[0].CreatedAt,
	}
}

func extractResponseTimes(results []*types.HealthCheckResult) []float64 {
	var times []float64
	for _, r := range results {
		if r.ResponseTime > 0 {
			times = append(times, float64(r.ResponseTime))
		}
	}
	return times
}

func meanStdDev(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var varianceSum float64
	for _, v := range values {
		d := v - mean
		varianceSum += d * d
	}
	stddev := math.Sqrt(varianceSum / float64(len(values)))

	return mean, stddev
}
