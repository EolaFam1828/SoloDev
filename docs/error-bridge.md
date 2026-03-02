# Error-to-AI Bridge Module

## Overview

The Error Bridge is the connective tissue between the Error Tracker, Pipeline Execution, and the AI Auto-Remediation system. It watches for failures and automatically creates pending remediation tasks with the full error context — stack traces, file paths, severity, and commit information — ready for an AI agent to pick up and fix.

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

When `ReportError()` is called on the error tracker controller:
1. Error group is created/updated (existing flow)
2. Event is reported (existing flow)
3. **NEW:** If `errorBridge` is set, `OnErrorReported()` is called
4. Bridge creates a `Remediation` with:
   - `TriggerSource = "error_tracker"`
   - `TriggerRef = errorGroup.Identifier`
   - `ErrorLog = occurrence.StackTrace`
   - `FilePath = errorGroup.FilePath`
   - `Title = "[Auto] Fix: {errorGroup.Title}"`

### From Pipeline Failures

`OnPipelineFailed()` can be called from pipeline completion hooks:
1. Bridge creates a `Remediation` with:
   - `TriggerSource = "pipeline"`
   - `TriggerRef = executionNumber`
   - `ErrorLog = build log output`
   - `Branch`, `CommitSHA` from the failed execution

## Configuration

The bridge is configurable per-instance:
- `enabled` — master on/off toggle (passed to `NewBridge()`)
- `autoTriggerFatal` — when true, only auto-trigger on fatal/error severity (skip warnings). Hardcoded to `true` in `NewBridge()`; not externally configurable

## Filtering

The bridge skips creating remediations when:
- Bridge is disabled (`enabled = false`)
- Error severity is warning and `autoTriggerFatal` is true
- Error group status is already `resolved` or `ignored`

## Architecture

### Bridge (`app/services/errorbridge/bridge.go`)

```go
type Bridge struct {
    remediationStore store.RemediationStore
    enabled          bool
    autoTriggerFatal bool
}
```

**Methods:**
- `NewBridge(store, enabled)` — constructor
- `OnErrorReported(ctx, errorGroup, occurrence)` — called from error tracker controller
- `OnPipelineFailed(ctx, spaceID, repoID, executionNumber, branch, commitSHA, errorLog, createdBy)` — called from pipeline hooks

### Error Tracker Integration

The `errortracker.Controller` was modified to include:
- `errorBridge *errorbridge.Bridge` field
- `SetErrorBridge(bridge)` method for optional injection
- Auto-trigger call in `ReportError()` after event reporting

## File Structure

```
├── app/services/errorbridge/
│   └── bridge.go                       # Bridge service
└── app/api/controller/errortracker/
    └── controller.go                   # Modified: SetErrorBridge + auto-trigger
```

## Usage

```go
// Create the bridge
bridge := errorbridge.NewBridge(remediationStore, true)

// Attach to error tracker controller
errorTrackerCtrl.SetErrorBridge(bridge)

// Now every error report will auto-create a remediation task
// No manual intervention needed!
```
