// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package types

// AnomalyType classifies the kind of detected anomaly.
type AnomalyType string

const (
	AnomalyTypeLatencySpike AnomalyType = "latency_spike" // Sudden increase in response time
	AnomalyTypeLatencyDrift AnomalyType = "latency_drift" // Gradual upward trend in response time
	AnomalyTypeErrorBurst   AnomalyType = "error_burst"   // Cluster of errors in a short window
	AnomalyTypeFlapping     AnomalyType = "flapping"      // Rapid up/down status oscillation
)

// AnomalySeverity indicates how concerning the anomaly is.
type AnomalySeverity string

const (
	AnomalySeverityLow      AnomalySeverity = "low"
	AnomalySeverityMedium   AnomalySeverity = "medium"
	AnomalySeverityHigh     AnomalySeverity = "high"
	AnomalySeverityCritical AnomalySeverity = "critical"
)

// Anomaly represents a detected anomaly in health check behavior.
type Anomaly struct {
	HealthCheckID   int64           `json:"health_check_id"`
	Identifier      string          `json:"identifier"`
	Name            string          `json:"name"`
	Type            AnomalyType     `json:"type"`
	Severity        AnomalySeverity `json:"severity"`
	Description     string          `json:"description"`
	CurrentValue    float64         `json:"current_value"`
	BaselineValue   float64         `json:"baseline_value"`
	DeviationFactor float64         `json:"deviation_factor"`
	DetectedAt      int64           `json:"detected_at"`
}

// AnomalyReport is the response from the anomaly detection endpoint.
type AnomalyReport struct {
	SpaceID       int64     `json:"space_id"`
	ChecksScanned int       `json:"checks_scanned"`
	Anomalies     []Anomaly `json:"anomalies"`
	AnalyzedAt    int64     `json:"analyzed_at"`
}
