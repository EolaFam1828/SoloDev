// Copyright 2023 Harness, Inc.
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
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ store.ErrorTrackerStore = (*ErrorTrackerStore)(nil)

// NewErrorTrackerStore returns a new ErrorTrackerStore.
func NewErrorTrackerStore(db *sqlx.DB) *ErrorTrackerStore {
	return &ErrorTrackerStore{
		db: db,
	}
}

// ErrorTrackerStore implements store.ErrorTrackerStore backed by a relational database.
type ErrorTrackerStore struct {
	db *sqlx.DB
}

const (
	errorGroupColumns = `
		 eg_id
		,eg_space_id
		,eg_repo_id
		,eg_identifier
		,eg_title
		,eg_message
		,eg_fingerprint
		,eg_status
		,eg_severity
		,eg_first_seen
		,eg_last_seen
		,eg_occurrence_count
		,eg_file_path
		,eg_line_number
		,eg_function_name
		,eg_language
		,eg_tags
		,eg_assigned_to
		,eg_resolved_at
		,eg_resolved_by
		,eg_created_by
		,eg_created
		,eg_updated
		,eg_version`

	errorGroupSelectBase = `
    SELECT` + errorGroupColumns + `
	FROM error_groups`

	errorOccurrenceColumns = `
		 eo_id
		,eo_error_group_id
		,eo_stack_trace
		,eo_environment
		,eo_runtime
		,eo_os
		,eo_arch
		,eo_metadata
		,eo_created_at`

	errorOccurrenceSelectBase = `
    SELECT` + errorOccurrenceColumns + `
	FROM error_occurrences`
)

type errorGroup struct {
	ID              int64                `db:"eg_id"`
	SpaceID         int64                `db:"eg_space_id"`
	RepoID          int64                `db:"eg_repo_id"`
	Identifier      string               `db:"eg_identifier"`
	Title           string               `db:"eg_title"`
	Message         string               `db:"eg_message"`
	Fingerprint     string               `db:"eg_fingerprint"`
	Status          types.ErrorGroupStatus `db:"eg_status"`
	Severity        types.ErrorSeverity  `db:"eg_severity"`
	FirstSeen       int64                `db:"eg_first_seen"`
	LastSeen        int64                `db:"eg_last_seen"`
	OccurrenceCount int64                `db:"eg_occurrence_count"`
	FilePath        string               `db:"eg_file_path"`
	LineNumber      int                  `db:"eg_line_number"`
	FunctionName    string               `db:"eg_function_name"`
	Language        string               `db:"eg_language"`
	Tags            json.RawMessage      `db:"eg_tags"`
	AssignedTo      *int64               `db:"eg_assigned_to"`
	ResolvedAt      *int64               `db:"eg_resolved_at"`
	ResolvedBy      *int64               `db:"eg_resolved_by"`
	CreatedBy       int64                `db:"eg_created_by"`
	Created         int64                `db:"eg_created"`
	Updated         int64                `db:"eg_updated"`
	Version         int64                `db:"eg_version"`
}

type errorOccurrence struct {
	ID           int64           `db:"eo_id"`
	ErrorGroupID int64           `db:"eo_error_group_id"`
	StackTrace   string          `db:"eo_stack_trace"`
	Environment  string          `db:"eo_environment"`
	Runtime      string          `db:"eo_runtime"`
	OS           string          `db:"eo_os"`
	Arch         string          `db:"eo_arch"`
	Metadata     json.RawMessage `db:"eo_metadata"`
	CreatedAt    int64           `db:"eo_created_at"`
}

