# Code Security Scanner Module Implementation

This document describes the Code Security Scanner module implementation for the SoloDev platform.

## Overview

The Code Security Scanner module provides capabilities for conducting security scans on code repositories, detecting vulnerabilities, secrets, and other security issues. It supports multiple scan types (SAST, SCA, Secret Detection) and tracks findings with severity levels and status management.

## Module Structure

### 1. Types Definition
**File:** `types/securityscan.go`

Defines core data structures:
- `ScanResult`: Represents a security scan execution
- `ScanResultInput`: Input for creating/updating scans
- `ScanResultFilter`: Filtering options for listing scans
- `ScanFinding`: Represents a security finding from a scan
- `ScanFindingInput`: Input for creating findings
- `ScanFindingFilter`: Filtering options for findings
- `SecuritySummary`: Aggregate security posture data
- `ScanFindingStatusUpdate`: Update payload for finding status

### 2. Enumerations
**File:** `types/enum/securityscan.go`

Defines enumeration types:
- `SecurityScanType`: sast, sca, secret_detection
- `SecurityScanStatus`: pending, running, completed, failed
- `SecurityScanTrigger`: manual, pipeline, webhook
- `SecurityScanAttr`: Sortable attributes
- `SecurityFindingSeverity`: critical, high, medium, low, info
- `SecurityFindingCategory`: vulnerability, secret, code_smell, bug
- `SecurityFindingStatus`: open, resolved, ignored, false_positive
- `SecurityFindingAttr`: Sortable attributes

### 3. Database Layer
**File:** `app/store/database/securityscan.go`

Implements store interfaces:
- `SecurityScanStore`: CRUD operations for security scans
- `ScanFindingStore`: CRUD operations for scan findings

Key methods:
- `Create()`: Insert new scan/finding
- `Find()`: Retrieve by ID
- `FindByIdentifier()`: Retrieve by unique identifier
- `List()`: Query with pagination and filtering
- `Update()`: Modify existing records
- `Delete()`: Remove records

### 4. Controller
**File:** `app/api/controller/securityscan/controller.go`

Business logic layer implementing:
- `TriggerScan()`: Initiate a new security scan
- `FindScan()`: Retrieve scan details
- `ListScans()`: List scans with filtering
- `UpdateScan()`: Update scan status/results
- `FindFinding()`: Retrieve finding details
- `ListFindings()`: List findings for a scan
- `UpdateFindingStatus()`: Change finding status
- `GetSecuritySummary()`: Get security posture summary

### 5. HTTP Handlers
**Files:** `app/api/handler/securityscan/`

Handler files and functions:
- `trigger_scan.go` - `HandleTriggerScan()`: POST /security-scans
- `find_scan.go` - `HandleFindScan()`: GET /security-scans/{scan_identifier}
- `list_scans.go` - `HandleListScans()`: GET /security-scans
- `list_findings.go` - `HandleListFindings()`: GET /security-scans/{scan_identifier}/findings
- `update_finding_status.go` - `HandleUpdateFindingStatus()`: **Not currently registered in the router** (the controller method exists but no route is defined)
- `security_summary.go` - `HandleGetSecuritySummary()`: GET /security-scans/{scan_identifier}/summary

### 6. Events
**Files:** `app/events/securityscan/`
- `events.go` - Event type and payload definitions
- `reporter.go` - Event publishing interface
- `scan_triggered.go` - Fired when a scan is triggered
- `scan_completed.go` - Fired when a scan completes successfully
- `scan_failed.go` - Fired when a scan fails

### 7. Database Migrations
**Files:**
- `app/store/database/migrate/postgres/0102_create_table_security_scans.up.sql`
- `app/store/database/migrate/postgres/0102_create_table_scan_findings.up.sql`
- `app/store/database/migrate/sqlite/0102_create_table_security_scans.up.sql`
- `app/store/database/migrate/sqlite/0102_create_table_scan_findings.up.sql`

Creates two tables:
- `security_scans`: Stores scan execution records
- `scan_findings`: Stores detected security findings

### 8. Request Helpers
**File:** `app/api/request/securityscan.go`

Provides HTTP request parameter extraction:
- `GetScanIdentifierFromPath()`: Extract scan_identifier from URL

## API Endpoints

### Scan Management
- **POST** `/api/v1/spaces/{space_ref}/security-scans`
  - Trigger a new security scan
  - Body: `ScanResultInput`
  - Returns: `ScanResult`

- **GET** `/api/v1/spaces/{space_ref}/security-scans`
  - List security scans for a repository
  - Query params: `page`, `limit`, `sort`, `order`, `status`, `scan_type`, `triggered_by`
  - Returns: List of `ScanResult` with count

- **GET** `/api/v1/spaces/{space_ref}/security-scans/{scan_identifier}`
  - Get details of a specific scan
  - Returns: `ScanResult`

### Finding Management
- **GET** `/api/v1/spaces/{space_ref}/security-scans/{scan_identifier}/findings`
  - List findings for a scan
  - Query params: `page`, `limit`, `sort`, `order`, `severity`, `category`, `status`
  - Returns: List of `ScanFinding` with count

- **PATCH** `/api/v1/spaces/{space_ref}/security-scans/findings/{id}/status`
  - Update finding status (open, resolved, ignored, false_positive)
  - Body: `ScanFindingStatusUpdate`
  - Returns: Updated `ScanFinding`
  - **Note:** This endpoint is defined in the handler (`update_finding_status.go`) and controller but is **not currently registered** in the router (`app/router/api_modules.go`). A route must be added before this endpoint will be accessible.

### Security Posture
- **GET** `/api/v1/spaces/{space_ref}/security-scans/{scan_identifier}/summary`
  - Get security posture summary for a specific scan
  - Query params: `repo_ref` (optional)
  - Returns: `SecuritySummary`

