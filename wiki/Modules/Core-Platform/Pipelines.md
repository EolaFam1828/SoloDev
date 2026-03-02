# Pipelines

## Purpose

Provides CI/CD pipeline execution using the Drone-based engine inherited from Gitness. Pipelines build, test, and deploy code. In SoloDev, pipeline failures are a primary signal source that feeds the remediation loop.

## Inputs

- Pipeline YAML definitions (`.harness/pipeline.yaml` in repositories)
- Git push events triggering pipeline execution
- Auto-Pipeline generated YAML (from stack detection)
- Manual trigger requests

## Processing

- Parses pipeline YAML definitions
- Schedules and executes pipeline stages in Docker containers
- Captures build logs, test output, and exit codes
- Tracks execution status (pending, running, success, failure, killed)
- Publishes completion events

## Outputs

- Pipeline execution results (success/failure, logs, duration)
- Pipeline failure events consumed by the Error Bridge
- Build artifacts stored in the Registry
- Status badges and webhooks

## Integration with SoloDev

Pipelines are the primary trigger source for the SoloDev Loop:

1. **Auto-Pipeline** — The Auto-Pipeline module generates pipeline YAML from detected tech stacks, eliminating manual configuration
2. **Error Bridge** — When a pipeline fails, the Error Bridge creates a remediation task with the build log, branch, and commit SHA
3. **Quality Gates** — Pipeline completion can trigger quality gate evaluation
4. **Self-healing loop** (planned) — After a remediation patch is applied, the failed pipeline is re-triggered to verify the fix

### Pipeline YAML Example

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

## Key Paths

| Purpose | Path |
|---------|------|
| Pipeline engine | `app/pipeline/` |
| Auto-Pipeline | `app/pipeline/autopipeline/` |
| Pipeline types | `types/pipeline.go` |

## Status

**Implemented** — Drone-based CI/CD inherited from Gitness. Auto-Pipeline stack detection and YAML generation added by SoloDev.

## Future Work

- Self-healing pipeline loop (detect failure → generate fix → re-run)
- Pipeline failure → remediation → re-run flow
- Pipeline metrics and success rate tracking
