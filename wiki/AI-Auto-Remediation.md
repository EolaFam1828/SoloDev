# AI Auto-Remediation

## Overview

The AI Auto-Remediation module provides an automated error-to-fix pipeline. When a build fails, a test breaks, a security scan flags a vulnerability, or a runtime error appears, this module captures the full error context and creates a remediation task. An AI agent can then pick up that task, generate a code fix, create a patch, and optionally open a pull request — all without manual intervention.

## Remediation Lifecycle

```
  Error / Pipeline Failure
          │
          ▼
      ┌─────────┐
      │ PENDING  │  ← Task created (manually or via Error Bridge)
      └────┬─────┘
           │
           ▼
      ┌────────────┐
      │ PROCESSING │  ← AI agent picks up the task
      └─────┬──────┘
            │
      ┌─────┴──────┐
      ▼             ▼
┌───────────┐  ┌────────┐
│ COMPLETED │  │ FAILED │  ← AI could not generate a fix
└─────┬─────┘  └────────┘
      │
      ▼
  ┌─────────┐
  │ APPLIED │  ← Fix was pushed / PR merged
  └─────────┘
```

Status values: `pending` → `processing` → `completed` → `applied` → `failed` → `dismissed`

Trigger sources: `error_tracker`, `pipeline`, `security_scan`, `quality_gate`, `manual`

## API Endpoints

All endpoints are under `/api/v1/spaces/{space_ref}/remediations`.

### Trigger Remediation

**POST** `/api/v1/spaces/{space_ref}/remediations`

```json
{
  "title": "Fix: null pointer in handler",
  "description": "NPE raised in GET /api/users",
  "trigger_source": "error_tracker",
  "trigger_ref": "err-npe-12345",
  "error_log": "panic: runtime error: invalid memory address...",
  "file_path": "app/api/handler/users/find.go",
  "branch": "main",
  "commit_sha": "abc123"
}
```

### List Remediations

**GET** `/api/v1/spaces/{space_ref}/remediations`

Query parameters: `status`, `trigger_source`, `page`, `limit`

### Get Remediation

**GET** `/api/v1/spaces/{space_ref}/remediations/{remediation_identifier}`

### Update Remediation

**PATCH** `/api/v1/spaces/{space_ref}/remediations/{remediation_identifier}`

```json
{
  "status": "completed",
  "ai_response": "Analysis: NPE due to unchecked nil receiver...",
  "patch_diff": "--- a/handler.go\n+++ b/handler.go\n...",
  "fix_branch": "fix/npe-handler-12345",
  "pr_link": "https://github.com/org/repo/pull/42"
}
```

### Get Summary

**GET** `/api/v1/spaces/{space_ref}/remediations/summary`

```json
{
  "total": 24,
  "pending": 3,
  "processing": 1,
  "completed": 5,
  "applied": 12,
  "failed": 2,
  "dismissed": 1
}
```

## Database Schema

Table: `remediations` (migration `0172`)

| Column | Type | Description |
|--------|------|-------------|
| `rem_id` | SERIAL PK | Primary key |
| `rem_space_id` | INT | FK to spaces |
| `rem_repo_id` | INT | FK to repositories |
| `rem_identifier` | TEXT | Space-scoped unique identifier |
| `rem_title` | TEXT | Human-readable title |
| `rem_description` | TEXT | Detailed description |
| `rem_status` | TEXT | Lifecycle status |
| `rem_trigger_source` | TEXT | What triggered the task |
| `rem_trigger_ref` | TEXT | Error identifier or pipeline number |
| `rem_error_log` | TEXT | Full stack trace or build log |
| `rem_file_path` | TEXT | Source file involved |
| `rem_source_code` | TEXT | Relevant source code snippet |
| `rem_branch` | TEXT | Target branch |
| `rem_commit_sha` | TEXT | Commit that triggered failure |
| `rem_ai_model` | TEXT | AI model used |
| `rem_ai_prompt` | TEXT | Prompt sent to the LLM |
| `rem_ai_response` | TEXT | Raw AI/LLM response |
| `rem_patch_diff` | TEXT | Generated unified diff |
| `rem_fix_branch` | TEXT | Branch where fix was pushed |
| `rem_pr_link` | TEXT | Link to pull request |
| `rem_confidence` | REAL | Confidence score 0.0–1.0 |
| `rem_tokens_used` | BIGINT | Total LLM tokens consumed |
| `rem_duration_ms` | BIGINT | Wall-clock time of AI generation |
| `rem_metadata` | TEXT | Arbitrary JSON blob |
| `rem_created_by` | INT | FK to principals |
| `rem_created` | BIGINT | Unix milliseconds |
| `rem_updated` | BIGINT | Unix milliseconds |
| `rem_version` | BIGINT | Optimistic locking version |

**Indexes:** `(rem_space_id, rem_status)`, `(rem_space_id, rem_trigger_source)`, `(rem_created DESC)`

## Events

Three events are published:
- `RemediationTriggered` — when a task is created
- `RemediationCompleted` — when AI generation finishes
- `RemediationApplied` — when the fix is pushed/merged

## Integration with Error Bridge

The [Error Bridge](Error-Bridge) automatically creates remediation tasks when:
- A runtime error is reported via the Error Tracker (severity ≥ `error`)
- A pipeline execution fails

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/ai_remediation.go` |
| Database store | `app/store/database/ai_remediation.go` |
| Controller | `app/api/controller/airemediation/controller.go` |
| Handlers | `app/api/handler/airemediation/` |
| Events | `app/events/airemediation/` |
| Request helpers | `app/api/request/airemediation.go` |
| Migrations | `app/store/database/migrate/*/0172_create_table_remediations.*` |
| Router registration | `app/router/api_modules.go` (`setupRemediations()`) |
