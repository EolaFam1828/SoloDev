# Prompt Templates

## Purpose

Prompt Templates define the structured prompts used for AI analysis and remediation tasks. They encode SoloDev's domain expertise into reusable formats that produce consistent, parseable output from LLMs.

## Inputs

- Error context (log, stack trace, file path, branch, commit)
- Source code from the affected file
- Task metadata (trigger source, severity, language)

## Processing

### Remediation Prompt

The AI Worker constructs a system prompt that instructs the LLM to:

1. Analyze the error context and source code
2. Identify the root cause of the failure
3. Produce a unified diff (`patch -p1` compatible) that fixes the issue
4. Provide a confidence score between 0.0 and 1.0

The system prompt enforces output format:
- The diff must be enclosed in a code block
- The confidence score must be a decimal number
- The response must include an analysis section explaining the fix

### MCP Prompt Templates

The MCP server includes 5 pre-built prompt templates that encode domain expertise for common agent workflows:

| Prompt Name | Purpose |
|-------------|---------|
| `solodev_review` | Code review incorporating security, quality, and tech debt context |
| `solodev_incident` | Incident investigation correlating errors, health, and security |
| `solodev_pipeline_debug` | Pipeline debugging with generation context and validation |
| `solodev_security_audit` | Security audit across all scan findings for a repository |
| `solodev_debt_sprint` | Tech debt sprint planning with prioritized remediation items |

## Outputs

- Formatted prompt strings ready for LLM submission
- Template parameters documented for each prompt type

## Key Paths

| Purpose | Path |
|---------|------|
| AI Worker prompt construction | `app/services/aiworker/worker.go` |
| MCP prompt templates | `mcp/prompts.go` |

## Status

**Implemented** — Remediation prompts are constructed in the AI Worker. MCP prompt templates are served via the `prompts/list` and `prompts/get` protocol methods.

## Future Work

- Prompt versioning and A/B testing
- Language-specific prompt variants (Go, Python, TypeScript, Rust)
- User-customizable prompt templates
- Prompt effectiveness tracking (which prompts produce applied patches)
