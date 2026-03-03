# Insights

## Purpose

Insights provide analytics and feedback surfaced to developers about their project's health, quality, security posture, and remediation effectiveness. Insights aggregate data from across the platform to give a solo developer a single view of what matters.

## Inputs

- Remediation task statistics (pending, completed, applied, failed)
- Error group counts and trends
- Security scan finding summaries
- Quality gate evaluation results
- Health check uptime percentages
- Tech debt item counts by severity

## Processing

### Current Implementation

Insights are currently provided through per-module summary endpoints:

| Endpoint | Data |
|----------|------|
| `GET /remediations/summary` | Total, pending, processing, completed, applied, failed, dismissed |
| `GET /errors/summary` | Total errors, counts by status and severity |
| `GET /security-scans/{id}/summary` | Finding counts by severity |
| `GET /quality-gates/summary` | Repository pass/fail rates, average coverage |
| `GET /health-checks/{id}/summary` | Uptime percentage, response times |

### SoloDev Dashboard

The web dashboard (`web/src/pages/SoloDevDashboard/`) aggregates these summaries into a single-page view with cards for each module. This is the primary insights surface today.

### MCP Resources

AI agents can access live platform state through MCP resources:
- `solodev://errors/active` — Active errors
- `solodev://remediations/pending` — Pending remediations
- `solodev://quality/summary` — Quality summary
- `solodev://security/open-findings` — Open security findings
- `solodev://health/status` — Health check statuses (cataloged as coming-soon; may be unavailable depending runtime wiring)
- `solodev://tech-debt/hotspots` — Top debt hotspots

## Outputs

- Dashboard summary cards with module-level metrics
- Per-module summary statistics via REST API
- Real-time platform state via MCP resources

## Status

**Implemented** — Per-module summary endpoints and the consolidated web dashboard are working. MCP resources provide real-time insights to agents.

## Future Work

- Cross-module trend analysis (error rate over time, remediation success rate)
- Personalized daily digest for the solo developer
- Predictive insights (e.g., "this area of code is likely to produce errors based on complexity trends")
- Remediation effectiveness metrics (what percentage of AI patches were applied)
