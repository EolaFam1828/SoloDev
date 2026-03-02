# Architecture

## Overview

SoloDev is built on top of [Gitness by Harness](https://github.com/harness/gitness), an Apache-2.0 licensed open-source DevOps platform. SoloDev adds an AI-native intelligence layer on top of that foundation: automated error-to-fix pipelines, zero-config CI generation, enforcement gates, security scanning, runtime error tracking, and a first-class MCP server for AI agent access.

## High-Level Component Map

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              SoloDev Platform                                │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Web UI (React / TypeScript)                     │   │
│  │  Dashboard · Repos · Pipelines · Errors · Security · Quality ·       │   │
│  │  Health Monitor · Feature Flags · Tech Debt · MCP Setup              │   │
│  └───────────────────────────┬──────────────────────────────────────────┘   │
│                              │ REST / JSON-RPC                               │
│  ┌───────────────────────────▼──────────────────────────────────────────┐   │
│  │                     HTTP API Layer (chi router)                      │   │
│  │  /api/v1/spaces/{space_ref}/                                         │   │
│  │    remediations/  errors/  quality-gates/  security-scans/           │   │
│  │    health-checks/  tech-debt/  feature-flags/  auto-pipeline/        │   │
│  │  /mcp  (MCP JSON-RPC endpoint)                                       │   │
│  └───────────────────────────┬──────────────────────────────────────────┘   │
│                              │                                               │
│  ┌───────────────────────────▼──────────────────────────────────────────┐   │
│  │                       Controller Layer (Go)                          │   │
│  │  airemediation · autopipeline · errortracker · securityscan          │   │
│  │  qualitygate · healthcheck · featureflag · techdebt                  │   │
│  │  (+ upstream: repo · pipeline · pullreq · space · user · etc.)      │   │
│  └──────┬────────────────────┬─────────────────────────────────────────┘   │
│         │                    │                                               │
│  ┌──────▼──────┐   ┌─────────▼──────────────────────────────────────────┐  │
│  │   Services  │   │               Store Layer (SQLite / PostgreSQL)     │  │
│  │  sologate   │   │  RemediationStore · ErrorTrackerStore               │  │
│  │  errorbridge│   │  SecurityScanStore · QualityRuleStore               │  │
│  │  autopipeline│  │  HealthCheckStore · FeatureFlagStore · TechDebtStore│  │
│  └─────────────┘   └────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                          MCP Server (Go)                             │   │
│  │  Stdio transport · Streamable HTTP transport (/mcp)                  │   │
│  │  16 atomic tools · 5 compound tools · 7 resources · 5 prompts        │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Repository Layout

```
SoloDev/
├── app/
│   ├── api/
│   │   ├── controller/          # Business logic per domain
│   │   │   ├── airemediation/
│   │   │   ├── autopipeline/
│   │   │   ├── errortracker/
│   │   │   ├── featureflag/
│   │   │   ├── healthcheck/
│   │   │   ├── qualitygate/
│   │   │   ├── securityscan/
│   │   │   ├── techdebt/
│   │   │   └── ... (upstream controllers)
│   │   ├── handler/             # HTTP handler wrappers
│   │   └── request/             # Path/query parameter helpers
│   ├── events/                  # Domain event publishers/subscribers
│   │   ├── airemediation/
│   │   ├── errortracker/
│   │   ├── healthcheck/
│   │   ├── qualitygate/
│   │   └── securityscan/
│   ├── pipeline/autopipeline/   # Stack detection + YAML generation
│   ├── router/                  # HTTP router (chi) + route registration
│   ├── services/
│   │   ├── errorbridge/         # Error-to-remediation bridge
│   │   └── sologate/            # Gate evaluation engine
│   └── store/database/          # SQLite + PostgreSQL implementations
│       └── migrate/             # Schema migrations (numbered)
├── docs/                        # Module reference docs
├── mcp/                         # MCP server implementation
│   ├── server.go
│   ├── tools_atomic.go
│   ├── tools_compound.go
│   ├── resources.go
│   ├── prompts.go
│   ├── transport_stdio.go
│   └── transport_sse.go
├── types/                       # Domain types (structs, enums)
│   ├── ai_remediation.go
│   ├── auto_pipeline.go
│   ├── errortracker.go
│   ├── featureflag.go
│   ├── healthcheck.go
│   ├── qualitygate.go
│   ├── securityscan.go
│   ├── solo_gate.go
│   └── techdebt.go
└── web/src/                     # React TypeScript frontend
    ├── pages/
    │   ├── SoloDevDashboard/
    │   ├── ErrorList/
    │   ├── FeatureFlagList/
    │   ├── McpSetup/
    │   ├── MonitorList/
    │   ├── QualityRuleList/
    │   ├── SecurityScanList/
    │   └── TechDebtList/
    └── layouts/menu/ModuleMenu.tsx
```

## Database

SoloDev supports two database backends, both via numbered schema migrations:

| Backend | Used For |
|---------|---------|
| **SQLite** | Local development and single-node deployments |
| **PostgreSQL** | Production multi-user deployments |

SoloDev-specific migration numbers:
- `0102` — `security_scans`, `scan_findings`
- `0103` — `health_checks`
- `0104` — `health_check_results`
- `0105` — `quality_rules`, `quality_evaluations`
- `0171` — `error_groups`, `error_occurrences`
- `0172` — `remediations`

## Authentication

All API requests are authenticated via Bearer token using the Personal Access Token (PAT) system inherited from the Gitness foundation. The MCP server uses the same PAT system, reading the token from:
- `Authorization: Bearer <token>` header (HTTP transport)
- `SOLODEV_TOKEN` environment variable (stdio transport)

## Authorization

Role-based access control is enforced at every controller entry point using the `authz.Authorizer` interface and `apiauth.Check*` helpers. Permissions checked include `PermissionSpaceView`, `PermissionSpaceEdit`, `PermissionRepoView`, `PermissionRepoPush`, `PermissionQualityGateView`, `PermissionQualityGateCreate`, `PermissionQualityGateEdit`, `PermissionQualityGateDelete`, and `PermissionQualityGateEvaluate`.

## Events

All SoloDev modules emit domain events through a shared event system:
- `airemediation` — RemediationTriggered, RemediationCompleted, RemediationApplied
- `errortracker` — ErrorReported, ErrorStatusChanged, ErrorAssigned
- `healthcheck` — Created, Updated, Deleted, StatusChanged, ResultCreated
- `qualitygate` — RuleCreated, RuleUpdated, RuleDeleted, RuleEnabled, RuleDisabled, EvaluationCreated
- `securityscan` — ScanTriggered, ScanCompleted, ScanFailed

## Web Dashboard

The SoloDev web dashboard (`/pages/SoloDevDashboard`) shows a summary card grid for all eight product domains:

| Domain Card | Status Tracked |
|-------------|---------------|
| Pipelines | Active pipeline count |
| Security | Finding count |
| Quality Gates | Rule count |
| Error Tracker | Unresolved error count |
| Remediation | Pending task count |
| Health Monitor | Overall check status |
| Feature Flags | Flag count |
| Tech Debt | Item count |

The dashboard also shows an MCP connection status banner with a **Connect Client** button that navigates to the MCP setup page (`/pages/McpSetup`).

## Module Navigation (Sidebar)

The module navigation sidebar (`ModuleMenu`) in the web UI provides links to:
- Feature Flags
- Technical Debt
- Security Scanner
- Uptime Monitor
- Error Tracker
- Quality Gates
