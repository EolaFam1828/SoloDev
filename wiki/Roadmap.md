# Roadmap

SoloDev is an open-source DevOps platform that extends [Gitness by Harness](https://github.com/harness/gitness) with an AI intelligence layer for solo developers. This page tracks component maturity and planned milestones.

## Component Maturity

| Component | Status | Notes |
|-----------|--------|-------|
| Core platform (Gitness fork) | **Forked / compiling** | SCM, Pipelines (Drone-based), Gitspaces, Registry, all inherited |
| Web dashboard | **Working** | React/TypeScript UI with module pages and summary dashboard |
| AI remediation API | **Working** | CRUD endpoints, status tracking, trigger sources, patch storage |
| AI Worker (LLM integration) | **Working / early** | Background job polls pending tasks, calls LLM, parses diff + confidence |
| LLM provider support | **Working** | Anthropic, OpenAI, Google Gemini, Ollama (local) |
| Error Tracker | **Working** | Error groups, occurrences, fingerprinting, severity classification |
| Error-to-AI Bridge | **Working** | Auto-creates remediation tasks from errors and pipeline failures |
| Security Scanner | **Working** | SAST, SCA, secret scanning, finding tracking |
| Security → remediation bridge | **Working** | Scan findings auto-trigger remediation tasks |
| Quality Gates | **Working** | Rule CRUD, evaluation engine, coverage/complexity/style policies |
| Solo Gate Engine | **Working** | Strict/balanced/prototype enforcement modes |
| Health Monitor | **Working** | HTTP endpoint checks, uptime tracking, result history |
| Feature Flags | **Working** | Boolean and multivariate flags per space |
| Tech Debt Tracker | **Working** | Item CRUD, severity, categorization |
| Auto-Pipeline | **Working** | Stack detection, CI/CD YAML generation from file paths |
| MCP Server | **Working** | 16 atomic tools, 5 compound tools, 7 resources, 5 prompts, stdio + HTTP |
| Vector retrieval (Code Context Engine) | **Planned** | Embedding-based search over repository for richer LLM context |
| Auto-PR creation | **Planned** | Automatically open PRs from generated patches |
| Auto-merge (Agent Controller) | **Planned** | Configurable auto-merge for high-confidence patches |
| Signal Correlator | **Planned** | Cross-signal correlation between errors, failures, and metrics |
| Multi-agent MCP orchestration | **Planned** | Coordinate multiple AI agents across platform operations |
| Self-healing pipelines | **Planned** | Detect failure → generate fix → re-run pipeline automatically |
| Plugin ecosystem | **Planned** | Extend platform capabilities via community plugins |
| Hosted offering | **Planned** | Managed SoloDev deployment |

## Phased Milestones

### Phase 1 — Foundation

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

### Phase 2 — AI Remediation Loop

Close the loop from failure detection to applied fix.

| Milestone | Status |
|-----------|--------|
| Vector retrieval over repository for richer LLM context | Planned |
| Auto-PR creation from generated patches | Planned |
| Confidence-based auto-merge for trivial fixes | Planned |
| Pipeline failure → remediation → re-run flow | Planned |
| Remediation metrics and success rate tracking | Planned |

### Phase 3 — Autonomous Operations

Enable AI agents to operate the platform with minimal human intervention.

| Milestone | Status |
|-----------|--------|
| Multi-agent MCP orchestration | Planned |
| Self-healing pipeline loops | Planned |
| Signal Correlator (cross-domain error analysis) | Planned |
| Agent-driven deployment decisions | Planned |
| Anomaly detection from health check trends | Planned |

### Phase 4 — Platform Maturity

Scale beyond single-developer use and build ecosystem.

| Milestone | Status |
|-----------|--------|
| Plugin architecture for community extensions | Planned |
| Hosted / managed SoloDev offering | Planned |
| Community marketplace for pipeline templates | Planned |
| Multi-tenant workspace support | Planned |
| Enterprise SSO and audit logging | Planned |
