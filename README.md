# SoloDev

<p align="center">
  <img src="web/src/components/HarnessLogo/solodev-logo.svg" alt="SoloDev logo" width="164" />
</p>

<p align="center">
  <strong>An AI-native DevOps platform for solo builders, indie hackers, and tiny product teams.</strong>
</p>

<p align="center">
  Open source, Apache-2.0, and built from a fork of
  <a href="https://github.com/harness/gitness">Gitness by Harness</a>.
</p>

SoloDev takes the backbone of a modern DevOps platform and layers in the part I felt was missing: the brain. It combines code hosting, pipelines, registries, Gitspaces, security scanning, quality controls, runtime error tracking, AI remediation, and MCP-native agent access in one project aimed at people moving fast without a giant platform team behind them.

## Why SoloDev

I did not come into this through a traditional software engineering path. I came from recruiting, spent a lot of time around technical teams, and then dove hard into AI-assisted building. The more I built, the more obvious the gap became: the tools were powerful, but they were still missing the good stuff for the person trying to do everything alone.

I wanted the platform layer to do more than host code and run pipelines. I wanted it to understand what was breaking, surface the right problems, connect those problems to AI, and help me actually move from "something is wrong" to "here is the fix." That is why SoloDev exists.

SoloDev is my attempt to build the kind of DevOps product I would have wanted from day one: opinionated, agent-friendly, open source, and designed for the reality of modern one-person software teams.

## Who It Is For

- Solo founders shipping real products without dedicated platform, security, or DevOps staff.
- Technical operators who want one place to manage code, build pipelines, findings, errors, and AI-assisted fixes.
- Open-source contributors interested in AI-native developer tooling, MCP workflows, security/remediation loops, and better product ergonomics for small teams.

## What SoloDev Adds

| Area | What SoloDev focuses on | Docs |
|------|--------------------------|------|
| AI Auto-Remediation | Turn failures and findings into structured AI fix tasks with patch output | [docs/ai-remediation.md](docs/ai-remediation.md) |
| Zero-Config Pipelines | Detect stack and generate pipeline YAML automatically | [docs/auto-pipeline.md](docs/auto-pipeline.md) |
| Solo Gate Engine | Enforce quality through strict, balanced, or prototype-friendly modes | [docs/solo-gate.md](docs/solo-gate.md) |
| Error-to-AI Bridge | Feed runtime issues into remediation workflows | [docs/error-bridge.md](docs/error-bridge.md) |
| Error Tracker | Group, fingerprint, and manage application errors | [docs/error-tracker.md](docs/error-tracker.md) |
| Security Scanner | Run SAST, SCA, and secret scanning with finding tracking | [docs/security-scanner.md](docs/security-scanner.md) |
| Quality Gates | Evaluate coverage, complexity, and style policies | [docs/quality-gates.md](docs/quality-gates.md) |
| Health Monitor | Track uptime and service health | [docs/health-monitor.md](docs/health-monitor.md) |
| MCP Server | Expose SoloDev capabilities to AI agents over MCP | [docs/mcp-server.md](docs/mcp-server.md) |

## Project Status

SoloDev is being built in the open. The product direction is clear, the core platform is real, and the SoloDev-specific intelligence layer is actively expanding. Some internal command paths and package names still reflect upstream Gitness heritage while decoupling continues, but the project direction is SoloDev.

## Quick Start

### Run with Docker

Build a local image from this repo:

```bash
docker build -t solodev:local .

docker run -d \
  -p 3000:3000 \
  -p 3022:3022 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v solodev-data:/data \
  --name solodev \
  --restart unless-stopped \
  solodev:local
```

Visit `http://localhost:3000`.

### Run from Source

Prerequisites:

- Go 1.20+
- Node.js (latest stable)
- `protoc` 3.21.11

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

make dep && make tools

cd web && yarn install && yarn build && cd ..

make build
./gitness server .local.env
```

Current note: while SoloDev branding is the product direction, some local build and runtime entrypoints in the repo still use the inherited `gitness` command name.

### Authentication

```bash
./gitness login
./gitness user pat "my-pat" 2592000

curl http://localhost:3000/api/v1/user \
  -H "Authorization: Bearer $TOKEN"
```

### MCP

```bash
./gitness mcp stdio
./gitness mcp sse --port 3001
```

See [docs/mcp-server.md](docs/mcp-server.md) for setup details.

## Core API Surface

SoloDev capabilities are exposed under `/api/v1/spaces/{space_ref}/`:

```text
remediations/
auto-pipeline/
errors/
quality-gates/
security-scans/
health-checks/
tech-debt/
feature-flags/
```

Swagger and generated client support are available at:

- Swagger UI: `http://localhost:3000/swagger`
- OpenAPI spec: `http://localhost:3000/openapi.yaml`
- Registry API: `http://localhost:3000/registry/swagger/`

To refresh the frontend API client:

```bash
./gitness swagger > web/src/services/code/swagger.yaml
cd web && yarn services
```

## Open Source Provenance

SoloDev is a fully open-source fork built from the Apache-2.0 licensed
[Gitness repository by Harness](https://github.com/harness/gitness).

That provenance is intentional and explicit:

- upstream attribution remains preserved
- Apache-2.0 licensing remains in place
- derivative-work notices are maintained in [NOTICE](NOTICE)
- the project is being extended in public as SoloDev, not presented as an unrelated clean-room rewrite

If you are evaluating the project, that is the honest framing: SoloDev stands on top of a serious open-source foundation and pushes it in a more AI-native, solo-builder-focused direction.

## Contributing

If you care about AI-native DevOps, MCP tooling, remediation systems, solo-builder workflows, or making open-source developer infrastructure more useful, this project wants collaborators.

Good contribution areas:

- backend product surfaces and API contracts
- MCP tools and roadmap completion
- security and remediation workflows
- frontend/product UX polish
- docs, onboarding, and contributor ergonomics

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0. See [LICENSE](LICENSE).

> SoloDev is what happens when you give a non-traditional developer the tools to move fast and the stubbornness to keep going.
