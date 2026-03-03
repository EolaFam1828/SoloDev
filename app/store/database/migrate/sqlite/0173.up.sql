CREATE TABLE IF NOT EXISTS solo_gate_configs (
    sgc_id              INTEGER PRIMARY KEY AUTOINCREMENT,
    sgc_space_id        INTEGER NOT NULL UNIQUE,
    sgc_enforcement_mode TEXT NOT NULL DEFAULT 'strict',
    sgc_auto_ignore_low  BOOLEAN NOT NULL DEFAULT 0,
    sgc_auto_triage_known BOOLEAN NOT NULL DEFAULT 0,
    sgc_ai_auto_fix      BOOLEAN NOT NULL DEFAULT 0,
    sgc_log_tech_debt    BOOLEAN NOT NULL DEFAULT 0,
    sgc_created          INTEGER NOT NULL DEFAULT 0,
    sgc_updated          INTEGER NOT NULL DEFAULT 0
);
