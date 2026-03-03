# SoloDev MCP Server

Model Context Protocol (MCP) server exposing SoloDev platform capabilities as MCP
tools, resources, and prompts for AI agent consumption.

## Overview

The MCP server implements the [Model Context Protocol 2024-11-05](https://modelcontextprotocol.io)
specification over two transports, allowing any MCP-compatible AI client (Claude
Desktop, Cursor, custom agents) to interact with the full SoloDev platform
programmatically.

## Architecture

```
mcp/
  auth.go            — Bearer-token authenticator (wraps existing PAT system)
  server.go          — JSON-RPC 2.0 dispatch (initialize, ping, tools/*, resources/*, prompts/*)
  types.go           — MCP wire-format types (requests, responses, JSON-RPC envelope)
  wire.go            — Controllers struct + NewServer() constructor
  tools_atomic.go    — Tier 1: atomic controller wrappers (runtime-active subset depends on module availability)
  tools_compound.go  — Tier 2: compound power tools
  resources.go       — Tier 3: live resource URIs
  prompts.go         — Tier 4: expert prompts (runtime-gated by catalog state)
  transport_stdio.go — Stdio transport for Claude Desktop / local agents
  transport_sse.go   — Streamable HTTP transport mounted at /mcp
  server_test.go     — 16 unit tests covering protocol, tools, resources, prompts

cli/operations/mcpcmd/
  mcpcmd.go          — CLI subcommands: `gitness mcp stdio`, `gitness mcp sse`
```

## Transports

### Stdio (Claude Desktop)

```bash
gitness mcp stdio
```

Reads JSON-RPC requests from stdin, writes responses to stdout. Designed for
use as a Claude Desktop MCP server configuration:

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

Or, when running the full Gitness server, the MCP endpoint is automatically
mounted at `/mcp` on the main API router. Clients send `POST /mcp` with
`Content-Type: application/json` and receive streaming JSON-RPC responses.

```bash
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"curl","version":"1.0"},"capabilities":{}}}'
```

## Authentication

All requests are authenticated via Bearer token using the existing Gitness
Personal Access Token (PAT) system. The MCP authenticator extracts the token
from either:

- `Authorization: Bearer <token>` header (HTTP transport)
- `SOLODEV_TOKEN` environment variable (stdio transport)

## Tools (Tier 1 — Atomic)

Direct wrappers around individual SoloDev controllers. Runtime availability is
not a fixed number: the active subset depends on which controllers are enabled
and ready. The catalog still tracks blocked and coming-soon surfaces separately.

| Tool Name | Module | Description |
|---|---|---|
| `pipeline_generate` | Auto-Pipeline | Generate CI/CD pipeline from repo analysis |
| `security_scan` | Security Scanner | Trigger a security scan |
| `security_findings` | Security Scanner | List security scan findings |
| `security_fix_finding` | Security Scanner | Create or reuse a remediation for a finding |
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

## Tools (Tier 2 — Compound Power Tools)

Multi-step orchestrations that combine multiple controllers into high-level
workflows.

| Tool Name | Description |
|---|---|
| `fix_this` | Analyze an error, find root cause, trigger remediation, report |
| `pr_ready` | Run security scan + quality gate + tech debt check, return PR readiness verdict |
| `pipeline_validate` | Generate pipeline, dry-run validate, return warnings |
| `onboard_repo` | Full repository onboarding: scan + pipeline + quality rules + health checks |
| `incident_triage` | Error + health + security correlation, incident severity assessment |

## Resources (Tier 3 — Live Context)

Real-time, read-only data URIs that AI agents can subscribe to for contextual
awareness. As with tools, the active set varies at runtime; blocked and
coming-soon resources remain visible through the catalog API.

| URI | Description |
|---|---|
| `solodev://errors/active` | Currently active error events |
| `solodev://remediations/pending` | Pending remediation attempts |
| `solodev://quality/rules` | Configured quality gate rules |
| `solodev://quality/summary` | Quality gate evaluation summary |
| `solodev://security/open-findings` | Open security scan findings |
| `solodev://health/status` | Current health check statuses |
| `solodev://tech-debt/hotspots` | Top tech debt hotspots |

## Prompts (Tier 4 — Expert Reasoning Chains)

Pre-baked prompt templates that encode SoloDev domain expertise for common
AI agent workflows.

| Prompt Name | Description |
|---|---|
| `solodev_review` | Code review with security, quality, and tech debt context |
| `solodev_incident` | Incident investigation correlating errors, health, and security |
| `solodev_pipeline_debug` | Pipeline debugging with generation context and validation |
| `solodev_security_audit` | Security audit across all scan findings for a repo |
| `solodev_debt_sprint` | Tech debt sprint planning with prioritized remediation items |

## Runtime Truthfulness

- Use `/api/v1/system/mcp/catalog` to inspect cataloged, blocked, and coming-soon surfaces.
- Use `/api/v1/spaces/{space_ref}/solodev/overview` for the live dashboard counts that back the UI.
- `fix_this` still stops at the remediation object; `remediation_apply` is the explicit draft-PR handoff.

## Testing

```bash
go test ./mcp/... -v
```

Runs 16 tests covering:
- Protocol handshake (initialize, ping)
- Tool/resource/prompt listing and invocation
- Error handling (unknown tools, invalid JSON, bad JSON-RPC version)
- HTTP transport (streaming handler, CORS preflight, method filtering)
- Resource reads and prompt rendering

## License

Copyright 2026 EolaFam1828. All rights reserved.

Licensed under the Apache License, Version 2.0. See [LICENSE](../LICENSE) for details.

Built upon [Harness Open Source](https://github.com/harness/gitness).
