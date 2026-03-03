// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/store/database"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"

	"github.com/jmoiron/sqlx"
)

var _ store.SoloGateConfigStore = (*soloGateConfigStore)(nil)

type soloGateConfigStore struct {
	db *sqlx.DB
}

// ProvideSoloGateConfigStore creates a new solo gate config store.
func ProvideSoloGateConfigStore(db *sqlx.DB) store.SoloGateConfigStore {
	return &soloGateConfigStore{db: db}
}

type soloGateConfigDB struct {
	ID              int64  `db:"sgc_id"`
	SpaceID         int64  `db:"sgc_space_id"`
	EnforcementMode string `db:"sgc_enforcement_mode"`
	AutoIgnoreLow   bool   `db:"sgc_auto_ignore_low"`
	AutoTriageKnown bool   `db:"sgc_auto_triage_known"`
	AIAutoFix       bool   `db:"sgc_ai_auto_fix"`
	LogTechDebt     bool   `db:"sgc_log_tech_debt"`
	Created         int64  `db:"sgc_created"`
	Updated         int64  `db:"sgc_updated"`
}

func (d *soloGateConfigDB) toType() *types.SoloGateConfig {
	return &types.SoloGateConfig{
		ID:              d.ID,
		SpaceID:         d.SpaceID,
		EnforcementMode: types.EnforcementMode(d.EnforcementMode),
		AutoIgnoreLow:   d.AutoIgnoreLow,
		AutoTriageKnown: d.AutoTriageKnown,
		AIAutoFix:       d.AIAutoFix,
		LogTechDebt:     d.LogTechDebt,
		Created:         d.Created,
		Updated:         d.Updated,
	}
}

func (s *soloGateConfigStore) FindBySpaceID(ctx context.Context, spaceID int64) (*types.SoloGateConfig, error) {
	const query = `SELECT sgc_id, sgc_space_id, sgc_enforcement_mode,
		sgc_auto_ignore_low, sgc_auto_triage_known, sgc_ai_auto_fix, sgc_log_tech_debt,
		sgc_created, sgc_updated
		FROM solo_gate_configs WHERE sgc_space_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	var dst soloGateConfigDB
	if err := db.GetContext(ctx, &dst, query, spaceID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, database.ProcessSQLErrorf(ctx, err, "find solo gate config by space_id %d", spaceID)
	}
	return dst.toType(), nil
}

func (s *soloGateConfigStore) Upsert(ctx context.Context, config *types.SoloGateConfig) error {
	const query = `
		INSERT INTO solo_gate_configs (
			sgc_space_id, sgc_enforcement_mode, sgc_auto_ignore_low, sgc_auto_triage_known,
			sgc_ai_auto_fix, sgc_log_tech_debt, sgc_created, sgc_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (sgc_space_id) DO UPDATE SET
			sgc_enforcement_mode = EXCLUDED.sgc_enforcement_mode,
			sgc_auto_ignore_low = EXCLUDED.sgc_auto_ignore_low,
			sgc_auto_triage_known = EXCLUDED.sgc_auto_triage_known,
			sgc_ai_auto_fix = EXCLUDED.sgc_ai_auto_fix,
			sgc_log_tech_debt = EXCLUDED.sgc_log_tech_debt,
			sgc_updated = EXCLUDED.sgc_updated`

	now := types.NowMillis()
	if config.Created == 0 {
		config.Created = now
	}
	config.Updated = now

	db := dbtx.GetAccessor(ctx, s.db)
	_, err := db.ExecContext(ctx, query,
		config.SpaceID, string(config.EnforcementMode),
		config.AutoIgnoreLow, config.AutoTriageKnown,
		config.AIAutoFix, config.LogTechDebt,
		config.Created, config.Updated,
	)
	if err != nil {
		return fmt.Errorf("upsert solo gate config: %w", err)
	}
	return nil
}
