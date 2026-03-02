// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Removed unused imports.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"

	"github.com/jmoiron/sqlx"
)

var _ store.SecurityScanStore = (*SecurityScanStore)(nil)
var _ store.ScanFindingStore = (*ScanFindingStore)(nil)

// NewSecurityScanStore returns a new SecurityScanStore.
func NewSecurityScanStore(db *sqlx.DB) *SecurityScanStore {
	return &SecurityScanStore{
		db: db,
	}
}

// SecurityScanStore implements store.SecurityScanStore backed by a relational database.
type SecurityScanStore struct {
	db *sqlx.DB
}

// scanResult is an internal representation used to store security scan data in the database.
type scanResult struct {
	ID            int64                    `db:"ss_id"`
	SpaceID       int64                    `db:"ss_space_id"`
	RepoID        int64                    `db:"ss_repo_id"`
	Identifier    string                   `db:"ss_identifier"`
	ScanType      enum.SecurityScanType    `db:"ss_scan_type"`
	Status        enum.SecurityScanStatus  `db:"ss_status"`
	CommitSHA     string                   `db:"ss_commit_sha"`
	Branch        string                   `db:"ss_branch"`
	TotalIssues   int                      `db:"ss_total_issues"`
	CriticalCount int                      `db:"ss_critical_count"`
	HighCount     int                      `db:"ss_high_count"`
	MediumCount   int                      `db:"ss_medium_count"`
	LowCount      int                      `db:"ss_low_count"`
	FailureReason string                   `db:"ss_failure_reason"`
	Duration      int64                    `db:"ss_duration"`
	TriggeredBy   enum.SecurityScanTrigger `db:"ss_triggered_by"`
	CreatedBy     int64                    `db:"ss_created_by"`
	Created       int64                    `db:"ss_created"`
	Updated       int64                    `db:"ss_updated"`
	Version       int64                    `db:"ss_version"`
}

const (
	scanResultColumns = `
		 ss_id
		,ss_space_id
		,ss_repo_id
		,ss_identifier
		,ss_scan_type
		,ss_status
		,ss_commit_sha
		,ss_branch
		,ss_total_issues
		,ss_critical_count
		,ss_high_count
		,ss_medium_count
		,ss_low_count
		,ss_failure_reason
		,ss_duration
		,ss_triggered_by
		,ss_created_by
		,ss_created
		,ss_updated
		,ss_version`

	scanResultSelectBase = `
	SELECT` + scanResultColumns + `
	FROM security_scans`
)

// Create creates a new security scan.
func (s *SecurityScanStore) Create(ctx context.Context, scan *types.ScanResult) error {
	const sqlQuery = `
	INSERT INTO security_scans (
		 ss_space_id
		,ss_repo_id
		,ss_identifier
		,ss_scan_type
		,ss_status
		,ss_commit_sha
		,ss_branch
		,ss_total_issues
		,ss_critical_count
		,ss_high_count
		,ss_medium_count
		,ss_low_count
		,ss_failure_reason
		,ss_duration
		,ss_triggered_by
		,ss_created_by
		,ss_created
		,ss_updated
		,ss_version
	) VALUES (
		 $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12, $13, $14, $15, $16, $17, $18, $19
	) RETURNING ss_id`

	db := dbtx.GetAccessor(ctx, s.db)

	if err := db.QueryRowxContext(
		ctx,
		sqlQuery,
		scan.SpaceID,
		scan.RepoID,
		scan.Identifier,
		scan.ScanType,
		scan.Status,
		scan.CommitSHA,
		scan.Branch,
		scan.TotalIssues,
		scan.CriticalCount,
		scan.HighCount,
		scan.MediumCount,
		scan.LowCount,
		scan.FailureReason,
		scan.Duration,
		scan.TriggeredBy,
		scan.CreatedBy,
		scan.Created,
		scan.Updated,
		scan.Version,
	).Scan(&scan.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Insert query failed")
	}

	return nil
}

