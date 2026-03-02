# Getting Started

## Prerequisites

- Go 1.20+
- Node.js (latest stable)
- `protoc` 3.21.11

## Run with Docker

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

## Run from Source

Install protobuf code-generation tools:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

Install dependencies, build the frontend, and build the binary:

```bash
make dep && make tools

cd web && yarn install && yarn build && cd ..

make build
```

Start the server:

```bash
./gitness server .local.env
```

> **Note:** While SoloDev branding is the product direction, some local build and runtime entrypoints in the repo still use the inherited `gitness` command name from the upstream Gitness fork.

Visit `http://localhost:3000`.

## Authentication

Log in and create a personal access token (PAT):

```bash
./gitness login
./gitness user pat "my-pat" 2592000
```

Use the token with the REST API:

```bash
curl http://localhost:3000/api/v1/user \
  -H "Authorization: Bearer $TOKEN"
```

## MCP Server

Start the MCP server in stdio mode (for Claude Desktop):

```bash
./gitness mcp stdio
```

Start the MCP server in SSE mode (for remote clients):

```bash
./gitness mcp sse --port 3001
```

See [MCP Server](MCP-Server) for full configuration details.

## Refresh the Frontend API Client

```bash
./gitness swagger > web/src/services/code/swagger.yaml
cd web && yarn services
```

## Environment Configuration

The server is configured via an environment file passed as the first argument (`./gitness server .local.env`). Common variables:

| Variable | Purpose |
|----------|---------|
| `SOLODEV_TOKEN` | Bearer token used by MCP stdio transport |

## Next Steps

- Explore the [Architecture](Architecture) to understand how modules fit together
- Set up the [MCP Server](MCP-Server) to connect an AI agent
- Review [API Reference](API-Reference) for the full endpoint listing
