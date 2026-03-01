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
	"github.com/harness/gitness/types/enum"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	_ store.QualityRuleStore       = (*QualityRuleStore)(nil)
	_ store.QualityEvaluationStore = (*QualityEvaluationStore)(nil)
)

// NewQualityRuleStore returns a new QualityRuleStore.
func NewQualityRuleStore(db *sqlx.DB) *QualityRuleStore {
	return &QualityRuleStore{db: db}
}

// QualityRuleStore implements store.QualityRuleStore backed by a relational database.
type QualityRuleStore struct {
	db *sqlx.DB
}

const (
	qualityRuleColumns = `
		 qr_id
		,qr_space_id
		,qr_identifier
		,qr_name
		,qr_description
		,qr_category
		,qr_enforcement
		,qr_condition
		,qr_target_repo_ids
		,qr_target_branches
		,qr_enabled
		,qr_tags
		,qr_created_by
		,qr_created
		,qr_updated
		,qr_version`

	qualityRuleSelectBase = `
    SELECT` + qualityRuleColumns + `
	FROM quality_rules`
)

type qualityRule struct {
	ID             int64                    `db:"qr_id"`
	SpaceID        int64                    `db:"qr_space_id"`
	Identifier     string                   `db:"qr_identifier"`
	Name           string                   `db:"qr_name"`
	Description    string                   `db:"qr_description"`
	Category       enum.QualityRuleCategory `db:"qr_category"`
	Enforcement    enum.QualityEnforcement  `db:"qr_enforcement"`
	Condition      string                   `db:"qr_condition"`
	TargetRepoIDs  json.RawMessage          `db:"qr_target_repo_ids"`
	TargetBranches json.RawMessage          `db:"qr_target_branches"`
	Enabled        bool                     `db:"qr_enabled"`
	Tags           json.RawMessage          `db:"qr_tags"`
	CreatedBy      int64                    `db:"qr_created_by"`
	Created        int64                    `db:"qr_created"`
	Updated        int64                    `db:"qr_updated"`
	Version        int64                    `db:"qr_version"`
}

