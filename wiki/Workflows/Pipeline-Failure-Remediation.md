# Pipeline Failure Remediation

A step-by-step walkthrough of how SoloDev detects and fixes a pipeline failure, demonstrating the full SoloDev Loop.

## Scenario

A Go project has a pipeline that runs `go test ./...`. A developer pushes a commit that introduces a nil pointer dereference in `app/handler.go`. The pipeline fails.

## Step 1: Detect

The pipeline executes and the test step fails:

```
--- FAIL: TestHandleRequest (0.00s)
    handler_test.go:25: runtime error: invalid memory address or nil pointer dereference
FAIL    app/handler  0.003s
```

Pipeline status changes to `failed`. The failure event includes the build log, branch (`main`), and commit SHA (`abc123`).

## Step 2: Analyze

The Error Bridge receives the pipeline failure event via `OnPipelineFailed()`:

```go
bridge.OnPipelineFailed(ctx, spaceID, repoID, "42", "main", "abc123", buildLog, createdBy)
```

The bridge creates a remediation task:

| Field | Value |
|-------|-------|
| `title` | `[Auto] Fix: Pipeline #42 failure` |
| `trigger_source` | `pipeline` |
| `trigger_ref` | `42` |
| `error_log` | Full build log output |
| `branch` | `main` |
| `commit_sha` | `abc123` |
| `status` | `pending` |

## Step 3: Propose

The AI Worker poller picks up the pending task (within 15 seconds):

1. Marks status → `processing`
2. Builds prompt with error log, file path, and source code
3. Sends to configured LLM (e.g., Anthropic Claude)
4. LLM returns:

```
Analysis: The test fails because GetUser() can return nil when the user
is not found, but the handler does not check for nil before calling
user.Process(). Adding a nil check resolves the panic.

Confidence: 0.92

```diff
--- a/app/handler.go
+++ b/app/handler.go
@@ -40,6 +40,10 @@ func HandleRequest(w http.ResponseWriter, r *http.Request) {
     user := GetUser(r)
+    if user == nil {
+        http.Error(w, "user not found", http.StatusNotFound)
+        return
+    }
     user.Process()
```

## Step 4: Validate

The parser extracts:
- Unified diff patch
- Confidence score: 0.92

Status updated to `completed`. The patch is stored in `rem_patch_diff`.

## Step 5: Apply

The developer (or agent) views the remediation via dashboard, API, or MCP:

```bash
# Via API
curl http://localhost:3000/api/v1/spaces/my-space/remediations/rem-42 \
  -H "Authorization: Bearer $TOKEN"

# Via MCP
remediation_get(identifier="rem-42")
```

The developer reviews the diff, applies it to the codebase, and pushes. The pipeline re-runs and passes.

## Current vs. Planned

| Step | Current | Planned |
|------|---------|---------|
| Detect | Automatic | Automatic |
| Analyze | Automatic (Error Bridge) | Automatic |
| Propose | Automatic (AI Worker) | Automatic |
| Validate | Confidence score only | Full validation (tests, security scan) |
| Apply | Manual (developer copies diff) | Auto-PR, auto-merge, auto-re-run |
