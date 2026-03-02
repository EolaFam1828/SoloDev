# AI Auto-Remediation Module

## Overview

The AI Auto-Remediation module provides an automated error-to-fix pipeline for the SoloDev platform. When a build fails, a test breaks, a security scan flags a vulnerability, or a runtime error appears, this module captures the full error context and creates a remediation task that an AI agent can pick up, generate a code fix, create a patch, and optionally open a pull request — all without manual intervention.

## Architecture

### Components

#### 1. Types (`types/ai_remediation.go`)
Core data structures:
- **Remediation**: Full remediation task with status tracking, error context, and AI-generated outputs
- **RemediationStatus**: Lifecycle enum (`pending` → `analyzing` → `fix_generated` → `applied` → `failed` → `skipped`)
- **RemediationTriggerSource**: What triggered the remediation (`error`, `pipeline`, `security`, `quality`, `manual`)
- **TriggerRemediationInput**: API request body for manually triggering a remediation
- **UpdateRemediationInput**: API request body for updating status, AI response, patch diff, etc.
- **RemediationListFilter**: Filtering options for listing remediations
- **RemediationSummary**: Aggregate statistics

#### 2. Database Store (`app/store/database/ai_remediation.go`)
Implements CRUD operations using squirrel query builder:
- `Create()`: Insert new remediation task
- `Find()`: Retrieve by ID
- `FindByIdentifier()`: Retrieve by space-scoped identifier
- `List()`: Paginated listing with filtering by status, trigger source, repo
- `Count()`: Count with filter
- `Update()`: Full update with version increment
- `UpdateStatus()`: Status-only update
- `Summary()`: Aggregate statistics (counts by status)

#### 3. Store Interface (`app/store/database.go`)
Added `RemediationStore` interface to the central store definitions.

#### 4. Controller (`app/api/controller/airemediation/controller.go`)
Business logic with space-scoped authorization:
- `Trigger()`: Manually create a remediation task
- `List()`: List remediations with filtering
- `Get()`: Get single remediation with full detail
- `Update()`: Update status, AI response, patch diff, fix branch, PR link
- `Summary()`: Get aggregate statistics

#### 5. Handlers (`app/api/handler/airemediation/`)
- `trigger.go`: POST handler for creating remediations
- `list.go`: GET handler for listing
- `get.go`: GET handler for detail
- `update.go`: PATCH handler for updates
- `summary.go`: GET handler for summary statistics

#### 6. Events (`app/events/airemediation/`)
- `events.go`: Event data structures (RemediationTriggered, RemediationCompleted, RemediationApplied)
- `reporter.go`: Event reporter (3 methods)
- `reader.go`: Event reader

#### 7. Request Helpers (`app/api/request/airemediation.go`)
Path parameter extraction for `remediation_identifier`.

#### 8. Migrations (0172)
- PostgreSQL: `0172_create_table_remediations.up.sql` / `.down.sql`
- SQLite: Same migration number

## Database Schema

### remediations Table
```sql
- rem_id              (SERIAL PRIMARY KEY)
- rem_space_id        (INT, FK to spaces)
- rem_repo_id         (INT, FK to repositories)
- rem_identifier      (TEXT, unique per space)
- rem_title           (TEXT)
- rem_description     (TEXT)
- rem_status          (TEXT: pending|analyzing|fix_generated|applied|failed|skipped)
- rem_trigger_source  (TEXT: error|pipeline|security|quality|manual)
- rem_trigger_ref     (TEXT, e.g. error identifier or pipeline number)
- rem_error_log       (TEXT, full stack trace or build log)
- rem_file_path       (TEXT, source file involved)
- rem_source_code     (TEXT, relevant source code snippet)
- rem_branch          (TEXT, target branch)
- rem_commit_sha      (TEXT, commit that triggered failure)
- rem_ai_response     (TEXT, raw AI/LLM response)
- rem_patch_diff      (TEXT, generated unified diff)
- rem_fix_branch      (TEXT, branch where fix was pushed)
- rem_pr_link         (TEXT, link to pull request)
- rem_created_by      (INT, FK to principals)
- rem_created         (BIGINT, unix milliseconds)
- rem_updated         (BIGINT, unix milliseconds)
- rem_version         (BIGINT, optimistic locking)
```

