# FAQ

## General

**What is SoloDev?**
An open-source, AI-native DevOps platform for solo builders. It combines code hosting, pipelines, error tracking, security scanning, and AI-driven remediation into a single deployable system.

**How is SoloDev different from GitHub or GitLab?**
SoloDev closes the loop between failure detection and code fixes. When something breaks, the platform generates a patch — not just an alert. It also provides enforcement modes designed for solo developers, not teams. See [Why SoloDev](../Why-SoloDev) and [Comparisons](Comparisons) for details.

**Is SoloDev production-ready?**
The core modules are implemented and functional. The AI remediation loop generates real patches. However, some planned features (auto-PR, auto-merge, self-healing pipelines) are not yet implemented. See [Current State](../Roadmap/Current-State) for an honest assessment.

**What license is SoloDev under?**
Apache License 2.0, same as the upstream Gitness project.

## Technical

**What LLM providers are supported?**
Anthropic, OpenAI, Google Gemini, and Ollama (local). The provider is configured at the server level.

**Can I run SoloDev without an LLM provider?**
Yes. All modules except AI patch generation work without an LLM. Remediation tasks will be created but not processed.

**What database does SoloDev use?**
SQLite for local development, PostgreSQL for production.

**Why does the binary say `gitness`?**
SoloDev is built from a fork of Gitness by Harness. Some internal names still reflect this heritage. Renaming is in progress.

**How do AI agents connect to SoloDev?**
Via the Model Context Protocol (MCP). SoloDev supports stdio transport (for local agents like Claude Desktop) and HTTP transport (for remote agents). See [MCP Server](../Modules/Agent-System/MCP-Server).

## Operations

**How do I set up the Error Bridge?**
The Error Bridge is automatically configured at startup when enabled. Errors with severity ≥ `error` automatically trigger remediation tasks.

**How do I configure enforcement modes?**
Set the `SoloGateConfig` for your space with the desired enforcement mode (strict, balanced, or prototype).

**What happens when a pipeline fails?**
If the Error Bridge is enabled, it creates a pending remediation task with the build log. The AI Worker picks it up, generates a patch, and stores it for review.

**Can I use SoloDev with an existing GitHub repository?**
SoloDev hosts its own Git repositories. You can push an existing project to SoloDev's Git hosting, or mirror from GitHub.
