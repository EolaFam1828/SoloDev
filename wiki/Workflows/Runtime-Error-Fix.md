# Runtime Error Fix

A step-by-step example showing how a runtime error reported to the Error Tracker leads to AI-generated remediation.

## Scenario

A production application reports a database connection timeout via the Error Tracker API. The error includes a stack trace pointing to `internal/db/connect.go:42`.

## Step 1: Report the Error

The application (or monitoring system) sends an error report:

```bash
curl -X POST http://localhost:3000/api/v1/spaces/my-space/errors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "db-connection-timeout",
    "title": "Database Connection Timeout",
    "message": "Failed to connect to database within 30s",
    "severity": "error",
    "file_path": "internal/db/connect.go",
    "line_number": 42,
    "function_name": "Connect",
    "language": "go",
    "tags": ["database", "critical"],
    "stack_trace": "goroutine 1 [running]:\ninternal/db/connect.go:42 +0x1a\nmain.go:15 +0x2b",
    "environment": "production",
    "runtime": "go1.24"
  }'
```

## Step 2: Error Group Created

The Error Tracker:
1. Generates a fingerprint from `SHA256(title + stack_trace)`
2. Creates a new error group (or increments occurrence count if fingerprint matches)
3. Stores the occurrence with full stack trace and metadata
4. Publishes `ErrorReported` event

## Step 3: Error Bridge Triggers

The Error Bridge receives `OnErrorReported()`:
1. Checks severity — `error` meets the threshold (≥ `error`)
2. Checks status — group is `open` (not resolved or ignored)
3. Creates a pending remediation task:

| Field | Value |
|-------|-------|
| `title` | `[Auto] Fix: Database Connection Timeout` |
| `trigger_source` | `error_tracker` |
| `trigger_ref` | `db-connection-timeout` |
| `error_log` | Stack trace from occurrence |
| `file_path` | `internal/db/connect.go` |

## Step 4: AI Worker Processes

Within 15 seconds, the AI Worker:
1. Picks up the pending task
2. Builds prompt with the error context and source code
3. Sends to the LLM
4. Receives a patch that adds connection retry logic with exponential backoff
5. Stores the diff with confidence score 0.78

## Step 5: Developer Reviews

The developer sees the remediation in the dashboard:
- Error: Database Connection Timeout
- Proposed fix: Add retry logic with backoff
- Confidence: 0.78

The developer reviews the patch, adjusts the retry count, applies it, and pushes. The error group can be marked as `resolved`.

## Error Bridge Filtering

The bridge would have skipped remediation creation if:
- Severity was `warning` (threshold is `error` or `fatal`)
- The error group was already `resolved` or `ignored`
- The bridge was disabled

## MCP Agent Alternative

An MCP-connected agent can handle this workflow:

```
1. Agent reads solodev://errors/active → sees new error
2. Agent calls fix_this(error_id="db-connection-timeout")
3. fix_this triggers remediation, monitors progress
4. Agent retrieves the patch via remediation_get()
5. Agent applies the patch to the codebase
```
