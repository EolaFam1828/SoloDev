# Harness Integration

SoloDev is built from a fork of [Gitness by Harness](https://github.com/harness/gitness), an open-source DevOps platform licensed under Apache 2.0. This page clearly documents which components are inherited from Harness/Gitness and which are added by SoloDev.

## Inherited from Gitness (Layer 1 — Core)

These components exist in the upstream Gitness repository and are used by SoloDev with minimal modification:

| Component | Gitness Package | Description |
|-----------|----------------|-------------|
| Git hosting (SCM) | `app/api/controller/repo/` | Repository management, commits, branches, tags, pull requests, code review, webhooks |
| Pipelines | `app/pipeline/` | Drone-based CI/CD execution engine |
| Gitspaces | `app/api/controller/gitspace/` | Cloud-based developer environments |
| Registry | `registry/` | OCI artifact and container image storage |
| Authentication | `app/auth/` | User management, sessions, personal access tokens |
| Web UI framework | `web/` | React/TypeScript application shell, routing, common components |
| Database framework | `app/store/database/` | Migration system, store interfaces, connection management |
| CLI framework | `cli/`, `cmd/gitness/` | Cobra-based command-line interface |

## Added by SoloDev (Layers 2–3)

These components are entirely new, built on top of the Gitness foundation:

### Layer 2 — Observability

| Component | SoloDev Package | Migration |
|-----------|----------------|-----------|
| Error Tracker | `app/api/controller/errortracker/` | 0171 |
| Security Scanner | `app/api/controller/securityscan/` | 0102 |
| Quality Gates | `app/api/controller/qualitygate/` | 0105 |
| Health Monitor | `app/api/controller/healthcheck/` | 0103, 0104 |
| Feature Flags | `app/api/controller/featureflag/` | — |
| Tech Debt Tracker | `app/api/controller/techdebt/` | — |

### Layer 3 — Intelligence

| Component | SoloDev Package |
|-----------|----------------|
| AI Auto-Remediation | `app/api/controller/airemediation/` (migration 0172) |
| AI Worker | `app/services/aiworker/` |
| Error Bridge | `app/services/errorbridge/` |
| Solo Gate Engine | `app/services/sologate/` |
| Auto-Pipeline | `app/pipeline/autopipeline/` |
| Security Remediation | `app/services/securityremediation/` |
| MCP Server | `mcp/` |

### Layer 4 — UI Additions

| Component | SoloDev Path |
|-----------|-------------|
| SoloDev Dashboard | `web/src/pages/SoloDevDashboard/` |
| Module-specific pages | `web/src/pages/{ErrorList,SecurityScanList,QualityGateList,...}/` |
| MCP Setup page | `web/src/pages/McpSetup/` |

## Modified Gitness Components

Some upstream components have been modified to support SoloDev features:

| Component | Modification |
|-----------|-------------|
| Router | `app/router/api_modules.go` added to register all SoloDev module routes |
| Wire setup | `cmd/gitness/wire.go` extended with SoloDev dependency injection |
| Web sidebar | Navigation items added for SoloDev modules |
| Docker configuration | `docker-compose.yml` and `Dockerfile` updated for SoloDev build |

## Naming Heritage

Some internal command paths and package names still use the `gitness` name inherited from the upstream fork. The binary is currently built as `./gitness` and some configuration keys reference the upstream naming. This is a known decoupling task — the product direction is SoloDev.

| Current Name | Context | Plan |
|-------------|---------|------|
| `./gitness` | Binary name | Will be renamed to `./solodev` |
| `gitness server` | CLI command | Will be renamed |
| `gitness mcp` | CLI command | Will be renamed |

## Upstream Synchronization

SoloDev maintains awareness of upstream Gitness changes but does not automatically merge upstream commits. The fork is extended in a direction (AI-native, solo-builder-focused) that diverges from the upstream roadmap. Upstream changes are evaluated and cherry-picked when beneficial.

## License Compliance

- Apache 2.0 license is preserved from upstream
- Derivative-work notices are maintained in [NOTICE](https://github.com/EolaFam1828/SoloDev/blob/main/NOTICE)
- Upstream attribution remains in all inherited files
- SoloDev additions are also licensed under Apache 2.0
