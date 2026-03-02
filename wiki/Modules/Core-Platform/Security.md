# Security (Baseline)

## Purpose

Provides baseline security capabilities inherited from the Gitness platform. This page covers the foundational security model. For SoloDev's added security scanning capabilities, see [Security Scanner](../Signal-System/Error-Tracker).

## Inputs

- Authentication credentials (username/password, personal access tokens)
- Authorization requests against spaces and repositories
- Session management events

## Processing

- User authentication via local accounts
- Personal Access Token (PAT) generation and validation
- Role-based access control on spaces and repositories
- Session lifecycle management

## Outputs

- Authenticated user sessions
- Bearer tokens for API access
- Authorization decisions (allow/deny) on resource operations

## Authentication Model

SoloDev uses Bearer token authentication for all API and MCP requests:

```bash
# Create a PAT
./gitness login
./gitness user pat "my-pat" 2592000

# Use in API requests
curl http://localhost:3000/api/v1/user \
  -H "Authorization: Bearer $TOKEN"
```

For MCP stdio transport, set `SOLODEV_TOKEN` in the environment.

## Access Control

SoloDev modules use the following permissions:

| Module | View | Create | Edit | Delete |
|--------|------|--------|------|--------|
| Repositories | PermissionRepoView | PermissionRepoPush | PermissionRepoPush | PermissionRepoDelete |
| Security Scans | PermissionRepoView | PermissionRepoPush | PermissionSpaceEdit | — |
| Quality Gates | PermissionQualityGateView | PermissionQualityGateCreate | PermissionQualityGateEdit | PermissionQualityGateDelete |
| Health Checks | PermissionRepoView | PermissionRepoView | PermissionRepoView | PermissionRepoView |

## Key Paths

| Purpose | Path |
|---------|------|
| Authentication | `app/auth/` |
| Crypto | `crypto/` |
| Secret management | `secret/` |

## Status

**Implemented** — Inherited from Gitness. Local authentication, PAT-based API access, and role-based authorization.

## Future Work

- Enterprise SSO integration (planned for Phase 4)
- Audit logging for security-relevant operations
