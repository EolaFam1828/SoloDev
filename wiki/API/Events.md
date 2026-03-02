# Events

## Purpose

SoloDev modules emit events when state changes occur. Events enable loose coupling between modules — the Error Bridge subscribes to error events, the Security Remediation service subscribes to scan events, and the control plane coordinates responses.

## Event Sources

### Error Tracker Events

| Event | Trigger |
|-------|---------|
| `ErrorReported` | New error or new occurrence of existing error |
| `ErrorStatusChanged` | Status transition (open → resolved, etc.) |
| `ErrorAssigned` | Error assigned to a user |

### Security Scanner Events

| Event | Trigger |
|-------|---------|
| `ScanTriggered` | Security scan started |
| `ScanCompleted` | Scan finished successfully |
| `ScanFailed` | Scan failed |

### Remediation Events

| Event | Trigger |
|-------|---------|
| `RemediationTriggered` | Remediation task created |
| `RemediationCompleted` | AI generation finished |
| `RemediationApplied` | Fix was pushed/merged |

### Quality Gate Events

| Event | Trigger |
|-------|---------|
| `RuleCreatedEvent` | Quality rule created |
| `RuleUpdatedEvent` | Quality rule updated |
| `RuleDeletedEvent` | Quality rule deleted |
| `RuleEnabledEvent` | Quality rule enabled |
| `RuleDisabledEvent` | Quality rule disabled |
| `EvaluationCreatedEvent` | Quality evaluation completed |

### Health Check Events

| Event | Trigger |
|-------|---------|
| `healthcheck_created` | Monitor created |
| `healthcheck_updated` | Monitor configuration changed |
| `healthcheck_deleted` | Monitor deleted |
| `healthcheck_status_changed` | Status changed (e.g., up → down) |
| `healthcheck_result_created` | Individual check completed |

## Event Flow

Events are published synchronously within the Go process after successful database operations. Subscribers are registered at startup via dependency injection.

```
Controller ──▶ Database ──▶ Event Publisher ──▶ Subscriber
                                                    │
                                    ┌───────────────┼───────────────┐
                                    ▼               ▼               ▼
                              Error Bridge    Sec Remediation   (future: webhooks)
```

## Event Infrastructure

Events use SoloDev's internal pub/sub system (`pubsub/`). There is no external message queue — events are dispatched in-process.

| Purpose | Path |
|---------|------|
| Pub/sub system | `pubsub/` |
| Error Tracker events | `app/events/errortracker/` |
| Security Scanner events | `app/events/securityscan/` |
| Remediation events | `app/events/airemediation/` |
| Quality Gate events | `app/events/qualitygate/` |
| Health Check events | `app/events/healthcheck/` |

## Future Work

- External webhook delivery for events
- Event replay for debugging
- Event-driven notifications (Slack, email, Discord)