// Create creates a new quality rule.
func (s *QualityRuleStore) Create(ctx context.Context, in *types.QualityRule) error {
	in.Created = database.GetCurrentTimestamp()
	in.Updated = database.GetCurrentTimestamp()

	const sqlQuery = `
		INSERT INTO quality_rules (
			 qr_space_id
			,qr_identifier
			,qr_name
			,qr_description
			,qr_category
			,qr_enforcement
			,qr_condition
			,qr_target_repo_ids
			,qr_target_branches
			,qr_enabled
			,qr_tags
			,qr_created_by
			,qr_created
			,qr_updated
			,qr_version
		) VALUES (
			 :qr_space_id
			,:qr_identifier
			,:qr_name
			,:qr_description
			,:qr_category
			,:qr_enforcement
			,:qr_condition
			,:qr_target_repo_ids
			,:qr_target_branches
			,:qr_enabled
			,:qr_tags
			,:qr_created_by
			,:qr_created
			,:qr_updated
			,:qr_version
		) RETURNING qr_id`

	args := map[string]interface{}{
		"qr_space_id":        in.SpaceID,
		"qr_identifier":      in.Identifier,
		"qr_name":            in.Name,
		"qr_description":     in.Description,
		"qr_category":        in.Category,
		"qr_enforcement":     in.Enforcement,
		"qr_condition":       in.Condition,
		"qr_target_repo_ids": in.TargetRepoIDs,
		"qr_target_branches": in.TargetBranches,
		"qr_enabled":         in.Enabled,
		"qr_tags":            in.Tags,
		"qr_created_by":      in.CreatedBy,
		"qr_created":         in.Created,
		"qr_updated":         in.Updated,
		"qr_version":         in.Version,
	}

	query, values, err := squirrel.Dialect("postgres").
		Insert("quality_rules").
		Columns(
			"qr_space_id",
			"qr_identifier",
			"qr_name",
			"qr_description",
			"qr_category",
			"qr_enforcement",
			"qr_condition",
			"qr_target_repo_ids",
			"qr_target_branches",
			"qr_enabled",
			"qr_tags",
			"qr_created_by",
			"qr_created",
			"qr_updated",
			"qr_version",
		).
		Values(
			in.SpaceID,
			in.Identifier,
			in.Name,
			in.Description,
			in.Category,
			in.Enforcement,
			in.Condition,
			in.TargetRepoIDs,
			in.TargetBranches,
			in.Enabled,
			in.Tags,
			in.CreatedBy,
			in.Created,
			in.Updated,
			in.Version,
		).
		Suffix("RETURNING qr_id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	return s.db.QueryRowContext(ctx, query, values...).Scan(&in.ID)
}

// Update updates an existing quality rule.
func (s *QualityRuleStore) Update(ctx context.Context, in *types.QualityRule) error {
	in.Updated = database.GetCurrentTimestamp()
	in.Version++

	const sqlQuery = `
		UPDATE quality_rules SET
			 qr_name = :qr_name
			,qr_description = :qr_description
			,qr_category = :qr_category
			,qr_enforcement = :qr_enforcement
			,qr_condition = :qr_condition
			,qr_target_repo_ids = :qr_target_repo_ids
			,qr_target_branches = :qr_target_branches
			,qr_enabled = :qr_enabled
			,qr_tags = :qr_tags
			,qr_updated = :qr_updated
			,qr_version = :qr_version
		WHERE qr_id = :qr_id AND qr_version = :qr_old_version`

	query, values, err := squirrel.Dialect("postgres").
		Update("quality_rules").
		Set("qr_name", in.Name).
		Set("qr_description", in.Description).
		Set("qr_category", in.Category).
		Set("qr_enforcement", in.Enforcement).
		Set("qr_condition", in.Condition).
		Set("qr_target_repo_ids", in.TargetRepoIDs).
		Set("qr_target_branches", in.TargetBranches).
		Set("qr_enabled", in.Enabled).
		Set("qr_tags", in.Tags).
		Set("qr_updated", in.Updated).
		Set("qr_version", in.Version).
		Where(squirrel.Eq{"qr_id": in.ID, "qr_version": in.Version - 1}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return store.ErrVersionConflict
	}

	return nil
}

// Find finds a quality rule by id.
func (s *QualityRuleStore) Find(ctx context.Context, id int64) (*types.QualityRule, error) {
	query, values, err := squirrel.Dialect("postgres").
		Select(qualityRuleColumns).
		From("quality_rules").
		Where(squirrel.Eq{"qr_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var rule qualityRule
	err = s.db.GetContext(ctx, &rule, query, values...)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality rule")
	}

	return mapQualityRule(&rule), nil
}

// FindByIdentifier finds a quality rule by space id and identifier.
func (s *QualityRuleStore) FindByIdentifier(ctx context.Context, spaceID int64, identifier string) (*types.QualityRule, error) {
	query, values, err := squirrel.Dialect("postgres").
		Select(qualityRuleColumns).
		From("quality_rules").
		Where(squirrel.Eq{"qr_space_id": spaceID, "qr_identifier": identifier}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var rule qualityRule
	err = s.db.GetContext(ctx, &rule, query, values...)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality rule")
	}

	return mapQualityRule(&rule), nil
}

// List lists quality rules in a space.
func (s *QualityRuleStore) List(ctx context.Context, spaceID int64, filter *types.QualityRuleFilter) ([]*types.QualityRule, error) {
	query := squirrel.Dialect("postgres").
		Select(qualityRuleColumns).
		From("quality_rules").
		Where(squirrel.Eq{"qr_space_id": spaceID})

	if filter.Category != nil {
		query = query.Where(squirrel.Eq{"qr_category": *filter.Category})
	}

	if filter.Enforcement != nil {
		query = query.Where(squirrel.Eq{"qr_enforcement": *filter.Enforcement})
	}

	if filter.Enabled != nil {
		query = query.Where(squirrel.Eq{"qr_enabled": *filter.Enabled})
	}

	query = query.
		OrderBy("qr_created DESC").
		Limit(uint64(filter.Size)).
		Offset(uint64(filter.Page * filter.Size))

	sql, values, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var rules []qualityRule
	err = s.db.SelectContext(ctx, &rules, sql, values...)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality rules")
	}

	return mapQualityRules(rules), nil
}

// Count returns the count of quality rules in a space.
func (s *QualityRuleStore) Count(ctx context.Context, spaceID int64, filter *types.QualityRuleFilter) (int64, error) {
	query := squirrel.Dialect("postgres").
		Select("COUNT(*)").
		From("quality_rules").
		Where(squirrel.Eq{"qr_space_id": spaceID})

	if filter.Category != nil {
		query = query.Where(squirrel.Eq{"qr_category": *filter.Category})
	}

	if filter.Enforcement != nil {
		query = query.Where(squirrel.Eq{"qr_enforcement": *filter.Enforcement})
	}

	if filter.Enabled != nil {
		query = query.Where(squirrel.Eq{"qr_enabled": *filter.Enabled})
	}

	sql, values, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var count int64
	err = s.db.GetContext(ctx, &count, sql, values...)
	return count, database.ProcessError(ctx, err, "count quality rules")
}

// Delete deletes a quality rule.
func (s *QualityRuleStore) Delete(ctx context.Context, id int64) error {
	query, values, err := squirrel.Dialect("postgres").
		Delete("quality_rules").
		Where(squirrel.Eq{"qr_id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, values...)
	return database.ProcessError(ctx, err, "quality rule")
}

// NewQualityEvaluationStore returns a new QualityEvaluationStore.
func NewQualityEvaluationStore(db *sqlx.DB) *QualityEvaluationStore {
	return &QualityEvaluationStore{db: db}
}

// QualityEvaluationStore implements store.QualityEvaluationStore backed by a relational database.
type QualityEvaluationStore struct {
	db *sqlx.DB
}

const (
	qualityEvaluationColumns = `
		 qe_id
		,qe_space_id
		,qe_repo_id
		,qe_identifier
		,qe_commit_sha
		,qe_branch
		,qe_overall_status
		,qe_rules_evaluated
		,qe_rules_passed
		,qe_rules_failed
		,qe_rules_warned
		,qe_results
		,qe_triggered_by
		,qe_pipeline_id
		,qe_duration_ms
		,qe_created_by
		,qe_created`

	qualityEvaluationSelectBase = `
    SELECT` + qualityEvaluationColumns + `
	FROM quality_evaluations`
)

type qualityEvaluation struct {
	ID              int64               `db:"qe_id"`
	SpaceID         int64               `db:"qe_space_id"`
	RepoID          int64               `db:"qe_repo_id"`
	Identifier      string              `db:"qe_identifier"`
	CommitSHA       string              `db:"qe_commit_sha"`
	Branch          string              `db:"qe_branch"`
	OverallStatus   enum.QualityStatus  `db:"qe_overall_status"`
	RulesEvaluated  int                 `db:"qe_rules_evaluated"`
	RulesPassed     int                 `db:"qe_rules_passed"`
	RulesFailed     int                 `db:"qe_rules_failed"`
	RulesWarned     int                 `db:"qe_rules_warned"`
	Results         json.RawMessage     `db:"qe_results"`
	TriggeredBy     enum.QualityTrigger `db:"qe_triggered_by"`
	PipelineID      int64               `db:"qe_pipeline_id"`
	Duration        int64               `db:"qe_duration_ms"`
	CreatedBy       int64               `db:"qe_created_by"`
	Created         int64               `db:"qe_created"`
}

// Create creates a new quality evaluation.
func (s *QualityEvaluationStore) Create(ctx context.Context, in *types.QualityEvaluation) error {
	in.Created = database.GetCurrentTimestamp()

	query, values, err := squirrel.Dialect("postgres").
		Insert("quality_evaluations").
		Columns(
			"qe_space_id",
			"qe_repo_id",
			"qe_identifier",
			"qe_commit_sha",
			"qe_branch",
			"qe_overall_status",
			"qe_rules_evaluated",
			"qe_rules_passed",
			"qe_rules_failed",
			"qe_rules_warned",
			"qe_results",
			"qe_triggered_by",
			"qe_pipeline_id",
			"qe_duration_ms",
			"qe_created_by",
			"qe_created",
		).
		Values(
			in.SpaceID,
			in.RepoID,
			in.Identifier,
			in.CommitSHA,
			in.Branch,
			in.OverallStatus,
			in.RulesEvaluated,
			in.RulesPassed,
			in.RulesFailed,
			in.RulesWarned,
			in.Results,
			in.TriggeredBy,
			in.PipelineID,
			in.Duration,
			in.CreatedBy,
			in.Created,
		).
		Suffix("RETURNING qe_id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	return s.db.QueryRowContext(ctx, query, values...).Scan(&in.ID)
}

// Find finds a quality evaluation by id.
func (s *QualityEvaluationStore) Find(ctx context.Context, id int64) (*types.QualityEvaluation, error) {
	query, values, err := squirrel.Dialect("postgres").
		Select(qualityEvaluationColumns).
		From("quality_evaluations").
		Where(squirrel.Eq{"qe_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var eval qualityEvaluation
	err = s.db.GetContext(ctx, &eval, query, values...)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality evaluation")
	}

	return mapQualityEvaluation(&eval), nil
}

// FindByIdentifier finds a quality evaluation by identifier.
func (s *QualityEvaluationStore) FindByIdentifier(ctx context.Context, identifier string) (*types.QualityEvaluation, error) {
	query, values, err := squirrel.Dialect("postgres").
		Select(qualityEvaluationColumns).
		From("quality_evaluations").
		Where(squirrel.Eq{"qe_identifier": identifier}).
		OrderBy("qe_created DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var eval qualityEvaluation
	err = s.db.GetContext(ctx, &eval, query, values...)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality evaluation")
	}

	return mapQualityEvaluation(&eval), nil
}

// List lists quality evaluations for a space.
func (s *QualityEvaluationStore) List(ctx context.Context, spaceID int64, filter *types.QualityEvaluationFilter) ([]*types.QualityEvaluation, error) {
	query := squirrel.Dialect("postgres").
		Select(qualityEvaluationColumns).
		From("quality_evaluations").
		Where(squirrel.Eq{"qe_space_id": spaceID})

	if filter.OverallStatus != nil {
		query = query.Where(squirrel.Eq{"qe_overall_status": *filter.OverallStatus})
	}

	if filter.TriggeredBy != nil {
		query = query.Where(squirrel.Eq{"qe_triggered_by": *filter.TriggeredBy})
	}

	query = query.
		OrderBy("qe_created DESC").
		Limit(uint64(filter.Size)).
		Offset(uint64(filter.Page * filter.Size))

	sql, values, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var evals []qualityEvaluation
	err = s.db.SelectContext(ctx, &evals, sql, values...)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality evaluations")
	}

	return mapQualityEvaluations(evals), nil
}

// Count returns the count of quality evaluations for a space.
func (s *QualityEvaluationStore) Count(ctx context.Context, spaceID int64, filter *types.QualityEvaluationFilter) (int64, error) {
	query := squirrel.Dialect("postgres").
		Select("COUNT(*)").
		From("quality_evaluations").
		Where(squirrel.Eq{"qe_space_id": spaceID})

	if filter.OverallStatus != nil {
		query = query.Where(squirrel.Eq{"qe_overall_status": *filter.OverallStatus})
	}

	if filter.TriggeredBy != nil {
		query = query.Where(squirrel.Eq{"qe_triggered_by": *filter.TriggeredBy})
	}

	sql, values, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var count int64
	err = s.db.GetContext(ctx, &count, sql, values...)
	return count, database.ProcessError(ctx, err, "count quality evaluations")
}

// Summary returns aggregate quality statistics for a space.
func (s *QualityEvaluationStore) Summary(ctx context.Context, spaceID int64) (*types.QualitySummary, error) {
	query := `
		SELECT
			? AS space_id,
			COUNT(DISTINCT qe_repo_id) AS total_repositories,
			COUNT(DISTINCT CASE WHEN qe_overall_status = 'passed' THEN qe_repo_id END) AS repositories_passed,
			COUNT(DISTINCT CASE WHEN qe_overall_status = 'failed' THEN qe_repo_id END) AS repositories_failed,
			COUNT(DISTINCT CASE WHEN qe_overall_status = 'warning' THEN qe_repo_id END) AS repositories_warned,
			COUNT(*) AS total_evaluations,
			COUNT(CASE WHEN qe_overall_status = 'failed' THEN 1 END) AS failed_evaluations,
			MAX(qe_created) AS last_evaluation_time
		FROM quality_evaluations
		WHERE qe_space_id = ?`

	var summary types.QualitySummary
	err := s.db.GetContext(ctx, &summary, query, spaceID, spaceID)
	if err != nil {
		return nil, database.ProcessError(ctx, err, "quality summary")
	}

	return &summary, nil
}

// Helper functions

func mapQualityRule(in *qualityRule) *types.QualityRule {
	return &types.QualityRule{
		ID:             in.ID,
		SpaceID:        in.SpaceID,
		Identifier:     in.Identifier,
		Name:           in.Name,
		Description:    in.Description,
		Category:       in.Category,
		Enforcement:    in.Enforcement,
		Condition:      in.Condition,
		TargetRepoIDs:  in.TargetRepoIDs,
		TargetBranches: in.TargetBranches,
		Enabled:        in.Enabled,
		Tags:           in.Tags,
		CreatedBy:      in.CreatedBy,
		Created:        in.Created,
		Updated:        in.Updated,
		Version:        in.Version,
	}
}

func mapQualityRules(in []qualityRule) []*types.QualityRule {
	out := make([]*types.QualityRule, len(in))
	for i := range in {
		out[i] = mapQualityRule(&in[i])
	}
	return out
}

func mapQualityEvaluation(in *qualityEvaluation) *types.QualityEvaluation {
	return &types.QualityEvaluation{
		ID:             in.ID,
		SpaceID:        in.SpaceID,
		RepoID:         in.RepoID,
		Identifier:     in.Identifier,
		CommitSHA:      in.CommitSHA,
		Branch:         in.Branch,
		OverallStatus:  in.OverallStatus,
		RulesEvaluated: in.RulesEvaluated,
		RulesPassed:    in.RulesPassed,
		RulesFailed:    in.RulesFailed,
		RulesWarned:    in.RulesWarned,
		Results:        in.Results,
		TriggeredBy:    in.TriggeredBy,
		PipelineID:     in.PipelineID,
		Duration:       in.Duration,
		CreatedBy:      in.CreatedBy,
		Created:        in.Created,
	}
}

func mapQualityEvaluations(in []qualityEvaluation) []*types.QualityEvaluation {
	out := make([]*types.QualityEvaluation, len(in))
	for i := range in {
		out[i] = mapQualityEvaluation(&in[i])
	}
	return out
}