// CreateOrUpdateErrorGroup creates a new or updates an existing error group.
func (s *ErrorTrackerStore) CreateOrUpdateErrorGroup(
	ctx context.Context,
	errorGroup *types.ErrorGroup,
) error {
	const sqlQuery = `
	INSERT INTO error_groups (
		 eg_space_id
		,eg_repo_id
		,eg_identifier
		,eg_title
		,eg_message
		,eg_fingerprint
		,eg_status
		,eg_severity
		,eg_first_seen
		,eg_last_seen
		,eg_occurrence_count
		,eg_file_path
		,eg_line_number
		,eg_function_name
		,eg_language
		,eg_tags
		,eg_assigned_to
		,eg_resolved_at
		,eg_resolved_by
		,eg_created_by
		,eg_created
		,eg_updated
		,eg_version
	) VALUES (
		 :eg_space_id
		,:eg_repo_id
		,:eg_identifier
		,:eg_title
		,:eg_message
		,:eg_fingerprint
		,:eg_status
		,:eg_severity
		,:eg_first_seen
		,:eg_last_seen
		,:eg_occurrence_count
		,:eg_file_path
		,:eg_line_number
		,:eg_function_name
		,:eg_language
		,:eg_tags
		,:eg_assigned_to
		,:eg_resolved_at
		,:eg_resolved_by
		,:eg_created_by
		,:eg_created
		,:eg_updated
		,:eg_version
	)
	ON CONFLICT (eg_fingerprint) DO
	UPDATE SET
		 eg_last_seen = :eg_last_seen
		,eg_occurrence_count = eg_occurrence_count + 1
		,eg_status = CASE WHEN eg_status = 'resolved' THEN 'regressed' ELSE eg_status END
		,eg_updated = :eg_updated
		,eg_version = eg_version + 1
	RETURNING eg_id, eg_created_by, eg_created, eg_version`

	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(sqlQuery, mapInternalErrorGroup(errorGroup))
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind error group object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(
		&errorGroup.ID,
		&errorGroup.CreatedBy,
		&errorGroup.Created,
		&errorGroup.Version,
	); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Create or update error group query failed")
	}

	return nil
}

// CreateErrorOccurrence creates a new error occurrence.
func (s *ErrorTrackerStore) CreateErrorOccurrence(
	ctx context.Context,
	occurrence *types.ErrorOccurrence,
) error {
	const sqlQuery = `
	INSERT INTO error_occurrences (
		 eo_error_group_id
		,eo_stack_trace
		,eo_environment
		,eo_runtime
		,eo_os
		,eo_arch
		,eo_metadata
		,eo_created_at
	) VALUES (
		 :eo_error_group_id
		,:eo_stack_trace
		,:eo_environment
		,:eo_runtime
		,:eo_os
		,:eo_arch
		,:eo_metadata
		,:eo_created_at
	)
	RETURNING eo_id`

	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(sqlQuery, mapInternalErrorOccurrence(occurrence))
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind error occurrence object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&occurrence.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Create error occurrence query failed")
	}

	return nil
}

// FindByIdentifier returns an error group by identifier.
func (s *ErrorTrackerStore) FindByIdentifier(
	ctx context.Context,
	spaceID int64,
	identifier string,
) (*types.ErrorGroup, error) {
	const sqlQuery = errorGroupSelectBase + `
		WHERE eg_space_id = $1 AND eg_identifier = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(errorGroup)
	if err := db.GetContext(ctx, dst, sqlQuery, spaceID, identifier); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find error group")
	}

	return mapErrorGroup(dst), nil
}

// FindByFingerprint returns an error group by fingerprint.
func (s *ErrorTrackerStore) FindByFingerprint(
	ctx context.Context,
	spaceID int64,
	fingerprint string,
) (*types.ErrorGroup, error) {
	const sqlQuery = errorGroupSelectBase + `
		WHERE eg_space_id = $1 AND eg_fingerprint = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(errorGroup)
	if err := db.GetContext(ctx, dst, sqlQuery, spaceID, fingerprint); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find error group by fingerprint")
	}

	return mapErrorGroup(dst), nil
}

// Count returns the count of error groups matching the options.
func (s *ErrorTrackerStore) Count(
	ctx context.Context,
	spaceID int64,
	opts types.ErrorTrackerListOptions,
) (int64, error) {
	stmt := database.Builder.
		Select("count(*)").
		From("error_groups").
		Where("eg_space_id = ?", spaceID)

	stmt = s.applyFilter(stmt, opts)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to convert query to sql: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed to execute count error groups query")
	}

	return count, nil
}

// List returns a paginated list of error groups.
func (s *ErrorTrackerStore) List(
	ctx context.Context,
	spaceID int64,
	opts types.ErrorTrackerListOptions,
) ([]types.ErrorGroup, error) {
	stmt := database.Builder.
		Select(errorGroupColumns).
		From("error_groups").
		Where("eg_space_id = ?", spaceID)

	stmt = s.applyFilter(stmt, opts)

	stmt = stmt.
		Limit(database.Limit(opts.Size)).
		Offset(database.Offset(opts.Page, opts.Size)).
		OrderBy("eg_last_seen desc")

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to convert query to sql: %w", err)
	}

	dst := make([]*errorGroup, 0)
	db := dbtx.GetAccessor(ctx, s.db)

	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to execute list error groups query")
	}

	result := make([]types.ErrorGroup, len(dst))
	for i, eg := range dst {
		result[i] = *mapErrorGroup(eg)
	}

	return result, nil
}

