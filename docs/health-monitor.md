# Uptime & Health Monitor Module

This document describes the Uptime & Health Monitor module implementation for the Harness platform.

## Overview

The Health Monitor module provides comprehensive health check and uptime monitoring capabilities for the Harness platform. It allows users to configure monitors for HTTP endpoints, track their status over time, and view uptime statistics.

## Architecture

The module follows the standard Harness architecture patterns:

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
  - All handlers follow standard Harness patterns with authentication and error handling

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
POST /api/v1/spaces/{space_ref}/+/monitors
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
GET /api/v1/spaces/{space_ref}/+/monitors?page=0&size=50&query=api
```

**Get Monitor**
```
GET /api/v1/spaces/{space_ref}/+/monitors/{identifier}
```

**Update Monitor**
```
PATCH /api/v1/spaces/{space_ref}/+/monitors/{identifier}
Content-Type: application/json

{
  "name": "Updated Name",
  "interval_seconds": 600
}
```

**Toggle Monitor**
```
PATCH /api/v1/spaces/{space_ref}/+/monitors/{identifier}/toggle
```

**Delete Monitor**
```
DELETE /api/v1/spaces/{space_ref}/+/monitors/{identifier}
```

### Results & Statistics

**Get Recent Results**
```
GET /api/v1/spaces/{space_ref}/+/monitors/{identifier}/results?limit=100
```

Returns array of HealthCheckResult objects with latest checks first.

**Get Uptime Summary**
```
GET /api/v1/spaces/{space_ref}/+/monitors/summary
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
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/types/healthcheck.go`

### Store (Database)
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database.go` (interfaces)
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/healthcheck.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/healthcheck_result.go`

### Controller
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/controller.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/create.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/find.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/update.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/delete.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/results.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/wire.go`

### Handlers
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/handler/healthcheck/create.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/handler/healthcheck/find.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/controller/healthcheck/update.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/api/handler/healthcheck/delete.go`

### Events
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/events/healthcheck/category.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/events/healthcheck/events.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/events/healthcheck/reporter.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/events/healthcheck/reader.go`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/events/healthcheck/wire.go`

### Migrations

PostgreSQL:
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/postgres/0103_create_table_health_checks.up.sql`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/postgres/0103_create_table_health_checks.down.sql`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/postgres/0104_create_table_health_check_results.up.sql`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/postgres/0104_create_table_health_check_results.down.sql`

SQLite:
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/sqlite/0103_create_table_health_checks.up.sql`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/sqlite/0103_create_table_health_checks.down.sql`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/sqlite/0104_create_table_health_check_results.up.sql`
- `/sessions/fervent-eloquent-cerf/mnt/Harness-io/harness/app/store/database/migrate/sqlite/0104_create_table_health_check_results.down.sql`

## Integration Points

To integrate this module into the Harness application:

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

3. **Router Setup**: Add routes:
   ```go
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors",
      healthcheckhandler.HandleCreate(controller)).Methods(http.MethodPost)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors",
      healthcheckhandler.HandleList(controller)).Methods(http.MethodGet)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors/{identifier}",
      healthcheckhandler.HandleFind(controller)).Methods(http.MethodGet)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors/{identifier}",
      healthcheckhandler.HandleUpdate(controller)).Methods(http.MethodPatch)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors/{identifier}",
      healthcheckhandler.HandleDelete(controller)).Methods(http.MethodDelete)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors/{identifier}/toggle",
      healthcheckhandler.HandleToggle(controller)).Methods(http.MethodPatch)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors/{identifier}/results",
      healthcheckhandler.HandleListResults(controller)).Methods(http.MethodGet)
   router.Handle("/api/v1/spaces/{space_ref}/+/monitors/summary",
      healthcheckhandler.HandleSummary(controller)).Methods(http.MethodGet)
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
- Authorization checks follow Harness conventions using the authz package
- Error handling uses the usererror package for user-friendly messages
- The module is fully multi-tenant and space-aware
