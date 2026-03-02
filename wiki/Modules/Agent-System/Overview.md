# Agent System Overview

## Purpose

The Agent System enables AI agents to interact with SoloDev as first-class operators. Through the Model Context Protocol (MCP), agents can monitor platform state, trigger operations, diagnose problems, and apply fixes — performing the same actions a developer would, but autonomously.

## Role Inside SoloDev

The Agent System is the external interface of the SoloDev Loop. While the internal loop (Detect → Analyze → Propose → Validate → Apply) runs within the platform, the Agent System allows external AI agents to:

1. **Observe** — Read platform state via MCP resources
2. **Diagnose** — Use compound tools like `fix_this` and `incident_triage`
3. **Act** — Trigger remediations, toggle flags, report errors
4. **Plan** — Use prompt templates for structured reasoning

## Components

| Component | Description | Page |
|-----------|-------------|------|
| MCP Server | Model Context Protocol implementation with tools, resources, and prompts | [MCP Server](MCP-Server) |
| Agent Workflows | Autonomous and semi-autonomous development workflows | [Agent Workflows](Agent-Workflows) |
| External Tooling | Integrations with third-party services and tools | [External Tooling](External-Tooling) |

## MCP Capability Summary

| Tier | Type | Count | Description |
|------|------|-------|-------------|
| Tier 1 | Atomic Tools | 16 | Direct wrappers around individual controllers |
| Tier 2 | Compound Tools | 5 | Multi-step orchestrations combining multiple controllers |
| Tier 3 | Resources | 7 | Real-time, read-only data URIs |
| Tier 4 | Prompts | 5 | Pre-built reasoning chains |

## Agent Architecture

```
┌──────────────────────┐
│    AI Agent           │
│  (Claude, Cursor,     │
│   custom agent)       │
└──────────┬───────────┘
           │ MCP Protocol (JSON-RPC 2.0)
           │
┌──────────┴───────────┐
│    MCP Server         │
│                       │
│  ┌─────────────────┐ │
│  │ Atomic Tools    │ │ ──▶ Individual controller methods
│  │ Compound Tools  │ │ ──▶ Multi-step orchestrations
│  │ Resources       │ │ ──▶ Live platform state
│  │ Prompts         │ │ ──▶ Domain expertise templates
│  └─────────────────┘ │
│                       │
│  Auth: Bearer token   │
│  Transport: stdio/HTTP│
└───────────────────────┘
```

## Status

| Component | Status |
|-----------|--------|
| MCP Server (stdio + HTTP) | Implemented |
| Atomic Tools (16) | Implemented |
| Compound Tools (5) | Implemented |
| Resources (7) | Implemented |
| Prompts (5) | Implemented |
| Multi-agent orchestration | Planned |

## Future Work

- Multi-agent MCP orchestration (coordinate multiple agents)
- Agent-driven deployment decisions
- Agent memory and context persistence across sessions
- Agent authentication and authorization scoping