// UpdateStatus updates the status of an error group.
func (s *ErrorTrackerStore) UpdateStatus(
	ctx context.Context,
	errorGroupID int64,
	status types.ErrorGroupStatus,
	updatedBy int64,
) error {
	const sqlQuery = `
	UPDATE error_groups
	SET eg_status = $1,
		eg_updated = $2,
		eg_version = eg_version + 1,
		eg_resolved_by = $3,
		eg_resolved_at = CASE WHEN $1 = 'resolved' THEN $2 ELSE eg_resolved_at END
	WHERE eg_id = $4`

	db := dbtx.GetAccessor(ctx, s.db)

	result, err := db.ExecContext(ctx, sqlQuery, status, types.NowMillis(), updatedBy, errorGroupID)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to update error group status")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to get rows affected")
	}

	if rows == 0 {
		return fmt.Errorf("error group not found")
	}

	return nil
}

// UpdateAssignment updates the assigned user of an error group.
func (s *ErrorTrackerStore) UpdateAssignment(
	ctx context.Context,
	errorGroupID int64,
	assignedTo *int64,
) error {
	const sqlQuery = `
	UPDATE error_groups
	SET eg_assigned_to = $1,
		eg_updated = $2,
		eg_version = eg_version + 1
	WHERE eg_id = $3`

	db := dbtx.GetAccessor(ctx, s.db)

	result, err := db.ExecContext(ctx, sqlQuery, assignedTo, types.NowMillis(), errorGroupID)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to update error group assignment")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to get rows affected")
	}

	if rows == 0 {
		return fmt.Errorf("error group not found")
	}

	return nil
}

// ListOccurrences returns a paginated list of error occurrences for an error group.
func (s *ErrorTrackerStore) ListOccurrences(
	ctx context.Context,
	errorGroupID int64,
	limit int,
	offset int,
) ([]types.ErrorOccurrence, error) {
	const sqlQuery = errorOccurrenceSelectBase + `
		WHERE eo_error_group_id = $1
		ORDER BY eo_created_at DESC
		LIMIT $2 OFFSET $3`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := make([]*errorOccurrence, 0)

	if err := db.SelectContext(ctx, &dst, sqlQuery, errorGroupID, limit, offset); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to list error occurrences")
	}

	result := make([]types.ErrorOccurrence, len(dst))
	for i, eo := range dst {
		result[i] = *mapErrorOccurrence(eo)
	}

	return result, nil
}

