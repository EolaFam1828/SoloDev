CREATE TABLE remediations (
    rem_id          SERIAL PRIMARY KEY,
    rem_space_id    INTEGER NOT NULL REFERENCES spaces(space_id) ON DELETE CASCADE,
    rem_repo_id     INTEGER NOT NULL DEFAULT 0,
    rem_identifier  TEXT    NOT NULL,
    rem_title       TEXT    NOT NULL,
    rem_description TEXT    NOT NULL DEFAULT '',
    rem_status      TEXT    NOT NULL DEFAULT 'pending',
    rem_trigger_source TEXT NOT NULL,
    rem_trigger_ref TEXT    NOT NULL DEFAULT '',
    rem_branch      TEXT    NOT NULL,
    rem_commit_sha  TEXT    NOT NULL DEFAULT '',
    rem_error_log   TEXT    NOT NULL,
    rem_source_code TEXT    NOT NULL DEFAULT '',
    rem_file_path   TEXT    NOT NULL DEFAULT '',
    rem_ai_model    TEXT    NOT NULL DEFAULT '',
    rem_ai_prompt   TEXT    NOT NULL DEFAULT '',
    rem_ai_response TEXT    NOT NULL DEFAULT '',
    rem_patch_diff  TEXT    NOT NULL DEFAULT '',
    rem_fix_branch  TEXT    NOT NULL DEFAULT '',
    rem_pr_link     TEXT    NOT NULL DEFAULT '',
    rem_confidence  REAL    NOT NULL DEFAULT 0,
    rem_tokens_used BIGINT  NOT NULL DEFAULT 0,
    rem_duration_ms BIGINT  NOT NULL DEFAULT 0,
    rem_metadata    TEXT,
    rem_created_by  INTEGER NOT NULL,
    rem_created     BIGINT  NOT NULL,
    rem_updated     BIGINT  NOT NULL,
    rem_version     BIGINT  NOT NULL DEFAULT 1,
    UNIQUE (rem_space_id, rem_identifier)
);

CREATE INDEX idx_remediations_space_status ON remediations (rem_space_id, rem_status);
CREATE INDEX idx_remediations_space_trigger ON remediations (rem_space_id, rem_trigger_source);
CREATE INDEX idx_remediations_created ON remediations (rem_created DESC);
