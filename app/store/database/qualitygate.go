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
	"time"

	gitness_store "github.com/harness/gitness/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"

	"github.com/harness/gitness/app/store"

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
	now := time.Now().UnixMilli()
	in.Created = now
	in.Updated = now

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
			 $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING qr_id`

	db := dbtx.GetAccessor(ctx, s.db)
	if err := db.QueryRowContext(ctx, sqlQuery,
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
	).Scan(&in.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Insert quality rule failed")
	}

	return nil
}

// Update updates an existing quality rule.
func (s *QualityRuleStore) Update(ctx context.Context, in *types.QualityRule) error {
	in.Updated = time.Now().UnixMilli()
	in.Version++

	const sqlQuery = `
		UPDATE quality_rules SET
			 qr_name = $1
			,qr_description = $2
			,qr_category = $3
			,qr_enforcement = $4
			,qr_condition = $5
			,qr_target_repo_ids = $6
			,qr_target_branches = $7
			,qr_enabled = $8
			,qr_tags = $9
			,qr_updated = $10
			,qr_version = $11
		WHERE qr_id = $12 AND qr_version = $13`

	db := dbtx.GetAccessor(ctx, s.db)
	result, err := db.ExecContext(ctx, sqlQuery,
		in.Name,
		in.Description,
		in.Category,
		in.Enforcement,
		in.Condition,
		in.TargetRepoIDs,
		in.TargetBranches,
		in.Enabled,
		in.Tags,
		in.Updated,
		in.Version,
		in.ID,
		in.Version-1,
	)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Update quality rule failed")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Update quality rule rows affected")
	}

	if rows == 0 {
		return gitness_store.ErrVersionConflict
	}

	return nil
}

// Find finds a quality rule by id.
func (s *QualityRuleStore) Find(ctx context.Context, id int64) (*types.QualityRule, error) {
	const sqlQuery = qualityRuleSelectBase + ` WHERE qr_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := &qualityRule{}
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Find quality rule failed")
	}

	return mapQualityRule(dst), nil
}

// FindByIdentifier finds a quality rule by space id and identifier.
func (s *QualityRuleStore) FindByIdentifier(ctx context.Context, spaceID int64, identifier string) (*types.QualityRule, error) {
	const sqlQuery = qualityRuleSelectBase + ` WHERE qr_space_id = $1 AND qr_identifier = $2`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := &qualityRule{}
	if err := db.GetContext(ctx, dst, sqlQuery, spaceID, identifier); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Find quality rule by identifier failed")
	}

	return mapQualityRule(dst), nil
}

// List lists quality rules in a space.
func (s *QualityRuleStore) List(ctx context.Context, spaceID int64, filter *types.QualityRuleFilter) ([]*types.QualityRule, error) {
	const sqlQuery = qualityRuleSelectBase + `
		WHERE qr_space_id = $1
		ORDER BY qr_created DESC
		LIMIT $2 OFFSET $3`

	db := dbtx.GetAccessor(ctx, s.db)
	var rules []qualityRule
	if err := db.SelectContext(ctx, &rules, sqlQuery,
		spaceID, filter.Size, filter.Page*filter.Size); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List quality rules failed")
	}

	return mapQualityRules(rules), nil
}

// Count returns the count of quality rules in a space.
func (s *QualityRuleStore) Count(ctx context.Context, spaceID int64, filter *types.QualityRuleFilter) (int64, error) {
	const sqlQuery = `SELECT COUNT(*) FROM quality_rules WHERE qr_space_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)
	var count int64
	if err := db.GetContext(ctx, &count, sqlQuery, spaceID); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Count quality rules failed")
	}

	return count, nil
}

// Delete deletes a quality rule.
func (s *QualityRuleStore) Delete(ctx context.Context, id int64) error {
	const sqlQuery = `DELETE FROM quality_rules WHERE qr_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sqlQuery, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Delete quality rule failed")
	}

	return nil
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
	ID             int64               `db:"qe_id"`
	SpaceID        int64               `db:"qe_space_id"`
	RepoID         int64               `db:"qe_repo_id"`
	Identifier     string              `db:"qe_identifier"`
	CommitSHA      string              `db:"qe_commit_sha"`
	Branch         string              `db:"qe_branch"`
	OverallStatus  enum.QualityStatus  `db:"qe_overall_status"`
	RulesEvaluated int                 `db:"qe_rules_evaluated"`
	RulesPassed    int                 `db:"qe_rules_passed"`
	RulesFailed    int                 `db:"qe_rules_failed"`
	RulesWarned    int                 `db:"qe_rules_warned"`
	Results        json.RawMessage     `db:"qe_results"`
	TriggeredBy    enum.QualityTrigger `db:"qe_triggered_by"`
	PipelineID     int64               `db:"qe_pipeline_id"`
	Duration       int64               `db:"qe_duration_ms"`
	CreatedBy      int64               `db:"qe_created_by"`
	Created        int64               `db:"qe_created"`
}

