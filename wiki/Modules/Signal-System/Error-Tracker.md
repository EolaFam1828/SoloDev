# Error Tracker

## Purpose

The Error Tracker module provides runtime error tracking and management. Applications report errors via REST API; the module groups similar errors by fingerprint, tracks individual occurrences, and manages error lifecycle status. It is the primary signal source feeding the remediation loop.

## Inputs

- Error reports via `POST /api/v1/spaces/{space_ref}/errors`
- Each report includes: title, message, severity, file path, line number, function name, language, stack trace, environment, runtime, OS, architecture, tags, and metadata

## Processing

### Error Fingerprinting

Errors are automatically fingerprinted using a SHA256 hash of `title + stack_trace`. Duplicate errors are grouped together and the occurrence count is incremented rather than creating new records.

### Automatic Status Regression

When an error is reported that was previously resolved, its status automatically changes to `regressed`.

### Optimistic Locking

The `ErrorGroup.Version` field implements optimistic locking to prevent concurrent update conflicts.

### Events

The module emits three event types:
- `ErrorReported` — new error or new occurrence of existing error
- `ErrorStatusChanged` — status transition (open, resolved, ignored, regressed)
- `ErrorAssigned` — error assigned to a user

## Outputs

- Stored error groups with occurrence counts, severity, and status
- Error events consumed by the Error Bridge for automated remediation
- Summary statistics via the API
- Error data exposed to MCP agents via `solodev://errors/active` resource

## API Endpoints

Base path: `/api/v1/spaces/{space_ref}/errors`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/errors` | Report an error |
| GET | `/errors` | List error groups (`?status=&severity=&language=&page=&limit=&query=`) |
| GET | `/errors/{id}` | Get error group detail (includes sample occurrences) |
| PATCH | `/errors/{id}` | Update error group (status, assigned_to, title, tags) |
| GET | `/errors/{id}/occurrences` | List occurrences for an error group |
| GET | `/errors/summary` | Get aggregate statistics |

Severity values: `fatal`, `error`, `warning`
Status values: `open`, `resolved`, `ignored`, `regressed`

## Database Schema

### `error_groups` table (migration `0171`)

| Column | Type | Description |
|--------|------|-------------|
| `eg_id` | SERIAL PK | Primary key |
| `eg_space_id` | INT | FK to spaces |
| `eg_repo_id` | INT | FK to repositories |
| `eg_identifier` | TEXT | Space-scoped identifier |
| `eg_fingerprint` | TEXT UNIQUE | SHA256 of title+stack_trace |
| `eg_status` | TEXT | `open`, `resolved`, `ignored`, `regressed` |
| `eg_severity` | TEXT | `fatal`, `error`, `warning` |
| `eg_occurrence_count` | BIGINT | Total occurrence count |
| `eg_file_path` | TEXT | Source file (optional) |
| `eg_language` | TEXT | Programming language (optional) |

### `error_occurrences` table (migration `0171`)

| Column | Type | Description |
|--------|------|-------------|
| `eo_id` | SERIAL PK | Primary key |
| `eo_error_group_id` | INT | FK to error_groups |
| `eo_stack_trace` | TEXT | Full stack trace |
| `eo_environment` | TEXT | `production`, `staging`, `development` |
| `eo_metadata` | JSON | Arbitrary context |

## Key Paths

| Purpose | Path |
|---------|------|
| Types | `types/errortracker.go` |
| Database store | `app/store/database/errortracker.go` |
| Controller | `app/api/controller/errortracker/controller.go` |
| Handlers | `app/api/handler/errortracker/` |
| Events | `app/events/errortracker/` |
| Migrations | `app/store/database/migrate/*/0171_create_table_error_groups.*` |

## Status

**Implemented** — Error groups, occurrences, fingerprinting, severity classification, status management, and Error Bridge integration are all working.

## Future Work

- Source map support for frontend error stack traces
- Error rate alerting thresholds
- Automatic resolution when a fix is applied and verified
