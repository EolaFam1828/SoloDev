# Approval Workflow

## Purpose

The Approval Workflow defines the gates between AI-generated patches and their application to the codebase. It determines whether a fix proceeds automatically, requires human review, or is discarded.

## Inputs

- Completed remediation task with patch diff and confidence score
- Solo Gate enforcement mode for the space
- Configurable approval thresholds (planned)

## Processing

### Current Implementation

The current approval workflow is manual:

1. AI Worker generates a patch and stores it with status `completed`
2. Developer views the patch via web dashboard, REST API, or MCP tools
3. Developer manually applies the diff to the codebase
4. Developer updates the remediation status to `applied`

### Planned: Automated Approval

The planned workflow introduces automation gates:

```
Completed Remediation
        │
        ▼
┌─────────────────┐
│ Confidence Check │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
   High      Low
    │         │
    ▼         ▼
 Auto-PR   Flag for
 + merge   review
```

1. **High confidence + strict mode** — Create PR, require manual merge
2. **High confidence + balanced mode** — Create PR, auto-merge if tests pass
3. **High confidence + prototype mode** — Auto-merge directly
4. **Low confidence (any mode)** — Create PR, require manual review
5. **Very low confidence** — Mark as failed, do not create PR

### Solo Gate Integration

The Solo Gate enforcement mode influences the approval workflow:

| Mode | High Confidence | Low Confidence |
|------|----------------|----------------|
| Strict | PR with manual merge | PR with manual review |
| Balanced | PR with auto-merge | PR with manual review |
| Prototype | Direct apply | PR with auto-merge |

## Outputs

- Approval decision (auto-apply, create PR, require review, discard)
- PR creation with AI analysis as description (planned)
- Status update on the remediation task

## Key Paths

| Purpose | Path |
|---------|------|
| Remediation status management | `app/api/controller/airemediation/controller.go` |
| Solo Gate engine | `app/services/sologate/engine.go` |

## Status

**Concept** — Manual approval workflow is implemented (developer reviews and applies patches). Automated approval with confidence thresholds and Solo Gate integration is planned.

## Future Work

- Auto-PR creation from completed remediations
- Confidence-based auto-merge
- Approval history and audit trail
- Notification system for pending reviews
