# Installation

Step-by-step setup instructions for SoloDev.

## Option 1: Docker Compose (Recommended)

```bash
git clone https://github.com/EolaFam1828/SoloDev.git
cd SoloDev
docker compose up -d
```

The first build takes several minutes (compiling Go, building the React frontend). Subsequent starts are fast.

The container exposes:

| Port | Service |
|------|---------|
| 3000 | Web UI + REST API |
| 3022 | SSH (Git over SSH) |

Data is persisted to a Docker volume (`solodev-data`).

```bash
# Stop
docker compose down

# Stop and delete data
docker compose down -v
```

## Option 2: Docker Image

Build and run a local image:

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

## Option 3: Build from Source

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

## Authentication Setup

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

## MCP Setup

Connect an AI agent via the Model Context Protocol:

```bash
# Stdio transport (for local AI clients)
./gitness mcp stdio

# HTTP transport (for remote clients)
./gitness mcp sse --port 3001
```

For MCP stdio transport, set `SOLODEV_TOKEN` in your environment.

## Verify Installation

Open [http://localhost:3000](http://localhost:3000). You should see the SoloDev login page. Register a new account — the first user becomes the admin.
