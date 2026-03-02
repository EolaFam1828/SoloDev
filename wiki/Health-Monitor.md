# Health Monitor

## Overview

The Health Monitor module provides uptime monitoring for HTTP endpoints. Users configure monitors with check intervals and expected responses; the module tracks check results over time and provides uptime statistics and status summaries.

## Monitor Configuration

When creating a monitor, the following validation rules apply:

| Field | Constraint |
|-------|-----------|
| `identifier` | Pattern `^[a-zA-Z0-9_-]+$`, 3–200 chars |
| `name` | Required, 1–255 characters |
| `url` | Must be a valid HTTP/HTTPS URL |
| `method` | Must be `GET`, `POST`, or `HEAD` |
| `expected_status` | 100–599 |
| `interval_seconds` | 60–86400 |
| `timeout_seconds` | 1–300 |

## API Endpoints

### Create Monitor

**POST** `/api/v1/spaces/{space_ref}/health-checks`

```json
{
  "identifier": "api-health",
  "name": "API Health Check",
  "description": "Monitor main API endpoint",
  "url": "https://api.example.com/health",
  "method": "GET",
  "expected_status": 200,
  "interval_seconds": 300,
  "timeout_seconds": 10,
  "enabled": true,
  "headers": "{}",
  "body": "",
  "tags": "[\"production\"]"
}
```

### List Monitors

**GET** `/api/v1/spaces/{space_ref}/health-checks`

Query parameters: `page`, `size`, `query`

### Get Monitor

**GET** `/api/v1/spaces/{space_ref}/health-checks/{identifier}`

### Update Monitor

**PATCH** `/api/v1/spaces/{space_ref}/health-checks/{identifier}`

```json
{
  "name": "Updated Name",
  "interval_seconds": 600
}
```

> **Note:** There is no dedicated toggle endpoint in the router. Enabling or disabling a monitor is done through the standard Update endpoint (`PATCH`) by sending `{"enabled": true}` or `{"enabled": false}`. This differs from the Quality Gates module, which has a dedicated `/toggle` route.

### Delete Monitor

**DELETE** `/api/v1/spaces/{space_ref}/health-checks/{identifier}`

### Get Recent Results

**GET** `/api/v1/spaces/{space_ref}/health-checks/{identifier}/results`

Query parameters: `limit` (returns latest results first)

### Get Uptime Summary

**GET** `/api/v1/spaces/{space_ref}/health-checks/{identifier}/summary`

Returns an array of `HealthCheckSummary` objects with aggregated statistics:

| Field | Description |
|-------|-------------|
| `uptime_percentage` | Percentage of successful checks in past 24h |
| `total_checks` | Total check count |
| `successful_checks` | Count of successful checks |
| `failed_checks` | Count of failed checks |
| `average_response_time` | Average response time in milliseconds |

## Status Values

`up` · `down` · `degraded` · `unknown`

## Database Schema

### `health_checks` table (migration `0103`)

| Column | Description |
|--------|-------------|
| `hc_id` | Primary key (auto-increment) |
| `hc_space_id` | FK to spaces |
| `hc_identifier` | Unique within space (indexed) |
| `hc_name` | Display name |
| `hc_description` | Optional description |
| `hc_url` | HTTP endpoint to monitor |
| `hc_method` | `GET`, `POST`, or `HEAD` |
| `hc_expected_status` | Expected HTTP status code |
| `hc_interval_seconds` | Check frequency (60–86400) |
| `hc_timeout_seconds` | Request timeout (1–300) |
| `hc_enabled` | Enable/disable monitor |
| `hc_headers` | JSON custom headers |
| `hc_body` | Request body for POST |
| `hc_tags` | JSON array of tags |
| `hc_last_status` | Latest status |
| `hc_last_checked_at` | Timestamp of last check (unix millis) |
| `hc_last_response_time` | Last response time in ms |
| `hc_consecutive_failures` | Count of consecutive failures |
| `hc_created_by` | FK to principals |
| `hc_created` | Unix milliseconds |
| `hc_updated` | Unix milliseconds |
| `hc_version` | Optimistic locking version |

### `health_check_results` table (migration `0104`)

| Column | Description |
|--------|-------------|
| `hcr_id` | Primary key |
| `hcr_health_check_id` | FK to health_checks (indexed) |
| `hcr_status` | `up`, `down`, `degraded` |
| `hcr_response_time` | Response time in milliseconds |
| `hcr_status_code` | HTTP status code received |
| `hcr_error_message` | Error details if failed |
| `hcr_created_at` | Result timestamp (unix millis, indexed) |

## Optimistic Locking

Health checks use optimistic locking via `hc_version`. The version is incremented on each update; if there is a version mismatch, the update fails with `ErrVersionConflict` and the controller automatically retries with the latest version.

## Events

Five events are published:

| Event | Trigger |
|-------|---------|
| `healthcheck_created` | Monitor created |
| `healthcheck_updated` | Monitor configuration changed |
| `healthcheck_deleted` | Monitor deleted |
| `healthcheck_status_changed` | Status changed (e.g., `up` → `down`) |
| `healthcheck_result_created` | Individual check completed |

## Access Control

All endpoints check authorization against the parent space using `PermissionRepoView`.

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/healthcheck.go` |
| Store interfaces | `app/store/database.go` |
| Health check store | `app/store/database/healthcheck.go` |
| Results store | `app/store/database/healthcheck_result.go` |
| Controller | `app/api/controller/healthcheck/` |
| Handlers | `app/api/handler/healthcheck/` |
| Events | `app/events/healthcheck/` |
| Migrations | `app/store/database/migrate/*/0103_create_table_health_checks.*` and `0104_create_table_health_check_results.*` |
