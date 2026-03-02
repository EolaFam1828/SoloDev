# Signal System Overview

## Purpose

The Signal System is responsible for ingesting, storing, and normalizing runtime and pipeline signals from all platform sources. It is the "Detect" stage of the SoloDev Loop — every problem that SoloDev can fix must first be captured as a signal.

## What Is a Signal

A signal is any platform event indicating something requires attention:

| Signal Type | Source | Example |
|-------------|--------|---------|
| Runtime error | Error Tracker API | `panic: nil pointer dereference in handler.go:42` |
| Pipeline failure | Pipeline Runner | Build step exited with code 1 |
| Security finding | Security Scanner | SQL injection in `db/query.go` |
| Quality violation | Quality Gates | Coverage dropped below 80% |
| Health check failure | Health Monitor | API endpoint returned 503 |

## Components

| Component | Description | Page |
|-----------|-------------|------|
| Error Tracker | Ingests runtime errors, groups by fingerprint, tracks occurrences | [Error Tracker](Error-Tracker) |
| Error Bridge | Normalizes errors into structured data for AI modules | [Error Bridge](Error-Bridge) |
| Health Monitor | Tracks HTTP endpoint uptime and response times | [Health Monitor](Health-Monitor) |
| Telemetry | Logs, metrics, and traces used for analysis | [Telemetry](Telemetry) |

## Signal Flow

```
External Event ──▶ Signal Module ──▶ Database ──▶ Bridge Service ──▶ Remediation Task
                                                       │
                                          Error Bridge / Sec Remediation
```

All signals follow the same pattern:
1. An external event arrives via REST API or internal hook
2. The appropriate module stores the structured signal
3. A bridge service evaluates whether the signal should trigger remediation
4. If triggered, a pending remediation task is created for the AI Layer

## Status

| Component | Status |
|-----------|--------|
| Error Tracker | Implemented |
| Error Bridge | Implemented |
| Health Monitor | Implemented |
| Telemetry (structured logging) | Concept |
| Signal Correlator | Planned |

## Future Work

- Signal Correlator: cross-domain correlation between errors, pipeline failures, health degradation, and security findings
- Anomaly detection from health check trends
- Signal deduplication across multiple sources
