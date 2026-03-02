# SoloDev

**An open-source, autonomous feedback system for software development.**

SoloDev extends [Gitness by Harness](https://github.com/harness/gitness) with an AI intelligence layer that connects runtime failures, security findings, and pipeline errors to automated code fixes. It targets indie hackers, solo SaaS founders, and small teams who need DevOps without the DevOps team.

SoloDev is not just DevOps tooling. It is a closed-loop system: when something breaks, the platform detects it, analyzes it, proposes a fix, validates the result, and applies the change — with the developer reviewing rather than performing each step.

## Core Loop

```
Detect → Analyze → Propose → Validate → Apply
```

A pipeline fails. The Error Bridge captures the failure context. The AI Worker generates a patch. The Validation Engine scores it. The developer reviews and merges — or an MCP-connected agent handles it autonomously.

## Who It Is For

- Solo founders shipping real products without dedicated platform, security, or DevOps staff.
- Technical operators who want one place to manage code, pipelines, findings, errors, and AI-assisted fixes.
- Open-source contributors interested in AI-native developer tooling, MCP workflows, and autonomous remediation loops.

## Documentation

| Section | Description |
|---------|-------------|
| [Why SoloDev](Why-SoloDev) | Technical motivation and gap analysis |
| [Core Concepts](Core-Concepts/Design-Principles) | Design principles, the SoloDev Loop, and terminology |
| [Architecture](Architecture/Overview) | System layers, data flow, and deployment topology |
| [Modules](Modules/Core-Platform/Source-Control) | Platform components from Git hosting to AI remediation |
| [Workflows](Workflows/Pipeline-Failure-Remediation) | End-to-end examples of the feedback loop in action |
| [API](API/Overview) | REST endpoints and event streams |
| [Getting Started](Getting-Started/Quick-Start) | Install, configure, and run your first pipeline |
| [Deployment](Deployment/Local-Deployment) | Local, self-hosted, and scaled deployment options |
| [Roadmap](Roadmap/Current-State) | What exists today and what is planned |
| [Contributing](Contributing/How-To-Contribute) | How to contribute to SoloDev |
| [Reference](Reference/Glossary) | Glossary, FAQ, and platform comparisons |

## Quick Start

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

Open [http://localhost:3000](http://localhost:3000), create an account, and start building.

## License

Apache License 2.0 — [LICENSE](https://github.com/EolaFam1828/SoloDev/blob/main/LICENSE).
