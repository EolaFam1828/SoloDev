# Design Principles

SoloDev follows a set of architectural philosophies that shape every component and decision. These principles distinguish SoloDev from general-purpose DevOps platforms.

## Solo-First UX

Every feature is designed for a single operator. There are no phantom team members required. Branch protection does not need reviewers. Approval gates do not need approvers. Enforcement modes (strict, balanced, prototype) replace multi-person governance with single-developer controls. If a feature requires organizational overhead that a solo developer cannot provide, it is redesigned or removed.

## Self-Hostability

SoloDev deploys as a single binary with Docker Compose. There is no mandatory cloud dependency, no required SaaS subscription, and no phone-home telemetry. The platform runs on a developer's laptop, a VPS, or a home server. External LLM providers are optional — local inference via Ollama is supported.

## Automation Over Configuration

The platform should do things, not ask the developer to configure things. Auto-Pipeline detects the tech stack and generates CI/CD YAML. The Error Bridge creates remediation tasks without manual intervention. The AI Worker polls and processes tasks in the background. Quality gates evaluate automatically on pipeline completion. The default state is "the platform handles it."

## Minimal Cognitive Overhead

A solo developer has limited attention. SoloDev consolidates code hosting, pipelines, error tracking, security scanning, quality enforcement, health monitoring, and remediation into one place. Context switching between six separate tools is replaced with a single dashboard. Alerts are actionable — they include proposed fixes, not just notifications.

## Closed-Loop Feedback

SoloDev is an autonomous feedback system, not a monitoring dashboard. Detection flows into analysis, analysis flows into proposals, proposals flow into validated patches, and patches flow back into the codebase. The platform closes the loop between "something broke" and "here is the fix." See [The SoloDev Loop](The-SoloDev-Loop) for the full cycle.

## AI as Platform Infrastructure

AI is not a feature bolted onto the side. It is embedded in the platform's operational model. The AI Worker is a background service that processes remediation tasks. The MCP server exposes the full platform to AI agents. Prompt templates encode domain expertise. The platform assumes AI agents will operate alongside — and eventually instead of — manual human actions.

## Mechanical, Not Marketing

Documentation describes what the system does and how it works. Every module page includes Purpose, Inputs, Processing, Outputs, Status, and Future Work. Capabilities are reported honestly: "Working" means the code compiles and runs. "Planned" means the feature does not exist yet. There are no implied promises.
