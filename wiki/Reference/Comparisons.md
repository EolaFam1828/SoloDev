# Comparisons

How SoloDev compares with other platforms.

## Feature Comparison

| Feature | GitHub | GitLab | Harness OSS (Gitness) | SoloDev |
|---------|--------|--------|----------------------|---------|
| Git hosting + code review | Yes | Yes | Yes | Yes (inherited) |
| CI/CD pipelines | Actions | Built-in | Drone-based | Drone-based (inherited) |
| AI failure remediation | Copilot (code suggestions, not failure-driven) | No | No | Implemented — LLM-driven patch generation from error logs |
| Error tracking with auto-fix bridge | No (requires Sentry/Datadog) | Partial (error tracking exists) | No | Implemented — errors auto-create remediation tasks |
| Zero-config pipeline generation | No | Auto DevOps (heavyweight) | No | Implemented — stack detection + YAML generation |
| Enforcement gates for solo devs | Branch protection (team-oriented) | Approval rules (team-oriented) | No | Implemented — strict/balanced/prototype modes |
| MCP-native agent access | No | No | No | Implemented — 16 atomic tools, 5 compound tools |
| Security scan → auto-remediation | Dependabot (updates only) | SAST/DAST (no auto-fix) | No | Implemented — findings trigger remediation tasks |
| Lightweight local deploy | Codespaces (cloud) | Heavy self-hosted | Yes | Yes (inherited) |
| Self-healing pipelines | No | No | No | Planned |
| Multi-agent orchestration | No | No | No | Planned |

## Key Differences

### vs. GitHub

GitHub provides excellent code hosting and CI/CD (Actions) but treats failures as notifications. The developer is responsible for diagnosing, fixing, and re-running. SoloDev automates the diagnosis-to-fix path. GitHub Copilot assists with code generation at the editor level; SoloDev's AI operates at the platform level, responding to operational failures.

### vs. GitLab

GitLab is the closest in scope (integrated DevOps platform) but is designed for teams. Branch protection, approval rules, and governance features assume multiple users. GitLab has error tracking but no automated remediation path. SoloDev adapts enforcement to solo developers and closes the detection-to-fix loop.

### vs. Harness OSS (Gitness)

SoloDev is a direct fork of Gitness. Gitness provides the core platform (SCM, pipelines, Gitspaces, registry) but has no AI layer, no error tracking, no security scanning, no quality gates, and no MCP integration. SoloDev adds all of these.

## Status Key

- **Implemented** — Code exists, compiles, and exposes API endpoints
- **Planned** — On the roadmap but not yet built
- **Inherited** — Capability comes from the upstream Gitness fork
