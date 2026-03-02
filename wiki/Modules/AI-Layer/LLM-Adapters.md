# LLM Adapters

## Purpose

LLM Adapters define how SoloDev communicates with different AI model providers. The adapter layer abstracts the provider-specific API details so the AI Worker can be configured to use any supported provider without changing the core logic.

## Inputs

- System prompt (instructions for the LLM)
- User prompt (error context, source code, task description)
- Provider configuration (API key, model name, endpoint URL)

## Processing

The AI Worker (`app/services/aiworker/worker.go`) selects the appropriate provider based on configuration and sends the assembled prompt. Each provider adapter handles:

1. **Request formatting** — Converts the prompt into the provider's API format
2. **API communication** — Sends the request and handles response streaming
3. **Error handling** — Retries on transient failures, reports permanent failures
4. **Response extraction** — Returns the raw text response for parsing

## Supported Providers

| Provider | Model Examples | Local/Remote | Notes |
|----------|---------------|-------------|-------|
| **Anthropic** | Claude Sonnet, Claude Opus | Remote | API key required |
| **OpenAI** | GPT-4, GPT-4o | Remote | API key required |
| **Google Gemini** | Gemini Pro, Gemini Flash | Remote | API key required |
| **Ollama** | Llama, CodeLlama, Mistral | Local | No API key, runs on developer's machine |

## Outputs

- Raw LLM response text containing a unified diff and confidence score
- Token usage metrics (stored in `rem_tokens_used`)
- Processing duration (stored in `rem_duration_ms`)
- Model identifier (stored in `rem_ai_model`)

## Configuration

The provider is configured at the server level. The AI Worker reads the configuration to determine which provider and model to use for remediation tasks.

## Swapping Providers

Changing the LLM provider requires updating the server configuration. No code changes are needed. The prompt format and response parsing are provider-agnostic — all providers receive the same prompt and are expected to return a unified diff with a confidence score.

## Key Paths

| Purpose | Path |
|---------|------|
| AI Worker (provider selection) | `app/services/aiworker/worker.go` |
| Response parser | `app/services/aiworker/parser.go` |

## Status

**Implemented** — Four providers are supported: Anthropic, OpenAI, Google Gemini, and Ollama. Provider selection is configuration-driven.

## Future Work

- Provider fallback chain (try primary, fall back to secondary on failure)
- Cost tracking per provider and model
- Multi-model consensus (send to multiple providers, compare outputs)
- Custom provider adapter interface for community extensions
