# Solo Gate

## Purpose

The Solo Gate module replaces enterprise-grade multi-team governance with fast, solopreneur-friendly enforcement modes for quality and security gates. Instead of rigid block-or-pass rules, solo developers choose how strictly pipelines are enforced depending on their current development phase.

## Inputs

- List of findings (security, quality, or combined)
- `SoloGateConfig` for the space (enforcement mode, auto-fix settings, tech debt logging)

## Processing

### Enforcement Modes

| Mode | Behavior | Best For |
|------|----------|---------|
| **Strict** | Block on any finding (critical, high, medium, or low) | Release branches, production deploy |
| **Balanced** | Block critical+high, warn medium, auto-ignore low | Feature development |
| **Prototype** | Never block — log everything as tech debt | Hackathons, rapid prototyping, spikes |

### Decision Matrix

```
                  ┌─────────┬──────────┬───────────┐
                  │ Strict  │ Balanced │ Prototype │
     ┌────────────┼─────────┼──────────┼───────────┤
     │ Critical   │ BLOCK   │ BLOCK    │ WARN      │
     │ High       │ BLOCK   │ BLOCK    │ WARN      │
     │ Medium     │ BLOCK   │ WARN     │ WARN      │
     │ Low        │ BLOCK   │ PASS     │ PASS      │
     └────────────┴─────────┴──────────┴───────────┘
```

### Gate Evaluation

`Evaluate(config, findings) → EvaluateResult`

The gate engine takes a list of findings and the space's gate config, and returns:

| Field | Type | Description |
|-------|------|-------------|
| `Action` | string | `"block"`, `"warn"`, or `"pass"` |
| `Reasons` | []string | Human-readable explanation |
| `AutoRemediate` | bool | True if AI should auto-fix |
| `TechDebtLogged` | bool | True if issues were logged as debt |

### AI Remediation Integration

When `AIAutoFix` is enabled and critical/high findings are detected:
1. Gate engine sets `AutoRemediate = true` in the result
2. Caller triggers AI remediation via the remediation controller
3. AI generates fix

### Tech Debt Integration

When `LogTechDebt` is enabled and the gate passes despite findings:
1. Each skipped finding is logged as a tech debt item
2. Debt is tracked per-space for later review

## Outputs

- Gate evaluation result (block/warn/pass) with reasons
- Auto-remediation trigger flag
- Tech debt items created for skipped findings

## Configuration

| Field | Type | Description |
|-------|------|-------------|
| `EnforcementMode` | enum | `strict`, `balanced`, or `prototype` |
| `AutoIgnoreLow` | bool | Auto-dismiss low severity findings |
| `AutoTriageKnown` | bool | Auto-dismiss repeat false positives |
| `AIAutoFix` | bool | Trigger AI remediation on critical/high findings |
| `LogTechDebt` | bool | Log skipped issues as tech debt items |

## Key Paths

| Purpose | Path |
|---------|------|
| Types | `types/solo_gate.go` |
| Gate evaluation engine | `app/services/sologate/engine.go` |

## Status

**Implemented** — Strict/balanced/prototype enforcement modes, finding evaluation, AI auto-fix trigger, and tech debt logging are all working.

## Future Work

- Per-branch enforcement mode overrides
- Custom severity thresholds per mode
- Gate evaluation history and trend analysis
