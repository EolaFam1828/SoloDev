// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed database column mapping and query builder.
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

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	featureFlagColumns = `
		 ff_space_id
		,ff_identifier
		,ff_name
		,ff_description
		,ff_kind
		,ff_default_on_variation
		,ff_default_off_variation
		,ff_enabled
		,ff_variations
		,ff_tags
		,ff_permanent
		,ff_created_by
		,ff_created
		,ff_updated
		,ff_version`

	featureFlagColumnsWithID = featureFlagColumns + `,ff_id`
)

type featureFlag struct {
	ID                  int64  `db:"ff_id"`
	SpaceID             int64  `db:"ff_space_id"`
	Identifier          string `db:"ff_identifier"`
	Name                string `db:"ff_name"`
	Description         string `db:"ff_description"`
	Kind                string `db:"ff_kind"`
	DefaultOnVariation  string `db:"ff_default_on_variation"`
	DefaultOffVariation string `db:"ff_default_off_variation"`
	Enabled             bool   `db:"ff_enabled"`
	Variations          string `db:"ff_variations"`
	Tags                string `db:"ff_tags"`
	Permanent           bool   `db:"ff_permanent"`
	CreatedBy           int64  `db:"ff_created_by"`
	Created             int64  `db:"ff_created"`
	Updated             int64  `db:"ff_updated"`
	Version             int64  `db:"ff_version"`
}

var _ store.FeatureFlagStore = (*FeatureFlagStore)(nil)

type FeatureFlagStore struct {
	db *sqlx.DB
}

func NewFeatureFlagStore(db *sqlx.DB) store.FeatureFlagStore {
	return &FeatureFlagStore{
		db: db,
	}
}

func (s *FeatureFlagStore) Create(ctx context.Context, featureFlag *types.FeatureFlag) error {
	const sqlQuery = `
		INSERT INTO feature_flags (` + featureFlagColumns + `)
		values (
			:ff_space_id
			,:ff_identifier
			,:ff_name
			,:ff_description
			,:ff_kind
			,:ff_default_on_variation
			,:ff_default_off_variation
			,:ff_enabled
			,:ff_variations
			,:ff_tags
			,:ff_permanent
			,:ff_created_by
			,:ff_created
			,:ff_updated
			,:ff_version
		)
		RETURNING ff_id`

	db := dbtx.GetAccessor(ctx, s.db)
	dbFeatureFlag := mapInternalFeatureFlag(featureFlag)

	query, args, err := db.BindNamed(sqlQuery, dbFeatureFlag)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind query")
	}

	if err = db.QueryRowContext(ctx, query, args...).Scan(&featureFlag.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to create feature flag")
	}

	return nil
}

func (s *FeatureFlagStore) Update(ctx context.Context, featureFlag *types.FeatureFlag) error {
	const sqlQuery = `
		UPDATE feature_flags
		SET
			 ff_name = :ff_name
			,ff_description = :ff_description
			,ff_kind = :ff_kind
			,ff_default_on_variation = :ff_default_on_variation
			,ff_default_off_variation = :ff_default_off_variation
			,ff_enabled = :ff_enabled
			,ff_variations = :ff_variations
			,ff_tags = :ff_tags
			,ff_permanent = :ff_permanent
			,ff_updated = :ff_updated
			,ff_version = ff_version + 1
		WHERE ff_id = :ff_id AND ff_version = :ff_version`

	dbFeatureFlag := mapInternalFeatureFlag(featureFlag)
	dbFeatureFlag.Updated = time.Now().UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)
	query, args, err := db.BindNamed(sqlQuery, dbFeatureFlag)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind query")
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to update feature flag")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to get rows affected")
	}

	if rowsAffected == 0 {
		return database.ProcessSQLErrorf(ctx, errors.New("version conflict"), "Feature flag version mismatch")
	}

	// Update the version in the object
	featureFlag.Version++
	featureFlag.Updated = dbFeatureFlag.Updated

	return nil
}

func (s *FeatureFlagStore) Find(ctx context.Context, id int64) (*types.FeatureFlag, error) {
	stmt := database.Builder.
		Select(featureFlagColumnsWithID).
		From("feature_flags").
		Where("ff_id = ?", id)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	var dst featureFlag
	if err := db.GetContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get feature flag")
	}

	return mapFeatureFlag(&dst), nil
}

