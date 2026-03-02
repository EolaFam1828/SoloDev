# API Overview

## Philosophy

SoloDev exposes all platform capabilities through a REST API and the MCP protocol. The API is designed for both human-driven clients (dashboards, CLI tools) and AI agent consumption. Every module has CRUD endpoints scoped to a space, and the MCP server provides tool-based access for AI agents.

## Authentication

All API requests require a Bearer token (Personal Access Token) in the `Authorization` header:

```
Authorization: Bearer <personal-access-token>
```

Create a PAT:
```bash
./gitness login
./gitness user pat "my-pat" 2592000
```

## Base URL

```
http://localhost:3000/api/v1
```

## Space Scoping

All SoloDev module endpoints are scoped to a space:

```
/api/v1/spaces/{space_ref}/
```

Where `{space_ref}` is the unique identifier for the space.

## Major Interfaces

| Module | Base Path | Description |
|--------|-----------|-------------|
| Remediations | `/spaces/{space_ref}/remediations` | AI remediation CRUD, summary |
| Errors | `/spaces/{space_ref}/errors` | Error group management, occurrences |
| Security Scans | `/spaces/{space_ref}/security-scans` | Scan triggers, findings |
| Quality Gates | `/spaces/{space_ref}/quality-gates` | Rules, evaluations, summary |
| Health Checks | `/spaces/{space_ref}/health-checks` | Monitor CRUD, results, summary |
| Feature Flags | `/spaces/{space_ref}/feature-flags` | Flag CRUD |
| Tech Debt | `/spaces/{space_ref}/tech-debt` | Debt item CRUD |
| Auto-Pipeline | `/spaces/{space_ref}/auto-pipeline` | Pipeline generation |
| MCP | `/mcp` | JSON-RPC 2.0 endpoint |

## Discovery

| URL | Description |
|-----|-------------|
| `http://localhost:3000/swagger` | Swagger UI (interactive API browser) |
| `http://localhost:3000/openapi.yaml` | Raw OpenAPI spec |
| `http://localhost:3000/registry/swagger/` | Registry API swagger |

## Response Format

All responses use JSON. List endpoints support pagination via `page` and `limit` query parameters. Error responses include a `message` field.

## Related Pages

- [Events](Events) — Event streams emitted by the platform
- [Remediation API](Remediation-API) — Detailed remediation endpoint reference
- [Integration](Integration) — How external tools interact with SoloDev
