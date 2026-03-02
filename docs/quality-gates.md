# Code Quality Gates Module

This document describes the Code Quality Gates module for the SoloDev platform, which enables enforcement of code quality rules and evaluation of repositories against quality metrics.

## Overview

The Code Quality Gates module provides:

1. **Quality Rules Management** - Create, update, delete, and manage quality rules for a space
2. **Quality Evaluations** - Evaluate repositories against quality rules and track results
3. **Quality Metrics** - Aggregate statistics and summaries of code quality across a space

## Architecture

### Types (`types/qualitygate.go`, `types/enum/qualitygate.go`)

#### QualityRule
Represents an individual code quality rule/policy with the following fields:
- `ID`, `SpaceID` - Unique identifiers
- `Identifier`, `Name`, `Description` - Display information
- `Category` - One of: coverage, complexity, documentation, naming, testing, security, style, custom
- `Enforcement` - One of: block (fails pipeline), warn (warning only), info (informational)
- `Condition` - Rule expression (e.g., "coverage >= 80", "max_function_lines <= 50")
- `TargetRepoIDs` - JSON array of repo IDs (empty means all repos)
- `TargetBranches` - JSON array of branch patterns (e.g., ["main", "release/*"])
- `Enabled` - Boolean to enable/disable the rule
- `Tags` - JSON array for categorization
- `CreatedBy`, `Created`, `Updated`, `Version` - Audit fields

#### QualityEvaluation
Represents the result of evaluating rules against code:
- `ID`, `SpaceID`, `RepoID` - Unique identifiers
- `Identifier` - Unique evaluation identifier
- `CommitSHA`, `Branch` - Git references
- `OverallStatus` - One of: passed, failed, warning
- `RulesEvaluated`, `RulesPassed`, `RulesFailed`, `RulesWarned` - Statistics
- `Results` - JSON array of individual rule results
- `TriggeredBy` - One of: pipeline, manual, pr
- `PipelineID` - Optional pipeline reference
- `Duration` - Evaluation duration in milliseconds
- `CreatedBy`, `Created` - Audit fields

#### QualitySummary
Aggregate statistics for a space:
- `TotalRepositories`, `RepositoriesPassed`, `RepositoriesFailed`, `RepositoriesWarned`
- `AverageCoverage` - Optional coverage metric
- `TotalEvaluations`, `FailedEvaluations`
- `LastEvaluationTime` - Timestamp of most recent evaluation

### Database (`app/store/database/qualitygate.go`)

#### Tables
- `quality_rules` - Stores quality rule definitions (prefix: `qr_`)
- `quality_evaluations` - Stores evaluation results (prefix: `qe_`)

#### Store Interfaces
Defined in `app/store/database.go`:
- `QualityRuleStore` - CRUD operations for quality rules
- `QualityEvaluationStore` - CRUD and query operations for evaluations

### API Controller (`app/api/controller/qualitygate/controller.go`)

The main controller implements:

**Rule Management:**
- `CreateRule(ctx, session, spaceRef, input)` - Create new rule
- `GetRule(ctx, session, spaceRef, identifier)` - Retrieve rule
- `ListRules(ctx, session, spaceRef, filter)` - List rules with filtering
- `UpdateRule(ctx, session, spaceRef, identifier, input)` - Update rule
- `ToggleRule(ctx, session, spaceRef, identifier, input)` - Enable/disable rule
- `DeleteRule(ctx, session, spaceRef, identifier)` - Delete rule

**Evaluation:**
- `Evaluate(ctx, session, spaceRef, input)` - Trigger evaluation
- `ListEvaluations(ctx, session, spaceRef, filter)` - List evaluations
- `GetEvaluation(ctx, session, spaceRef, identifier)` - Retrieve evaluation detail
- `GetSummary(ctx, session, spaceRef)` - Get aggregate statistics

### HTTP Handlers (`app/api/handler/qualitygate/`)

#### Rule Endpoints
- `HandleRuleCreate` - POST /api/v1/spaces/{space_ref}/quality-gates/rules
- `HandleRuleList` - GET /api/v1/spaces/{space_ref}/quality-gates/rules
- `HandleRuleGet` - GET /api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}
- `HandleRuleUpdate` - PATCH /api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}
- `HandleRuleToggle` - POST /api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}/toggle
- `HandleRuleDelete` - DELETE /api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}

#### Evaluation Endpoints
- `HandleEvaluate` - POST /api/v1/spaces/{space_ref}/quality-gates/evaluate
- `HandleEvaluationList` - GET /api/v1/spaces/{space_ref}/quality-gates/evaluations
- `HandleEvaluationGet` - GET /api/v1/spaces/{space_ref}/quality-gates/evaluations/{identifier}
- `HandleSummaryGet` - GET /api/v1/spaces/{space_ref}/quality-gates/summary

