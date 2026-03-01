CREATE TABLE health_checks (
    hc_id SERIAL PRIMARY KEY,
    hc_space_id INTEGER NOT NULL,
    hc_identifier TEXT NOT NULL,
    hc_name TEXT NOT NULL,
    hc_description TEXT NOT NULL DEFAULT '',
    hc_url TEXT NOT NULL,
    hc_method TEXT NOT NULL DEFAULT 'GET',
    hc_expected_status INTEGER NOT NULL DEFAULT 200,
    hc_interval_seconds INTEGER NOT NULL DEFAULT 300,
    hc_timeout_seconds INTEGER NOT NULL DEFAULT 10,
    hc_enabled BOOLEAN NOT NULL DEFAULT true,
    hc_headers TEXT NOT NULL DEFAULT '{}',
    hc_body TEXT NOT NULL DEFAULT '',
    hc_tags TEXT NOT NULL DEFAULT '[]',
    hc_last_status TEXT NOT NULL DEFAULT 'unknown',
    hc_last_checked_at BIGINT NOT NULL DEFAULT 0,
    hc_last_response_time BIGINT NOT NULL DEFAULT 0,
    hc_consecutive_failures INTEGER NOT NULL DEFAULT 0,
    hc_created_by INTEGER NOT NULL,
    hc_created BIGINT NOT NULL,
    hc_updated BIGINT NOT NULL,
    hc_version INTEGER NOT NULL DEFAULT 0,
    UNIQUE (hc_space_id, hc_identifier),
    CONSTRAINT fk_health_checks_space_id FOREIGN KEY (hc_space_id)
        REFERENCES spaces (space_id) ON DELETE CASCADE,
    CONSTRAINT fk_health_checks_created_by FOREIGN KEY (hc_created_by)
        REFERENCES principals (principal_id) ON DELETE NO ACTION
);

CREATE INDEX idx_health_checks_space_id ON health_checks(hc_space_id);
CREATE INDEX idx_health_checks_enabled ON health_checks(hc_enabled);
