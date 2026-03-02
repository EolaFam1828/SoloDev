// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
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

	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ store.RemediationStore = (*RemediationStore)(nil)

// NewRemediationStore returns a new RemediationStore.
func NewRemediationStore(db *sqlx.DB) *RemediationStore {
	return &RemediationStore{
		db: db,
	}
}

// RemediationStore implements store.RemediationStore backed by a relational database.
type RemediationStore struct {
	db *sqlx.DB
}

const (
	remediationColumns = `
		 rem_id
		,rem_space_id
		,rem_repo_id
		,rem_identifier
		,rem_title
		,rem_description
		,rem_status
		,rem_trigger_source
		,rem_trigger_ref
		,rem_branch
		,rem_commit_sha
		,rem_error_log
		,rem_source_code
		,rem_file_path
		,rem_ai_model
		,rem_ai_prompt
		,rem_ai_response
		,rem_patch_diff
		,rem_fix_branch
		,rem_pr_link
		,rem_confidence
		,rem_tokens_used
		,rem_duration_ms
		,rem_metadata
		,rem_created_by
		,rem_created
		,rem_updated
		,rem_version`
)

type remediation struct {
	ID            int64   `db:"rem_id"`
	SpaceID       int64   `db:"rem_space_id"`
	RepoID        int64   `db:"rem_repo_id"`
	Identifier    string  `db:"rem_identifier"`
	Title         string  `db:"rem_title"`
	Description   string  `db:"rem_description"`
	Status        string  `db:"rem_status"`
	TriggerSource string  `db:"rem_trigger_source"`
	TriggerRef    string  `db:"rem_trigger_ref"`
	Branch        string  `db:"rem_branch"`
	CommitSHA     string  `db:"rem_commit_sha"`
	ErrorLog      string  `db:"rem_error_log"`
	SourceCode    string  `db:"rem_source_code"`
	FilePath      string  `db:"rem_file_path"`
	AIModel       string  `db:"rem_ai_model"`
	AIPrompt      string  `db:"rem_ai_prompt"`
	AIResponse    string  `db:"rem_ai_response"`
	PatchDiff     string  `db:"rem_patch_diff"`
	FixBranch     string  `db:"rem_fix_branch"`
	PRLink        string  `db:"rem_pr_link"`
	Confidence    float64 `db:"rem_confidence"`
	TokensUsed    int64   `db:"rem_tokens_used"`
	DurationMs    int64   `db:"rem_duration_ms"`
	Metadata      []byte  `db:"rem_metadata"`
	CreatedBy     int64   `db:"rem_created_by"`
	Created       int64   `db:"rem_created"`
	Updated       int64   `db:"rem_updated"`
	Version       int64   `db:"rem_version"`
}

func mapRemediation(r *remediation) *types.Remediation {
	return &types.Remediation{
		ID:            r.ID,
		SpaceID:       r.SpaceID,
		RepoID:        r.RepoID,
		Identifier:    r.Identifier,
		Title:         r.Title,
		Description:   r.Description,
		Status:        types.RemediationStatus(r.Status),
		TriggerSource: types.RemediationTriggerSource(r.TriggerSource),
		TriggerRef:    r.TriggerRef,
		Branch:        r.Branch,
		CommitSHA:     r.CommitSHA,
		ErrorLog:      r.ErrorLog,
		SourceCode:    r.SourceCode,
		FilePath:      r.FilePath,
		AIModel:       r.AIModel,
		AIPrompt:      r.AIPrompt,
		AIResponse:    r.AIResponse,
		PatchDiff:     r.PatchDiff,
		FixBranch:     r.FixBranch,
		PRLink:        r.PRLink,
		Confidence:    r.Confidence,
		TokensUsed:    r.TokensUsed,
		DurationMs:    r.DurationMs,
		Metadata:      r.Metadata,
		CreatedBy:     r.CreatedBy,
		Created:       r.Created,
		Updated:       r.Updated,
		Version:       r.Version,
	}
}

