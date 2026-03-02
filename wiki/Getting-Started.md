# Getting Started

SoloDev is an open-source DevOps platform that extends Gitness with AI-driven remediation, error tracking, security scanning, and MCP-based agent access. This guide gets you running locally in under 5 minutes.

## Prerequisites

**Docker path** (recommended): Docker Engine 20.10+ with Docker Compose V2.

**Source path**: Go 1.20+, Node.js 16+ with Yarn, `protoc` 3.21.11.

## Quick Start with Docker Compose

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

The first build takes several minutes (compiling Go, building the React frontend). Subsequent starts are fast. Once running, open [http://localhost:3000](http://localhost:3000).

The container exposes two ports:

| Port | Service |
|------|---------|
| 3000 | Web UI + REST API |
| 3022 | SSH (Git over SSH) |

Data is persisted to a Docker volume (`solodev-data`). To stop: `docker compose down`. To stop and delete data: `docker compose down -v`.

## First-Use Walkthrough

### 1. Create an account

Open `http://localhost:3000` and register. The first user becomes the admin.

### 2. Create a space

A space is an organizational container (like a GitHub organization). Click **New Space**, give it a name, and enter it.

### 3. Create a repository

Inside the space, click **New Repository**. Initialize it with a README or push an existing project:

```bash
git remote add solodev http://localhost:3000/git/<space>/<repo>.git
git push solodev main
```

### 4. Trigger a pipeline

Create a `.harness/pipeline.yaml` in your repo (or use Auto-Pipeline to generate one — see [Auto-Pipeline](Auto-Pipeline)):

```yaml
kind: pipeline
spec:
  stages:
    - type: ci
      spec:
        steps:
          - name: test
            type: run
            spec:
              container: golang:1.21
              script: go test ./...
```

Push the file and the pipeline runs automatically.

### 5. View the SoloDev dashboard

Navigate to the SoloDev dashboard from the sidebar menu. It shows summary cards for all modules: Pipelines, Security, Quality Gates, Error Tracker, Remediation, Health Monitor, Feature Flags, and Tech Debt.

### 6. Report an error and see AI remediation

Post an error to the Error Tracker API:

```bash
TOKEN="<your-personal-access-token>"

curl -X POST http://localhost:3000/api/v1/spaces/<space>/errors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "nil pointer in handler",
    "message": "runtime error: invalid memory address",
    "severity": "error",
    "language": "go",
    "file_path": "app/handler.go",
    "stack_trace": "goroutine 1:\napp/handler.go:42 +0x1a"
  }'
```

If the Error Bridge is enabled, a remediation task is automatically created. If the AI Worker is configured with an LLM provider, it picks up the task, generates a patch, and stores the result. Check the remediation via the API or dashboard.

## Build from Source

```bash
# Install protobuf tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Install dependencies
make dep && make tools

# Build frontend
cd web && yarn install && yarn build && cd ..

# Build and run
make build
./gitness server .local.env
```

> **Note:** The binary and some internal paths use the `gitness` name inherited from the upstream fork. The product is SoloDev.

## MCP Server

Connect an AI agent (Claude Desktop, Cursor, or custom) via the Model Context Protocol:

```bash
# Stdio transport (for local AI clients)
./gitness mcp stdio

# HTTP transport (for remote clients)
./gitness mcp sse --port 3001
```

See [MCP Server](MCP-Server) for full configuration and tool reference.

## Authentication

Create a personal access token for API and MCP access:

```bash
./gitness login
./gitness user pat "my-pat" 2592000
```

Use it in API requests:

```bash
curl http://localhost:3000/api/v1/user \
  -H "Authorization: Bearer $TOKEN"
```

For MCP stdio transport, set `SOLODEV_TOKEN` in your environment.

## Troubleshooting

**Build fails during Docker Compose up**
The Go build requires significant memory. Ensure Docker has at least 4 GB RAM allocated. On Apple Silicon, the cross-compilation step may take longer.

**Port 3000 already in use**
Stop the conflicting service, or change the port mapping in `docker-compose.yml`: `"3001:3000"`.

**Pipeline does not trigger**
Ensure the pipeline file is at `.harness/pipeline.yaml` in the repository root and that the file was pushed to a branch the pipeline is configured to watch.

**AI Worker not generating patches**
The AI Worker requires an LLM provider to be configured. Check that the provider (Anthropic, OpenAI, Gemini, or Ollama) is set up in the server configuration and that the API key is valid.

**Permission denied on `/var/run/docker.sock`**
The container needs access to the Docker socket for Gitspaces and pipeline execution. On Linux, add your user to the `docker` group or run with appropriate permissions.

## Next Steps

- [Architecture](Architecture) — understand how modules connect
- [MCP Server](MCP-Server) — connect an AI agent to the platform
- [API Reference](API-Reference) — explore all endpoints
- [Roadmap](Roadmap) — see what is working today vs. planned
