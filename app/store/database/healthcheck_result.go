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

	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/store/database"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.HealthCheckResultStore = (*healthCheckResultStore)(nil)

const (
	healthCheckResultQueryBase = `
		SELECT` + healthCheckResultColumns + `
		FROM health_check_results`

	healthCheckResultColumns = `
	hcr_id,
	hcr_health_check_id,
	hcr_status,
	hcr_response_time,
	hcr_status_code,
	hcr_error_message,
	hcr_created_at
	`
)

func NewHealthCheckResultStore(db *sqlx.DB) store.HealthCheckResultStore {
	return &healthCheckResultStore{
		db: db,
	}
}

type healthCheckResultStore struct {
	db *sqlx.DB
}

func (s *healthCheckResultStore) Find(ctx context.Context, id int64) (*types.HealthCheckResult, error) {
	const findQueryStmt = healthCheckResultQueryBase + `
		WHERE hcr_id = $1`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(types.HealthCheckResult)
	if err := db.GetContext(ctx, dst, findQueryStmt, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to find health check result")
	}
	return dst, nil
}

func (s *healthCheckResultStore) Create(ctx context.Context, result *types.HealthCheckResult) error {
	const healthCheckResultInsertStmt = `
	INSERT INTO health_check_results (
		hcr_health_check_id,
		hcr_status,
		hcr_response_time,
		hcr_status_code,
		hcr_error_message,
		hcr_created_at
	) VALUES (
		:hcr_health_check_id,
		:hcr_status,
		:hcr_response_time,
		:hcr_status_code,
		:hcr_error_message,
		:hcr_created_at
	) RETURNING hcr_id`
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(healthCheckResultInsertStmt, result)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind health check result object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&result.ID); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "health check result query failed")
	}

	return nil
}

func (s *healthCheckResultStore) ListByHealthCheckID(ctx context.Context, healthCheckID int64, limit int) ([]*types.HealthCheckResult, error) {
	stmt := database.Builder.
		Select(healthCheckResultColumns).
		From("health_check_results").
		Where("hcr_health_check_id = ?", fmt.Sprint(healthCheckID)).
		OrderBy("hcr_created_at DESC").
		Limit(uint64(limit))

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	dst := []*types.HealthCheckResult{}
	if err = db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed executing list query")
	}

	return dst, nil
}

func (s *healthCheckResultStore) CountByStatus(ctx context.Context, healthCheckID int64, status string) (int64, error) {
	stmt := database.Builder.
		Select("count(*)").
		From("health_check_results").
		Where("hcr_health_check_id = ?", fmt.Sprint(healthCheckID)).
		Where("hcr_status = ?", status)

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

func (s *healthCheckResultStore) CountTotal(ctx context.Context, healthCheckID int64) (int64, error) {
	stmt := database.Builder.
		Select("count(*)").
		From("health_check_results").
		Where("hcr_health_check_id = ?", fmt.Sprint(healthCheckID))

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

func (s *healthCheckResultStore) GetAverageResponseTime(ctx context.Context, healthCheckID int64, hours int) (int64, error) {
	const getAvgStmt = `
	SELECT COALESCE(CAST(AVG(hcr_response_time) AS BIGINT), 0)
	FROM health_check_results
	WHERE hcr_health_check_id = $1
		AND hcr_created_at > EXTRACT(EPOCH FROM NOW() - INTERVAL '1' HOUR) * 1000 * $2`

	db := dbtx.GetAccessor(ctx, s.db)

	var avg int64
	err := db.QueryRowContext(ctx, getAvgStmt, healthCheckID, hours).Scan(&avg)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed executing average query")
	}
	return avg, nil
}

func (s *healthCheckResultStore) DeleteOlderThan(ctx context.Context, olderThanMillis int64) error {
	const deleteStmt = `
	DELETE FROM health_check_results
	WHERE hcr_created_at < $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, deleteStmt, olderThanMillis); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to delete old health check results")
	}

	return nil
}
