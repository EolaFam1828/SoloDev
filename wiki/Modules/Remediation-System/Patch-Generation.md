# Patch Generation

## Purpose

Patch Generation is the process by which the AI Worker transforms error context and source code into a concrete code change. The output is a unified diff that can be applied to the repository.

## Inputs

- Error log or stack trace
- Source code of the affected file
- File path, branch, and commit SHA
- Language and framework metadata

## Processing

### Prompt Construction

The AI Worker builds a prompt containing:

1. **System instructions** — Directs the LLM to analyze the error, identify the root cause, and produce a unified diff with a confidence score
2. **Error context** — The full error log or stack trace
3. **Source code** — The content of the affected file
4. **Metadata** — File path, branch, commit SHA, language

### LLM Generation

The assembled prompt is sent to the configured LLM provider. The LLM produces:

1. **Analysis** — Explanation of the root cause
2. **Unified diff** — A `patch -p1` compatible code change
3. **Confidence score** — A 0.0–1.0 estimate of fix reliability

### Response Parsing

The parser (`app/services/aiworker/parser.go`) extracts:

- The diff block from markdown code fences
- The confidence score as a floating-point number
- The full AI response text

### Diff Format

Generated patches follow the standard unified diff format:

```diff
--- a/app/handler.go
+++ b/app/handler.go
@@ -40,6 +40,9 @@ func HandleRequest(w http.ResponseWriter, r *http.Request) {
     user := GetUser(r)
+    if user == nil {
+        http.Error(w, "user not found", http.StatusNotFound)
+        return
+    }
     user.Process()
```

## Outputs

- Unified diff stored in `rem_patch_diff`
- AI analysis stored in `rem_ai_response`
- Confidence score stored in `rem_confidence`
- Model identifier stored in `rem_ai_model`

## Key Paths

| Purpose | Path |
|---------|------|
| AI Worker | `app/services/aiworker/worker.go` |
| Response parser | `app/services/aiworker/parser.go` |

## Status

**Implemented** — LLM-driven patch generation from error context is working across all four supported providers (Anthropic, OpenAI, Gemini, Ollama).

## Future Work

- Multi-file patch generation for cross-file bugs
- Iterative patch refinement (generate → test → refine)
- Patch quality metrics and tracking
