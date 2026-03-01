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
	"fmt"
	"time"

	"github.com/harness/gitness/app/store"
	gitness_store "github.com/harness/gitness/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.HealthCheckStore = (*healthCheckStore)(nil)

const (
	healthCheckQueryBase = `
		SELECT` + healthCheckColumns + `
		FROM health_checks`

	healthCheckColumns = `
	hc_id,
	hc_space_id,
	hc_identifier,
	hc_name,
	hc_description,
	hc_url,
	hc_method,
	hc_expected_status,
	hc_interval_seconds,
	hc_timeout_seconds,
	hc_enabled,
	hc_headers,
	hc_body,
	hc_tags,
	hc_last_status,
	hc_last_checked_at,
	hc_last_response_time,
	hc_consecutive_failures,
	hc_created_by,
	hc_created,
	hc_updated,
	hc_version
	`
)

func NewHealthCheckStore(db *sqlx.DB) store.HealthCheckStore {
	return &healthCheckStore{
		db: db,
	}
}

type healthCheckStore struct {
	db *sqlx.DB
}

func (s *healthCheckStore) Find(ctx context.Context, id int64) (*types.HealthCheck, error) {
	const findQueryStmt = healthCheckQueryBase + `
		WHERE hc_id = $1`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(types.HealthCheck)
	if err := db.GetContext(ctx, dst, findQueryStmt, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find health check")
	}
	return dst, nil
}

func (s *healthCheckStore) FindByIdentifier(ctx context.Context, spaceID int64, identifier string) (*types.HealthCheck, error) {
	const findQueryStmt = healthCheckQueryBase + `
		WHERE hc_space_id = $1 AND hc_identifier = $2`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(types.HealthCheck)
	if err := db.GetContext(ctx, dst, findQueryStmt, spaceID, identifier); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find health check by identifier")
	}
	return dst, nil
}

func (s *healthCheckStore) Create(ctx context.Context, hc *types.HealthCheck) error {
	const healthCheckInsertStmt = `
	INSERT INTO health_checks (
		hc_space_id,
		hc_identifier,
		hc_name,
		hc_description,
		hc_url,
		hc_method,
		hc_expected_status,
		hc_interval_seconds,
		hc_timeout_seconds,
		hc_enabled,
		hc_headers,
		hc_body,
		hc_tags,
		hc_last_status,
		hc_last_checked_at,
		hc_last_response_time,
		hc_consecutive_failures,
		hc_created_by,
		hc_created,
		hc_updated,
		hc_version
	) VALUES (
		:hc_space_id,
		:hc_identifier,
		:hc_name,
		:hc_description,
		:hc_url,
		:hc_method,
		:hc_expected_status,
		:hc_interval_seconds,
		:hc_timeout_seconds,
		:hc_enabled,
		:hc_headers,
		:hc_body,
		:hc_tags,
		:hc_last_status,
		:hc_last_checked_at,
		:hc_last_response_time,
		:hc_consecutive_failures,
		:hc_created_by,
		:hc_created,
		:hc_updated,
		:hc_version
	) RETURNING hc_id`
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(healthCheckInsertStmt, hc)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind health check object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&hc.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "health check query failed")
	}

	return nil
}

func (s *healthCheckStore) Update(ctx context.Context, hc *types.HealthCheck) error {
	const healthCheckUpdateStmt = `
	UPDATE health_checks
	SET
		hc_identifier = :hc_identifier,
		hc_name = :hc_name,
		hc_description = :hc_description,
		hc_url = :hc_url,
		hc_method = :hc_method,
		hc_expected_status = :hc_expected_status,
		hc_interval_seconds = :hc_interval_seconds,
		hc_timeout_seconds = :hc_timeout_seconds,
		hc_enabled = :hc_enabled,
		hc_headers = :hc_headers,
		hc_body = :hc_body,
		hc_tags = :hc_tags,
		hc_last_status = :hc_last_status,
		hc_last_checked_at = :hc_last_checked_at,
		hc_last_response_time = :hc_last_response_time,
		hc_consecutive_failures = :hc_consecutive_failures,
		hc_updated = :hc_updated,
		hc_version = :hc_version
	WHERE hc_id = :hc_id AND hc_version = :hc_version - 1`
	updatedAt := time.Now()
	hcCopy := *hc

	hcCopy.Version++
	hcCopy.Updated = updatedAt.UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(healthCheckUpdateStmt, hcCopy)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind health check object")
	}

	result, err := db.ExecContext(ctx, query, arg...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to update health check")
	}

	count, err := result.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to get number of updated rows")
	}

	if count == 0 {
		return gitness_store.ErrVersionConflict
	}

	hc.Version = hcCopy.Version
	hc.Updated = hcCopy.Updated
	return nil
}

func (s *healthCheckStore) UpdateOptLock(ctx context.Context,
	hc *types.HealthCheck,
	mutateFn func(hc *types.HealthCheck) error,
) (*types.HealthCheck, error) {
	for {
		dup := *hc

		err := mutateFn(&dup)
		if err != nil {
			return nil, err
		}

		err = s.Update(ctx, &dup)
		if err == nil {
			return &dup, nil
		}
		if !errors.Is(err, gitness_store.ErrVersionConflict) {
			return nil, err
		}

		hc, err = s.Find(ctx, hc.ID)
		if err != nil {
			return nil, err
		}
	}
}

func (s *healthCheckStore) List(ctx context.Context, spaceID int64, filter types.ListQueryFilter) ([]*types.HealthCheck, error) {
	stmt := database.Builder.
		Select(healthCheckColumns).
		From("health_checks").
		Where("hc_space_id = ?", fmt.Sprint(spaceID))

	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("hc_identifier", filter.Query))
	}

	stmt = stmt.Limit(database.Limit(filter.Size))
	stmt = stmt.Offset(database.Offset(filter.Page, filter.Size))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dst := []*types.HealthCheck{}
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed executing custom list query")
	}

	return dst, nil
}

func (s *healthCheckStore) ListAll(ctx context.Context, spaceID int64) ([]*types.HealthCheck, error) {
	stmt := database.Builder.
		Select(healthCheckColumns).
		From("health_checks").
		Where("hc_space_id = ?", fmt.Sprint(spaceID))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dst := []*types.HealthCheck{}
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed executing custom list query")
	}

	return dst, nil
}

func (s *healthCheckStore) Delete(ctx context.Context, id int64) error {
	const healthCheckDeleteStmt = `
		DELETE FROM health_checks
		WHERE hc_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, healthCheckDeleteStmt, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Could not delete health check")
	}

	return nil
}

func (s *healthCheckStore) DeleteByIdentifier(ctx context.Context, spaceID int64, identifier string) error {
	const healthCheckDeleteStmt = `
	DELETE FROM health_checks
	WHERE hc_space_id = $1 AND hc_identifier = $2`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, healthCheckDeleteStmt, spaceID, identifier); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Could not delete health check by identifier")
	}

	return nil
}

func (s *healthCheckStore) Count(ctx context.Context, spaceID int64, filter types.ListQueryFilter) (int64, error) {
	stmt := database.Builder.
		Select("count(*)").
		From("health_checks").
		Where("hc_space_id = ?", spaceID)

	if filter.Query != "" {
		stmt = stmt.Where(PartialMatch("hc_identifier", filter.Query))
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed executing count query")
	}
	return count, nil
}
