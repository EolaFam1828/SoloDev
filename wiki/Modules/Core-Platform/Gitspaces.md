# Gitspaces

## Purpose

Provides cloud-based developer environments that can be launched from a repository. Gitspaces allow a developer to spin up a pre-configured development environment without local setup. Inherited from Gitness.

## Inputs

- Repository reference (which repo to create the environment from)
- Gitspace configuration (IDE type, resource allocation)
- Docker socket access (required for container management)

## Processing

- Provisions a containerized development environment
- Clones the target repository into the environment
- Configures IDE access (VS Code, JetBrains, etc.)
- Manages environment lifecycle (create, start, stop, delete)

## Outputs

- Running development environment accessible via browser or IDE
- Environment status and resource usage metrics
- Git operations from within the environment flow back to the SCM

## Integration with SoloDev

Gitspaces play a role in the feedback loop by providing the environment where:
- Developers can review and test AI-generated patches
- MCP-connected agents can operate within a sandboxed environment
- Remediation patches can be validated before merging

## Key Paths

| Purpose | Path |
|---------|------|
| Gitspace controller | `app/api/controller/gitspace/` |
| Infrastructure provider | `infraprovider/` |

## Status

**Implemented** — Inherited from Gitness. Standard developer environment provisioning.

## Future Work

- Agent-operated Gitspaces for autonomous testing of remediation patches
- Ephemeral environments for patch validation