// Find finds a security scan by ID.
func (s *SecurityScanStore) Find(ctx context.Context, id int64) (*types.ScanResult, error) {
	const sqlQuery = scanResultSelectBase + ` WHERE ss_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := &scanResult{}
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select query failed")
	}

	return mapToScanResult(dst), nil
}

// FindByIdentifier finds a security scan by identifier for a repo.
func (s *SecurityScanStore) FindByIdentifier(
	ctx context.Context,
	repoID int64,
	identifier string,
) (*types.ScanResult, error) {
	const sqlQuery = scanResultSelectBase + `
	WHERE ss_repo_id = $1 AND LOWER(ss_identifier) = LOWER($2)`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := &scanResult{}
	if err := db.GetContext(ctx, dst, sqlQuery, repoID, identifier); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select query failed")
	}

	return mapToScanResult(dst), nil
}

// List lists security scans with pagination and filtering.
func (s *SecurityScanStore) List(
	ctx context.Context,
	repoID int64,
	filter *types.ScanResultFilter,
) ([]*types.ScanResult, int64, error) {
	if filter == nil {
		filter = &types.ScanResultFilter{}
	}

	stmt := database.Builder.
		Select(scanResultColumns).
		From("security_scans").
		Where("ss_repo_id = ?", repoID)

	if filter.Status != "" {
		stmt = stmt.Where("ss_status = ?", filter.Status)
	}

	if filter.ScanType != "" {
		stmt = stmt.Where("ss_scan_type = ?", filter.ScanType)
	}

	if filter.TriggeredBy != "" {
		stmt = stmt.Where("ss_triggered_by = ?", filter.TriggeredBy)
	}

	countStmt := database.Builder.
		Select("COUNT(*)").
		From("security_scans").
		Where("ss_repo_id = ?", repoID)

	if filter.Status != "" {
		countStmt = countStmt.Where("ss_status = ?", filter.Status)
	}

	if filter.ScanType != "" {
		countStmt = countStmt.Where("ss_scan_type = ?", filter.ScanType)
	}

	if filter.TriggeredBy != "" {
		countStmt = countStmt.Where("ss_triggered_by = ?", filter.TriggeredBy)
	}

	// Apply sorting
	if filter.Sort != enum.SecurityScanAttrNone {
		stmt = stmt.OrderBy(fmt.Sprintf("%s %s", filter.Sort.String(), filter.Order.String()))
	} else {
		stmt = stmt.OrderBy(fmt.Sprintf("ss_created %s", enum.OrderDesc.String()))
	}

	// Apply pagination
	if filter.Size > 0 {
		stmt = stmt.Limit(uint64(filter.Size))
		if filter.Page > 0 {
			stmt = stmt.Offset(uint64((filter.Page - 1) * filter.Size))
		}
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert query to sql: %w", err)
	}

	countSQL, countArgs, err := countStmt.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert count query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	if err := db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return nil, 0, database.ProcessSQLErrorf(ctx, err, "Count query failed")
	}

	dsts := make([]*scanResult, 0)
	if err := db.SelectContext(ctx, &dsts, sql, args...); err != nil {
		return nil, 0, database.ProcessSQLErrorf(ctx, err, "Select query failed")
	}

	results := make([]*types.ScanResult, len(dsts))
	for i := range dsts {
		results[i] = mapToScanResult(dsts[i])
	}

	return results, count, nil
}

// ListByStatus lists scans across all repos by status, ordered by created time.
func (s *SecurityScanStore) ListByStatus(
	ctx context.Context,
	status enum.SecurityScanStatus,
	limit int,
) ([]*types.ScanResult, error) {
	if limit <= 0 {
		limit = 10
	}

	stmt := database.Builder.
		Select(scanResultColumns).
		From("security_scans").
		Where("ss_status = ?", status).
		OrderBy("ss_created ASC").
		Limit(uint64(limit))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to convert query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dsts := make([]*scanResult, 0)
	if err := db.SelectContext(ctx, &dsts, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "ListByStatus query failed")
	}

	results := make([]*types.ScanResult, len(dsts))
	for i := range dsts {
		results[i] = mapToScanResult(dsts[i])
	}

	return results, nil
}

// Summary returns the latest completed scan summary for a repo or a space-wide aggregate.
func (s *SecurityScanStore) Summary(ctx context.Context, spaceID int64, repoID *int64) (*types.SecuritySummary, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	if repoID != nil {
		const sqlQuery = `
WITH latest_scan AS (
	SELECT ss_id, ss_space_id, ss_repo_id, ss_created
	FROM security_scans
	WHERE ss_space_id = $1
	  AND ss_repo_id = $2
	  AND ss_status = 'completed'
	ORDER BY ss_created DESC, ss_id DESC
	LIMIT 1
)
SELECT
	ls.ss_space_id AS space_id,
	ls.ss_repo_id AS repo_id,
	ls.ss_id AS last_scan_id,
	ls.ss_created AS last_scan_time,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' THEN 1 ELSE 0 END), 0) AS total_findings,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'critical' THEN 1 ELSE 0 END), 0) AS critical_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'high' THEN 1 ELSE 0 END), 0) AS high_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'medium' THEN 1 ELSE 0 END), 0) AS medium_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'low' THEN 1 ELSE 0 END), 0) AS low_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'info' THEN 1 ELSE 0 END), 0) AS info_issues
FROM latest_scan ls
LEFT JOIN scan_findings sf ON sf.sf_scan_id = ls.ss_id
GROUP BY ls.ss_space_id, ls.ss_repo_id, ls.ss_id, ls.ss_created`

		var summary types.SecuritySummary
		if err := db.GetContext(ctx, &summary, sqlQuery, spaceID, *repoID); err != nil {
			return nil, nil
		}
		return &summary, nil
	}

	const sqlQuery = `
WITH ranked_scans AS (
	SELECT
		ss_id,
		ss_space_id,
		ss_repo_id,
		ss_created,
		ROW_NUMBER() OVER (PARTITION BY ss_repo_id ORDER BY ss_created DESC, ss_id DESC) AS rn
	FROM security_scans
	WHERE ss_space_id = $1
	  AND ss_status = 'completed'
),
latest_scans AS (
	SELECT ss_id, ss_space_id, ss_repo_id, ss_created
	FROM ranked_scans
	WHERE rn = 1
)
SELECT
	$1 AS space_id,
	0 AS repo_id,
	COALESCE(MAX(ls.ss_id), 0) AS last_scan_id,
	COALESCE(MAX(ls.ss_created), 0) AS last_scan_time,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' THEN 1 ELSE 0 END), 0) AS total_findings,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'critical' THEN 1 ELSE 0 END), 0) AS critical_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'high' THEN 1 ELSE 0 END), 0) AS high_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'medium' THEN 1 ELSE 0 END), 0) AS medium_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'low' THEN 1 ELSE 0 END), 0) AS low_issues,
	COALESCE(SUM(CASE WHEN sf.sf_status = 'open' AND sf.sf_severity = 'info' THEN 1 ELSE 0 END), 0) AS info_issues
FROM latest_scans ls
LEFT JOIN scan_findings sf ON sf.sf_scan_id = ls.ss_id`

	var summary types.SecuritySummary
	if err := db.GetContext(ctx, &summary, sqlQuery, spaceID); err != nil {
		return nil, fmt.Errorf("failed to summarize security scans: %w", err)
	}
	return &summary, nil
}

// Update updates a security scan.
func (s *SecurityScanStore) Update(ctx context.Context, scan *types.ScanResult) error {
	scan.Updated = time.Now().UnixMilli()
	scan.Version++

	stmt := database.Builder.
		Update("security_scans").
		Set("ss_status", scan.Status).
		Set("ss_total_issues", scan.TotalIssues).
		Set("ss_critical_count", scan.CriticalCount).
		Set("ss_high_count", scan.HighCount).
		Set("ss_medium_count", scan.MediumCount).
		Set("ss_low_count", scan.LowCount).
		Set("ss_failure_reason", scan.FailureReason).
		Set("ss_duration", scan.Duration).
		Set("ss_updated", scan.Updated).
		Set("ss_version", scan.Version).
		Where("ss_id = ?", scan.ID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to convert query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Update query failed")
	}

	return nil
}

// Delete deletes a security scan.
func (s *SecurityScanStore) Delete(ctx context.Context, id int64) error {
	const sqlQuery = `DELETE FROM security_scans WHERE ss_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sqlQuery, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete query failed")
	}

	return nil
}