func (s *FeatureFlagStore) FindByIdentifier(ctx context.Context, spaceID int64, identifier string) (*types.FeatureFlag, error) {
	stmt := database.Builder.
		Select(featureFlagColumnsWithID).
		From("feature_flags").
		Where("ff_space_id = ? AND ff_identifier = ?", spaceID, identifier)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	var dst featureFlag
	if err := db.GetContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get feature flag by identifier")
	}

	return mapFeatureFlag(&dst), nil
}

func (s *FeatureFlagStore) List(
	ctx context.Context,
	spaceID int64,
	filter *types.FeatureFlagFilter,
) ([]*types.FeatureFlag, error) {
	stmt := database.Builder.
		Select(featureFlagColumnsWithID).
		From("feature_flags").
		Where("ff_space_id = ?", spaceID).
		OrderBy("ff_created asc").
		Limit(database.Limit(filter.Size)).
		Offset(database.Offset(filter.Page, filter.Size))

	if filter.Query != "" {
		stmt = stmt.Where("LOWER(ff_identifier) LIKE '%' || LOWER(?) || '%' OR LOWER(ff_name) LIKE '%' || LOWER(?) || '%'", filter.Query, filter.Query)
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var dst []*featureFlag
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to list feature flags")
	}

	return mapFeatureFlags(dst), nil
}

func (s *FeatureFlagStore) Count(
	ctx context.Context,
	spaceID int64,
	filter *types.FeatureFlagFilter,
) (int64, error) {
	stmt := database.Builder.
		Select("COUNT(*)").
		From("feature_flags").
		Where("ff_space_id = ?", spaceID)

	if filter.Query != "" {
		stmt = stmt.Where("LOWER(ff_identifier) LIKE '%' || LOWER(?) || '%' OR LOWER(ff_name) LIKE '%' || LOWER(?) || '%'", filter.Query, filter.Query)
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	if err := db.GetContext(ctx, &count, sql, args...); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed to count feature flags")
	}

	return count, nil
}

func (s *FeatureFlagStore) Delete(ctx context.Context, id int64) error {
	stmt := database.Builder.
		Delete("feature_flags").
		Where("ff_id = ?", id)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, sql, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to delete feature flag")
	}

	return nil
}

func mapInternalFeatureFlag(in *types.FeatureFlag) *featureFlag {
	variationsJSON, _ := json.Marshal(in.Variations)
	tagsJSON, _ := json.Marshal(in.Tags)

	return &featureFlag{
		ID:                  in.ID,
		SpaceID:             in.SpaceID,
		Identifier:          in.Identifier,
		Name:                in.Name,
		Description:         in.Description,
		Kind:                in.Kind,
		DefaultOnVariation:  in.DefaultOnVariation,
		DefaultOffVariation: in.DefaultOffVariation,
		Enabled:             in.Enabled,
		Variations:          string(variationsJSON),
		Tags:                string(tagsJSON),
		Permanent:           in.Permanent,
		CreatedBy:           in.CreatedBy,
		Created:             in.Created,
		Updated:             in.Updated,
		Version:             in.Version,
	}
}

func mapFeatureFlag(internal *featureFlag) *types.FeatureFlag {
	variations := []types.Variation{}
	json.Unmarshal([]byte(internal.Variations), &variations)

	tags := []string{}
	json.Unmarshal([]byte(internal.Tags), &tags)

	return &types.FeatureFlag{
		ID:                  internal.ID,
		SpaceID:             internal.SpaceID,
		Identifier:          internal.Identifier,
		Name:                internal.Name,
		Description:         internal.Description,
		Kind:                internal.Kind,
		DefaultOnVariation:  internal.DefaultOnVariation,
		DefaultOffVariation: internal.DefaultOffVariation,
		Enabled:             internal.Enabled,
		Variations:          variations,
		Tags:                tags,
		Permanent:           internal.Permanent,
		CreatedBy:           internal.CreatedBy,
		Created:             internal.Created,
		Updated:             internal.Updated,
		Version:             internal.Version,
	}
}

func mapFeatureFlags(dbFeatureFlags []*featureFlag) []*types.FeatureFlag {
	result := make([]*types.FeatureFlag, len(dbFeatureFlags))

	for i, ff := range dbFeatureFlags {
		result[i] = mapFeatureFlag(ff)
	}

	return result
}
