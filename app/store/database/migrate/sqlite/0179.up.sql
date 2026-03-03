CREATE TABLE IF NOT EXISTS quality_rules (
    qr_id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    qr_space_id             INTEGER NOT NULL,
    qr_identifier           TEXT NOT NULL,
    qr_name                 TEXT NOT NULL,
    qr_description          TEXT,
    qr_category             TEXT NOT NULL,
    qr_enforcement          TEXT NOT NULL,
    qr_condition            TEXT NOT NULL,
    qr_target_repo_ids      TEXT,
    qr_target_branches      TEXT,
    qr_enabled              INTEGER NOT NULL DEFAULT 1,
    qr_tags                 TEXT,
    qr_created_by           INTEGER NOT NULL,
    qr_created              INTEGER NOT NULL,
    qr_updated              INTEGER NOT NULL,
    qr_version              INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (qr_space_id) REFERENCES spaces(space_id) ON DELETE CASCADE,
    UNIQUE(qr_space_id, qr_identifier)
);

CREATE INDEX idx_quality_rules_space_id ON quality_rules(qr_space_id);
CREATE INDEX idx_quality_rules_category ON quality_rules(qr_category);
CREATE INDEX idx_quality_rules_enforcement ON quality_rules(qr_enforcement);
CREATE INDEX idx_quality_rules_enabled ON quality_rules(qr_enabled);

CREATE TABLE IF NOT EXISTS quality_evaluations (
    qe_id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    qe_space_id             INTEGER NOT NULL,
    qe_repo_id              INTEGER NOT NULL,
    qe_identifier           TEXT NOT NULL,
    qe_commit_sha           TEXT NOT NULL,
    qe_branch               TEXT,
    qe_overall_status       TEXT NOT NULL,
    qe_rules_evaluated      INTEGER NOT NULL,
    qe_rules_passed         INTEGER NOT NULL,
    qe_rules_failed         INTEGER NOT NULL,
    qe_rules_warned         INTEGER NOT NULL,
    qe_results              TEXT,
    qe_triggered_by         TEXT NOT NULL,
    qe_pipeline_id          INTEGER,
    qe_duration_ms          INTEGER,
    qe_created_by           INTEGER NOT NULL,
    qe_created              INTEGER NOT NULL,
    FOREIGN KEY (qe_space_id) REFERENCES spaces(space_id) ON DELETE CASCADE,
    FOREIGN KEY (qe_repo_id) REFERENCES repositories(repo_id) ON DELETE CASCADE
);

CREATE INDEX idx_quality_evaluations_space_id ON quality_evaluations(qe_space_id);
CREATE INDEX idx_quality_evaluations_repo_id ON quality_evaluations(qe_repo_id);
CREATE INDEX idx_quality_evaluations_identifier ON quality_evaluations(qe_identifier);
CREATE INDEX idx_quality_evaluations_overall_status ON quality_evaluations(qe_overall_status);
CREATE INDEX idx_quality_evaluations_triggered_by ON quality_evaluations(qe_triggered_by);
CREATE INDEX idx_quality_evaluations_created ON quality_evaluations(qe_created);
