# First Remediation

See the AI remediation loop in action.

## Prerequisites

- SoloDev running (see [Quick Start](Quick-Start))
- A personal access token (see [Installation](Installation))
- An LLM provider configured (Anthropic, OpenAI, Gemini, or Ollama)

## Step 1: Report an Error

Post an error to the Error Tracker:

```bash
TOKEN="<your-personal-access-token>"

curl -X POST http://localhost:3000/api/v1/spaces/<space>/errors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "nil pointer in handler",
    "message": "runtime error: invalid memory address",
    "severity": "error",
    "language": "go",
    "file_path": "app/handler.go",
    "stack_trace": "goroutine 1:\napp/handler.go:42 +0x1a"
  }'
```

## Step 2: Verify Bridge Trigger

If the Error Bridge is enabled, a remediation task is automatically created. Check:

```bash
curl http://localhost:3000/api/v1/spaces/<space>/remediations?status=pending \
  -H "Authorization: Bearer $TOKEN"
```

You should see a remediation task with `trigger_source: "error_tracker"`.

## Step 3: Wait for AI Processing

The AI Worker polls every 15 seconds. Within a minute, the remediation should move from `pending` → `processing` → `completed`.

```bash
curl http://localhost:3000/api/v1/spaces/<space>/remediations?status=completed \
  -H "Authorization: Bearer $TOKEN"
```

## Step 4: Review the Patch

Get the completed remediation:

```bash
curl http://localhost:3000/api/v1/spaces/<space>/remediations/<identifier> \
  -H "Authorization: Bearer $TOKEN"
```

The response includes:
- `ai_response` — Analysis of the root cause
- `patch_diff` — The proposed code fix as a unified diff
- `confidence` — How confident the AI is in the fix (0.0–1.0)

## Step 5: Apply the Fix

Review the diff. If it looks correct, apply it to your codebase:

```bash
# Apply the diff
echo '<patch_diff_content>' | patch -p1
```

Then update the remediation status:

```bash
curl -X PATCH http://localhost:3000/api/v1/spaces/<space>/remediations/<identifier> \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "applied"}'
```

## Troubleshooting

**Remediation stays `pending`**
The AI Worker requires an LLM provider to be configured. Check that the provider API key is set in the server configuration.

**No remediation created**
The Error Bridge may not be enabled, or the error severity may be `warning` (threshold is `error` or `fatal`).

**Patch quality is poor**
Try a more capable model. Ensure the error report includes a detailed stack trace and file path for better context.
