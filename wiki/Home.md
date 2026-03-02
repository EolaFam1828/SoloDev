# SoloDev Wiki

<p align="center">
  <strong>An AI-native DevOps platform for solo builders, indie hackers, and tiny product teams.</strong>
</p>

<p align="center">
  Open source · Apache-2.0 · Built from a fork of <a href="https://github.com/harness/gitness">Gitness by Harness</a>
</p>

---

SoloDev takes the backbone of a modern DevOps platform and layers in the intelligence layer that was missing: automated error-to-fix pipelines, zero-config CI, enforced quality gates, security scanning, runtime error tracking, and first-class MCP support for AI agent access — all in one project aimed at people moving fast without a platform team.

## Platform Modules

| Module | Summary | Wiki Page |
|--------|---------|-----------|
| **AI Auto-Remediation** | Turn failures and findings into structured AI fix tasks with patch output | [AI-Auto-Remediation](AI-Auto-Remediation) |
| **Zero-Config Auto-Pipeline** | Detect stack, generate CI/CD pipeline YAML automatically | [Auto-Pipeline](Auto-Pipeline) |
| **Solo Gate Engine** | Enforce quality via strict, balanced, or prototype modes | [Solo-Gate](Solo-Gate) |
| **Error Tracker** | Group, fingerprint, and manage application runtime errors | [Error-Tracker](Error-Tracker) |
| **Error-to-AI Bridge** | Feed runtime errors into the remediation workflow automatically | [Error-Bridge](Error-Bridge) |
| **Security Scanner** | Run SAST, SCA, and secret scanning with finding tracking | [Security-Scanner](Security-Scanner) |
| **Quality Gates** | Evaluate coverage, complexity, and style policies | [Quality-Gates](Quality-Gates) |
| **Health Monitor** | Track uptime and HTTP endpoint health | [Health-Monitor](Health-Monitor) |
| **MCP Server** | Expose all SoloDev capabilities to AI agents over Model Context Protocol | [MCP-Server](MCP-Server) |
| **Feature Flags** | Create and toggle boolean and multivariate feature flags per space | [Feature-Flags](Feature-Flags) |
| **Tech Debt** | Log, triage, and track technical debt items | [Tech-Debt](Tech-Debt) |

## Core API Surface

All SoloDev-specific capabilities are exposed under `/api/v1/spaces/{space_ref}/`:

```
remediations/
auto-pipeline/
errors/
quality-gates/
security-scans/
health-checks/
tech-debt/
feature-flags/
```

- Swagger UI: `http://localhost:3000/swagger`
- OpenAPI spec: `http://localhost:3000/openapi.yaml`
- Registry API: `http://localhost:3000/registry/swagger/`

## Quick Links

- [Getting Started](Getting-Started) — installation, Docker, run from source
- [Architecture](Architecture) — platform architecture and component map
- [API Reference](API-Reference) — full API endpoint reference
- [MCP Server](MCP-Server) — connect AI agents (Claude Desktop, Cursor, custom)
- [Contributing](Contributing) — how to contribute

## Project Status

SoloDev is being built in the open. The product direction is clear, the core platform is real, and the SoloDev-specific intelligence layer is actively expanding. Some internal command paths and package names still reflect upstream Gitness heritage while decoupling continues.

## Open Source Provenance

SoloDev is a fully open-source fork built from the Apache-2.0 licensed [Gitness repository by Harness](https://github.com/harness/gitness). Upstream attribution is preserved, Apache-2.0 licensing remains in place, and derivative-work notices are maintained in [NOTICE](https://github.com/EolaFam1828/SoloDev/blob/main/NOTICE).

## License

Apache License 2.0.
