# Control Plane

The control plane is the set of orchestration components that coordinate pipelines, signals, and AI modules within SoloDev. It determines how events flow between detection and remediation.

## Components

### Event Router

When a signal is produced (error reported, pipeline failed, scan completed), the control plane routes it to the appropriate handler:

- **Error Tracker events** → Error Bridge → AI Worker
- **Pipeline failure events** → Error Bridge → AI Worker
- **Security scan events** → Security Remediation service → AI Worker
- **Quality gate events** → Solo Gate Engine → Tech Debt Tracker (or Error Bridge)
- **Health check events** → Stored for trend analysis (remediation path planned)

### Background Services

| Service | Interval | Function |
|---------|----------|----------|
| AI Worker Poller | 15 seconds | Queries for `pending` remediation tasks and schedules worker jobs |
| Health Check Runner | Per-monitor interval (60–86400s) | Executes configured HTTP checks and stores results |

### Dependency Injection

SoloDev uses Wire-based dependency injection to compose the control plane at startup. Controllers, stores, services, and bridges are wired together in `cmd/gitness/wire.go`. The Error Bridge is attached to the Error Tracker controller. The Security Remediation service is attached to the Security Scanner controller. The Solo Gate engine receives both the remediation store and the Error Bridge.

## Coordination Flow

```
Pipeline Runner ──────────────────┐
                                  ▼
Error Tracker ──▶ Error Bridge ──▶ Remediation Store ──▶ AI Worker Poller ──▶ AI Worker ──▶ LLM
                                  ▲
Security Scanner ──▶ Sec Remediation ─┘
```

The control plane does not implement a message queue or distributed event bus. Events are dispatched synchronously within the Go process and stored in the database. The AI Worker Poller reads from the database on a fixed interval. This design is intentional: SoloDev targets single-instance deployment where in-process coordination is simpler and more reliable than distributed messaging.

## Status

| Aspect | Status |
|--------|--------|
| Error → Bridge → AI Worker flow | Implemented |
| Security → Remediation → AI Worker flow | Implemented |
| Quality Gate → Solo Gate → Remediation flow | Implemented |
| Health Check → Remediation flow | Concept |
| Distributed event bus | Not planned for current scope |

## Future Work

- Health check failure → remediation pipeline
- Signal Correlator: cross-domain event correlation (errors + health + security)
- Event replay for debugging remediation flows
