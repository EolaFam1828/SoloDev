# Milestones

Planned development phases for SoloDev.

## Phase 1 — Foundation (Complete)

Establish a stable, deployable fork with the core SoloDev modules operational.

| Milestone | Status |
|-----------|--------|
| Fork Gitness, build and run locally | Done |
| Docker image builds and serves on port 3000 | Done |
| Docker Compose for single-command deployment | Done |
| AI remediation module with CRUD API | Done |
| Error Tracker with grouping and fingerprinting | Done |
| Error Bridge auto-creates remediation tasks | Done |
| Security Scanner with finding tracking | Done |
| Quality Gates with rule evaluation | Done |
| Solo Gate enforcement engine | Done |
| Health Monitor with uptime checks | Done |
| Feature Flags and Tech Debt modules | Done |
| Auto-Pipeline stack detection and YAML generation | Done |
| MCP server with full tool/resource/prompt coverage | Done |
| Web dashboard with summary cards for all modules | Done |
| AI Worker prototype with LLM-driven patch generation | Done |

## Phase 2 — AI Remediation Loop

Close the loop from failure detection to applied fix.

| Milestone | Status |
|-----------|--------|
| Vector retrieval over repository for richer LLM context | In Progress (prototype indexing/search APIs available) |
| Auto-PR creation from generated patches | Implemented (opt-in via `SOLODEV_AI_REMEDIATION_CREATE_FIX_BRANCH`) |
| Confidence-based auto-merge for trivial fixes | Planned |
| Pipeline failure → remediation → re-run flow | Planned |
| Remediation metrics and success rate tracking | In Progress (space-level metrics endpoint available) |

## Phase 3 — Autonomous Operations

Enable AI agents to operate the platform with minimal human intervention.

| Milestone | Status |
|-----------|--------|
| Multi-agent MCP orchestration | Planned |
| Self-healing pipeline loops | Planned |
| Signal Correlator (cross-domain error analysis) | Planned |
| Agent-driven deployment decisions | Planned |
| Anomaly detection from health check trends | Planned |

## Phase 4 — Platform Maturity

Scale beyond single-developer use and build ecosystem.

| Milestone | Status |
|-----------|--------|
| Plugin architecture for community extensions | Planned |
| Hosted / managed SoloDev offering | Planned |
| Community marketplace for pipeline templates | Planned |
| Multi-tenant workspace support | Planned |
| Enterprise SSO and audit logging | Planned |
