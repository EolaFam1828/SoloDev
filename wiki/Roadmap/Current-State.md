# Current State

An honest description of what exists in SoloDev today.

## Component Maturity

| Component | Status | Notes |
|-----------|--------|-------|
| Core platform (Gitness fork) | **Forked / compiling** | SCM, Pipelines (Drone-based), Gitspaces, Registry — all inherited |
| Web dashboard | **Implemented** | React/TypeScript UI with module pages and summary dashboard |
| AI remediation API | **Implemented** | CRUD endpoints, status tracking, trigger sources, patch storage, manual apply endpoint |
| AI Worker (LLM integration) | **Implemented** | Background poller, source enrichment, diff generation, optional auto-delivery to draft PR |
| LLM provider support | **Implemented** | Anthropic, OpenAI, Google Gemini, Ollama (local) |
| Error Tracker | **Implemented** | Error groups, occurrences, fingerprinting, severity classification |
| Error-to-AI Bridge | **Implemented** | Auto-creates remediation tasks from errors and pipeline failures |
| Security Scanner | **Implemented** | SAST, SCA, secret scanning, finding tracking |
| Security → remediation bridge | **Implemented** | Scan findings auto-trigger remediation tasks |
| Quality Gates | **Implemented** | Rule CRUD, evaluation engine, coverage/complexity/style policies |
| Solo Gate Engine | **Implemented** | Strict/balanced/prototype enforcement modes |
| Health Monitor | **Implemented** | HTTP endpoint checks, uptime tracking, result history |
| Feature Flags | **Implemented** | Boolean and multivariate flags per space |
| Tech Debt Tracker | **Implemented** | Item CRUD, severity, categorization |
| Auto-Pipeline | **Implemented** | Stack detection, CI/CD YAML generation from file paths |
| MCP Server | **Implemented** | Runtime-active surfaces are live; broader catalog also tracks blocked and coming-soon tools/resources/prompts |
| Vector retrieval (Code Context Engine) | **In Progress** | Prototype embedding-based indexing/search exists; ongoing hardening for broader remediation coverage |
| Draft PR creation from remediations | **Implemented** | Manual apply path plus opt-in worker auto-delivery via `SOLODEV_AI_REMEDIATION_CREATE_FIX_BRANCH` |
| Auto-merge (Agent Controller) | **Planned** | Configurable auto-merge for high-confidence patches |
| Signal Correlator | **Planned** | Cross-signal correlation between errors, failures, and metrics |
| Multi-agent MCP orchestration | **Planned** | Coordinate multiple AI agents across platform operations |
| Self-healing pipelines | **Planned** | Detect failure → generate fix → re-run pipeline automatically |

## What Works End-to-End Today

1. Report an error → Error Bridge creates remediation → AI Worker generates patch → Manual or opt-in automatic draft PR delivery
2. Run a security scan → Findings trigger remediation → AI Worker generates patch → Manual or opt-in automatic draft PR delivery
3. Detect stack → Generate pipeline YAML → Execute pipeline
4. AI agent connects via MCP → Uses tools to interact with all modules
5. Configure quality gates → Evaluate code → Solo Gate enforces based on mode

## What Does Not Work Yet

1. Auto-merge for high-confidence patches
2. Self-healing pipeline loop (fix → re-run → verify)
3. Production-grade vector-based code retrieval for richer AI context
4. Cross-signal correlation
5. Multi-agent orchestration

## Known Limitations

- Binary and some CLI commands still use the `gitness` name from upstream
- Health check → remediation path is not yet connected
- Remediation webhooks exist, but broader cross-module notification routing is still limited