// Create creates a new quality evaluation.
func (s *QualityEvaluationStore) Create(ctx context.Context, in *types.QualityEvaluation) error {
	in.Created = time.Now().UnixMilli()

	const sqlQuery = `
		INSERT INTO quality_evaluations (
			 qe_space_id
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
			,qe_created
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING qe_id`

	db := dbtx.GetAccessor(ctx, s.db)
	if err := db.QueryRowContext(ctx, sqlQuery,
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
	).Scan(&in.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Insert quality evaluation failed")
	}

	return nil
}

// Find finds a quality evaluation by id.
func (s *QualityEvaluationStore) Find(ctx context.Context, id int64) (*types.QualityEvaluation, error) {
	const sqlQuery = qualityEvaluationSelectBase + ` WHERE qe_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := &qualityEvaluation{}
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Find quality evaluation failed")
	}

	return mapQualityEvaluation(dst), nil
}

// FindByIdentifier finds a quality evaluation by identifier.
func (s *QualityEvaluationStore) FindByIdentifier(ctx context.Context, identifier string) (*types.QualityEvaluation, error) {
	const sqlQuery = qualityEvaluationSelectBase + `
		WHERE qe_identifier = $1
		ORDER BY qe_created DESC
		LIMIT 1`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := &qualityEvaluation{}
	if err := db.GetContext(ctx, dst, sqlQuery, identifier); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Find quality evaluation by identifier failed")
	}

	return mapQualityEvaluation(dst), nil
}

// List lists quality evaluations for a space.
func (s *QualityEvaluationStore) List(ctx context.Context, spaceID int64, filter *types.QualityEvaluationFilter) ([]*types.QualityEvaluation, error) {
	const sqlQuery = qualityEvaluationSelectBase + `
		WHERE qe_space_id = $1
		ORDER BY qe_created DESC
		LIMIT $2 OFFSET $3`

	db := dbtx.GetAccessor(ctx, s.db)
	var evals []qualityEvaluation
	if err := db.SelectContext(ctx, &evals, sqlQuery,
		spaceID, filter.Size, filter.Page*filter.Size); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List quality evaluations failed")
	}

	return mapQualityEvaluations(evals), nil
}

// Count returns the count of quality evaluations for a space.
func (s *QualityEvaluationStore) Count(ctx context.Context, spaceID int64, filter *types.QualityEvaluationFilter) (int64, error) {
	const sqlQuery = `SELECT COUNT(*) FROM quality_evaluations WHERE qe_space_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)
	var count int64
	if err := db.GetContext(ctx, &count, sqlQuery, spaceID); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Count quality evaluations failed")
	}

	return count, nil
}

// Summary returns aggregate quality statistics for a space.
func (s *QualityEvaluationStore) Summary(ctx context.Context, spaceID int64) (*types.QualitySummary, error) {
	const sqlQuery = `
		SELECT
			$1 AS space_id,
			COUNT(DISTINCT qe_repo_id) AS total_repositories,
			COUNT(DISTINCT CASE WHEN qe_overall_status = 'passed' THEN qe_repo_id END) AS repositories_passed,
			COUNT(DISTINCT CASE WHEN qe_overall_status = 'failed' THEN qe_repo_id END) AS repositories_failed,
			COUNT(DISTINCT CASE WHEN qe_overall_status = 'warning' THEN qe_repo_id END) AS repositories_warned,
			COUNT(*) AS total_evaluations,
			COUNT(CASE WHEN qe_overall_status = 'failed' THEN 1 END) AS failed_evaluations,
			MAX(qe_created) AS last_evaluation_time
		FROM quality_evaluations
		WHERE qe_space_id = $2`

	db := dbtx.GetAccessor(ctx, s.db)
	var summary types.QualitySummary
	if err := db.GetContext(ctx, &summary, sqlQuery, spaceID, spaceID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Quality summary failed")
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
