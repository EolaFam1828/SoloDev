# SoloDev

AI-accelerated DevOps platform for solo developers. Code hosting, pipelines,
artifact registries, Gitspaces — plus AI auto-remediation, zero-config
pipelines, quality/security gates, and a live error-to-AI feedback loop.

Built on [Harness Open Source](https://github.com/harness/gitness).

## Modules

| Module | What it does | Docs |
|--------|-------------|------|
| AI Auto-Remediation | Captures failures, generates patches, opens PRs | [docs/ai-remediation.md](docs/ai-remediation.md) |
| Zero-Config Pipelines | Detects stack, generates pipeline YAML | [docs/auto-pipeline.md](docs/auto-pipeline.md) |
| Solo Gate Engine | Strict / Balanced / Prototype enforcement modes | [docs/solo-gate.md](docs/solo-gate.md) |
| Error-to-AI Bridge | Auto-creates remediation tasks from runtime errors | [docs/error-bridge.md](docs/error-bridge.md) |
| Error Tracker | Error grouping, fingerprinting, occurrence tracking | [docs/error-tracker.md](docs/error-tracker.md) |
| Security Scanner | SAST, SCA, secret detection with finding management | [docs/security-scanner.md](docs/security-scanner.md) |
| Quality Gates | Rule-based quality eval (coverage, complexity, style) | [docs/quality-gates.md](docs/quality-gates.md) |
| Health Monitor | HTTP endpoint uptime monitoring | [docs/health-monitor.md](docs/health-monitor.md) |
| MCP Server | Model Context Protocol interface for AI agents | [docs/mcp-server.md](docs/mcp-server.md) |

## API Routes

All SoloDev endpoints live under `/api/v1/spaces/{space_ref}/`:

```
remediations/       POST  GET  GET /summary  GET /{id}  PATCH /{id}
auto-pipeline/      POST /generate
errors/             POST  GET  GET /summary  GET /{id}  PATCH /{id}  GET /{id}/occurrences
quality-gates/      POST /evaluate  GET /summary  rules/*  evaluations/*
security-scans/     POST  GET  GET /{id}  GET /{id}/findings  GET /{id}/summary
health-checks/      POST  GET  CRUD /{id}  GET /{id}/results  GET /{id}/summary
tech-debt/          POST  GET  CRUD /{id}  GET /{id}/summary
feature-flags/      POST  GET  CRUD /{id}  POST /{id}/toggle
```

The MCP Server is mounted at `/mcp` (Streamable HTTP) or via `gitness mcp stdio`
for Claude Desktop. See [docs/mcp-server.md](docs/mcp-server.md).

## Quick Start

### Docker

```bash
docker run -d \
  -p 3000:3000 \
  -p 3022:3022 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp/harness:/data \
  --name solodev \
  --restart always \
  harness/harness
```

Visit `http://localhost:3000`. Use a bind mount or named volume — data is lost
if the container stops without one.

### Build from Source

Prerequisites: Go 1.20+, Node.js (latest stable), protoc 3.21.11.

```bash
# Protobuf tooling
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
make dep && make tools

# UI
pushd web && yarn install && yarn build && popd

# Backend
make build

# Run
./gitness server .local.env
```

### Auth

```bash
./gitness login                        # user: admin, pw: changeit
./gitness user pat "my-pat" 2592000    # 30-day PAT

curl http://localhost:3000/api/v1/user \
  -H "Authorization: Bearer $TOKEN"
```

### MCP Server

```bash
gitness mcp stdio               # Claude Desktop / local agents
gitness mcp sse --port 3001     # Remote HTTP clients
```

## Docker Socket Config

| Runtime | Socket | Fix |
|---------|--------|-----|
| Docker Desktop | `/var/run/docker.sock` | Default |
| Rancher Desktop | `~/.rd/docker.sock` | `sudo ln -sf ~/.rd/docker.sock /var/run/docker.sock` |
| Colima | `~/.colima/default/docker.sock` | `sudo ln -sf ~/.colima/default/docker.sock /var/run/docker.sock` |

Or set `GITNESS_DOCKER_HOST` in `.local.env`.

## Swagger / API Client

```bash
./gitness swagger > web/src/services/code/swagger.yaml
cd web && yarn services
```

- Swagger UI: `http://localhost:3000/swagger`
- OpenAPI spec: `http://localhost:3000/openapi.yaml`
- Registry API: `http://localhost:3000/registry/swagger/`

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