// GetSummary returns summary statistics for error groups in a space.
func (s *ErrorTrackerStore) GetSummary(
	ctx context.Context,
	spaceID int64,
) (*types.ErrorTrackerSummary, error) {
	const sqlQuery = `
	SELECT
		COUNT(*) as total,
		COALESCE(SUM(CASE WHEN eg_status = 'open' THEN 1 ELSE 0 END), 0) as open_count,
		COALESCE(SUM(CASE WHEN eg_status = 'resolved' THEN 1 ELSE 0 END), 0) as resolved_count,
		COALESCE(SUM(CASE WHEN eg_status = 'ignored' THEN 1 ELSE 0 END), 0) as ignored_count,
		COALESCE(SUM(CASE WHEN eg_status = 'regressed' THEN 1 ELSE 0 END), 0) as regression_count,
		COALESCE(SUM(CASE WHEN eg_severity = 'fatal' THEN 1 ELSE 0 END), 0) as fatal_count,
		COALESCE(SUM(CASE WHEN eg_severity = 'error' THEN 1 ELSE 0 END), 0) as error_count,
		COALESCE(SUM(CASE WHEN eg_severity = 'warning' THEN 1 ELSE 0 END), 0) as warning_count,
		COALESCE(MAX(eg_last_seen), 0) as last_updated
	FROM error_groups
	WHERE eg_space_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	summary := &types.ErrorTrackerSummary{}
	err := db.QueryRowContext(ctx, sqlQuery, spaceID).Scan(
		&summary.TotalErrors,
		&summary.OpenErrors,
		&summary.ResolvedErrors,
		&summary.IgnoredErrors,
		&summary.RegressionErrors,
		&summary.FatalCount,
		&summary.ErrorCount,
		&summary.WarningCount,
		&summary.LastUpdated,
	)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get error tracker summary")
	}

	return summary, nil
}

// applyFilter applies filter options to the query builder.
func (s *ErrorTrackerStore) applyFilter(
	stmt squirrel.SelectBuilder,
	opts types.ErrorTrackerListOptions,
) squirrel.SelectBuilder {
	if opts.Status != nil {
		stmt = stmt.Where("eg_status = ?", *opts.Status)
	}

	if opts.Severity != nil {
		stmt = stmt.Where("eg_severity = ?", *opts.Severity)
	}

	if opts.Language != nil {
		stmt = stmt.Where("eg_language = ?", *opts.Language)
	}

	if opts.Query != "" {
		searchTerm := "%" + opts.Query + "%"
		stmt = stmt.Where(
			"(eg_title ILIKE ? OR eg_message ILIKE ? OR eg_identifier ILIKE ?)",
			searchTerm, searchTerm, searchTerm,
		)
	}

	return stmt
}

// Helper functions for mapping between internal and external types.

func mapErrorGroup(eg *errorGroup) *types.ErrorGroup {
	return &types.ErrorGroup{
		ID:              eg.ID,
		SpaceID:         eg.SpaceID,
		RepoID:          eg.RepoID,
		Identifier:      eg.Identifier,
		Title:           eg.Title,
		Message:         eg.Message,
		Fingerprint:     eg.Fingerprint,
		Status:          eg.Status,
		Severity:        eg.Severity,
		FirstSeen:       eg.FirstSeen,
		LastSeen:        eg.LastSeen,
		OccurrenceCount: eg.OccurrenceCount,
		FilePath:        eg.FilePath,
		LineNumber:      eg.LineNumber,
		FunctionName:    eg.FunctionName,
		Language:        eg.Language,
		Tags:            eg.Tags,
		AssignedTo:      eg.AssignedTo,
		ResolvedAt:      eg.ResolvedAt,
		ResolvedBy:      eg.ResolvedBy,
		CreatedBy:       eg.CreatedBy,
		Created:         eg.Created,
		Updated:         eg.Updated,
		Version:         eg.Version,
	}
}

func mapInternalErrorGroup(eg *types.ErrorGroup) *errorGroup {
	return &errorGroup{
		ID:              eg.ID,
		SpaceID:         eg.SpaceID,
		RepoID:          eg.RepoID,
		Identifier:      eg.Identifier,
		Title:           eg.Title,
		Message:         eg.Message,
		Fingerprint:     eg.Fingerprint,
		Status:          eg.Status,
		Severity:        eg.Severity,
		FirstSeen:       eg.FirstSeen,
		LastSeen:        eg.LastSeen,
		OccurrenceCount: eg.OccurrenceCount,
		FilePath:        eg.FilePath,
		LineNumber:      eg.LineNumber,
		FunctionName:    eg.FunctionName,
		Language:        eg.Language,
		Tags:            eg.Tags,
		AssignedTo:      eg.AssignedTo,
		ResolvedAt:      eg.ResolvedAt,
		ResolvedBy:      eg.ResolvedBy,
		CreatedBy:       eg.CreatedBy,
		Created:         eg.Created,
		Updated:         eg.Updated,
		Version:         eg.Version,
	}
}

func mapErrorOccurrence(eo *errorOccurrence) *types.ErrorOccurrence {
	return &types.ErrorOccurrence{
		ID:          eo.ID,
		ErrorGroupID: eo.ErrorGroupID,
		StackTrace:  eo.StackTrace,
		Environment: eo.Environment,
		Runtime:     eo.Runtime,
		OS:          eo.OS,
		Arch:        eo.Arch,
		Metadata:    eo.Metadata,
		CreatedAt:   eo.CreatedAt,
	}
}

func mapInternalErrorOccurrence(eo *types.ErrorOccurrence) *errorOccurrence {
	return &errorOccurrence{
		ID:           eo.ID,
		ErrorGroupID: eo.ErrorGroupID,
		StackTrace:   eo.StackTrace,
		Environment:  eo.Environment,
		Runtime:      eo.Runtime,
		OS:           eo.OS,
		Arch:         eo.Arch,
		Metadata:     eo.Metadata,
		CreatedAt:    eo.CreatedAt,
	}
}