// mapToScanResult maps a scanResult to types.ScanResult.
func mapToScanResult(sr *scanResult) *types.ScanResult {
	return &types.ScanResult{
		ID:            sr.ID,
		SpaceID:       sr.SpaceID,
		RepoID:        sr.RepoID,
		Identifier:    sr.Identifier,
		ScanType:      sr.ScanType,
		Status:        sr.Status,
		CommitSHA:     sr.CommitSHA,
		Branch:        sr.Branch,
		TotalIssues:   sr.TotalIssues,
		CriticalCount: sr.CriticalCount,
		HighCount:     sr.HighCount,
		MediumCount:   sr.MediumCount,
		LowCount:      sr.LowCount,
		FailureReason: sr.FailureReason,
		Duration:      sr.Duration,
		TriggeredBy:   sr.TriggeredBy,
		CreatedBy:     sr.CreatedBy,
		Created:       sr.Created,
		Updated:       sr.Updated,
		Version:       sr.Version,
	}
}

// NewScanFindingStore returns a new ScanFindingStore.
func NewScanFindingStore(db *sqlx.DB) *ScanFindingStore {
	return &ScanFindingStore{
		db: db,
	}
}

// ScanFindingStore implements store.ScanFindingStore backed by a relational database.
type ScanFindingStore struct {
	db *sqlx.DB
}

