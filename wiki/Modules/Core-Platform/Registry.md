# Registry

## Purpose

Provides OCI-compliant artifact and container image storage. The Registry stores build outputs from pipelines, making them available for deployment. Inherited from Gitness.

## Inputs

- Docker push operations from pipeline steps
- OCI artifact uploads
- Image pull requests from deployment targets

## Processing

- Stores container images and OCI artifacts
- Manages image tags and manifests
- Handles garbage collection of unused layers
- Enforces access control on push and pull operations

## Outputs

- Container images available for deployment
- Artifact metadata accessible via the Registry API
- Image layer data served to Docker clients

## Integration with SoloDev

The Registry is the artifact endpoint of the pipeline:

1. Pipeline builds produce container images
2. Images are pushed to the Registry
3. Deployment steps pull from the Registry
4. When a remediation patch is applied and the pipeline re-runs, new artifacts are produced

## Key Paths

| Purpose | Path |
|---------|------|
| Registry implementation | `registry/` |
| Registry API swagger | `http://localhost:3000/registry/swagger/` |

## Status

**Implemented** — OCI registry inherited from Gitness.

## Future Work

- Artifact scanning integration with the Security Scanner
- Image provenance tracking for remediation-generated builds
