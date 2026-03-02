# Quality Gates

## Purpose

The Quality Gates module provides code quality rule management and evaluation. It enables enforcement of quality policies, evaluation of repositories against rules, and aggregated quality metrics per space.

## Inputs

- Quality rule definitions (identifier, category, enforcement level, condition, targets)
- Evaluation trigger requests (manual, pipeline completion, PR event)
- Repository reference, commit SHA, and branch

## Processing

### Quality Rules

Each rule defines a single quality policy:

| Field | Description |
|-------|-------------|
| `Category` | `coverage`, `complexity`, `documentation`, `naming`, `testing`, `security`, `style`, `custom` |
| `Enforcement` | `block` (fails pipeline), `warn` (warning only), `info` (informational) |
| `Condition` | Rule expression, e.g. `coverage >= 80`, `max_function_lines <= 50` |
| `TargetRepoIDs` | Repos to apply to (empty = all) |
| `TargetBranches` | Branch patterns, e.g. `["main", "release/*"]` |

### Evaluation

Evaluates all enabled rules against a specific commit:
- `failed` if any rule fails
- `warning` if no failures but warnings exist
- `passed` if all rules pass

Results include per-rule pass/fail status with counts.

## Outputs

- Rule evaluation results (passed/failed/warning per rule)
- Overall evaluation status
- Quality summary statistics per space
- Events: `RuleCreatedEvent`, `RuleUpdatedEvent`, `RuleDeletedEvent`, `RuleEnabledEvent`, `RuleDisabledEvent`, `EvaluationCreatedEvent`

## API Endpoints

Base path: `/api/v1/spaces/{space_ref}/quality-gates`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/rules` | Create a quality rule |
| GET | `/rules` | List rules |
| GET | `/rules/{id}` | Get rule detail |
| PATCH | `/rules/{id}` | Update rule |
| POST | `/rules/{id}/toggle` | Enable or disable rule |
| DELETE | `/rules/{id}` | Delete rule |
| POST | `/evaluate` | Trigger evaluation |
| GET | `/evaluations` | List evaluations |
| GET | `/evaluations/{id}` | Get evaluation detail |
| GET | `/summary` | Get quality summary |

## Database Schema

Tables created in migration `0105`: `quality_rules`, `quality_evaluations`

## Key Paths

| Purpose | Path |
|---------|------|
| Types | `types/qualitygate.go` |
| Enums | `types/enum/qualitygate.go` |
| Controller | `app/api/controller/qualitygate/controller.go` |
| Handlers | `app/api/handler/qualitygate/` |
| Events | `app/events/qualitygate/` |

## Status

**Implemented** — Rule CRUD, evaluation engine, coverage/complexity/style policies, and aggregated quality metrics are working.

## Future Work

- Integration with CI pipeline results for automated coverage data
- Custom rule expression language
- Quality trend tracking over time
