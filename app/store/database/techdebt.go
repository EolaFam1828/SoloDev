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

	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/store/database"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ store.TechDebtStore = (*TechDebtStore)(nil)

// NewTechDebtStore returns a new TechDebtStore.
func NewTechDebtStore(db *sqlx.DB) *TechDebtStore {
	return &TechDebtStore{
		db: db,
	}
}

// TechDebtStore implements store.TechDebtStore backed by a relational database.
type TechDebtStore struct {
	db *sqlx.DB
}

const (
	techDebtColumns = `
		 td_id
		,td_space_id
		,td_repo_id
		,td_identifier
		,td_title
		,td_description
		,td_severity
		,td_status
		,td_category
		,td_file_path
		,td_line_start
		,td_line_end
		,td_estimated_effort
		,td_tags
		,td_due_date
		,td_resolved_at
		,td_resolved_by
		,td_created_by
		,td_created
		,td_updated
		,td_version`

	techDebtSelectBase = `
    SELECT` + techDebtColumns + `
	FROM tech_debts`
)

type techDebt struct {
	ID               int64          `db:"td_id"`
	SpaceID          int64          `db:"td_space_id"`
	RepoID           int64          `db:"td_repo_id"`
	Identifier       string         `db:"td_identifier"`
	Title            string         `db:"td_title"`
	Description      string         `db:"td_description"`
	Severity         string         `db:"td_severity"`
	Status           string         `db:"td_status"`
	Category         string         `db:"td_category"`
	FilePath         string         `db:"td_file_path"`
	LineStart        int            `db:"td_line_start"`
	LineEnd          int            `db:"td_line_end"`
	EstimatedEffort  string         `db:"td_estimated_effort"`
	Tags             json.RawMessage `db:"td_tags"`
	DueDate          int64          `db:"td_due_date"`
	ResolvedAt       int64          `db:"td_resolved_at"`
	ResolvedBy       int64          `db:"td_resolved_by"`
	CreatedBy        int64          `db:"td_created_by"`
	Created          int64          `db:"td_created"`
	Updated          int64          `db:"td_updated"`
	Version          int64          `db:"td_version"`
}

// Create creates a new technical debt item.
func (s *TechDebtStore) Create(ctx context.Context, td *types.TechDebt) error {
	const sqlQuery = `
	INSERT INTO tech_debts (
		 td_space_id
		,td_repo_id
		,td_identifier
		,td_title
		,td_description
		,td_severity
		,td_status
		,td_category
		,td_file_path
		,td_line_start
		,td_line_end
		,td_estimated_effort
		,td_tags
		,td_due_date
		,td_resolved_at
		,td_resolved_by
		,td_created_by
		,td_created
		,td_updated
		,td_version
	) VALUES (
		 :td_space_id
		,:td_repo_id
		,:td_identifier
		,:td_title
		,:td_description
		,:td_severity
		,:td_status
		,:td_category
		,:td_file_path
		,:td_line_start
		,:td_line_end
		,:td_estimated_effort
		,:td_tags
		,:td_due_date
		,:td_resolved_at
		,:td_resolved_by
		,:td_created_by
		,:td_created
		,:td_updated
		,:td_version
	) RETURNING td_id`

	db := dbtx.GetAccessor(ctx, s.db)

	tags, err := json.Marshal(td.Tags)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to marshal tags")
	}

	dbTd := &techDebt{
		SpaceID:         td.SpaceID,
		RepoID:          td.RepoID,
		Identifier:      td.Identifier,
		Title:           td.Title,
		Description:     td.Description,
		Severity:        string(td.Severity),
		Status:          string(td.Status),
		Category:        string(td.Category),
		FilePath:        td.FilePath,
		LineStart:       td.LineStart,
		LineEnd:         td.LineEnd,
		EstimatedEffort: td.EstimatedEffort,
		Tags:            tags,
		DueDate:         td.DueDate,
		ResolvedAt:      td.ResolvedAt,
		ResolvedBy:      td.ResolvedBy,
		CreatedBy:       td.CreatedBy,
		Created:         td.Created,
		Updated:         td.Updated,
		Version:         td.Version,
	}

	if err := db.QueryRowContext(ctx, sqlQuery, dbTd).Scan(&td.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to create technical debt")
	}

	return nil
}

