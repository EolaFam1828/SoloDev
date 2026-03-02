# Scaling

How SoloDev can scale with usage.

## Current Architecture

SoloDev runs as a single process with all components in one binary:
- Web server (API + UI)
- AI Worker (background poller)
- MCP server
- Health check runner

This is sufficient for the target audience (solo developers and small teams) and simplifies deployment.

## Vertical Scaling

The simplest scaling path is increasing server resources:

| Resource | Impact |
|----------|--------|
| RAM | More concurrent pipeline executions, larger codebase analysis |
| CPU | Faster builds, more concurrent health checks |
| Disk I/O | Faster Git operations and database queries |

## Database Scaling

- **SQLite** — Suitable for single-user local development
- **PostgreSQL** — Recommended for production; supports concurrent access, better durability, standard backup tools

## AI Worker Scaling (Planned)

The AI Worker is the most resource-intensive component due to LLM API calls. Planned scaling:

1. **Multiple worker instances** — Separate the AI Worker from the API server
2. **Worker pool** — Multiple workers reading from the same remediation queue
3. **Provider load balancing** — Distribute requests across multiple LLM providers

## Distributed Topology (Planned)

```
Load Balancer
      │
      ├──▶ API Server 1 (stateless)
      ├──▶ API Server 2 (stateless)
      │
      ├──▶ AI Worker 1
      ├──▶ AI Worker 2
      │
      └──▶ PostgreSQL (shared state)
```

API servers are stateless and can be horizontally scaled. AI Workers read from the shared database queue.

## Status

| Scaling Path | Status |
|-------------|--------|
| Single binary (vertical scaling) | Implemented |
| PostgreSQL for production | Implemented |
| Separated AI Worker | Planned |
| Horizontal API scaling | Planned |
| Kubernetes Helm chart | Planned |

## Future Work

- Kubernetes Helm chart for container orchestration
- Redis-based job queue for AI Worker distribution
- Horizontal API server scaling behind a load balancer
- Managed / hosted SoloDev offering
