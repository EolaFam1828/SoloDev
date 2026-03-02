# SoloDev

**An AI-native DevOps platform for solo builders, indie hackers, and tiny product teams.**

Open source, Apache-2.0, and built from a fork of [Gitness by Harness](https://github.com/harness/gitness).

SoloDev combines code hosting, pipelines, registries, Gitspaces, security scanning, quality controls, runtime error tracking, AI remediation, and MCP-native agent access into a single autonomous feedback system for software development.

## Features

| Area | Description |
|------|-------------|
| AI Auto-Remediation | Failures and findings become structured AI fix tasks with patch output |
| Zero-Config Pipelines | Detect tech stack and generate pipeline YAML automatically |
| Solo Gate Engine | Enforce quality through strict, balanced, or prototype-friendly modes |
| Error-to-AI Bridge | Feed runtime issues into remediation workflows |
| Error Tracker | Group, fingerprint, and manage application errors |
| Security Scanner | Run SAST, SCA, and secret scanning with finding tracking |
| Quality Gates | Evaluate coverage, complexity, and style policies |
| Health Monitor | Track uptime and service health |
| MCP Server | Expose SoloDev capabilities to AI agents over Model Context Protocol |
| Tech Debt Tracker | Detect and track technical debt over time |

## Quick Install

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

Visit `http://localhost:3000`.

## Documentation

- [Wiki Home](wiki/Home.md) — Executive overview
- [Why SoloDev](wiki/Why-SoloDev.md) — Technical motivation
- [Architecture Overview](wiki/Architecture/Overview.md) — System layers and data flow
- [Getting Started](wiki/Getting-Started/Quick-Start.md) — Install and run in minutes
- [API Overview](wiki/API/Overview.md) — REST endpoints and MCP tools
- [Roadmap](wiki/Roadmap/Current-State.md) — What exists and what is planned
- [Contributing](wiki/Contributing/How-To-Contribute.md) — How to help

## Project Status

SoloDev is being built in the open. The core platform is real, the AI remediation loop produces patches, and the MCP server exposes the full platform to AI agents. Some internal command paths and package names still reflect upstream Gitness heritage while decoupling continues.

## Open Source Provenance

SoloDev is a fully open-source fork built from the Apache-2.0 licensed [Gitness repository by Harness](https://github.com/harness/gitness). Upstream attribution is preserved and derivative-work notices are maintained in [NOTICE](NOTICE).

## License

Apache License 2.0. See [LICENSE](LICENSE).
