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

package events

import (
	"github.com/harness/gitness/events"
	"github.com/harness/gitness/types"
)

const (
	EventHealthCheckCreated   events.EventType = "healthcheck_created"
	EventHealthCheckUpdated   events.EventType = "healthcheck_updated"
	EventHealthCheckDeleted   events.EventType = "healthcheck_deleted"
	EventHealthCheckStatusChanged events.EventType = "healthcheck_status_changed"
	EventHealthCheckResultCreated events.EventType = "healthcheck_result_created"
)

type EventHealthCheckCreated struct {
	HealthCheck *types.HealthCheck `json:"health_check,omitempty"`
	CreatedBy   int64              `json:"created_by,omitempty"`
}

type EventHealthCheckUpdated struct {
	HealthCheck *types.HealthCheck `json:"health_check,omitempty"`
	UpdatedBy   int64              `json:"updated_by,omitempty"`
}

type EventHealthCheckDeleted struct {
	HealthCheckID int64  `json:"health_check_id,omitempty"`
	SpaceID       int64  `json:"space_id,omitempty"`
	Identifier    string `json:"identifier,omitempty"`
	DeletedBy     int64  `json:"deleted_by,omitempty"`
}

type EventHealthCheckStatusChanged struct {
	HealthCheckID int64  `json:"health_check_id,omitempty"`
	SpaceID       int64  `json:"space_id,omitempty"`
	OldStatus     string `json:"old_status,omitempty"`
	NewStatus     string `json:"new_status,omitempty"`
	Timestamp     int64  `json:"timestamp,omitempty"`
}

type EventHealthCheckResultCreated struct {
	Result *types.HealthCheckResult `json:"result,omitempty"`
}
