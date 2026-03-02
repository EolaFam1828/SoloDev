# Architecture Decisions

Major technical decisions and reasoning behind SoloDev's architecture.

## AD-001: Fork Gitness Instead of Building from Scratch

**Decision:** Build SoloDev as a fork of Gitness by Harness.

**Reasoning:** Gitness provides a mature, Apache-2.0 licensed DevOps foundation (SCM, pipelines, Gitspaces, registry) that would take years to build from scratch. By forking, SoloDev inherits a working platform and focuses engineering effort on the intelligence layer.

**Trade-offs:**
- (+) Immediate access to Git hosting, CI/CD, and developer environments
- (+) Apache-2.0 license allows derivative works
- (-) Naming heritage requires ongoing decoupling
- (-) Upstream divergence means selective cherry-picking, not automatic merging

## AD-002: Single Binary Architecture

**Decision:** Package all components (API server, AI Worker, MCP server, health checker) into a single binary.

**Reasoning:** Solo developers are the target audience. A single binary deployed via Docker Compose is dramatically simpler than a multi-service architecture. Operational complexity is the enemy of adoption for a one-person team.

**Trade-offs:**
- (+) One command to deploy (`docker compose up -d`)
- (+) No inter-service communication failures
- (-) Vertical scaling only (addressed in planned distributed topology)
- (-) AI Worker shares resources with API server

## AD-003: SQLite Default, PostgreSQL for Production

**Decision:** Use SQLite for local development and PostgreSQL for production.

**Reasoning:** SQLite requires zero configuration and works perfectly for a single-developer instance. PostgreSQL provides durability and concurrent access for production. Dual-backend support lets developers start instantly and upgrade when needed.

## AD-004: In-Process Event System

**Decision:** Use synchronous in-process event dispatch instead of a message queue.

**Reasoning:** For a single-instance deployment, an in-process pub/sub system is simpler, faster, and has no external dependencies. The Error Bridge and Security Remediation service subscribe to events within the same process. An external message queue (Redis, RabbitMQ) adds operational complexity without benefit at the target scale.

**Trade-offs:**
- (+) No external dependencies
- (+) Guaranteed event delivery within the process
- (-) Events are lost if the process crashes mid-dispatch
- (-) Cannot distribute consumers across multiple instances

## AD-005: LLM Provider Abstraction

**Decision:** Support multiple LLM providers behind a common interface.

**Reasoning:** Solo developers have different preferences and constraints. Some have Anthropic API access, others prefer OpenAI, some need local inference via Ollama. The AI Worker should work with any provider without code changes.

## AD-006: MCP as Agent Interface

**Decision:** Implement the Model Context Protocol as the primary agent interface.

**Reasoning:** MCP is emerging as the standard for AI agent tool access. By implementing MCP, SoloDev is immediately compatible with Claude Desktop, Cursor, and any MCP-compliant client. Building a proprietary agent API would limit adoption.

## AD-007: Enforcement Modes Over Branch Protection

**Decision:** Replace team-oriented branch protection with Solo Gate enforcement modes.

**Reasoning:** Branch protection assumes reviewers and approvers exist. For a solo developer, strict/balanced/prototype modes provide equivalent safety guarantees without requiring phantom team members. The enforcement adapts to the development phase, not the team size.