// Update updates an existing technical debt item.
func (s *TechDebtStore) Update(ctx context.Context, td *types.TechDebt) error {
	const sqlQuery = `
	UPDATE tech_debts SET
		 td_title = :td_title
		,td_description = :td_description
		,td_severity = :td_severity
		,td_status = :td_status
		,td_category = :td_category
		,td_file_path = :td_file_path
		,td_line_start = :td_line_start
		,td_line_end = :td_line_end
		,td_estimated_effort = :td_estimated_effort
		,td_tags = :td_tags
		,td_due_date = :td_due_date
		,td_resolved_at = :td_resolved_at
		,td_resolved_by = :td_resolved_by
		,td_updated = :td_updated
		,td_version = :td_version
	WHERE td_id = :td_id AND td_version = :td_version - 1`

	db := dbtx.GetAccessor(ctx, s.db)

	tags, err := json.Marshal(td.Tags)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to marshal tags")
	}

	dbTd := &techDebt{
		ID:              td.ID,
		Title:           td.Title,
		Description:     td.Description,
		Severity:        string(td.Severity),
		Status:          string(td.Status),
		Category:        string(td.Category),
		FilePath:        td.FilePath,
		LineStart:       td.LineStart,
		LineEnd:         td.LineEnd,
		EstimatedEffort: td.EstimatedEffort,
		Tags:            tags,
		DueDate:         td.DueDate,
		ResolvedAt:      td.ResolvedAt,
		ResolvedBy:      td.ResolvedBy,
		Updated:         td.Updated,
		Version:         td.Version,
	}

	result, err := db.NamedExecContext(ctx, sqlQuery, dbTd)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to update technical debt")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to get rows affected")
	}

	if rows == 0 {
		return database.ProcessSQLErrorf(ctx, fmt.Errorf("no rows affected"), "Failed to update technical debt, version mismatch")
	}

	return nil
}

// Find finds a technical debt item by id.
func (s *TechDebtStore) Find(ctx context.Context, id int64) (*types.TechDebt, error) {
	const sqlQuery = techDebtSelectBase + ` WHERE td_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(techDebt)
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find technical debt")
	}

	return mapTechDebt(dst), nil
}

// FindByIdentifier finds a technical debt item by space ID and identifier.
func (s *TechDebtStore) FindByIdentifier(ctx context.Context, spaceID int64, identifier string) (*types.TechDebt, error) {
	const sqlQuery = techDebtSelectBase + ` WHERE td_space_id = $1 AND td_identifier = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(techDebt)
	if err := db.GetContext(ctx, dst, sqlQuery, spaceID, identifier); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find technical debt by identifier")
	}

	return mapTechDebt(dst), nil
}

// List lists technical debt items based on the provided filter.
func (s *TechDebtStore) List(ctx context.Context, spaceID int64, filter *types.TechDebtFilter) ([]*types.TechDebt, error) {
	builder := squirrel.Select(techDebtColumns).
		From("tech_debts").
		Where("td_space_id = ?", spaceID).
		OrderBy("td_created DESC")

	if filter != nil {
		if len(filter.Severity) > 0 {
			builder = builder.Where(squirrel.Eq{"td_severity": filter.Severity})
		}
		if len(filter.Status) > 0 {
			builder = builder.Where(squirrel.Eq{"td_status": filter.Status})
		}
		if len(filter.Category) > 0 {
			builder = builder.Where(squirrel.Eq{"td_category": filter.Category})
		}
		if filter.RepoID > 0 {
			builder = builder.Where("td_repo_id = ?", filter.RepoID)
		}

		if filter.Limit > 0 {
			builder = builder.Limit(uint64(filter.Limit))
		}
		if filter.Page > 0 {
			builder = builder.Offset(uint64((filter.Page - 1) * filter.Limit))
		}
	}

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to build query")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var dbTds []techDebt
	if err := db.SelectContext(ctx, &dbTds, sqlQuery, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to list technical debts")
	}

	result := make([]*types.TechDebt, len(dbTds))
	for i, td := range dbTds {
		result[i] = mapTechDebt(&td)
	}

	return result, nil
}

