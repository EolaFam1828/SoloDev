# Telemetry

## Purpose

Telemetry covers the logs, metrics, and traces produced by the platform and monitored applications. These signals provide the raw data that other modules consume for analysis and remediation.

## Inputs

- Application logs from pipeline execution
- Build output and test results
- Health check response data
- Error occurrence metadata (environment, runtime, OS, architecture)
- Pipeline execution timing and resource usage

## Processing

### Current Implementation

Telemetry data is currently captured within each module's own data model:

- **Pipeline logs** — Stored as part of pipeline execution records (inherited from Gitness)
- **Error metadata** — Environment, runtime, OS, and architecture stored per error occurrence
- **Health check metrics** — Response time, status code, and error messages stored per check result
- **Remediation metrics** — Token usage, duration, and model identifier per AI generation

There is no centralized telemetry pipeline or external metrics backend in the current implementation.

### Planned: Structured Telemetry

A centralized telemetry system would aggregate signals across modules:

1. **Logs** — Structured JSON logs from all platform operations
2. **Metrics** — Time-series data for health checks, pipeline durations, error rates, remediation success rates
3. **Traces** — Request-level tracing through the remediation pipeline (signal → bridge → AI → patch)

## Outputs

- Module-specific metrics (response times, error counts, remediation durations)
- Pipeline execution logs (build output, test results)
- Health check result history

## Key Paths

| Purpose | Path |
|---------|------|
| Logging utilities | `logging/` |
| Livelog (pipeline streaming) | `livelog/` |
| Health check results | `app/store/database/healthcheck_result.go` |
| Remediation metrics | `types/ai_remediation.go` (tokens_used, duration_ms) |

## Status

**Concept** — Telemetry data is captured per-module but there is no centralized telemetry pipeline. Structured logging exists but is not aggregated.

## Future Work

- Centralized metrics aggregation across all modules
- OpenTelemetry-compatible trace export
- Dashboard widgets for platform-wide health trends
- Signal correlation using telemetry data