// scanFinding is an internal representation used to store scan finding data in the database.
type scanFinding struct {
	ID          int64                        `db:"sf_id"`
	ScanID      int64                        `db:"sf_scan_id"`
	Identifier  string                       `db:"sf_identifier"`
	Severity    enum.SecurityFindingSeverity `db:"sf_severity"`
	Category    enum.SecurityFindingCategory `db:"sf_category"`
	Title       string                       `db:"sf_title"`
	Description string                       `db:"sf_description"`
	FilePath    string                       `db:"sf_file_path"`
	LineStart   int                          `db:"sf_line_start"`
	LineEnd     int                          `db:"sf_line_end"`
	Rule        string                       `db:"sf_rule"`
	Snippet     string                       `db:"sf_snippet"`
	Suggestion  string                       `db:"sf_suggestion"`
	Status      enum.SecurityFindingStatus   `db:"sf_status"`
	CWE         string                       `db:"sf_cwe"`
	Created     int64                        `db:"sf_created"`
	Updated     int64                        `db:"sf_updated"`
}

const (
	scanFindingColumns = `
		 sf_id
		,sf_scan_id
		,sf_identifier
		,sf_severity
		,sf_category
		,sf_title
		,sf_description
		,sf_file_path
		,sf_line_start
		,sf_line_end
		,sf_rule
		,sf_snippet
		,sf_suggestion
		,sf_status
		,sf_cwe
		,sf_created
		,sf_updated`

	scanFindingSelectBase = `
	SELECT` + scanFindingColumns + `
	FROM scan_findings`
)

// Create creates a new scan finding.
func (s *ScanFindingStore) Create(ctx context.Context, finding *types.ScanFinding) error {
	const sqlQuery = `
	INSERT INTO scan_findings (
		 sf_scan_id
		,sf_identifier
		,sf_severity
		,sf_category
		,sf_title
		,sf_description
		,sf_file_path
		,sf_line_start
		,sf_line_end
		,sf_rule
		,sf_snippet
		,sf_suggestion
		,sf_status
		,sf_cwe
		,sf_created
		,sf_updated
	) VALUES (
		 $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		,$11, $12, $13, $14, $15, $16
	) RETURNING sf_id`

	db := dbtx.GetAccessor(ctx, s.db)

	if err := db.QueryRowxContext(
		ctx,
		sqlQuery,
		finding.ScanID,
		finding.Identifier,
		finding.Severity,
		finding.Category,
		finding.Title,
		finding.Description,
		finding.FilePath,
		finding.LineStart,
		finding.LineEnd,
		finding.Rule,
		finding.Snippet,
		finding.Suggestion,
		finding.Status,
		finding.CWE,
		finding.Created,
		finding.Updated,
	).Scan(&finding.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Insert query failed")
	}

	return nil
}

