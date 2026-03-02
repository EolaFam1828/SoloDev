# The SoloDev Loop

The SoloDev Loop is the core feedback cycle that drives the platform. Every component in SoloDev exists to serve one or more stages of this loop.

## The Five Stages

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  DETECT  │────▶│ ANALYZE  │────▶│ PROPOSE  │────▶│ VALIDATE │────▶│  APPLY   │
└──────────┘     └──────────┘     └──────────┘     └──────────┘     └──────────┘
     ▲                                                                    │
     └────────────────────────────────────────────────────────────────────┘
```

### 1. Detect

The platform ingests signals indicating something requires attention. Sources include:

- **Pipeline failures** — build errors, test failures, deployment issues
- **Runtime errors** — application exceptions reported via the Error Tracker API
- **Security findings** — vulnerabilities, secrets, and code smells from the Security Scanner
- **Quality gate violations** — coverage drops, complexity increases, style violations
- **Health check failures** — HTTP endpoint monitoring detecting downtime or degradation

Each signal is structured and stored with full context: stack traces, file paths, commit SHAs, branch information, and severity classification.

### 2. Analyze

The Error Bridge and Signal System normalize detected signals into a common format. The system determines:

- What component failed
- What source code is involved
- What commit introduced the change
- What the severity and category are
- Whether the signal correlates with other recent signals

This stage bridges the gap between "something happened" and "here is the structured context an AI can work with."

### 3. Propose

The AI Worker picks up the analyzed signal and generates a remediation proposal. It:

- Builds a prompt from the error context, source code, and file path
- Sends the prompt to the configured LLM provider (Anthropic, OpenAI, Gemini, or Ollama)
- Receives a unified diff patch and confidence score
- Stores the proposal as a remediation record

The proposal is a concrete code change, not a suggestion or explanation.

### 4. Validate

Before application, the proposal is scored and checked:

- **Confidence scoring** — the LLM provides a 0.0–1.0 confidence value
- **Quality gate evaluation** — the proposed change can be checked against existing rules
- **Solo Gate enforcement** — the enforcement mode (strict/balanced/prototype) determines whether the proposal can proceed automatically

Proposals below the confidence threshold require human review. Proposals above it can proceed to automatic application.

### 5. Apply

The validated change is applied to the codebase:

- A fix branch is created from the patch diff
- A pull request is opened with the AI analysis as the description
- If auto-merge is enabled and confidence is high enough, the PR merges without intervention
- The failed pipeline is re-triggered to verify the fix resolves the original failure

The loop then returns to Detect: the re-triggered pipeline either succeeds (closing the loop) or fails again (starting a new iteration).

## Current Implementation Status

| Stage | Status |
|-------|--------|
| Detect | Implemented — Error Tracker, Security Scanner, Health Monitor, Quality Gates, Pipeline hooks |
| Analyze | Implemented — Error Bridge auto-creates structured remediation tasks |
| Propose | Implemented — AI Worker polls pending tasks, calls LLM, produces diff patches |
| Validate | Partial — Confidence scoring exists; auto-PR and auto-merge are planned |
| Apply | Planned — Manual application via dashboard, API, or MCP agent today |

## Why a Loop

Most DevOps tools are linear: they detect problems and stop. The developer handles the rest. SoloDev is designed as a closed loop because a solo developer cannot afford to be the glue between detection and resolution. The platform must complete the cycle autonomously, with the developer supervising rather than performing each step.
