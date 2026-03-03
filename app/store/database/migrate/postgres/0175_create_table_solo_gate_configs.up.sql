CREATE TABLE IF NOT EXISTS solo_gate_configs (
    sgc_id              SERIAL PRIMARY KEY,
    sgc_space_id        INTEGER NOT NULL UNIQUE,
    sgc_enforcement_mode TEXT NOT NULL DEFAULT 'strict',
    sgc_auto_ignore_low  BOOLEAN NOT NULL DEFAULT FALSE,
    sgc_auto_triage_known BOOLEAN NOT NULL DEFAULT FALSE,
    sgc_ai_auto_fix      BOOLEAN NOT NULL DEFAULT FALSE,
    sgc_log_tech_debt    BOOLEAN NOT NULL DEFAULT FALSE,
    sgc_created          BIGINT NOT NULL DEFAULT 0,
    sgc_updated          BIGINT NOT NULL DEFAULT 0
);