// Find finds a scan finding by ID.
func (s *ScanFindingStore) Find(ctx context.Context, id int64) (*types.ScanFinding, error) {
	const sqlQuery = scanFindingSelectBase + ` WHERE sf_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := &scanFinding{}
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select query failed")
	}

	return mapToScanFinding(dst), nil
}

// ListByScan lists all findings for a specific scan.
func (s *ScanFindingStore) ListByScan(
	ctx context.Context,
	scanID int64,
	filter *types.ScanFindingFilter,
) ([]*types.ScanFinding, int64, error) {
	if filter == nil {
		filter = &types.ScanFindingFilter{}
	}

	stmt := database.Builder.
		Select(scanFindingColumns).
		From("scan_findings").
		Where("sf_scan_id = ?", scanID)

	if filter.Severity != "" {
		stmt = stmt.Where("sf_severity = ?", filter.Severity)
	}

	if filter.Category != "" {
		stmt = stmt.Where("sf_category = ?", filter.Category)
	}

	if filter.Status != "" {
		stmt = stmt.Where("sf_status = ?", filter.Status)
	}

	countStmt := database.Builder.
		Select("COUNT(*)").
		From("scan_findings").
		Where("sf_scan_id = ?", scanID)

	if filter.Severity != "" {
		countStmt = countStmt.Where("sf_severity = ?", filter.Severity)
	}

	if filter.Category != "" {
		countStmt = countStmt.Where("sf_category = ?", filter.Category)
	}

	if filter.Status != "" {
		countStmt = countStmt.Where("sf_status = ?", filter.Status)
	}

	// Apply sorting
	if filter.Sort != enum.SecurityFindingAttrNone {
		stmt = stmt.OrderBy(fmt.Sprintf("%s %s", filter.Sort.String(), filter.Order.String()))
	} else {
		stmt = stmt.OrderBy(fmt.Sprintf("sf_severity %s", enum.OrderDesc.String()))
	}

	// Apply pagination
	if filter.Size > 0 {
		stmt = stmt.Limit(uint64(filter.Size))
		if filter.Page > 0 {
			stmt = stmt.Offset(uint64((filter.Page - 1) * filter.Size))
		}
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert query to sql: %w", err)
	}

	countSQL, countArgs, err := countStmt.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to convert count query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	if err := db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return nil, 0, database.ProcessSQLErrorf(ctx, err, "Count query failed")
	}

	dsts := make([]*scanFinding, 0)
	if err := db.SelectContext(ctx, &dsts, sql, args...); err != nil {
		return nil, 0, database.ProcessSQLErrorf(ctx, err, "Select query failed")
	}

	results := make([]*types.ScanFinding, len(dsts))
	for i := range dsts {
		results[i] = mapToScanFinding(dsts[i])
	}

	return results, count, nil
}

// Update updates a scan finding.
func (s *ScanFindingStore) Update(ctx context.Context, finding *types.ScanFinding) error {
	finding.Updated = time.Now().UnixMilli()

	stmt := database.Builder.
		Update("scan_findings").
		Set("sf_severity", finding.Severity).
		Set("sf_category", finding.Category).
		Set("sf_status", finding.Status).
		Set("sf_updated", finding.Updated).
		Where("sf_id = ?", finding.ID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to convert query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Update query failed")
	}

	return nil
}

// Delete deletes a scan finding.
func (s *ScanFindingStore) Delete(ctx context.Context, id int64) error {
	const sqlQuery = `DELETE FROM scan_findings WHERE sf_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sqlQuery, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete query failed")
	}

	return nil
}

// DeleteByScan deletes all findings for a scan.
func (s *ScanFindingStore) DeleteByScan(ctx context.Context, scanID int64) error {
	const sqlQuery = `DELETE FROM scan_findings WHERE sf_scan_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sqlQuery, scanID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete query failed")
	}

	return nil
}

// mapToScanFinding maps a scanFinding to types.ScanFinding.
func mapToScanFinding(sf *scanFinding) *types.ScanFinding {
	return &types.ScanFinding{
		ID:          sf.ID,
		ScanID:      sf.ScanID,
		Identifier:  sf.Identifier,
		Severity:    sf.Severity,
		Category:    sf.Category,
		Title:       sf.Title,
		Description: sf.Description,
		FilePath:    sf.FilePath,
		LineStart:   sf.LineStart,
		LineEnd:     sf.LineEnd,
		Rule:        sf.Rule,
		Snippet:     sf.Snippet,
		Suggestion:  sf.Suggestion,
		Status:      sf.Status,
		CWE:         sf.CWE,
		Created:     sf.Created,
		Updated:     sf.Updated,
	}
}
