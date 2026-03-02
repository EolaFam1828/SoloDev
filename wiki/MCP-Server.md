# MCP Server

## Overview

The MCP server implements the [Model Context Protocol 2024-11-05](https://modelcontextprotocol.io) specification over two transports, allowing any MCP-compatible AI client (Claude Desktop, Cursor, custom agents) to interact with the full SoloDev platform programmatically.

## Transports

### Stdio (Claude Desktop / Local Agents)

```bash
gitness mcp stdio
```

Reads JSON-RPC requests from stdin, writes responses to stdout.

**Claude Desktop configuration** (`~/.config/claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "solodev": {
      "command": "gitness",
      "args": ["mcp", "stdio"],
      "env": {
        "SOLODEV_TOKEN": "<your-personal-access-token>"
      }
    }
  }
}
```

### Streamable HTTP (Remote Clients)

```bash
gitness mcp sse --port 3001
```

When running the full server, the MCP endpoint is automatically mounted at `/mcp`. Clients send `POST /mcp` with `Content-Type: application/json`.

```bash
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"curl","version":"1.0"},"capabilities":{}}}'
```

## Authentication

All MCP requests are authenticated via Bearer token using the existing Personal Access Token (PAT) system. The token is read from:
- `Authorization: Bearer <token>` header (HTTP transport)
- `SOLODEV_TOKEN` environment variable (stdio transport)

## Tools — Tier 1: Atomic (16 tools)

Direct wrappers around individual SoloDev controllers. Each tool maps 1:1 to a controller method.

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
| `remediation_update` | Remediation | Update remediation status or details |
| `health_summary` | Health Monitor | Get health check summary |
| `feature_flag_toggle` | Feature Flags | Toggle a feature flag |
| `tech_debt_list` | Tech Debt | List tech debt items |

## Tools — Tier 2: Compound Power Tools (5 tools)

Multi-step orchestrations that combine multiple controllers into high-level workflows.

| Tool Name | Description |
|-----------|-------------|
| `fix_this` | Analyze an error, find root cause, trigger remediation, report |
| `pr_ready` | Run security scan + quality gate + tech debt check, return PR readiness verdict |
| `pipeline_validate` | Generate pipeline, dry-run validate, return warnings |
| `onboard_repo` | Full repository onboarding: scan + pipeline + quality rules + health checks |
| `incident_triage` | Error + health + security correlation, incident severity assessment |

## Resources — Tier 3: Live Context (7 resources)

Real-time, read-only data URIs that AI agents can subscribe to for contextual awareness.

| URI | Description |
|-----|-------------|
| `solodev://errors/active` | Currently active error events |
| `solodev://remediations/pending` | Pending remediation attempts |
| `solodev://quality/rules` | Configured quality gate rules |
| `solodev://quality/summary` | Quality gate evaluation summary |
| `solodev://security/open-findings` | Open security scan findings |
| `solodev://health/status` | Current health check statuses |
| `solodev://tech-debt/hotspots` | Top tech debt hotspots |

## Prompts — Tier 4: Expert Reasoning Chains (5 prompts)

Pre-baked prompt templates that encode SoloDev domain expertise for common AI agent workflows.

| Prompt Name | Description |
|-------------|-------------|
| `solodev_review` | Code review with security, quality, and tech debt context |
| `solodev_incident` | Incident investigation correlating errors, health, and security |
| `solodev_pipeline_debug` | Pipeline debugging with generation context and validation |
| `solodev_security_audit` | Security audit across all scan findings for a repo |
| `solodev_debt_sprint` | Tech debt sprint planning with prioritized remediation items |

## Protocol

The MCP server implements [JSON-RPC 2.0](https://www.jsonrpc.org/specification) dispatch supporting the following methods:
- `initialize`
- `ping`
- `tools/list`, `tools/call`
- `resources/list`, `resources/read`
- `prompts/list`, `prompts/get`

## Testing

```bash
go test ./mcp/... -v
```

The test suite (`mcp/server_test.go`) includes 16 tests covering:
- Protocol handshake (`initialize`, `ping`)
- Tool/resource/prompt listing and invocation
- Error handling (unknown tools, invalid JSON, bad JSON-RPC version)
- HTTP transport (streaming handler, CORS preflight, method filtering)
- Resource reads and prompt rendering

## Web UI Setup Page

The SoloDev web UI includes a dedicated MCP setup page (`/pages/McpSetup`) accessible from the main dashboard via the **Connect Client** button. The dashboard also shows an MCP connection status indicator that checks `/api/v1/system/config` to determine availability.

## File Locations

| Purpose | Path |
|---------|------|
| Auth | `mcp/auth.go` |
| Server (JSON-RPC dispatch) | `mcp/server.go` |
| Wire types | `mcp/types.go` |
| Constructor | `mcp/wire.go` |
| Atomic tools | `mcp/tools_atomic.go` |
| Compound tools | `mcp/tools_compound.go` |
| Resources | `mcp/resources.go` |
| Prompts | `mcp/prompts.go` |
| Stdio transport | `mcp/transport_stdio.go` |
| Streamable HTTP transport | `mcp/transport_sse.go` |
| Tests | `mcp/server_test.go` |
| CLI subcommands | `cli/operations/mcpcmd/mcpcmd.go` |
