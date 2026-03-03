# AI Auto-Remediation

## Purpose

The AI Auto-Remediation module provides the automated error-to-fix pipeline. When a build fails, a test breaks, a security scan flags a vulnerability, or a runtime error appears, this module captures the full error context, creates a remediation task, sends it to an LLM, stores a unified diff, and can deliver that diff into a fix branch plus draft PR.

## Inputs

- Remediation task creation requests (manual or via Error Bridge / Security Remediation)
- Each task includes: title, description, trigger source, trigger ref, error log, file path, source code, branch, commit SHA

## Processing

### Remediation Lifecycle

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
      │ PROCESSING │  ← AI Worker picks up the task
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

Status values: `pending` → `processing` → `completed` → `applied` / `failed` / `dismissed`

Trigger sources: `error_tracker`, `pipeline`, `security_scan`, `quality_gate`, `manual`

### AI Worker Flow

1. Background poller runs every 15 seconds
2. Queries for `pending` remediations
3. Marks task as `processing`
4. Builds prompt from error context + source code
5. Sends to configured LLM provider
6. Parses response for unified diff + confidence score
7. Stores result and updates status
8. Optionally applies the diff onto `solodev/rem-<identifier>` and opens a draft PR

## Outputs

- Patch diff (unified diff, `patch -p1` compatible)
- AI response (analysis and explanation)
- Confidence score (0.0–1.0)
- Token usage and processing duration metrics
- Fix branch and PR link once delivery runs
- Events: `RemediationTriggered`, `RemediationCompleted`, `RemediationApplied`

## API Endpoints

Base path: `/api/v1/spaces/{space_ref}/remediations`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/remediations` | Trigger a remediation task |
| POST | `/remediations/from-security-finding` | Trigger or reuse a remediation from a security finding |
| GET | `/remediations` | List remediations (`?limit=`) |
| GET | `/remediations/{remediation_identifier}` | Get remediation detail |
| PATCH | `/remediations/{remediation_identifier}` | Update remediation (status, AI response, patch diff, fix branch, PR link) |
| POST | `/remediations/{remediation_identifier}/apply` | Apply a completed remediation into a fix branch and draft PR |
| POST | `/remediations/{remediation_identifier}/validate` | Trigger remediation validation (optional `pipeline_identifier`) |
| GET | `/remediations/metrics` | Get time-windowed remediation metrics (`?window_days=`) |
| GET | `/remediations/summary` | Get aggregate statistics |

## Database Schema

Table: `remediations` (migration `0172`)

Key columns: `rem_status`, `rem_trigger_source`, `rem_error_log`, `rem_file_path`, `rem_source_code`, `rem_branch`, `rem_commit_sha`, `rem_ai_model`, `rem_ai_response`, `rem_patch_diff`, `rem_fix_branch`, `rem_pr_link`, `rem_metadata`, `rem_confidence`, `rem_tokens_used`, `rem_duration_ms`

## Key Paths

| Purpose | Path |
|---------|------|
| Types | `types/ai_remediation.go` |
| Database store | `app/store/database/ai_remediation.go` |
| Controller | `app/api/controller/airemediation/controller.go` |
| Handlers | `app/api/handler/airemediation/` |
| AI Worker | `app/services/aiworker/worker.go` |
| Response parser | `app/services/aiworker/parser.go` |
| Events | `app/events/airemediation/` |
| Migrations | `app/store/database/migrate/*/0172_create_table_remediations.*` |

## Status

**Implemented** — CRUD endpoints, AI Worker with LLM integration, source enrichment, Error Bridge and Security Remediation triggers, patch generation, confidence scoring, manual delivery to draft PR, and opt-in automatic draft PR delivery are all working.

## Future Work

- Auto-merge for high-confidence patches
- Self-healing pipeline loop (fix → re-run → verify)
- Remediation success rate tracking
