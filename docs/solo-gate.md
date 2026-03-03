# Solopreneur Gate Configuration Module

## Overview

The Solo Gate module replaces enterprise-grade multi-team governance with fast, solopreneur-friendly enforcement modes for quality and security gates. Instead of rigid block-or-pass rules, solo developers can choose how strictly they want pipelines enforced depending on their current development phase.

## Enforcement Modes

| Mode | Behavior | Best For |
|------|----------|----------|
| **Strict** | Block on any finding (critical, high, medium, or low) | Release branches, production deploy |
| **Balanced** | Block critical+high, warn medium, auto-ignore low | Feature development |
| **Prototype** | Never block — log everything as tech debt | Hackathons, rapid prototyping, spikes |

## Architecture

### Types (`types/solo_gate.go`)

#### SoloGateConfig
Per-space configuration:
- `EnforcementMode` — strict / balanced / prototype
- `AutoIgnoreLow` — auto-dismiss low severity findings
- `AutoTriageKnown` — auto-dismiss repeat false positives
- `AIAutoFix` — trigger AI remediation on critical/high findings
- `LogTechDebt` — log skipped issues as tech debt items

#### UpdateSoloGateConfigInput
Optional-field update struct for PATCH operations.

### Gate Engine (`app/services/sologate/engine.go`)

#### `Evaluate(config, findings) → EvaluateResult`

Takes a list of findings and the space's gate config, returns:
- `Action` — "block", "warn", or "pass"
- `Reasons` — human-readable explanation
- `AutoRemediate` — true if AI should auto-fix
- `TechDebtLogged` — true if issues were logged as debt

**Decision Matrix:**

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

## Integration Points

### With AI Remediation
When `AIAutoFix` is enabled and critical/high findings are detected:
1. Gate engine sets `AutoRemediate = true` in result
2. Caller triggers AI remediation via the remediation controller
3. AI generates fix → remediation reaches `completed`
4. Draft PR delivery can be triggered manually, or automatically when remediation delivery is enabled

### With Tech Debt Tracker
When `LogTechDebt` is enabled and passing despite findings:
1. Each skipped finding is logged as a tech debt item
2. Debt is tracked per-space for later review

## File Structure

```
├── types/
│   └── solo_gate.go                    # SoloGateConfig, EnforcementMode
└── app/services/sologate/
    └── engine.go                       # Gate evaluation engine
```

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
// result.Action = "block"
// result.AutoRemediate = true
// result.Reasons = ["1 critical, 0 high findings (blocking)"]
```
