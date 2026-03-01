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

import "database/sql/driver"

type HealthCheckStatus string

const (
	HealthCheckStatusUp       HealthCheckStatus = "up"
	HealthCheckStatusDown     HealthCheckStatus = "down"
	HealthCheckStatusDegraded HealthCheckStatus = "degraded"
	HealthCheckStatusUnknown  HealthCheckStatus = "unknown"
)

func (s HealthCheckStatus) String() string {
	return string(s)
}

func (s HealthCheckStatus) Value() (driver.Value, error) {
	return string(s), nil
}

type HealthCheck struct {
	ID                  int64                `db:"hc_id"                      json:"id"`
	SpaceID             int64                `db:"hc_space_id"                json:"space_id"`
	Identifier          string               `db:"hc_identifier"              json:"identifier"`
	Name                string               `db:"hc_name"                    json:"name"`
	Description         string               `db:"hc_description"             json:"description"`
	URL                 string               `db:"hc_url"                     json:"url"`
	Method              string               `db:"hc_method"                  json:"method"`
	ExpectedStatus      int                  `db:"hc_expected_status"         json:"expected_status"`
	IntervalSeconds     int                  `db:"hc_interval_seconds"        json:"interval_seconds"`
	TimeoutSeconds      int                  `db:"hc_timeout_seconds"         json:"timeout_seconds"`
	Enabled             bool                 `db:"hc_enabled"                 json:"enabled"`
	Headers             string               `db:"hc_headers"                 json:"headers"`
	Body                string               `db:"hc_body"                    json:"body"`
	Tags                string               `db:"hc_tags"                    json:"tags"`
	LastStatus          string               `db:"hc_last_status"             json:"last_status"`
	LastCheckedAt       int64                `db:"hc_last_checked_at"         json:"last_checked_at"`
	LastResponseTime    int64                `db:"hc_last_response_time"      json:"last_response_time"`
	ConsecutiveFailures int                  `db:"hc_consecutive_failures"    json:"consecutive_failures"`
	CreatedBy           int64                `db:"hc_created_by"              json:"created_by"`
	Created             int64                `db:"hc_created"                 json:"created"`
	Updated             int64                `db:"hc_updated"                 json:"updated"`
	Version             int64                `db:"hc_version"                 json:"-"`
}

type HealthCheckResult struct {
	ID              int64  `db:"hcr_id"                  json:"id"`
	HealthCheckID   int64  `db:"hcr_health_check_id"     json:"health_check_id"`
	Status          string `db:"hcr_status"              json:"status"`
	ResponseTime    int64  `db:"hcr_response_time"       json:"response_time"`
	StatusCode      int    `db:"hcr_status_code"         json:"status_code"`
	ErrorMessage    string `db:"hcr_error_message"       json:"error_message"`
	CreatedAt       int64  `db:"hcr_created_at"          json:"created_at"`
}

type HealthCheckSummary struct {
	HealthCheckID       int64   `json:"health_check_id"`
	Identifier          string  `json:"identifier"`
	Name                string  `json:"name"`
	CurrentStatus       string  `json:"current_status"`
	UptimePercentage    float64 `json:"uptime_percentage"`
	TotalChecks         int64   `json:"total_checks"`
	SuccessfulChecks    int64   `json:"successful_checks"`
	FailedChecks        int64   `json:"failed_checks"`
	AverageResponseTime int64   `json:"average_response_time"`
	LastCheckedAt       int64   `json:"last_checked_at"`
}
