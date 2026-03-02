# Uptime & Health Monitor Module

This document describes the Uptime & Health Monitor module implementation for the SoloDev platform.

## Overview

The Health Monitor module provides comprehensive health check and uptime monitoring capabilities for the SoloDev platform. It allows users to configure monitors for HTTP endpoints, track their status over time, and view uptime statistics.

## Architecture

The module follows the standard SoloDev architecture patterns:

### Database Layer
- **Store Interfaces**: `HealthCheckStore` and `HealthCheckResultStore` in `/app/store/database.go`
- **Store Implementations**:
  - `/app/store/database/healthcheck.go` - Health check CRUD operations
  - `/app/store/database/healthcheck_result.go` - Health check result operations
- **Migrations**:
  - PostgreSQL: `0103_create_table_health_checks.up.sql`, `0104_create_table_health_check_results.up.sql`
  - SQLite: Same files in `migrate/sqlite/`

### Domain Types
- `/types/healthcheck.go`:
  - `HealthCheck` - Monitor configuration entity
  - `HealthCheckResult` - Individual check result
  - `HealthCheckSummary` - Aggregated uptime statistics
  - `HealthCheckStatus` enum (up, down, degraded, unknown)

### API Layer

#### Controller
- `/app/api/controller/healthcheck/`
  - `controller.go` - Main controller with dependencies
  - `create.go` - Create health check validation and business logic
  - `find.go` - Find and list operations with authorization
  - `update.go` - Update and toggle operations
  - `delete.go` - Delete operations
  - `results.go` - Result retrieval and summary operations
  - `wire.go` - Dependency injection setup

#### Handlers
- `/app/api/handler/healthcheck/`
  - `create.go` - HTTP POST handler for creating monitors
  - `find.go` - HTTP GET handlers for retrieving monitors and results
  - `update.go` - HTTP PATCH handlers for updating and toggling
  - `delete.go` - HTTP DELETE handler
  - All handlers follow standard SoloDev patterns with authentication and error handling

### Events System
- `/app/events/healthcheck/`
  - `category.go` - Event category definition
  - `events.go` - Event type and payload definitions
  - `reporter.go` - Event publishing interface
  - `reader.go` - Event subscription interface
  - `wire.go` - Event system integration

## Database Schema

### health_checks Table
```
hc_id                  - Primary key (auto-increment)
hc_space_id            - FK to spaces (required for multi-tenancy)
hc_identifier          - Unique identifier within space (required, indexed)
hc_name                - Display name of the monitor
hc_description         - Optional description
hc_url                 - HTTP endpoint to monitor (required)
hc_method              - HTTP method (GET, POST, HEAD)
hc_expected_status     - Expected HTTP status code
hc_interval_seconds    - Check frequency (60-86400)
hc_timeout_seconds     - Request timeout (1-300)
hc_enabled             - Enable/disable monitor
hc_headers             - JSON custom headers
hc_body                - Request body for POST
hc_tags                - JSON array of tags
hc_last_status         - Latest status (up, down, degraded, unknown)
hc_last_checked_at     - Timestamp of last check (unix millis)
hc_last_response_time  - Last response time (ms)
hc_consecutive_failures- Count of consecutive failures
hc_created_by          - FK to principals
hc_created             - Creation timestamp (unix millis)
hc_updated             - Update timestamp (unix millis)
hc_version             - Optimistic locking version
```

### health_check_results Table
```
hcr_id                 - Primary key (auto-increment)
hcr_health_check_id    - FK to health_checks (indexed)
hcr_status             - Status (up, down, degraded)
hcr_response_time      - Response time in milliseconds
hcr_status_code        - HTTP status code received
hcr_error_message      - Error details if failed
hcr_created_at         - Result timestamp (unix millis, indexed)
```

## API Endpoints

All endpoints require authentication and authorization checks against the parent space.

### Health Check Operations

**Create Monitor**
```
POST /api/v1/spaces/{space_ref}/health-checks
Content-Type: application/json

{
  "identifier": "api-health",
  "name": "API Health Check",
  "description": "Monitor main API endpoint",
  "url": "https://api.example.com/health",
  "method": "GET",
  "expected_status": 200,
  "interval_seconds": 300,
  "timeout_seconds": 10,
  "enabled": true,
  "headers": "{}",
  "body": "",
  "tags": "[\"production\"]"
}
```

**List Monitors**
```
GET /api/v1/spaces/{space_ref}/health-checks?page=0&size=50&query=api
```

**Get Monitor**
```
GET /api/v1/spaces/{space_ref}/health-checks/{identifier}
```

**Update Monitor**
```
PATCH /api/v1/spaces/{space_ref}/health-checks/{identifier}
Content-Type: application/json

{
  "name": "Updated Name",
  "interval_seconds": 600
}
```

**Delete Monitor**
```
DELETE /api/v1/spaces/{space_ref}/health-checks/{identifier}
```

> **Note:** There is no dedicated toggle endpoint in the router. To enable/disable a monitor, use the Update endpoint with `{"enabled": true/false}`.

### Results & Statistics

**Get Recent Results**
```
GET /api/v1/spaces/{space_ref}/health-checks/{identifier}/results?limit=100
```

Returns array of HealthCheckResult objects with latest checks first.

**Get Uptime Summary**
```
GET /api/v1/spaces/{space_ref}/health-checks/{identifier}/summary
```

