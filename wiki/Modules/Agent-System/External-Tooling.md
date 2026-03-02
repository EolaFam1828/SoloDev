# External Tooling

## Purpose

External Tooling describes how SoloDev integrates with third-party services and tools beyond its core platform. These integrations extend the platform's capabilities without requiring the developer to configure separate systems.

## Inputs

- LLM provider API calls (Anthropic, OpenAI, Gemini, Ollama)
- MCP client connections from external AI agents
- Docker socket for container operations

## Processing

### Current Integrations

#### LLM Providers

SoloDev connects to external AI model providers for patch generation:

| Provider | Type | Configuration |
|----------|------|--------------|
| Anthropic | Remote API | API key |
| OpenAI | Remote API | API key |
| Google Gemini | Remote API | API key |
| Ollama | Local inference | Endpoint URL (default: localhost) |

#### MCP Clients

External AI agents connect to SoloDev via MCP:

| Client | Transport | Configuration |
|--------|-----------|--------------|
| Claude Desktop | stdio | `claude_desktop_config.json` |
| Cursor | stdio | MCP configuration file |
| Custom agents | HTTP | Bearer token + endpoint URL |

#### Docker

SoloDev uses Docker for:
- Pipeline execution (container-based steps)
- Gitspace provisioning
- Registry operations
- Development environment (Docker Compose)

### Planned Integrations

| Integration | Purpose |
|-------------|---------|
| Webhook notifications | Alert external systems on remediation events |
| Slack/Discord | Notification delivery for applied fixes |
| GitHub/GitLab | Remote repository mirroring and PR creation |
| OpenTelemetry | Export telemetry to external observability platforms |

## Outputs

- LLM-generated patches from provider APIs
- Agent actions triggered via MCP protocol
- Container execution results from Docker

## Status

**Implemented** — LLM provider integrations, MCP server, and Docker integration are working. Webhook and notification integrations are planned.

## Future Work

- Webhook notification system for remediation events
- Chat platform integrations (Slack, Discord)
- Remote repository mirroring (GitHub, GitLab)
- OpenTelemetry export
- Plugin architecture for community-contributed integrations
