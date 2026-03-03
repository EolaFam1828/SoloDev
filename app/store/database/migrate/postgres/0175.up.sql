CREATE TABLE tech_debts (
td_id                 BIGSERIAL PRIMARY KEY
,td_space_id          BIGINT NOT NULL
,td_repo_id           BIGINT
,td_identifier        TEXT NOT NULL
,td_title             TEXT NOT NULL
,td_description       TEXT
,td_severity          TEXT NOT NULL
,td_status            TEXT NOT NULL
,td_category          TEXT NOT NULL
,td_file_path         TEXT
,td_line_start        INTEGER
,td_line_end          INTEGER
,td_estimated_effort  TEXT NOT NULL
,td_tags              JSONB
,td_due_date          BIGINT
,td_resolved_at       BIGINT
,td_resolved_by       BIGINT
,td_created_by        BIGINT NOT NULL
,td_created           BIGINT NOT NULL
,td_updated           BIGINT NOT NULL
,td_version           BIGINT NOT NULL DEFAULT 1

,UNIQUE(td_space_id, td_identifier)
);

CREATE INDEX idx_tech_debts_space_id ON tech_debts(td_space_id);
CREATE INDEX idx_tech_debts_repo_id ON tech_debts(td_repo_id);
CREATE INDEX idx_tech_debts_status ON tech_debts(td_status);
CREATE INDEX idx_tech_debts_severity ON tech_debts(td_severity);
CREATE INDEX idx_tech_debts_category ON tech_debts(td_category);
