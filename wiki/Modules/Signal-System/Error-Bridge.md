# Error Bridge

## Purpose

The Error Bridge is the connective tissue between the Signal System and the AI Layer. It normalizes errors from the Error Tracker and pipeline failures into structured remediation tasks that the AI Worker can process. It is the "Analyze" stage of the SoloDev Loop.

## Inputs

- Error group and occurrence data from the Error Tracker (`OnErrorReported`)
- Pipeline failure data from the Pipeline Runner (`OnPipelineFailed`)
- Space ID, repo ID, branch, commit SHA, error log, and file path

## Processing

### From Error Tracker

When `ReportError()` is called on the Error Tracker controller:
1. Error group is created/updated (standard Error Tracker flow)
2. Events are published (standard flow)
3. If a bridge is attached, `OnErrorReported()` is called
4. Bridge creates a `Remediation` with:
   - `TriggerSource = "error_tracker"`
   - `TriggerRef = errorGroup.Identifier`
   - `ErrorLog = occurrence.StackTrace`
   - `FilePath = errorGroup.FilePath`
   - `Title = "[Auto] Fix: {errorGroup.Title}"`

### From Pipeline Failures

`OnPipelineFailed()` is called from pipeline completion hooks:
- `TriggerSource = "pipeline"`
- `TriggerRef = executionNumber`
- `ErrorLog = build log output`
- `Branch`, `CommitSHA` from the failed execution

### Filtering

The bridge skips creating remediations when:
- Bridge is disabled (`enabled = false`)
- Error severity is `warning` and `autoTriggerFatal` is `true`
- Error group status is already `resolved` or `ignored`

## Outputs

- Pending remediation tasks in the database, ready for the AI Worker
- Each task contains full error context for LLM processing

## Configuration

| Setting | Description |
|---------|-------------|
| `enabled` | Master on/off toggle (passed to `NewBridge()`) |
| `autoTriggerFatal` | When true, only auto-trigger on `fatal`/`error` severity. Hardcoded to `true`. |

## Key Paths

| Purpose | Path |
|---------|------|
| Bridge service | `app/services/errorbridge/bridge.go` |
| Error tracker controller (integration point) | `app/api/controller/errortracker/controller.go` |

## Status

**Implemented** — Auto-creates remediation tasks from runtime errors and pipeline failures. Severity filtering and status-based skip logic are working.

## Future Work

- Configurable severity threshold (expose `autoTriggerFatal` to users)
- Deduplication of remediation tasks for the same root cause
- Rate limiting to prevent remediation storms from repeated failures
