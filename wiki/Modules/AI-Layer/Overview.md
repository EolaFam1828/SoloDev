# AI Layer Overview

## Purpose

The AI Layer is SoloDev's intelligence infrastructure. It receives structured signals from the platform (errors, failures, findings), processes them through LLM-powered analysis, and produces concrete code fixes. The AI Layer is what transforms SoloDev from a DevOps tool into an autonomous feedback system.

## Relationship to Platform Signals

The AI Layer sits between the Signal System (which detects problems) and the Remediation System (which applies fixes):

```
Signal System ──▶ AI Layer ──▶ Remediation System
                     │
         ┌───────────┼───────────┐
         │           │           │
    Context      LLM Call    Prompt
    Engine       (Worker)    Templates
```

Every signal that enters the AI Layer has already been:
1. Detected by a platform module (Error Tracker, Security Scanner, Pipeline Runner)
2. Structured by a bridge service (Error Bridge, Security Remediation)
3. Stored as a pending remediation task in the database

The AI Layer's job is to turn that structured context into a working code patch.

## Components

| Component | Description | Page |
|-----------|-------------|------|
| Context Engine | Collects code, logs, and telemetry and prepares them for AI processing | [Context Engine](Context-Engine) |
| LLM Adapters | Defines how different LLM providers are integrated and swapped | [LLM Adapters](LLM-Adapters) |
| Prompt Templates | Structured prompts for analysis and remediation tasks | [Prompt Templates](Prompt-Templates) |
| Validation Engine | Verifies AI outputs before they are applied | [Validation-Engine](Validation-Engine) |

## Data Flow

```
Pending Remediation → Context Engine → Prompt Builder → LLM Adapter → Response Parser → Validated Patch
```

1. The AI Worker polls the database for pending remediation tasks
2. The Context Engine assembles error log, source code, file path, branch, and commit
3. The Prompt Builder constructs a structured prompt from a template
4. The LLM Adapter sends the prompt to the configured provider
5. The Response Parser extracts the unified diff and confidence score
6. The Validation Engine checks the output before storage

## Status

| Component | Status |
|-----------|--------|
| AI Worker (poll + process) | Implemented |
| LLM Adapters (Anthropic, OpenAI, Gemini, Ollama) | Implemented |
| Prompt construction from error context | Implemented |
| Response parsing (diff + confidence) | Implemented |
| Context Engine (vector-based retrieval) | Planned |
| Validation Engine (automated patch verification) | Concept |

## Future Work

- Vector-based code retrieval (embedding search over full repository)
- Multi-model consensus (run multiple LLMs and compare outputs)
- Feedback learning from applied vs. rejected patches
