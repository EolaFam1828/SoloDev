# Error-to-AI Bridge

## Overview

The Error Bridge is the connective tissue between the Error Tracker, Pipeline Execution, and the AI Auto-Remediation system. It watches for failures and automatically creates pending remediation tasks with full error context — stack traces, file paths, severity, and commit information — ready for an AI agent to pick up and fix.

## How It Works

```
  ┌──────────────────┐          ┌──────────────┐
  │  Error Tracker    │──────▶  │              │
  │  (ReportError)    │         │  Error Bridge │──────▶  Remediation
  └──────────────────┘          │              │           (pending)
                                │              │
  ┌──────────────────┐          │              │
  │  Pipeline Runner  │──────▶  │              │
  │  (OnFailure)      │         └──────────────┘
  └──────────────────┘
```

### From Error Tracker

When `ReportError()` is called:
1. Error group is created/updated (standard flow)
2. Events are published (standard flow)
3. If a bridge is attached, `OnErrorReported()` is called
4. Bridge creates a `Remediation` with:
   - `TriggerSource = "error_tracker"`
   - `TriggerRef = errorGroup.Identifier`
   - `ErrorLog = occurrence.StackTrace`
   - `FilePath = errorGroup.FilePath`
   - `Title = "[Auto] Fix: {errorGroup.Title}"`

### From Pipeline Failures

`OnPipelineFailed()` can be called from pipeline completion hooks. The bridge creates a `Remediation` with:
- `TriggerSource = "pipeline"`
- `TriggerRef = executionNumber`
- `ErrorLog = build log output`
- `Branch`, `CommitSHA` from the failed execution

## Filtering

The bridge skips creating remediations when:
- Bridge is disabled (`enabled = false`)
- Error severity is `warning` and `autoTriggerFatal` is `true`
- Error group status is already `resolved` or `ignored`

## Configuration

| Setting | Description |
|---------|-------------|
| `enabled` | Master on/off toggle (passed to `NewBridge()`) |
| `autoTriggerFatal` | When true, only auto-trigger on `fatal`/`error` severity (skip `warning`). Hardcoded to `true` in `NewBridge()`; not externally configurable. |

## Bridge Struct

```go
type Bridge struct {
    remediationStore store.RemediationStore
    enabled          bool
    autoTriggerFatal bool
}
```

### Methods

| Method | Description |
|--------|-------------|
| `NewBridge(store, enabled)` | Constructor |
| `OnErrorReported(ctx, errorGroup, occurrence)` | Called from error tracker controller |
| `OnPipelineFailed(ctx, spaceID, repoID, executionNumber, branch, commitSHA, errorLog, createdBy)` | Called from pipeline hooks |

## Setup

```go
// Create the bridge
bridge := errorbridge.NewBridge(remediationStore, true)

// Attach to error tracker controller
errorTrackerCtrl.SetErrorBridge(bridge)

// Now every error report will auto-create a remediation task
```

## File Locations

| Purpose | Path |
|---------|------|
| Bridge service | `app/services/errorbridge/bridge.go` |
| Error tracker controller (modified) | `app/api/controller/errortracker/controller.go` |

## Related Pages

- [Error Tracker](Error-Tracker) — the module that feeds into the bridge
- [AI Auto-Remediation](AI-Auto-Remediation) — the module that receives bridge output
