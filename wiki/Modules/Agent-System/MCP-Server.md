# MCP Server

## Purpose

The MCP server implements the [Model Context Protocol 2024-11-05](https://modelcontextprotocol.io) specification, allowing any MCP-compatible AI client (Claude Desktop, Cursor, custom agents) to interact with the full SoloDev platform programmatically.

## Inputs

- JSON-RPC 2.0 requests via stdio or HTTP transport
- Bearer token authentication (PAT)
- Tool call parameters, resource read requests, prompt get requests

## Processing

### Transports

**Stdio** (local agents):
```bash
gitness mcp stdio
```
Reads JSON-RPC from stdin, writes responses to stdout. Token from `SOLODEV_TOKEN` environment variable.

**Streamable HTTP** (remote agents):
```bash
gitness mcp sse --port 3001
```
When running the full server, the MCP endpoint is mounted at `/mcp`. Clients send `POST /mcp` with Bearer token.

### Protocol Methods

- `initialize` — Protocol handshake
- `ping` — Health check
- `tools/list`, `tools/call` — Tool discovery and execution
- `resources/list`, `resources/read` — Resource discovery and reading
- `prompts/list`, `prompts/get` — Prompt discovery and rendering

## Runtime Model

The MCP implementation exposes three classes of surface:

- **Active runtime surfaces**: tools/resources/prompts that are currently enabled and backed by live controllers
- **Blocked surfaces**: cataloged entries whose backing module is disabled or unavailable
- **Coming-soon surfaces**: documented but intentionally not active yet

Use the MCP catalog endpoint and the SoloDev dashboard to inspect the current runtime counts instead of relying on fixed totals in docs.

## Tools — Tier 1: Atomic (runtime-active subset)

| Tool Name | Module | Description |
|-----------|--------|-------------|
| `pipeline_generate` | Auto-Pipeline | Generate CI/CD pipeline from repo analysis |
| `security_scan` | Security Scanner | Trigger a security scan |
| `security_findings` | Security Scanner | List security scan findings |
| `quality_evaluate` | Quality Gate | Evaluate code against quality rules |
| `quality_rules_list` | Quality Gate | List quality gate rules |
| `quality_summary` | Quality Gate | Get quality gate evaluation summary |
| `error_report` | Error Tracker | Report a new error event |
| `error_list` | Error Tracker | List tracked errors |
| `error_get` | Error Tracker | Get a specific error by ID |
| `remediation_trigger` | Remediation | Trigger AI auto-remediation |
| `remediation_list` | Remediation | List remediation attempts |
| `remediation_get` | Remediation | Get a specific remediation by ID |
| `remediation_apply` | Remediation | Apply a completed remediation into a fix branch and draft PR |
| `remediation_update` | Remediation | Update remediation status or details |
| `health_summary` | Health Monitor | Get health check summary |
| `feature_flag_toggle` | Feature Flags | Toggle a feature flag |
| `tech_debt_list` | Tech Debt | List tech debt items |

## Tools — Tier 2: Compound (5 tools)

| Tool Name | Description |
|-----------|-------------|
| `fix_this` | Analyze error, find root cause, trigger remediation, report |
| `pr_ready` | Run security + quality + tech debt check, return PR readiness |
| `pipeline_validate` | Generate pipeline, dry-run validate, return warnings |
| `onboard_repo` | Full onboarding: scan + pipeline + quality rules + health checks |
| `incident_triage` | Error + health + security correlation, severity assessment |

## Resources — Tier 3: Live Context

| URI | Description |
|-----|-------------|
| `solodev://errors/active` | Currently active error events |
| `solodev://remediations/pending` | Pending remediation attempts |
| `solodev://quality/rules` | Configured quality gate rules |
| `solodev://quality/summary` | Quality gate evaluation summary |
| `solodev://security/open-findings` | Open security scan findings |
| `solodev://health/status` | Current health check statuses |
| `solodev://tech-debt/hotspots` | Top tech debt hotspots |

## Prompts — Tier 4: Expert Reasoning

| Prompt Name | Description |
|-------------|-------------|
| `solodev_review` | Code review with security, quality, and tech debt context |
| `solodev_incident` | Incident investigation correlating errors, health, and security |
| `solodev_pipeline_debug` | Pipeline debugging with generation context and validation |
| `solodev_security_audit` | Security audit across all scan findings |
| `solodev_debt_sprint` | Tech debt sprint planning with prioritized items |

## Outputs

- JSON-RPC 2.0 responses for all protocol methods
- Tool execution results (structured JSON)
- Resource content (platform state snapshots)
- Rendered prompt templates

## Key Paths

| Purpose | Path |
|---------|------|
| Server (JSON-RPC dispatch) | `mcp/server.go` |
| Auth | `mcp/auth.go` |
| Wire types | `mcp/types.go` |
| Atomic tools | `mcp/tools_atomic.go` |
| Compound tools | `mcp/tools_compound.go` |
| Resources | `mcp/resources.go` |
| Prompts | `mcp/prompts.go` |
| Stdio transport | `mcp/transport_stdio.go` |
| HTTP transport | `mcp/transport_sse.go` |
| Tests | `mcp/server_test.go` |
| CLI subcommands | `cli/operations/mcpcmd/mcpcmd.go` |

## Status

**Implemented** — Full MCP 2024-11-05 implementation with stdio and HTTP transports, live runtime gating, and a catalog that separates active, blocked, and coming-soon surfaces.

## Future Work

- Multi-agent orchestration protocol
- Agent session management and context persistence
- Tool authorization scoping (per-agent permissions)
- Webhook-based agent notifications
