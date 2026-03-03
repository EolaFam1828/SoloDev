# Context Engine

## Purpose

The Context Engine collects code, logs, and telemetry from the platform and prepares them for AI processing. It assembles the input context that the LLM needs to generate an accurate patch.

## Inputs

- Error log or stack trace from the remediation task
- File path of the affected source code
- Branch and commit SHA from the failure
- Source code stored in the remediation record
- Repository metadata (language, framework)

## Processing

### Current Implementation

The Context Engine currently operates inline within the AI Worker. When a remediation task is picked up:

1. **Error context extraction** — The error log, file path, branch, and commit SHA are read from the remediation record
2. **Source code inclusion** — The source code stored in the `rem_source_code` field is included directly in the prompt
3. **Metadata assembly** — Language detection from file extension, repository identifiers, and trigger source are added as context

### Current Prototype: Vector-Assisted Retrieval

The current implementation still prioritizes source code stored directly in the remediation record, and now also includes a prototype embedding-based search path:

1. Repository files are chunked and embedded into a vector store
2. The error context is used as a query against the vector store
3. The most relevant code chunks are retrieved and included in the prompt
4. Retrieved chunks are filtered by similarity and added as additional context fragments

## Outputs

- Assembled prompt context (error log + source code + metadata)
- Formatted input for the Prompt Template system

## Key Paths

| Purpose | Path |
|---------|------|
| AI Worker (contains current context assembly) | `app/services/aiworker/worker.go` |

## Status

**In Progress** — Context assembly is implemented inline within the AI Worker and supports prototype vector-assisted retrieval. Ongoing work is focused on robustness and broader repository coverage.

## Future Work

- Embedding-based vector retrieval over full repository content
- Multi-file context inclusion for cross-file bugs
- Dependency graph analysis for import chain context
- Historical remediation context (similar past fixes)
