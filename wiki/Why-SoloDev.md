# Why SoloDev

SoloDev is an open-source DevOps platform built on [Gitness by Harness](https://github.com/harness/gitness), extended with an AI intelligence layer for solo developers. This page describes the technical gaps in current DevOps tooling that motivated the project.

## The Solo Developer's DevOps Problem

A solo developer shipping a SaaS product needs the same operational capabilities as a team — CI/CD, error tracking, security scanning, uptime monitoring, quality enforcement — but has none of the headcount to configure, integrate, and maintain them. The typical result is a fragmented stack: GitHub for code, GitHub Actions for CI, Sentry for errors, UptimeRobot for monitoring, Snyk for security, and manual processes to connect insights across tools. Each service has its own billing, configuration, and alert fatigue. None of them can act on what they find.

## What Existing Tools Assume

GitHub, GitLab, and enterprise DevOps platforms are built around team workflows. Branch protection requires reviewers. Merge policies require approvers. CI/CD pipelines assume a platform engineer wrote the YAML and a separate team triages failures. These are reasonable assumptions for organizations with dedicated DevOps staff, but they create overhead without value for a developer who is the only reviewer, the only approver, and the only person who will fix what breaks. The platform should adapt to the operator, not the other way around.

## Failures Stay Disconnected from Fixes

When a pipeline fails on GitHub Actions or GitLab CI, the developer reads logs, context-switches to the codebase, diagnoses the problem, writes a fix, and pushes. The platform provides no bridge between "something broke" and "here is a patch." There is no closed-loop system where failure detection feeds directly into automated remediation. The developer is the glue.

## Observability and Code Live in Separate Tools

Most solo developers cobble together GitHub for code, Sentry or LogRocket for errors, Datadog or UptimeRobot for monitoring, and Snyk for security scanning. Each tool has its own account, billing, alert configuration, and context. None of them can take action on the code. A security finding in Snyk cannot trigger a code fix in GitHub. An error in Sentry cannot create a remediation task in the pipeline.

## AI Tooling Targets Code Generation, Not Operational Feedback Loops

GitHub Copilot, Cursor, and similar tools focus on code generation at the editor level. They suggest completions and refactors based on what the developer is currently typing. None of them watch runtime behavior, ingest error telemetry, or generate patches in response to production failures. The AI operates at the editor level, not the platform level.

## No Platform Treats AI Agents as First-Class Operators

Existing DevOps platforms expose REST APIs designed for human-driven dashboards and CLI tools. None provide native support for the Model Context Protocol (MCP), which is emerging as the standard for AI agent tool access. There is no way for an AI agent to monitor, diagnose, patch, and deploy within a single platform interaction.

## What SoloDev Does Differently

SoloDev consolidates code hosting, pipelines, error tracking, security scanning, quality enforcement, health monitoring, and AI-driven remediation into a single deployable binary. When an error is reported or a pipeline fails, the Error Bridge automatically creates a remediation task. The AI Worker sends the error context to an LLM and produces a unified diff patch. The MCP server exposes the entire platform to AI agents. Instead of switching between six tools and manually connecting insights, a solo developer gets a single system where failures flow directly to fixes.

SoloDev is not just DevOps tooling. It is an autonomous feedback system for software development. The platform detects problems, analyzes them, proposes solutions, validates confidence, and applies changes. The developer supervises the loop rather than performing every step manually.