// Count returns the count of technical debt items matching the filter.
func (s *TechDebtStore) Count(ctx context.Context, spaceID int64, filter *types.TechDebtFilter) (int64, error) {
	builder := squirrel.Select("COUNT(*)").
		From("tech_debts").
		Where("td_space_id = ?", spaceID)

	if filter != nil {
		if len(filter.Severity) > 0 {
			builder = builder.Where(squirrel.Eq{"td_severity": filter.Severity})
		}
		if len(filter.Status) > 0 {
			builder = builder.Where(squirrel.Eq{"td_status": filter.Status})
		}
		if len(filter.Category) > 0 {
			builder = builder.Where(squirrel.Eq{"td_category": filter.Category})
		}
		if filter.RepoID > 0 {
			builder = builder.Where("td_repo_id = ?", filter.RepoID)
		}
	}

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed to build query")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	if err := db.QueryRowContext(ctx, sqlQuery, args...).Scan(&count); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed to count technical debts")
	}

	return count, nil
}

// Delete deletes a technical debt item.
func (s *TechDebtStore) Delete(ctx context.Context, id int64) error {
	const sqlQuery = `DELETE FROM tech_debts WHERE td_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sqlQuery, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to delete technical debt")
	}

	return nil
}

// Summary returns aggregated statistics for technical debt items.
func (s *TechDebtStore) Summary(ctx context.Context, spaceID int64, filter *types.TechDebtFilter) (*types.TechDebtSummary, error) {
	summary := &types.TechDebtSummary{
		BySeverity: make(map[string]int),
		ByStatus:   make(map[string]int),
		ByCategory: make(map[string]int),
	}

	// Count by severity
	builder := squirrel.Select("td_severity", "COUNT(*) as count").
		From("tech_debts").
		Where("td_space_id = ?", spaceID).
		GroupBy("td_severity")

	if filter != nil && filter.RepoID > 0 {
		builder = builder.Where("td_repo_id = ?", filter.RepoID)
	}

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to build query")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to query severity summary")
	}
	defer rows.Close()

	for rows.Next() {
		var severity string
		var count int
		if err := rows.Scan(&severity, &count); err != nil {
			return nil, database.ProcessSQLErrorf(ctx, err, "Failed to scan severity")
		}
		summary.BySeverity[severity] = count
		summary.Total += count
	}

	// Count by status
	builder = squirrel.Select("td_status", "COUNT(*) as count").
		From("tech_debts").
		Where("td_space_id = ?", spaceID).
		GroupBy("td_status")

	if filter != nil && filter.RepoID > 0 {
		builder = builder.Where("td_repo_id = ?", filter.RepoID)
	}

	sqlQuery, args, err = builder.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to build query")
	}

	rows, err = db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to query status summary")
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, database.ProcessSQLErrorf(ctx, err, "Failed to scan status")
		}
		summary.ByStatus[status] = count
	}

	// Count by category
	builder = squirrel.Select("td_category", "COUNT(*) as count").
		From("tech_debts").
		Where("td_space_id = ?", spaceID).
		GroupBy("td_category")

	if filter != nil && filter.RepoID > 0 {
		builder = builder.Where("td_repo_id = ?", filter.RepoID)
	}

	sqlQuery, args, err = builder.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to build query")
	}

	rows, err = db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to query category summary")
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, database.ProcessSQLErrorf(ctx, err, "Failed to scan category")
		}
		summary.ByCategory[category] = count
	}

	return summary, nil
}

func mapTechDebt(dbTd *techDebt) *types.TechDebt {
	var tags []string
	if dbTd.Tags != nil {
		json.Unmarshal(dbTd.Tags, &tags)
	}
	if tags == nil {
		tags = []string{}
	}

	return &types.TechDebt{
		ID:              dbTd.ID,
		SpaceID:         dbTd.SpaceID,
		RepoID:          dbTd.RepoID,
		Identifier:      dbTd.Identifier,
		Title:           dbTd.Title,
		Description:     dbTd.Description,
		Severity:        types.TechDebtSeverity(dbTd.Severity),
		Status:          types.TechDebtStatus(dbTd.Status),
		Category:        types.TechDebtCategory(dbTd.Category),
		FilePath:        dbTd.FilePath,
		LineStart:       dbTd.LineStart,
		LineEnd:         dbTd.LineEnd,
		EstimatedEffort: dbTd.EstimatedEffort,
		Tags:            tags,
		DueDate:         dbTd.DueDate,
		ResolvedAt:      dbTd.ResolvedAt,
		ResolvedBy:      dbTd.ResolvedBy,
		CreatedBy:       dbTd.CreatedBy,
		Created:         dbTd.Created,
		Updated:         dbTd.Updated,
		Version:         dbTd.Version,
	}
}
