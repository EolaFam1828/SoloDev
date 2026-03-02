# Why SoloDev

SoloDev is an open-source DevOps platform built on [Gitness by Harness](https://github.com/harness/gitness), extended with an AI intelligence layer for solo developers. This page explains the problem space that motivated the project.

## The solo developer's DevOps problem

A solo developer shipping a SaaS product needs the same operational capabilities as a team — CI/CD, error tracking, security scanning, uptime monitoring, quality enforcement — but has none of the headcount to configure, integrate, and maintain them. The typical result is a fragmented stack: GitHub for code, GitHub Actions for CI, Sentry for errors, UptimeRobot for monitoring, Snyk for security, and manual processes to connect insights across tools. Each service has its own billing, configuration, and alert fatigue. None of them can act on what they find.

## What existing tools assume

GitHub, GitLab, and enterprise DevOps platforms are built around team workflows. Branch protection requires reviewers. Merge policies require approvers. CI/CD pipelines assume a platform engineer wrote the YAML and a separate team triages failures. These are reasonable assumptions for organizations with dedicated DevOps staff, but they create overhead without value for a developer who is the only reviewer, the only approver, and the only person who will fix what breaks. The platform should adapt to the operator, not the other way around.

## What SoloDev does differently

SoloDev consolidates code hosting, pipelines, error tracking, security scanning, quality enforcement, health monitoring, and AI-driven remediation into a single deployable binary. When an error is reported or a pipeline fails, the Error Bridge automatically creates a remediation task. The AI Worker sends the error context to an LLM and produces a unified diff patch. The MCP server exposes the entire platform to AI agents via the Model Context Protocol. Instead of switching between six tools and manually connecting insights, a solo developer gets a single system where failures flow directly to fixes.

The project is early. The core modules are working and the AI remediation loop produces real patches. Vector-based code retrieval, auto-PR creation, and self-healing pipelines are planned. See [Roadmap](Roadmap) for current maturity and [Differentiation](Differentiation) for a detailed comparison with existing platforms.