**Indexes:**
- `idx_remediations_space_status` on `(rem_space_id, rem_status)`
- `idx_remediations_trigger` on `(rem_space_id, rem_trigger_source, rem_trigger_ref)`

## API Endpoints

All under `/api/v1/spaces/{space_ref}/remediations`.

### 1. Trigger Remediation
**POST** `/api/v1/spaces/{space_ref}/remediations`

```json
{
  "title": "Fix: null pointer in handler",
  "description": "NPE raised in GET /api/users",
  "trigger_source": "error",
  "trigger_ref": "err-npe-12345",
  "error_log": "panic: runtime error: invalid memory address...",
  "file_path": "app/api/handler/users/find.go",
  "branch": "main",
  "commit_sha": "abc123"
}
```

### 2. List Remediations
**GET** `/api/v1/spaces/{space_ref}/remediations`

Query: `?status=pending&trigger_source=error&page=0&limit=50`

### 3. Get Remediation
**GET** `/api/v1/spaces/{space_ref}/remediations/{remediation_identifier}`

### 4. Update Remediation
**PATCH** `/api/v1/spaces/{space_ref}/remediations/{remediation_identifier}`

```json
{
  "status": "fix_generated",
  "ai_response": "Analysis: NPE due to unchecked nil receiver...",
  "patch_diff": "--- a/handler.go\n+++ b/handler.go\n@@ -42 +42 @@\n-  user.Name\n+  if user != nil { user.Name }",
  "fix_branch": "fix/npe-handler-12345",
  "pr_link": "https://github.com/org/repo/pull/42"
}
```

### 5. Get Summary
**GET** `/api/v1/spaces/{space_ref}/remediations/summary`

```json
{
  "total": 24,
  "pending": 3,
  "analyzing": 1,
  "fix_generated": 5,
  "applied": 12,
  "failed": 2,
  "skipped": 1
}
```

## Remediation Lifecycle

```
  Error/Pipeline Failure
          │
          ▼
      ┌─────────┐
      │ PENDING  │  ← Task created (manually or via error bridge)
      └────┬─────┘
           │
           ▼
      ┌──────────┐
      │ ANALYZING│  ← AI agent picks up the task
      └────┬─────┘
           │
     ┌─────┴──────┐
     ▼             ▼
┌────────────┐  ┌────────┐
│FIX_GENERATED│  │ FAILED │  ← AI could not generate a fix
└─────┬──────┘  └────────┘
      │
      ▼
  ┌─────────┐
  │ APPLIED │  ← Fix was pushed / PR merged
  └─────────┘
```

## Integration with Error Bridge

The Error Bridge (`app/services/errorbridge/bridge.go`) automatically creates remediation tasks when:
- A runtime error is reported via the Error Tracker (severity ≥ error)
- A pipeline execution fails

See `ERROR_BRIDGE_MODULE.md` for details.

## File Structure

```
├── types/
│   └── ai_remediation.go
├── app/
│   ├── api/
│   │   ├── controller/airemediation/
│   │   │   └── controller.go
│   │   ├── handler/airemediation/
│   │   │   ├── trigger.go
│   │   │   ├── list.go
│   │   │   ├── get.go
│   │   │   ├── update.go
│   │   │   └── summary.go
│   │   └── request/
│   │       └── airemediation.go
│   ├── events/airemediation/
│   │   ├── events.go
│   │   ├── reporter.go
│   │   └── reader.go
│   └── store/
│       └── database/
│           ├── ai_remediation.go
│           └── migrate/
│               ├── postgres/0172_create_table_remediations.{up,down}.sql
│               └── sqlite/0172_create_table_remediations.{up,down}.sql
└── app/router/
    └── api_modules.go          # setupRemediations()
```

## Dependencies

- `store.RemediationStore` (injected via wire)
- `authz.Authorizer` for space-scoped RBAC
- `refcache.SpaceFinder` for space resolution
- `errortrackerevents.Reporter` for event publishing (optional)
- `dbtx.Transactor` for transactions