## Database Schema

### security_scans table
```sql
- ss_id (BIGSERIAL): Primary key
- ss_space_id (BIGINT): Reference to space
- ss_repo_id (BIGINT): Reference to repository
- ss_identifier (TEXT): Unique identifier for the scan
- ss_scan_type (TEXT): Type of scan (sast, sca, secret_detection)
- ss_status (TEXT): Current status (pending, running, completed, failed)
- ss_commit_sha (TEXT): Git commit SHA
- ss_branch (TEXT): Git branch name
- ss_total_issues (INTEGER): Total issues found
- ss_critical_count (INTEGER): Critical severity count
- ss_high_count (INTEGER): High severity count
- ss_medium_count (INTEGER): Medium severity count
- ss_low_count (INTEGER): Low severity count
- ss_duration (BIGINT): Scan duration in milliseconds
- ss_triggered_by (TEXT): Who triggered the scan (manual, pipeline, webhook)
- ss_created_by (BIGINT): User who created the scan
- ss_created (BIGINT): Creation timestamp
- ss_updated (BIGINT): Last update timestamp
- ss_version (BIGINT): Optimistic locking version
```

### scan_findings table
```sql
- sf_id (BIGSERIAL): Primary key
- sf_scan_id (BIGINT): Reference to security scan
- sf_identifier (TEXT): Unique identifier for the finding
- sf_severity (TEXT): Severity level (critical, high, medium, low, info)
- sf_category (TEXT): Category (vulnerability, secret, code_smell, bug)
- sf_title (TEXT): Finding title
- sf_description (TEXT): Detailed description
- sf_file_path (TEXT): Path to affected file
- sf_line_start (INTEGER): Starting line number
- sf_line_end (INTEGER): Ending line number
- sf_rule (TEXT): Rule that triggered this finding
- sf_snippet (TEXT): Code snippet showing the issue
- sf_suggestion (TEXT): Fix recommendation
- sf_status (TEXT): Status (open, resolved, ignored, false_positive)
- sf_cwe (TEXT): CWE identifier if applicable
- sf_created (BIGINT): Creation timestamp
- sf_updated (BIGINT): Last update timestamp
```

## Access Control

All endpoints require proper authorization checked via `apiauth.CheckRepo` or `apiauth.CheckSpace`:
- **`enum.PermissionRepoView`**: Required for reading scans (`FindScan`, `ListScans`) and listing findings (`ListFindings`)
- **`enum.PermissionRepoPush`**: Required for triggering scans (`TriggerScan`) and updating scan results (`UpdateScan`)
- **`enum.PermissionSpaceView`**: Required for finding individual findings (`FindFinding`) and getting the security summary (`GetSecuritySummary`)
- **`enum.PermissionSpaceEdit`**: Required for updating finding status (`UpdateFindingStatus`)

## Integration Points

### With Existing Store Interfaces
The module integrates with the SoloDev store layer by defining new store interfaces in `app/store/database.go` (interfaces added):
- `SecurityScanStore`
- `ScanFindingStore`

### With Enum System
Uses consistent enum patterns from `types/enum` package for:
- Scan types and statuses
- Finding severities and categories
- Sorting and ordering options

### With Events System
Integrates with SoloDev events system for:
- Publishing scan lifecycle events
- Allowing subscribers to react to scan events

## Usage Example

```go
// Trigger a scan
scanInput := &types.ScanResultInput{
    ScanType:    enum.SecurityScanTypeSAST,
    CommitSHA:   "abc123def456",
    Branch:      "main",
    TriggeredBy: enum.SecurityScanTriggerManual,
}
scan, err := controller.TriggerScan(ctx, session, "my-space/my-repo", scanInput)

// List scans
filter := &types.ScanResultFilter{
    Page: 1,
    Size: 20,
    Status: enum.SecurityScanStatusCompleted,
}
scans, count, err := controller.ListScans(ctx, session, "my-space/my-repo", filter)

// Get findings
findings, count, err := controller.ListFindings(
    ctx, session, "my-space/my-repo", scan.Identifier, nil,
)

// Update finding status
update := &types.ScanFindingStatusUpdate{
    Status: enum.SecurityFindingStatusResolved,
}
finding, err := controller.UpdateFindingStatus(ctx, session, "my-space", findingID, update)

// Get security summary
summary, err := controller.GetSecuritySummary(ctx, session, "my-space", nil)
```

## File Locations Summary

| Purpose | File Path |
|---------|-----------|
| Types | `/types/securityscan.go` |
| Enums | `/types/enum/securityscan.go` |
| Store | `/app/store/database/securityscan.go` |
| Store Interfaces | `/app/store/database.go` (added) |
| Controller | `/app/api/controller/securityscan/controller.go` |
| Handlers | `/app/api/handler/securityscan/*.go` |
| Requests | `/app/api/request/securityscan.go` |
| Events | `/app/events/securityscan/*.go` |
| Migrations (Postgres) | `/app/store/database/migrate/postgres/0102_*.sql` |
| Migrations (SQLite) | `/app/store/database/migrate/sqlite/0102_*.sql` |

## Implementation Notes

1. **Identifier Generation**: ScanResult uses UUID for identifier to ensure uniqueness
2. **Timestamp Format**: All timestamps are Unix milliseconds for consistency
3. **Optimistic Locking**: ScanResult includes version field for concurrent updates
4. **Cascade Delete**: Deleting a scan automatically deletes associated findings
5. **Pagination**: All list operations support offset-based pagination
6. **Filtering**: List operations support multiple filter criteria simultaneously
7. **Authorization**: All operations are protected by role-based access control
8. **Events**: Key operations emit events for integration with other systems
