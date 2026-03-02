# Integration

## Purpose

This page explains how external tools and services can interact with SoloDev's APIs and MCP protocol.

## REST API Integration

Any HTTP client can interact with SoloDev's REST API. All endpoints are documented in the [API Overview](Overview) and available via Swagger UI at `http://localhost:3000/swagger`.

### Example: Report Errors from a Production App

```bash
curl -X POST http://localhost:3000/api/v1/spaces/my-space/errors \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Connection refused",
    "message": "Failed to connect to payment service",
    "severity": "error",
    "file_path": "services/payment.go",
    "language": "go",
    "stack_trace": "...",
    "environment": "production"
  }'
```

### Example: Check Remediation Status

```bash
curl http://localhost:3000/api/v1/spaces/my-space/remediations?status=completed \
  -H "Authorization: Bearer $TOKEN"
```

## MCP Integration

AI agents connect via the Model Context Protocol for programmatic access to all platform capabilities.

### Stdio Transport (Local Agents)

```bash
gitness mcp stdio
```

Configure in your AI client's MCP settings with the `SOLODEV_TOKEN` environment variable.

### HTTP Transport (Remote Agents)

```bash
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
```

## Webhook Integration (Planned)

SoloDev will emit webhook notifications for platform events:
- Remediation completed
- Security scan finished
- Health check status changed
- Quality gate evaluation failed

## Frontend API Client

To regenerate the frontend API client from the live OpenAPI spec:

```bash
./gitness swagger > web/src/services/code/swagger.yaml
cd web && yarn services
```

## Authentication

All integrations authenticate using Personal Access Tokens (PATs) via the `Authorization: Bearer <token>` header. For MCP stdio transport, use the `SOLODEV_TOKEN` environment variable.
