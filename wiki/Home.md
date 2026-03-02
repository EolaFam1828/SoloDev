# SoloDev

**An open-source, AI-native DevOps platform for solo builders.**

SoloDev extends [Gitness by Harness](https://github.com/harness/gitness) with an AI intelligence layer that connects runtime failures to automated code fixes. It targets indie hackers, solo SaaS founders, and small teams who need DevOps without the DevOps team. When a pipeline fails or an error is reported, SoloDev auto-creates a remediation task, sends it to an LLM, and produces a patch — see [Architecture](Architecture) for the full data flow.

## Documentation

| Page | Description |
|------|-------------|
| [Getting Started](Getting-Started) | Clone, run, and deploy in under 5 minutes |
| [Architecture](Architecture) | System layers, module map, AI remediation data flow |
| [Differentiation](Differentiation) | Comparison with GitHub, GitLab, Harness OSS |
| [Roadmap](Roadmap) | Component maturity and phased milestones |
| [API Reference](API-Reference) | Endpoint listing for all modules |
| [MCP Server](MCP-Server) | AI agent integration via Model Context Protocol |
| [Contributing](Contributing) | How to contribute |
| [Why SoloDev](Why-SoloDev) | The problem space SoloDev addresses |

## Quick Start

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

Open [http://localhost:3000](http://localhost:3000), create an account, and start building.

## License

Apache License 2.0 — [LICENSE](https://github.com/EolaFam1828/SoloDev/blob/main/LICENSE).
