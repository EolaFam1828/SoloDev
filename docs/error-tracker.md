# Error Tracker Module Implementation

## Overview

The Error Tracker module provides comprehensive error tracking and management capabilities for the Harness platform. It allows applications to report errors, group similar errors together, track occurrences, and manage error statuses.

## Architecture

### Components

#### 1. Types (`types/errortracker.go`)
Core data structures for the error tracker:
- **ErrorGroup**: Groups similar errors together based on fingerprint
- **ErrorOccurrence**: Individual error instances/occurrences
- **ErrorGroupStatus**: Status enum (open, resolved, ignored, regressed)
- **ErrorSeverity**: Severity enum (fatal, error, warning)
- **ReportErrorInput**: API request body for reporting errors
- **UpdateErrorGroupInput**: API request body for updating error groups
- **ErrorTrackerSummary**: Statistics summary for error groups

#### 2. Database Store (`app/store/database/errortracker.go`)
Implements database operations for error tracking:
- `CreateOrUpdateErrorGroup()`: Create new or update existing error groups with conflict handling
- `CreateErrorOccurrence()`: Create individual error occurrences
- `FindByIdentifier()`: Retrieve error group by identifier
- `FindByFingerprint()`: Retrieve error group by fingerprint
- `List()`: Paginated listing with filtering options
- `Count()`: Get count of error groups
- `UpdateStatus()`: Update error group status with version increment
- `UpdateAssignment()`: Assign error group to a user
- `ListOccurrences()`: Get occurrences for an error group
- `GetSummary()`: Get statistics summary

#### 3. Store Interface (`app/store/database.go`)
Public interface definition for ErrorTrackerStore that other components depend on.

#### 4. Controller (`app/api/controller/errortracker/controller.go`)
Business logic and orchestration:
- `ReportError()`: Handle error reporting with transaction management
- `ListErrors()`: List error groups with filtering
- `GetError()`: Get error group details with related occurrences and user information
- `UpdateError()`: Update error group status/assignment with event reporting
- `ListOccurrences()`: Get occurrences with pagination
- `GetSummary()`: Get statistics for a space

#### 5. Handlers (`app/api/handler/errortracker/`)
HTTP request handlers:
- `HandleErrorReport()`: POST endpoint for reporting errors
- `HandleErrorList()`: GET endpoint for listing error groups
- `HandleErrorDetail()`: GET endpoint for error group details
- `HandleErrorUpdate()`: PATCH endpoint for updating error groups
- `HandleErrorOccurrences()`: GET endpoint for listing occurrences
- `HandleErrorSummary()`: GET endpoint for summary statistics

#### 6. Events (`app/events/errortracker/`)
Event system for error tracking:
- `events.go`: Event data structures (ErrorReported, ErrorStatusChanged, ErrorAssigned)
- `reporter.go`: Reports events to event system
- `reader.go`: Provides read access to events

#### 7. API Request Helpers (`app/api/request/errortracker.go`)
Path parameter extraction utilities:
- `GetErrorIdentifierFromPath()`: Extract error identifier from URL path

#### 8. Migrations (`app/store/database/migrate/`)
Database schema migrations (0171):
- PostgreSQL and SQLite scripts
- Tables: `error_groups` and `error_occurrences`
- Indexed for efficient querying

## Database Schema

### error_groups Table
```sql
- eg_id (SERIAL PRIMARY KEY)
- eg_space_id (INT, FK to spaces)
- eg_repo_id (INT, FK to repositories)
- eg_identifier (TEXT)
- eg_title (TEXT)
- eg_message (TEXT)
- eg_fingerprint (TEXT, UNIQUE)
- eg_status (TEXT: open|resolved|ignored|regressed)
- eg_severity (TEXT: fatal|error|warning)
- eg_first_seen (BIGINT, unix milliseconds)
- eg_last_seen (BIGINT, unix milliseconds)
- eg_occurrence_count (BIGINT)
- eg_file_path (TEXT, optional)
- eg_line_number (INT, optional)
- eg_function_name (TEXT, optional)
- eg_language (TEXT, optional)
- eg_tags (JSON)
- eg_assigned_to (INT, FK to principals)
- eg_resolved_at (BIGINT, optional)
- eg_resolved_by (INT, FK to principals)
- eg_created_by (INT, FK to principals)
- eg_created (BIGINT)
- eg_updated (BIGINT)
- eg_version (BIGINT, for optimistic locking)
```

