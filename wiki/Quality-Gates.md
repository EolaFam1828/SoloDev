# Quality Gates

## Overview

The Quality Gates module provides code quality rule management and evaluation for the SoloDev platform. It enables enforcement of code quality policies, evaluation of repositories against rules, and aggregated quality metrics per space.

## Core Concepts

### Quality Rule

A rule defines a single quality policy:

| Field | Description |
|-------|-------------|
| `Identifier` | Space-unique identifier |
| `Name` | Display name |
| `Category` | `coverage`, `complexity`, `documentation`, `naming`, `testing`, `security`, `style`, `custom` |
| `Enforcement` | `block` (fails pipeline), `warn` (warning only), `info` (informational) |
| `Condition` | Rule expression, e.g. `coverage >= 80`, `max_function_lines <= 50` |
| `TargetRepoIDs` | JSON array of repo IDs (empty = all repos) |
| `TargetBranches` | JSON array of branch patterns, e.g. `["main", "release/*"]` |
| `Enabled` | Enable/disable the rule |
| `Tags` | JSON array for categorization |

### Quality Evaluation

Records the result of evaluating rules against a specific commit:

| Field | Description |
|-------|-------------|
| `CommitSHA`, `Branch` | Git references |
| `OverallStatus` | `passed`, `failed`, `warning` |
| `RulesEvaluated`, `RulesPassed`, `RulesFailed`, `RulesWarned` | Statistics |
| `Results` | JSON array of individual rule results |
| `TriggeredBy` | `pipeline`, `manual`, `pr` |

**Status determination:**
- `failed` if any rule fails
- `warning` if no failures but warnings exist
- `passed` if all rules pass

### Quality Summary

Aggregate statistics for a space:
- `TotalRepositories`, `RepositoriesPassed`, `RepositoriesFailed`, `RepositoriesWarned`
- `AverageCoverage` (optional coverage metric)
- `TotalEvaluations`, `FailedEvaluations`
- `LastEvaluationTime`

## API Endpoints

### Rule Management

#### Create a Rule

**POST** `/api/v1/spaces/{space_ref}/quality-gates/rules`

```json
{
  "identifier": "coverage-threshold",
  "name": "Code Coverage Threshold",
  "description": "Ensures minimum code coverage",
  "category": "coverage",
  "enforcement": "block",
  "condition": "coverage >= 80",
  "target_repo_ids": [1, 2, 3],
  "target_branches": ["main", "release/*"],
  "enabled": true,
  "tags": ["critical", "coverage"]
}
```

#### List Rules

**GET** `/api/v1/spaces/{space_ref}/quality-gates/rules`

Query parameters: `page`, `limit`, `category`, `enabled`

#### Get Rule

**GET** `/api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}`

#### Update Rule

**PATCH** `/api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}`

```json
{
  "enforcement": "warn",
  "condition": "coverage >= 75"
}
```

#### Toggle Rule

**POST** `/api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}/toggle`

```json
{
  "enabled": false
}
```

#### Delete Rule

**DELETE** `/api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}`

### Evaluations

#### Trigger Evaluation

**POST** `/api/v1/spaces/{space_ref}/quality-gates/evaluate`

```json
{
  "repo_ref": "my-repo",
  "commit_sha": "abc123def456",
  "branch": "main",
  "trigger": "pipeline",
  "pipeline_id": 12345
}
```

#### List Evaluations

**GET** `/api/v1/spaces/{space_ref}/quality-gates/evaluations`

#### Get Evaluation

**GET** `/api/v1/spaces/{space_ref}/quality-gates/evaluations/{identifier}`

#### Get Quality Summary

**GET** `/api/v1/spaces/{space_ref}/quality-gates/summary`

## Access Control

| Operation | Permission Required |
|-----------|-------------------|
| List / Get rules, evaluations, summary | `PermissionQualityGateView` |
| Create rule | `PermissionQualityGateCreate` |
| Update / Toggle rule | `PermissionQualityGateEdit` |
| Delete rule | `PermissionQualityGateDelete` |
| Trigger evaluation | `PermissionQualityGateEvaluate` |

## Database Schema

Tables created in migration `0105`:

### `quality_rules` table

Primary key: `qr_id`  
Unique constraint: `(qr_space_id, qr_identifier)`

Fields: space_id, identifier, name, description, category, enforcement, condition, target_repo_ids (JSON), target_branches (JSON), enabled, tags (JSON), created_by, created, updated, version

### `quality_evaluations` table

Primary key: `qe_id`  
FK to spaces and repositories

Fields: space_id, repo_id, identifier, commit_sha, branch, overall_status, rules_evaluated, rules_passed, rules_failed, rules_warned, results (JSON), triggered_by, pipeline_id, duration, created_by, created

## Events

Six events are published:
- `RuleCreatedEvent`
- `RuleUpdatedEvent`
- `RuleDeletedEvent`
- `RuleEnabledEvent`
- `RuleDisabledEvent`
- `EvaluationCreatedEvent`

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/qualitygate.go` |
| Enums | `types/enum/qualitygate.go` |
| Database store | `app/store/database/qualitygate.go` |
| Store interfaces | `app/store/database.go` |
| Controller | `app/api/controller/qualitygate/controller.go` |
| Handlers | `app/api/handler/qualitygate/` |
| Events | `app/events/qualitygate/` |
| Migrations | `app/store/database/migrate/postgres/0105.*` |
