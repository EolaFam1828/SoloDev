# Health Monitor

## Purpose

The Health Monitor provides uptime monitoring for HTTP endpoints. Users configure monitors with check intervals and expected responses; the module tracks results over time and provides uptime statistics and status summaries.

## Inputs

- Monitor configuration: URL, HTTP method, expected status code, interval, timeout, headers, body
- Check execution triggers (on configured interval)

## Processing

- Executes HTTP requests against configured endpoints on the specified interval
- Compares response status code against expected value
- Measures response time in milliseconds
- Tracks consecutive failures for degradation detection
- Stores individual check results with timestamps
- Computes 24-hour uptime percentage

## Outputs

- Monitor status: `up`, `down`, `degraded`, `unknown`
- Check result history with response times and status codes
- Uptime summary statistics (percentage, average response time, total/successful/failed counts)
- Events: `healthcheck_created`, `healthcheck_updated`, `healthcheck_deleted`, `healthcheck_status_changed`, `healthcheck_result_created`

## API Endpoints

Base path: `/api/v1/spaces/{space_ref}/health-checks`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/health-checks` | Create a health monitor |
| GET | `/health-checks` | List monitors |
| GET | `/health-checks/{id}` | Get monitor detail |
| PATCH | `/health-checks/{id}` | Update monitor (including enable/disable) |
| DELETE | `/health-checks/{id}` | Delete monitor |
| GET | `/health-checks/{id}/results` | Get recent check results |
| GET | `/health-checks/{id}/summary` | Get uptime summary |

## Monitor Constraints

| Field | Constraint |
|-------|-----------|
| `identifier` | Pattern `^[a-zA-Z0-9_-]+$`, 3–200 chars |
| `url` | Valid HTTP/HTTPS URL |
| `method` | `GET`, `POST`, or `HEAD` |
| `expected_status` | 100–599 |
| `interval_seconds` | 60–86400 |
| `timeout_seconds` | 1–300 |

## Database Schema

Tables: `health_checks` (migration 0103), `health_check_results` (migration 0104)

## Key Paths

| Purpose | Path |
|---------|------|
| Types | `types/healthcheck.go` |
| Store | `app/store/database/healthcheck.go`, `healthcheck_result.go` |
| Controller | `app/api/controller/healthcheck/` |
| Handlers | `app/api/handler/healthcheck/` |
| Events | `app/events/healthcheck/` |

## Status

**Implemented** — HTTP endpoint monitoring, uptime tracking, result history, and status change events are working.

## Future Work

- Health check failure → remediation trigger (connect to Error Bridge)
- Anomaly detection from response time trends
- Multi-protocol checks (TCP, DNS, ICMP)
- Alerting integrations (webhook, email, Slack)
