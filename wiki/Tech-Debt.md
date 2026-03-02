# Tech Debt

## Overview

The Tech Debt module provides tracking and triage of technical debt items within a space. Items can be linked to specific files and line ranges, assigned severities and categories, and managed through a lifecycle from open to resolved.

## Core Concepts

### `TechDebt`

| Field | Type | Description |
|-------|------|-------------|
| `ID` | int64 | Primary key |
| `SpaceID` | int64 | FK to spaces |
| `RepoID` | int64 | FK to repositories (optional) |
| `Identifier` | string | Space-unique identifier |
| `Title` | string | Item title |
| `Description` | string | Detailed description (optional) |
| `Severity` | TechDebtSeverity | `critical`, `high`, `medium`, `low` |
| `Status` | TechDebtStatus | `open`, `in_progress`, `resolved`, `accepted` |
| `Category` | TechDebtCategory | Category of the debt |
| `FilePath` | string | Affected file path (optional) |
| `LineStart` | int | Start line number (optional) |
| `LineEnd` | int | End line number (optional) |
| `EstimatedEffort` | string | Estimated remediation effort |
| `Tags` | []string | Tag array |
| `DueDate` | int64 | Due date in Unix milliseconds (optional) |
| `ResolvedAt` | int64 | Resolution timestamp (optional) |
| `ResolvedBy` | int64 | FK to principal who resolved (optional) |
| `CreatedBy` | int64 | FK to principals |
| `Created` | int64 | Unix milliseconds |
| `Updated` | int64 | Unix milliseconds |
| `Version` | int64 | Optimistic locking version |

### Severity Values

`critical` · `high` · `medium` · `low`

### Status Values

`open` · `in_progress` · `resolved` · `accepted`

### Category Values

| Category | Description |
|----------|-------------|
| `code_smell` | Code quality issues |
| `bug_risk` | Potential bugs or unreliable code |
| `performance` | Performance bottlenecks |
| `security` | Security weaknesses (non-critical) |
| `documentation` | Missing or outdated documentation |
| `test_coverage` | Insufficient test coverage |
| `dependency` | Outdated or risky dependencies |
| `architecture` | Architectural design issues |

## API Endpoints

### Create Item

**POST** `/api/v1/spaces/{space_ref}/tech-debt`

```json
{
  "identifier": "sql-n-plus-one-users",
  "title": "N+1 query in users listing",
  "description": "The users list endpoint issues one query per user to fetch roles",
  "severity": "medium",
  "category": "performance",
  "file_path": "app/api/handler/users/list.go",
  "line_start": 45,
  "line_end": 60,
  "estimated_effort": "2h",
  "tags": ["database", "performance"],
  "repo_id": 1
}
```

### List Items

**GET** `/api/v1/spaces/{space_ref}/tech-debt`

Query parameters: `severity`, `status`, `category`, `repo_id`, `page`, `limit`, `sort`

Response:
```json
{
  "items": [...],
  "total": 42,
  "page": 0,
  "limit": 20
}
```

### Get Item

**GET** `/api/v1/spaces/{space_ref}/tech-debt/{identifier}`

### Update Item

**PATCH** `/api/v1/spaces/{space_ref}/tech-debt/{identifier}`

```json
{
  "status": "in_progress",
  "estimated_effort": "4h"
}
```

### Delete Item

**DELETE** `/api/v1/spaces/{space_ref}/tech-debt/{identifier}`

## Summary / Aggregation

The `TechDebtSummary` type provides aggregated statistics:

```json
{
  "by_severity": {"critical": 2, "high": 5, "medium": 10, "low": 3},
  "by_status":   {"open": 15, "in_progress": 3, "resolved": 2},
  "by_category": {"performance": 5, "code_smell": 8, ...},
  "total": 20
}
```

## MCP Integration

The `tech_debt_list` atomic MCP tool wraps the tech debt controller, allowing AI agents to list tech debt items via the MCP protocol. The `solodev_debt_sprint` prompt provides AI-assisted sprint planning with prioritized remediation items, and the `solodev://tech-debt/hotspots` resource exposes the top tech debt hotspots as live context.

## Integration with Solo Gate

When the [Solo Gate](Solo-Gate) is configured with `LogTechDebt = true`, findings that are passed rather than blocked are automatically logged as tech debt items for later review.

## Web UI

The **Tech Debt** page (`/pages/TechDebtList`) and the **Technical Debt** navigation item in the module sidebar provide a web interface for managing debt items.

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/techdebt.go` |
| Controller | `app/api/controller/techdebt/controller.go` |
| List | `app/api/controller/techdebt/list.go` |
| Create | `app/api/controller/techdebt/create.go` |
| Update | `app/api/controller/techdebt/update.go` |
| Delete | `app/api/controller/techdebt/delete.go` |
| Find | `app/api/controller/techdebt/find.go` |
| Handlers | `app/api/handler/techdebt/` |
