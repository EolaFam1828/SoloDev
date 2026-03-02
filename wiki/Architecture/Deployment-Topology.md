# Deployment Topology

SoloDev can be deployed in multiple configurations depending on the operator's needs. All configurations run from the same codebase and binary.

## Local Development

Single binary running on the developer's machine. Suitable for development and evaluation.

```
┌─────────────────────────────────────┐
│          Developer Machine           │
│                                      │
│  ┌────────────────────────────────┐ │
│  │        SoloDev Binary          │ │
│  │                                │ │
│  │  Web UI (:3000)                │ │
│  │  REST API (:3000)              │ │
│  │  Git SSH (:3022)               │ │
│  │  MCP Server (stdio or :3001)   │ │
│  │  SQLite Database               │ │
│  │  AI Worker (background)        │ │
│  └────────────────────────────────┘ │
│                                      │
│  Volume: ./data/                     │
└─────────────────────────────────────┘
```

- Database: SQLite (file-based, zero configuration)
- LLM: Ollama (local) or external provider
- Storage: Local filesystem
- No external dependencies required

## Docker Compose (Recommended)

Single-command deployment using Docker Compose. This is the recommended path for most users.

```
┌─────────────────────────────────────┐
│          Docker Host                 │
│                                      │
│  ┌────────────────────────────────┐ │
│  │     SoloDev Container          │ │
│  │                                │ │
│  │  Ports: 3000, 3022             │ │
│  │  Volume: solodev-data          │ │
│  │  Docker socket: mounted        │ │
│  └────────────────────────────────┘ │
└─────────────────────────────────────┘
```

```bash
docker compose up -d
```

- Database: SQLite within the container volume
- Persistent storage via Docker volume (`solodev-data`)
- Docker socket mounted for Gitspaces and pipeline execution

## Self-Hosted (Production)

For production use, SoloDev can be deployed on a VPS or dedicated server with PostgreSQL.

```
┌─────────────────────────────────────────────────┐
│                 Production Host                   │
│                                                   │
│  ┌─────────────────┐   ┌──────────────────────┐ │
│  │ SoloDev Binary   │   │ PostgreSQL           │ │
│  │                  │──▶│                      │ │
│  │ Web UI + API     │   │ Production database  │ │
│  │ AI Worker        │   └──────────────────────┘ │
│  │ MCP Server       │                            │
│  └─────────────────┘                             │
│                                                   │
│  ┌─────────────────┐                             │
│  │ Reverse Proxy    │  (nginx, Caddy, etc.)      │
│  │ TLS termination  │                            │
│  └─────────────────┘                             │
└─────────────────────────────────────────────────┘
```

- Database: PostgreSQL (recommended for durability and concurrent access)
- TLS: Terminated at reverse proxy
- LLM: External provider (Anthropic, OpenAI, Gemini) or self-hosted Ollama

## Distributed (Planned)

Future topology for scaling beyond single-instance deployment.

```
┌──────────────────┐   ┌──────────────────┐   ┌──────────────────┐
│ SoloDev API      │   │ AI Worker Pool   │   │ PostgreSQL       │
│ (stateless)      │──▶│ (N instances)    │──▶│ (shared DB)      │
│ Web UI + REST    │   │ LLM processing   │   │                  │
│ MCP Server       │   │                  │   │                  │
└──────────────────┘   └──────────────────┘   └──────────────────┘
```

This topology separates the API server from the AI Worker, allowing independent scaling. Not yet implemented.

## Status

| Topology | Status |
|----------|--------|
| Local binary (SQLite) | Implemented |
| Docker Compose | Implemented |
| Self-hosted with PostgreSQL | Implemented |
| Distributed / multi-instance | Planned |

## Future Work

- Horizontal scaling of AI Worker instances
- Kubernetes Helm chart for container orchestration
- Managed / hosted SoloDev offering
