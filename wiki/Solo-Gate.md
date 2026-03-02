# Solo Gate Engine

## Overview

The Solo Gate module replaces enterprise-grade multi-team governance with fast, solopreneur-friendly enforcement modes for quality and security gates. Instead of rigid block-or-pass rules, solo developers choose how strictly pipelines are enforced depending on their current development phase.

## Enforcement Modes

| Mode | Behavior | Best For |
|------|----------|---------|
| **Strict** | Block on any finding (critical, high, medium, or low) | Release branches, production deploy |
| **Balanced** | Block critical+high, warn medium, auto-ignore low | Feature development |
| **Prototype** | Never block — log everything as tech debt | Hackathons, rapid prototyping, spikes |

## Decision Matrix

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

## Configuration

### `SoloGateConfig` (`types/solo_gate.go`)

Per-space configuration struct:

| Field | Type | Description |
|-------|------|-------------|
| `EnforcementMode` | enum | `strict`, `balanced`, or `prototype` |
| `AutoIgnoreLow` | bool | Auto-dismiss low severity findings |
| `AutoTriageKnown` | bool | Auto-dismiss repeat false positives |
| `AIAutoFix` | bool | Trigger AI remediation on critical/high findings |
| `LogTechDebt` | bool | Log skipped issues as tech debt items |

## Gate Evaluation

### `Evaluate(config, findings) → EvaluateResult`

The gate engine (`app/services/sologate/engine.go`) takes a list of findings and the space's gate config, and returns:

| Field | Type | Description |
|-------|------|-------------|
| `Action` | string | `"block"`, `"warn"`, or `"pass"` |
| `Reasons` | []string | Human-readable explanation |
| `AutoRemediate` | bool | True if AI should auto-fix |
| `TechDebtLogged` | bool | True if issues were logged as debt |

## Integration Points

### AI Remediation

When `AIAutoFix` is enabled and critical/high findings are detected:
1. Gate engine sets `AutoRemediate = true` in the result
2. Caller triggers AI remediation via the remediation controller
3. AI generates fix → PR created automatically

### Tech Debt Tracker

When `LogTechDebt` is enabled and the gate passes despite findings:
1. Each skipped finding is logged as a tech debt item
2. Debt is tracked per-space for later review

## Usage Example

```go
engine := sologate.NewEngine(remediationStore, bridge)

config := &types.SoloGateConfig{
    EnforcementMode: types.EnforcementModeBalanced,
    AutoIgnoreLow:   true,
    AIAutoFix:       true,
}

findings := []sologate.Finding{
    {ID: "CVE-2024-1234", Severity: "critical", Title: "SQL injection"},
    {ID: "LINT-005", Severity: "low", Title: "Unused variable"},
}

result := engine.Evaluate(ctx, config, findings)
// result.Action        = "block"
// result.AutoRemediate = true
// result.Reasons       = ["1 critical, 0 high findings (blocking)"]
```

## File Locations

| Purpose | Path |
|---------|------|
| Types | `types/solo_gate.go` |
| Gate evaluation engine | `app/services/sologate/engine.go` |
