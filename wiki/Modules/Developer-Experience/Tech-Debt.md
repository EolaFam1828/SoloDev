# Tech Debt

## Purpose

The Tech Debt module tracks and triages technical debt items within a space. Items can be linked to specific files and line ranges, assigned severities and categories, and managed through a lifecycle from open to resolved.

## Inputs

- Manual creation via API
- Automatic creation from the Solo Gate engine (when `LogTechDebt` is enabled and findings are passed rather than blocked)
- Security findings downgraded from block to debt
- Quality gate warnings in prototype mode

## Processing

- Stores debt items with severity, category, file location, and estimated effort
- Manages lifecycle: `open` → `in_progress` → `resolved` / `accepted`
- Aggregates statistics by severity, status, and category
- Links items to specific file paths and line ranges

### Categories

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

## Outputs

- Tech debt items accessible via API
- Summary statistics aggregated by severity, status, and category
- MCP resource: `solodev://tech-debt/hotspots` (top debt hotspots)
- MCP prompt: `solodev_debt_sprint` (AI-assisted sprint planning)

## API Endpoints

Base path: `/api/v1/spaces/{space_ref}/tech-debt`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/tech-debt` | Create a tech debt item |
| GET | `/tech-debt` | List items (`?severity=&status=&category=&repo_id=&page=&limit=&sort=`) |
| GET | `/tech-debt/{id}` | Get item detail |
| PATCH | `/tech-debt/{id}` | Update item |
| DELETE | `/tech-debt/{id}` | Delete item |

## Key Paths

| Purpose | Path |
|---------|------|
| Types | `types/techdebt.go` |
| Controller | `app/api/controller/techdebt/` |
| Handlers | `app/api/handler/techdebt/` |

## Status

**Implemented** — Item CRUD, severity/category classification, lifecycle management, Solo Gate integration, and MCP access are working.

## Future Work

- Debt age tracking and stale item alerting
- Automated debt detection from code analysis
- Debt reduction velocity metrics
