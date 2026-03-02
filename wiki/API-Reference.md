# API Reference

Complete reference for all SoloDev-specific API endpoints.

All endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <personal-access-token>
```

## Base URL

```
http://localhost:3000/api/v1
```

## Spaces

All SoloDev module endpoints are scoped to a space:

```
/api/v1/spaces/{space_ref}/
```

Where `{space_ref}` is the unique identifier for the space (e.g., `my-org/my-space` or just `my-space`).

---

## AI Auto-Remediation

Base path: `/api/v1/spaces/{space_ref}/remediations`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/remediations` | Trigger a remediation task |
| GET | `/remediations` | List remediations (`?status=&trigger_source=&page=&limit=`) |
| GET | `/remediations/{id}` | Get remediation detail |
| PATCH | `/remediations/{id}` | Update remediation (status, AI response, patch diff, fix branch, PR link) |
| GET | `/remediations/summary` | Get aggregate statistics |

See [AI-Auto-Remediation](AI-Auto-Remediation) for request/response schemas.

---

## Auto-Pipeline

Base path: `/api/v1/spaces/{space_ref}/auto-pipeline`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auto-pipeline/generate` | Generate a CI/CD pipeline from a list of file paths |

See [Auto-Pipeline](Auto-Pipeline) for request/response schemas.

---

## Error Tracker

Base path: `/api/v1/spaces/{space_ref}/errors`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/errors` | Report an error |
| GET | `/errors` | List error groups (`?status=&severity=&language=&page=&limit=&query=`) |
| GET | `/errors/{id}` | Get error group detail (includes sample occurrences) |
| PATCH | `/errors/{id}` | Update error group (status, assigned_to, title, tags) |
| GET | `/errors/{id}/occurrences` | List occurrences for an error group |
| GET | `/errors/summary` | Get aggregate statistics |

See [Error-Tracker](Error-Tracker) for request/response schemas.

---

## Security Scanner

Base path: `/api/v1/spaces/{space_ref}/security-scans`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/security-scans` | Trigger a security scan |
| GET | `/security-scans` | List security scans (`?status=&scan_type=&triggered_by=&page=&limit=`) |
| GET | `/security-scans/{id}` | Get scan detail |
| GET | `/security-scans/{id}/findings` | List findings (`?severity=&category=&status=&page=&limit=`) |
| GET | `/security-scans/{id}/summary` | Get security posture summary |
| PATCH | `/security-scans/findings/{id}/status` | Update finding status (**route not yet registered in router**) |

See [Security-Scanner](Security-Scanner) for request/response schemas.

---

## Quality Gates

Base path: `/api/v1/spaces/{space_ref}/quality-gates`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/quality-gates/rules` | Create a quality rule |
| GET | `/quality-gates/rules` | List rules (`?category=&enabled=&page=&limit=`) |
| GET | `/quality-gates/rules/{id}` | Get rule detail |
| PATCH | `/quality-gates/rules/{id}` | Update rule |
| POST | `/quality-gates/rules/{id}/toggle` | Enable or disable rule |
| DELETE | `/quality-gates/rules/{id}` | Delete rule |
| POST | `/quality-gates/evaluate` | Trigger an evaluation |
| GET | `/quality-gates/evaluations` | List evaluations |
| GET | `/quality-gates/evaluations/{id}` | Get evaluation detail |
| GET | `/quality-gates/summary` | Get quality summary |

See [Quality-Gates](Quality-Gates) for request/response schemas.

---

## Health Monitor

Base path: `/api/v1/spaces/{space_ref}/health-checks`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/health-checks` | Create a health monitor |
| GET | `/health-checks` | List monitors (`?page=&size=&query=`) |
| GET | `/health-checks/{id}` | Get monitor detail |
| PATCH | `/health-checks/{id}` | Update monitor |
| DELETE | `/health-checks/{id}` | Delete monitor |
| GET | `/health-checks/{id}/results` | Get recent check results (`?limit=`) |
| GET | `/health-checks/{id}/summary` | Get uptime summary |

See [Health-Monitor](Health-Monitor) for request/response schemas.

---

## Feature Flags

Base path: `/api/v1/spaces/{space_ref}/feature-flags`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/feature-flags` | Create a feature flag |
| GET | `/feature-flags` | List flags (`?query=`) |
| GET | `/feature-flags/{id}` | Get flag detail |
| PATCH | `/feature-flags/{id}` | Update flag |
| DELETE | `/feature-flags/{id}` | Delete flag |

See [Feature-Flags](Feature-Flags) for request/response schemas.

---

## Tech Debt

Base path: `/api/v1/spaces/{space_ref}/tech-debt`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/tech-debt` | Create a tech debt item |
| GET | `/tech-debt` | List items (`?severity=&status=&category=&repo_id=&page=&limit=&sort=`) |
| GET | `/tech-debt/{id}` | Get item detail |
| PATCH | `/tech-debt/{id}` | Update item |
| DELETE | `/tech-debt/{id}` | Delete item |

See [Tech-Debt](Tech-Debt) for request/response schemas.

---

## MCP Endpoint

| Method | Path | Description |
|--------|------|-------------|
| POST | `/mcp` | MCP JSON-RPC 2.0 endpoint (all methods) |

See [MCP-Server](MCP-Server) for the full protocol and tool reference.

---

## Discovery / Swagger

| URL | Description |
|-----|-------------|
| `http://localhost:3000/swagger` | Swagger UI (interactive API browser) |
| `http://localhost:3000/openapi.yaml` | Raw OpenAPI spec |
| `http://localhost:3000/registry/swagger/` | Registry (artifact) API swagger |

To regenerate the frontend API client from the live spec:

```bash
./gitness swagger > web/src/services/code/swagger.yaml
cd web && yarn services
```
