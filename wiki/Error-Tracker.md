# Error Tracker

## Overview

The Error Tracker module provides comprehensive error tracking and management for the SoloDev platform. Applications report errors via REST API; the module groups similar errors by fingerprint, tracks individual occurrences, and manages error lifecycle status.

## Key Features

### Error Fingerprinting

Errors are automatically fingerprinted using a SHA256 hash of `title + stack_trace`. Duplicate errors are grouped together and the occurrence count is incremented rather than creating new records.

### Automatic Status Regression

When an error is reported that was previously resolved, its status automatically changes to `regressed`.

### Optimistic Locking

The `ErrorGroup.Version` field implements optimistic locking to prevent concurrent update conflicts.

### Events

The module emits three event types:
- `ErrorReported`
- `ErrorStatusChanged`
- `ErrorAssigned`

## API Endpoints

### Report Error

**POST** `/api/v1/spaces/{space_ref}/errors`

```json
{
  "identifier": "db-connection-timeout",
  "title": "Database Connection Timeout",
  "message": "Failed to connect to database within 30s",
  "severity": "error",
  "file_path": "internal/db/connect.go",
  "line_number": 42,
  "function_name": "Connect",
  "language": "go",
  "tags": ["database", "critical"],
  "stack_trace": "goroutine 1 [running]:\n...",
  "environment": "production",
  "runtime": "go1.24",
  "os": "linux",
  "arch": "amd64",
  "metadata": {}
}
```

Severity values: `fatal`, `error`, `warning`
Environment values: `production`, `staging`, `development`

### List Error Groups

**GET** `/api/v1/spaces/{space_ref}/errors`

Query parameters:

| Parameter | Description |
|-----------|-------------|
| `page` | Page number (default 0) |
| `limit` | Results per page (default 50) |
| `query` | Search query |
| `status` | Filter: `open`, `resolved`, `ignored`, `regressed` |
| `severity` | Filter: `fatal`, `error`, `warning` |
| `language` | Filter by programming language |

### Get Error Details

**GET** `/api/v1/spaces/{space_ref}/errors/{error_identifier}`

Returns `ErrorGroupDetail` including sample occurrences and user information.

### Update Error Group

**PATCH** `/api/v1/spaces/{space_ref}/errors/{error_identifier}`

```json
{
  "status": "resolved",
  "assigned_to": 12345,
  "title": "Updated title",
  "tags": ["database"]
}
```

### List Error Occurrences

**GET** `/api/v1/spaces/{space_ref}/errors/{error_identifier}/occurrences`

Query parameters: `page`, `limit` (max 100)

### Get Summary Statistics

**GET** `/api/v1/spaces/{space_ref}/errors/summary`

```json
{
  "total_errors": 0,
  "open_errors": 0,
  "resolved_errors": 0,
  "ignored_errors": 0,
  "regression_errors": 0,
  "fatal_count": 0,
  "error_count": 0,
  "warning_count": 0,
  "last_updated": 0
}
```

## Database Schema

### `error_groups` table (migration `0171`)

| Column | Type | Description |
|--------|------|-------------|
| `eg_id` | SERIAL PK | Primary key |
| `eg_space_id` | INT | FK to spaces |
| `eg_repo_id` | INT | FK to repositories |
| `eg_identifier` | TEXT | Space-scoped identifier |
| `eg_title` | TEXT | Error title |
| `eg_message` | TEXT | Error message |
| `eg_fingerprint` | TEXT UNIQUE | SHA256 of title+stack_trace |
| `eg_status` | TEXT | `open`, `resolved`, `ignored`, `regressed` |
| `eg_severity` | TEXT | `fatal`, `error`, `warning` |
| `eg_first_seen` | BIGINT | Unix milliseconds |
| `eg_last_seen` | BIGINT | Unix milliseconds |
| `eg_occurrence_count` | BIGINT | Total occurrence count |
| `eg_file_path` | TEXT | Source file (optional) |
| `eg_line_number` | INT | Line number (optional) |
| `eg_function_name` | TEXT | Function name (optional) |
| `eg_language` | TEXT | Programming language (optional) |
| `eg_tags` | JSON | Tag array |
| `eg_assigned_to` | INT | FK to principals (optional) |
| `eg_resolved_at` | BIGINT | Resolution timestamp (optional) |
| `eg_resolved_by` | INT | FK to principals (optional) |
| `eg_created_by` | INT | FK to principals |
| `eg_created` | BIGINT | Creation timestamp |
| `eg_updated` | BIGINT | Update timestamp |
| `eg_version` | BIGINT | Optimistic locking version |

### `error_occurrences` table (migration `0171`)

| Column | Type | Description |
|--------|------|-------------|
| `eo_id` | SERIAL PK | Primary key |
| `eo_error_group_id` | INT | FK to error_groups |
| `eo_stack_trace` | TEXT | Full stack trace |
| `eo_environment` | TEXT | `production`, `staging`, `development` |
| `eo_runtime` | TEXT | Runtime version (optional) |
| `eo_os` | TEXT | Operating system (optional) |
| `eo_arch` | TEXT | CPU architecture (optional) |
| `eo_metadata` | JSON | Arbitrary context |
| `eo_created_at` | BIGINT | Occurrence timestamp |

## Error Bridge Integration

The Error Tracker controller has an optional integration with the [AI Auto-Remediation](AI-Auto-Remediation) system via the [Error Bridge](Error-Bridge). When `SetErrorBridge(bridge)` is called on the controller:

1. Every call to `ReportError()` will also call `bridge.OnErrorReported()`
2. The bridge creates a pending `Remediation` task (skips `warning`-level errors)

```go
bridge := errorbridge.NewBridge(remediationStore, true)
errorTrackerCtrl.SetErrorBridge(bridge)
```

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/errortracker.go` |
| Database store | `app/store/database/errortracker.go` |
| Store interface | `app/store/database.go` |
| Controller | `app/api/controller/errortracker/controller.go` |
| Handlers | `app/api/handler/errortracker/` |
| Request helpers | `app/api/request/errortracker.go` |
| Events | `app/events/errortracker/` |
| Migrations | `app/store/database/migrate/*/0171_create_table_error_groups.*` |