// Create creates a new remediation.
func (s *RemediationStore) Create(ctx context.Context, rem *types.Remediation) error {
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := squirrel.
		Insert("remediations").
		Columns(
			"rem_space_id",
			"rem_repo_id",
			"rem_identifier",
			"rem_title",
			"rem_description",
			"rem_status",
			"rem_trigger_source",
			"rem_trigger_ref",
			"rem_branch",
			"rem_commit_sha",
			"rem_error_log",
			"rem_source_code",
			"rem_file_path",
			"rem_ai_model",
			"rem_ai_prompt",
			"rem_ai_response",
			"rem_patch_diff",
			"rem_fix_branch",
			"rem_pr_link",
			"rem_confidence",
			"rem_tokens_used",
			"rem_duration_ms",
			"rem_metadata",
			"rem_created_by",
			"rem_created",
			"rem_updated",
			"rem_version",
		).
		Values(
			rem.SpaceID,
			rem.RepoID,
			rem.Identifier,
			rem.Title,
			rem.Description,
			string(rem.Status),
			string(rem.TriggerSource),
			rem.TriggerRef,
			rem.Branch,
			rem.CommitSHA,
			rem.ErrorLog,
			rem.SourceCode,
			rem.FilePath,
			rem.AIModel,
			rem.AIPrompt,
			rem.AIResponse,
			rem.PatchDiff,
			rem.FixBranch,
			rem.PRLink,
			rem.Confidence,
			rem.TokensUsed,
			rem.DurationMs,
			[]byte(rem.Metadata),
			rem.CreatedBy,
			rem.Created,
			rem.Updated,
			rem.Version,
		).
		Suffix("RETURNING rem_id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	err = db.QueryRowContext(ctx, query, args...).Scan(&rem.ID)
	if err != nil {
		return fmt.Errorf("failed to insert remediation: %w", err)
	}

	return nil
}

// Find finds a remediation by id.
func (s *RemediationStore) Find(ctx context.Context, id int64) (*types.Remediation, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := squirrel.
		Select(remediationColumns).
		From("remediations").
		Where(squirrel.Eq{"rem_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row remediation
	if err := db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, fmt.Errorf("failed to find remediation: %w", err)
	}
	return mapRemediation(&row), nil
}

// FindByIdentifier finds a remediation by space id and identifier.
func (s *RemediationStore) FindByIdentifier(ctx context.Context, spaceID int64, identifier string) (*types.Remediation, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := squirrel.
		Select(remediationColumns).
		From("remediations").
		Where(squirrel.Eq{"rem_space_id": spaceID, "rem_identifier": identifier}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row remediation
	if err := db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, fmt.Errorf("failed to find remediation: %w", err)
	}
	return mapRemediation(&row), nil
}

// List lists remediations in a space with optional filtering.
func (s *RemediationStore) List(ctx context.Context, spaceID int64, filter *types.RemediationListFilter) ([]*types.Remediation, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	qb := squirrel.
		Select(remediationColumns).
		From("remediations").
		Where(squirrel.Eq{"rem_space_id": spaceID}).
		OrderBy("rem_created DESC").
		PlaceholderFormat(squirrel.Dollar)

	if filter != nil {
		if filter.Status != nil {
			qb = qb.Where(squirrel.Eq{"rem_status": string(*filter.Status)})
		}
		if filter.TriggerSource != nil {
			qb = qb.Where(squirrel.Eq{"rem_trigger_source": string(*filter.TriggerSource)})
		}
		if filter.Size > 0 {
			qb = qb.Limit(uint64(filter.Size))
		}
		if filter.Page > 0 && filter.Size > 0 {
			qb = qb.Offset(uint64(filter.Page * filter.Size))
		}
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list query: %w", err)
	}

	var rows []remediation
	if err := db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list remediations: %w", err)
	}

	result := make([]*types.Remediation, len(rows))
	for i := range rows {
		result[i] = mapRemediation(&rows[i])
	}
	return result, nil
}

// Count returns the count of remediations matching the filter.
func (s *RemediationStore) Count(ctx context.Context, spaceID int64, filter *types.RemediationListFilter) (int64, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	qb := squirrel.
		Select("COUNT(*)").
		From("remediations").
		Where(squirrel.Eq{"rem_space_id": spaceID}).
		PlaceholderFormat(squirrel.Dollar)

	if filter != nil {
		if filter.Status != nil {
			qb = qb.Where(squirrel.Eq{"rem_status": string(*filter.Status)})
		}
		if filter.TriggerSource != nil {
			qb = qb.Where(squirrel.Eq{"rem_trigger_source": string(*filter.TriggerSource)})
		}
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var count int64
	if err := db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, fmt.Errorf("failed to count remediations: %w", err)
	}
	return count, nil
}

// Update updates a remediation.
func (s *RemediationStore) Update(ctx context.Context, rem *types.Remediation) error {
	db := dbtx.GetAccessor(ctx, s.db)

	now := types.NowMillis()

	query, args, err := squirrel.
		Update("remediations").
		Set("rem_status", string(rem.Status)).
		Set("rem_ai_response", rem.AIResponse).
		Set("rem_patch_diff", rem.PatchDiff).
		Set("rem_fix_branch", rem.FixBranch).
		Set("rem_pr_link", rem.PRLink).
		Set("rem_confidence", rem.Confidence).
		Set("rem_tokens_used", rem.TokensUsed).
		Set("rem_duration_ms", rem.DurationMs).
		Set("rem_updated", now).
		Set("rem_version", rem.Version+1).
		Where(squirrel.Eq{"rem_id": rem.ID, "rem_version": rem.Version}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update remediation: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("optimistic lock: version conflict on remediation %d", rem.ID)
	}

	rem.Updated = now
	rem.Version++
	return nil
}

// UpdateStatus updates only the status of a remediation.
func (s *RemediationStore) UpdateStatus(ctx context.Context, id int64, status types.RemediationStatus) error {
	db := dbtx.GetAccessor(ctx, s.db)

	now := types.NowMillis()
	query, args, err := squirrel.
		Update("remediations").
		Set("rem_status", string(status)).
		Set("rem_updated", now).
		Where(squirrel.Eq{"rem_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update status query: %w", err)
	}

	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update remediation status: %w", err)
	}
	return nil
}

