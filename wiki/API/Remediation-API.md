# Remediation API

Detailed reference for the AI Auto-Remediation endpoints.

## Base Path

```
/api/v1/spaces/{space_ref}/remediations
```

## Trigger Remediation

**POST** `/remediations`

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

Trigger source values: `error_tracker`, `pipeline`, `security_scan`, `quality_gate`, `health_check`, `manual`

## Trigger Remediation from Security Finding

**POST** `/remediations/from-security-finding`

```json
{
  "repo_ref": "my-repo",
  "scan_identifier": "scan-123",
  "finding_id": 42
}
```

Returns `201 Created` when a new remediation is created, or `200 OK` when an existing remediation is reused.

## List Remediations

**GET** `/remediations`

Query parameters:

| Parameter | Description |
|-----------|-------------|
| `limit` | Results per page |

## Get Remediation

**GET** `/remediations/{remediation_identifier}`

Returns full remediation detail including patch diff, AI response, confidence score, and metadata.

## Update Remediation

**PATCH** `/remediations/{remediation_identifier}`

```json
{
  "status": "completed",
  "ai_response": "Analysis: NPE due to unchecked nil receiver...",
  "patch_diff": "--- a/handler.go\n+++ b/handler.go\n...",
  "fix_branch": "fix/npe-handler-12345",
  "pr_link": "https://github.com/org/repo/pull/42"
}
```

## Get Summary

**GET** `/remediations/summary`

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

## Get Metrics

**GET** `/remediations/metrics`

Query parameters:

| Parameter | Description |
|-----------|-------------|
| `window_days` | Optional positive integer lookback window (defaults to `30`) |

## Apply Remediation

**POST** `/remediations/{remediation_identifier}/apply`

Applies a completed remediation and returns the updated remediation record.

## Validate Remediation

**POST** `/remediations/{remediation_identifier}/validate`

```json
{
  "pipeline_identifier": "pipeline-main"
}
```

Triggers remediation validation and returns the updated remediation record.

## Remediation Record Fields

| Field | Type | Description |
|-------|------|-------------|
| `identifier` | string | Space-scoped unique identifier |
| `title` | string | Human-readable title |
| `status` | string | Lifecycle status |
| `trigger_source` | string | What triggered the task |
| `trigger_ref` | string | Error identifier or pipeline number |
| `error_log` | string | Full stack trace or build log |
| `file_path` | string | Source file involved |
| `source_code` | string | Relevant source code snippet |
| `branch` | string | Target branch |
| `commit_sha` | string | Commit that triggered failure |
| `ai_model` | string | AI model used |
| `ai_response` | string | Raw AI/LLM response |
| `patch_diff` | string | Generated unified diff |
| `confidence` | float | Confidence score 0.0–1.0 |
| `tokens_used` | int | Total LLM tokens consumed |
| `duration_ms` | int | Wall-clock time of AI generation |