### Events (`app/events/qualitygate/`)

Event types published:
- `RuleCreatedEvent` - When a rule is created
- `RuleUpdatedEvent` - When a rule is updated
- `RuleDeletedEvent` - When a rule is deleted
- `RuleEnabledEvent` - When a rule is enabled
- `RuleDisabledEvent` - When a rule is disabled
- `EvaluationCreatedEvent` - When an evaluation is created

## Database Migrations

Migration 0105 creates the necessary tables:

```sql
-- quality_rules table with unique constraint on (space_id, identifier)
-- quality_evaluations table with foreign keys to spaces and repositories
```

Migration files:
- `0105.up.sql` - Creates tables and indexes
- `0105.down.sql` - Drops tables

## API Examples

### Create a Quality Rule
```bash
POST /api/v1/spaces/{space_ref}/quality-gates/rules
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

### List Quality Rules
```bash
GET /api/v1/spaces/{space_ref}/quality-gates/rules?page=0&limit=20&category=coverage&enabled=true
```

### Update a Quality Rule
```bash
PATCH /api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}
{
  "enforcement": "warn",
  "condition": "coverage >= 75"
}
```

### Toggle Rule
```bash
POST /api/v1/spaces/{space_ref}/quality-gates/rules/{rule_identifier}/toggle
{
  "enabled": false
}
```

### Trigger Evaluation
```bash
POST /api/v1/spaces/{space_ref}/quality-gates/evaluate
{
  "repo_ref": "my-repo",
  "commit_sha": "abc123def456",
  "branch": "main",
  "trigger": "pipeline",
  "pipeline_id": 12345
}
```

### Get Quality Summary
```bash
GET /api/v1/spaces/{space_ref}/quality-gates/summary
```

## Integration Points

### Authorization
All endpoints require the appropriate permission checked via `apiauth.CheckSpace`:
- `enum.PermissionQualityGateView` for GET operations (list, get, summary)
- `enum.PermissionQualityGateCreate` for creating rules
- `enum.PermissionQualityGateEdit` for updating and toggling rules
- `enum.PermissionQualityGateDelete` for deleting rules
- `enum.PermissionQualityGateEvaluate` for triggering evaluations

### Space Finder
Uses `refcache.SpaceFinder` to resolve space references

### Events System
Publishes events through the `events.System` interface for:
- Rule lifecycle changes
- Evaluation completions

## Implementation Notes

1. **Identifier Uniqueness** - Quality rules are uniquely identified by (space_id, identifier) combination
2. **Version Control** - Rules support optimistic locking via version field
3. **Flexible Targeting** - Rules can target specific repos or branches, or apply to all
4. **Status Determination** - Evaluation status is determined by rule results:
   - Failed if any rule fails
   - Warned if no failures but warnings exist
   - Passed if all rules pass
5. **JSON Storage** - Arrays and results are stored as JSON for flexibility

## Files Created

### Types
- `types/qualitygate.go`
- `types/enum/qualitygate.go`

### Database
- `app/store/database/qualitygate.go`
- `app/store/database/migrate/postgres/0105.up.sql`
- `app/store/database/migrate/postgres/0105.down.sql`

### Controller
- `app/api/controller/qualitygate/controller.go`
- `app/api/controller/qualitygate/wire.go`

### Handlers
- `app/api/handler/qualitygate/rule_create.go`
- `app/api/handler/qualitygate/rule_list.go`
- `app/api/handler/qualitygate/rule_get.go`
- `app/api/handler/qualitygate/rule_update.go`
- `app/api/handler/qualitygate/rule_delete.go`
- `app/api/handler/qualitygate/rule_toggle.go`
- `app/api/handler/qualitygate/eval_create.go`
- `app/api/handler/qualitygate/eval_list.go`
- `app/api/handler/qualitygate/eval_get.go`
- `app/api/handler/qualitygate/summary_get.go`
- `app/api/handler/qualitygate/wire.go`

### Events
- `app/events/qualitygate/category.go`
- `app/events/qualitygate/events.go`
- `app/events/qualitygate/reporter.go`
- `app/events/qualitygate/reader.go`
- `app/events/qualitygate/wire.go`

### Store Interfaces Update
- `app/store/database.go` (modified)

## Next Steps

To complete the integration:

1. **Register Routes** - Add routes to the router configuration in your HTTP server setup
2. **Wire Dependencies** - Add qualitygate.Controller and events to your wire configuration
3. **Database Migrations** - Run migration 0105 to create tables
4. **Testing** - Add unit and integration tests
5. **Documentation** - Add API documentation to your Swagger/OpenAPI specs
6. **Evaluation Engine** - Implement actual rule evaluation logic (currently placeholder)
