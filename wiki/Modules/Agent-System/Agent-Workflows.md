# Agent Workflows

## Purpose

Agent Workflows describe autonomous or semi-autonomous development workflows where AI agents interact with SoloDev to perform tasks that would otherwise require manual developer intervention.

## Inputs

- Agent decisions based on platform state (MCP resources)
- Tool call sequences defined by the agent's reasoning
- Developer approval for gated actions

## Processing

### Current Workflows

These workflows are possible today using the MCP server:

#### Fix-This Workflow

An agent uses the `fix_this` compound tool to:
1. Read the error context from the Error Tracker
2. Analyze the root cause
3. Trigger a remediation task
4. Monitor the remediation status
5. Report the result back

#### PR-Ready Check

An agent uses the `pr_ready` compound tool to:
1. Run a security scan
2. Evaluate quality gates
3. Check tech debt status
4. Return a consolidated readiness verdict

#### Repository Onboarding

An agent uses the `onboard_repo` compound tool to:
1. Trigger a security scan
2. Generate a pipeline from detected stack
3. Create quality gate rules
4. Set up health checks

#### Incident Triage

An agent uses the `incident_triage` compound tool to:
1. Correlate errors, health checks, and security findings
2. Assess incident severity
3. Recommend remediation actions

### Planned Workflows

#### Continuous Monitoring Agent

An agent runs continuously, watching MCP resources:
1. Polls `solodev://errors/active` and, when available, `solodev://health/status`
2. When a new error or health degradation is detected, triggers `fix_this`
3. Monitors the remediation through completion
4. Reports results via notification

#### Self-Healing Pipeline Agent

1. Detects pipeline failure
2. Triggers remediation
3. Waits for patch generation
4. Applies the patch
5. Re-triggers the pipeline
6. Loops until success or maximum retries

## Outputs

- Completed remediation tasks
- Applied code fixes
- Consolidated status reports
- PR readiness verdicts

## Status

**Implemented** — Compound tools enable the described workflows. Agents can execute them today via MCP. Continuous monitoring and self-healing agents are planned.

## Future Work

- Multi-agent coordination (multiple agents working on different aspects)
- Agent memory (persist learning across sessions)
- Workflow templates (pre-defined agent behaviors)
- Agent performance metrics and success rate tracking
