# Security Scanner

## Overview

The Security Scanner module provides capabilities for conducting security scans on code repositories, detecting vulnerabilities, secrets, and other security issues. It supports three scan types (SAST, SCA, Secret Detection) and tracks findings with severity levels and status management.

## Scan Types

| Type | Description |
|------|-------------|
| `sast` | Static Application Security Testing — analyze source code for vulnerabilities |
| `sca` | Software Composition Analysis — analyze dependencies for known vulnerabilities |
| `secret_detection` | Scan for accidentally committed secrets, tokens, and credentials |

## API Endpoints

### Scan Management

#### Trigger a Scan

**POST** `/api/v1/spaces/{space_ref}/security-scans`

```json
{
  "scan_type": "sast",
  "commit_sha": "abc123def456",
  "branch": "main",
  "triggered_by": "manual"
}
```

Trigger values: `manual`, `pipeline`, `webhook`

#### List Scans

**GET** `/api/v1/spaces/{space_ref}/security-scans`

Query parameters: `page`, `limit`, `sort`, `order`, `status`, `scan_type`, `triggered_by`

#### Get Scan Details

**GET** `/api/v1/spaces/{space_ref}/security-scans/{scan_identifier}`

### Finding Management

#### List Findings

**GET** `/api/v1/spaces/{space_ref}/security-scans/{scan_identifier}/findings`

Query parameters: `page`, `limit`, `sort`, `order`, `severity`, `category`, `status`

#### Update Finding Status

**PATCH** `/api/v1/spaces/{space_ref}/security-scans/findings/{id}/status`

```json
{
  "status": "resolved"
}
```

Finding status values: `open`, `resolved`, `ignored`, `false_positive`

> **Known limitation:** The handler (`update_finding_status.go`) and controller method (`UpdateFindingStatus`) exist in the codebase, but the route is not yet registered in the router (`app/router/api_modules.go`). The endpoint will return 404 until the route is added.

### Security Posture

#### Get Security Summary

**GET** `/api/v1/spaces/{space_ref}/security-scans/{scan_identifier}/summary`

Query parameters: `repo_ref` (optional)

Returns a `SecuritySummary` with aggregated counts by severity.

## Scan Status Lifecycle

```
pending → running → completed
                 ↘ failed
```

## Finding Severities

`critical` · `high` · `medium` · `low` · `info`

## Finding Categories

`vulnerability` · `secret` · `code_smell` · `bug`

## Database Schema

### `security_scans` table (migration `0102`)

| Column | Type | Description |
|--------|------|-------------|
| `ss_id` | BIGSERIAL PK | Primary key |
| `ss_space_id` | BIGINT | FK to spaces |
| `ss_repo_id` | BIGINT | FK to repositories |
| `ss_identifier` | TEXT | UUID unique identifier |
| `ss_scan_type` | TEXT | `sast`, `sca`, `secret_detection` |
| `ss_status` | TEXT | `pending`, `running`, `completed`, `failed` |
| `ss_commit_sha` | TEXT | Git commit SHA |
| `ss_branch` | TEXT | Git branch name |
| `ss_total_issues` | INTEGER | Total issues found |
| `ss_critical_count` | INTEGER | Critical severity count |
| `ss_high_count` | INTEGER | High severity count |
| `ss_medium_count` | INTEGER | Medium severity count |
| `ss_low_count` | INTEGER | Low severity count |
| `ss_duration` | BIGINT | Scan duration in milliseconds |
| `ss_triggered_by` | TEXT | `manual`, `pipeline`, `webhook` |
| `ss_created_by` | BIGINT | FK to principals |
| `ss_created` | BIGINT | Unix milliseconds |
| `ss_updated` | BIGINT | Unix milliseconds |
| `ss_version` | BIGINT | Optimistic locking version |

### `scan_findings` table (migration `0102`)

| Column | Type | Description |
|--------|------|-------------|
| `sf_id` | BIGSERIAL PK | Primary key |
| `sf_scan_id` | BIGINT | FK to security_scans |
| `sf_identifier` | TEXT | Unique finding identifier |
| `sf_severity` | TEXT | `critical`, `high`, `medium`, `low`, `info` |
| `sf_category` | TEXT | `vulnerability`, `secret`, `code_smell`, `bug` |
| `sf_title` | TEXT | Finding title |
| `sf_description` | TEXT | Detailed description |
| `sf_file_path` | TEXT | Path to affected file |
| `sf_line_start` | INTEGER | Starting line number |
| `sf_line_end` | INTEGER | Ending line number |
| `sf_rule` | TEXT | Rule that triggered the finding |
| `sf_snippet` | TEXT | Code snippet showing the issue |
| `sf_suggestion` | TEXT | Fix recommendation |
| `sf_status` | TEXT | `open`, `resolved`, `ignored`, `false_positive` |
| `sf_cwe` | TEXT | CWE identifier (if applicable) |
| `sf_created` | BIGINT | Unix milliseconds |
| `sf_updated` | BIGINT | Unix milliseconds |

## Access Control

| Operation | Permission Required |
|-----------|-------------------|
| List/Get scans | `PermissionRepoView` |
| Trigger scan / update results | `PermissionRepoPush` |
| Get individual finding | `PermissionSpaceView` |
| Get security summary | `PermissionSpaceView` |
| Update finding status | `PermissionSpaceEdit` |

## Events

Three events are published:
- `ScanTriggered` — when a scan starts
- `ScanCompleted` — when a scan finishes successfully
- `ScanFailed` — when a scan fails

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/securityscan.go` |
| Enums | `types/enum/securityscan.go` |
| Database store | `app/store/database/securityscan.go` |
| Store interfaces | `app/store/database.go` |
| Controller | `app/api/controller/securityscan/controller.go` |
| Handlers | `app/api/handler/securityscan/` |
| Request helpers | `app/api/request/securityscan.go` |
| Events | `app/events/securityscan/` |
| Migrations | `app/store/database/migrate/*/0102_create_table_security_scans.*` |