### error_occurrences Table
```sql
- eo_id (SERIAL PRIMARY KEY)
- eo_error_group_id (INT, FK to error_groups)
- eo_stack_trace (TEXT)
- eo_environment (TEXT: production|staging|development)
- eo_runtime (TEXT, optional)
- eo_os (TEXT, optional)
- eo_arch (TEXT, optional)
- eo_metadata (JSON, arbitrary context)
- eo_created_at (BIGINT)
```

## API Endpoints

### 1. Report Error
**POST** `/api/v1/spaces/{space_ref}/errors`

Request body:
```json
{
  "identifier": "string (required)",
  "title": "string (required)",
  "message": "string (required)",
  "severity": "fatal|error|warning",
  "file_path": "string",
  "line_number": 0,
  "function_name": "string",
  "language": "go|python|javascript",
  "tags": ["string"],
  "stack_trace": "string (required)",
  "environment": "production|staging|development",
  "runtime": "go1.24|node20",
  "os": "linux|darwin|windows",
  "arch": "amd64|arm64",
  "metadata": {}
}
```

Response: ErrorGroup

### 2. List Error Groups
**GET** `/api/v1/spaces/{space_ref}/errors`

Query parameters:
- `page`: Page number (default 0)
- `limit`: Results per page (default 50)
- `query`: Search query
- `status`: Filter by status (open|resolved|ignored|regressed)
- `severity`: Filter by severity (fatal|error|warning)
- `language`: Filter by language

Response: []ErrorGroup

### 3. Get Error Details
**GET** `/api/v1/spaces/{space_ref}/errors/{error_identifier}`

Response: ErrorGroupDetail (includes sample occurrences and user info)

### 4. Update Error Group
**PATCH** `/api/v1/spaces/{space_ref}/errors/{error_identifier}`

Request body:
```json
{
  "status": "open|resolved|ignored|regressed",
  "assigned_to": 12345,
  "title": "string",
  "tags": ["string"]
}
```

Response: ErrorGroup

### 5. List Error Occurrences
**GET** `/api/v1/spaces/{space_ref}/errors/{error_identifier}/occurrences`

Query parameters:
- `page`: Page number
- `limit`: Results per page (max 100)

Response: []ErrorOccurrence

### 6. Get Summary Statistics
**GET** `/api/v1/spaces/{space_ref}/errors/summary`

Response: ErrorTrackerSummary
```json
{
  "total_errors": 0,
  "open_errors": 0,
  "resolved_errors": 0,
  "ignored_errors": 0,
  "regression_errors": 0,
  "fatal_count": 0,
  "error_count": 0,
  "warning_count": 0,
  "last_updated": 0
}
```

## Key Features

### 1. Error Fingerprinting
Errors are automatically fingerprinted using SHA256 hash of title + stack trace. Duplicate errors are automatically grouped together and occurrence count is incremented.

### 2. Automatic Status Regression
When an error is reported that was previously resolved, its status automatically changes to "regressed".

### 3. Optimistic Locking
The ErrorGroup version field implements optimistic locking to prevent concurrent update conflicts.

### 4. Event Reporting
The system emits events for:
- Error reported
- Status changed
- Assignment changed

### 5. Role-Based Access Control
All endpoints require appropriate permissions (view/edit) on the space.

### 6. Transaction Management
Error reporting uses transactions to ensure atomic creation of error groups and occurrences.

## Integration Points

### 1. Space Finder (refcache.SpaceFinder)
Used to resolve space references and validate access.