// Summary returns aggregate remediation statistics for a space.
func (s *RemediationStore) Summary(ctx context.Context, spaceID int64) (*types.RemediationSummary, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN rem_status = 'pending'    THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(SUM(CASE WHEN rem_status = 'processing' THEN 1 ELSE 0 END), 0) as processing,
			COALESCE(SUM(CASE WHEN rem_status = 'completed'  THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(CASE WHEN rem_status = 'applied'    THEN 1 ELSE 0 END), 0) as applied,
			COALESCE(SUM(CASE WHEN rem_status = 'failed'     THEN 1 ELSE 0 END), 0) as failed,
			COALESCE(SUM(CASE WHEN rem_status = 'dismissed'  THEN 1 ELSE 0 END), 0) as dismissed
		FROM remediations
		WHERE rem_space_id = $1`

	var summary types.RemediationSummary
	if err := db.GetContext(ctx, &summary, query, spaceID); err != nil {
		return nil, fmt.Errorf("failed to get remediation summary: %w", err)
	}
	return &summary, nil
}

// ListPendingGlobal lists pending remediations across all spaces ordered by creation time.
func (s *RemediationStore) ListPendingGlobal(ctx context.Context, limit int) ([]*types.Remediation, error) {
	if limit <= 0 {
		limit = 10
	}

	query := fmt.Sprintf(`SELECT %s FROM remediations WHERE rem_status = 'pending' ORDER BY rem_created ASC LIMIT %d`,
		remediationColumns, limit)

	db := dbtx.GetAccessor(ctx, s.db)
	var dsts []*remediation
	if err := db.SelectContext(ctx, &dsts, query); err != nil {
		return nil, fmt.Errorf("failed to list pending remediations: %w", err)
	}

	results := make([]*types.Remediation, len(dsts))
	for i := range dsts {
		results[i] = mapRemediation(dsts[i])
	}

	return results, nil
}