Returns array of HealthCheckSummary objects with aggregated statistics:
- `uptime_percentage` - Percentage of successful checks in past 24h
- `total_checks` - Total check count
- `successful_checks` - Count of successful checks
- `failed_checks` - Count of failed checks
- `average_response_time` - Average response time in ms

## Validation Rules

### Health Check Creation/Update

- **Identifier**: Must match pattern `^[a-zA-Z0-9_-]+$` (3-200 chars)
- **Name**: Required, 1-255 characters
- **URL**: Must be valid HTTP/HTTPS URL
- **Method**: Must be GET, POST, or HEAD
- **Expected Status**: 100-599
- **Interval**: 60-86400 seconds
- **Timeout**: 1-300 seconds

## Events

Five event types are published by the module:

1. **healthcheck_created** - Monitor created
2. **healthcheck_updated** - Monitor configuration changed
3. **healthcheck_deleted** - Monitor deleted
4. **healthcheck_status_changed** - Status changed (up→down, etc.)
5. **healthcheck_result_created** - Individual check completed

See `/app/events/healthcheck/events.go` for payload structures.

## Optimistic Locking

Health checks use optimistic locking via the `hc_version` field:
- Version incremented on each update
- Update fails with `ErrVersionConflict` if version mismatch
- Controller automatically retries with latest version
- Prevents lost updates in concurrent scenarios

## Authorization

All endpoints check authorization against the parent space using:
- `enum.PermissionRepoView` - Required for all health check operations

## File Locations Summary

### Types
- `types/healthcheck.go`

### Store (Database)
- `app/store/database.go` (interfaces)
- `app/store/database/healthcheck.go`
- `app/store/database/healthcheck_result.go`

### Controller
- `app/api/controller/healthcheck/controller.go`
- `app/api/controller/healthcheck/create.go`
- `app/api/controller/healthcheck/find.go`
- `app/api/controller/healthcheck/update.go`
- `app/api/controller/healthcheck/delete.go`
- `app/api/controller/healthcheck/results.go`
- `app/api/controller/healthcheck/wire.go`

### Handlers
- `app/api/handler/healthcheck/create.go`
- `app/api/handler/healthcheck/find.go`
- `app/api/handler/healthcheck/update.go`
- `app/api/handler/healthcheck/delete.go`

### Events
- `app/events/healthcheck/category.go`
- `app/events/healthcheck/events.go`
- `app/events/healthcheck/reporter.go`
- `app/events/healthcheck/reader.go`
- `app/events/healthcheck/wire.go`

### Migrations

PostgreSQL:
- `app/store/database/migrate/postgres/0103_create_table_health_checks.up.sql`
- `app/store/database/migrate/postgres/0103_create_table_health_checks.down.sql`
- `app/store/database/migrate/postgres/0104_create_table_health_check_results.up.sql`
- `app/store/database/migrate/postgres/0104_create_table_health_check_results.down.sql`

SQLite:
- `app/store/database/migrate/sqlite/0103_create_table_health_checks.up.sql`
- `app/store/database/migrate/sqlite/0103_create_table_health_checks.down.sql`
- `app/store/database/migrate/sqlite/0104_create_table_health_check_results.up.sql`
- `app/store/database/migrate/sqlite/0104_create_table_health_check_results.down.sql`

## Integration Points

To integrate this module into the SoloDev application:

1. **Wire Setup**: Add to main wire.go:
   ```go
   healthcheckcontroller.WireSet,
   healthcheckevent.WireSet,
   ```

2. **Store Registration**: Add to database store setup:
   ```go
   HealthCheckStore: database.NewHealthCheckStore(db),
   HealthCheckResultStore: database.NewHealthCheckResultStore(db),
   ```

3. **Router Setup**: Add routes (see `app/router/api_modules.go`):
   ```go
   r.Route("/health-checks", func(r chi.Router) {
      r.Post("/", handlerhealthcheck.HandleCreate(healthCheckCtrl))
      r.Get("/", handlerhealthcheck.HandleList(healthCheckCtrl))
      r.Route("/{identifier}", func(r chi.Router) {
         r.Get("/", handlerhealthcheck.HandleFind(healthCheckCtrl))
         r.Patch("/", handlerhealthcheck.HandleUpdate(healthCheckCtrl))
         r.Delete("/", handlerhealthcheck.HandleDelete(healthCheckCtrl))
         r.Get("/results", handlerhealthcheck.HandleListResults(healthCheckCtrl))
         r.Get("/summary", handlerhealthcheck.HandleSummary(healthCheckCtrl))
      })
   })
   ```

4. **Migration Registration**: Add to migration list:
   ```go
   postgresqlMigrations = append(postgresqlMigrations,
      migrate.Postgre0103CreateHealthChecks,
      migrate.Postgre0104CreateHealthCheckResults)
   sqliteMigrations = append(sqliteMigrations,
      migrate.Sqlite0103CreateHealthChecks,
      migrate.Sqlite0104CreateHealthCheckResults)
   ```

## Implementation Notes

- All timestamps are stored as Unix milliseconds for consistency
- JSON fields (headers, body, tags) are stored as TEXT in database
- The module uses context-aware database transactions via dbtx.GetAccessor
- Authorization checks follow SoloDev conventions using the authz package
- Error handling uses the usererror package for user-friendly messages
- The module is fully multi-tenant and space-aware