### 2. Authorizer (authz.Authorizer)
Used for role-based access control checks.

### 3. Principal Info Cache (store.PrincipalInfoCache)
Used to fetch user information for display.

### 4. Transactor (dbtx.Transactor)
Used for transaction management during error reporting.

### 5. Event Reporter (errortrackerevents.Reporter)
Used to emit error tracking events.

## File Structure

```
harness/
├── types/
│   └── errortracker.go                    # Core data types
├── app/
│   ├── api/
│   │   ├── controller/
│   │   │   └── errortracker/
│   │   │       └── controller.go          # Business logic
│   │   ├── handler/
│   │   │   └── errortracker/
│   │   │       ├── error_report.go        # POST handler
│   │   │       ├── error_list.go          # GET list handler
│   │   │       ├── error_detail.go        # GET detail handler
│   │   │       ├── error_update.go        # PATCH handler
│   │   │       ├── error_occurrences.go   # GET occurrences handler
│   │   │       └── error_summary.go       # GET summary handler
│   │   └── request/
│   │       └── errortracker.go            # Request helpers
│   ├── events/
│   │   └── errortracker/
│   │       ├── events.go                  # Event structures
│   │       ├── reporter.go                # Event reporter
│   │       └── reader.go                  # Event reader
│   └── store/
│       ├── errortracker.go                # Store helper functions
│       └── database/
│           ├── errortracker.go            # Database implementation
│           └── migrate/
│               ├── postgres/
│               │   ├── 0171_create_table_error_groups.up.sql
│               │   └── 0171_create_table_error_groups.down.sql
│               └── sqlite/
│                   ├── 0171_create_table_error_groups.up.sql
│                   └── 0171_create_table_error_groups.down.sql
```

## Dependencies

### Go Packages
- `github.com/jmoiron/sqlx`: Database access
- `github.com/Masterminds/squirrel`: Query builder
- `github.com/harness/gitness`: Internal packages

### Services
- Space authentication and authorization
- Principal info caching
- Event publishing
- Database transactions

## Usage Example

```go
// Report an error
errorGroup, err := controller.ReportError(ctx, session, "my-space", &types.ReportErrorInput{
    Identifier:   "db-connection-timeout",
    Title:        "Database Connection Timeout",
    Message:      "Failed to connect to database within 30s",
    Severity:     types.ErrorSeverityError,
    StackTrace:   "...",
    Environment:  "production",
    Language:     "go",
    Tags:         []string{"database", "critical"},
})

// List errors
errors, err := controller.ListErrors(ctx, session, "my-space", types.ErrorTrackerListOptions{
    Status: &types.ErrorGroupStatusOpen,
})

// Update status
updated, err := controller.UpdateError(ctx, session, "my-space", "db-connection-timeout",
    &types.UpdateErrorGroupInput{
        Status: &types.ErrorGroupStatusResolved,
    })
```

## Notes

- Migration number: 0171
- Requires database schema migration
- All timestamps are stored in Unix milliseconds
- Fingerprints are SHA256 hashes (hex encoded)
- Foreign keys cascade on delete for repository and space
- Version field for optimistic locking

## Error Bridge Integration (SoloDev)

The Error Tracker controller now includes an optional integration with the AI Auto-Remediation system via the Error Bridge.

### How It Works

When `SetErrorBridge(bridge)` is called on the controller:
1. Every call to `ReportError()` will, after the normal event reporting, also call `bridge.OnErrorReported()`
2. The bridge creates a pending `Remediation` task with:
   - Stack trace from the error occurrence
   - File path from the error group
   - Severity context for prioritization
   - `[Auto] Fix: {title}` as the remediation title

### Filtering
The bridge skips remediation for:
- Warning-level errors (configurable)
- Already-resolved or ignored error groups

### Setup
```go
bridge := errorbridge.NewBridge(remediationStore, true)
errorTrackerCtrl.SetErrorBridge(bridge)
```

See [error-bridge.md](error-bridge.md) for full bridge documentation.

