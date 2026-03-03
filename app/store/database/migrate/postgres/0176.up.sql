CREATE TABLE security_scans (
 ss_id                   BIGSERIAL PRIMARY KEY
,ss_space_id             BIGINT NOT NULL
,ss_repo_id              BIGINT NOT NULL
,ss_identifier           TEXT NOT NULL
,ss_scan_type            TEXT NOT NULL
,ss_status               TEXT NOT NULL
,ss_commit_sha           TEXT NOT NULL
,ss_branch               TEXT NOT NULL
,ss_total_issues         INTEGER NOT NULL DEFAULT 0
,ss_critical_count       INTEGER NOT NULL DEFAULT 0
,ss_high_count           INTEGER NOT NULL DEFAULT 0
,ss_medium_count         INTEGER NOT NULL DEFAULT 0
,ss_low_count            INTEGER NOT NULL DEFAULT 0
,ss_duration             BIGINT NOT NULL DEFAULT 0
,ss_triggered_by         TEXT NOT NULL
,ss_created_by           BIGINT NOT NULL
,ss_created              BIGINT NOT NULL
,ss_updated              BIGINT NOT NULL
,ss_failure_reason       TEXT NOT NULL DEFAULT ''
,ss_version              BIGINT NOT NULL DEFAULT 0

,UNIQUE(ss_repo_id, ss_identifier)

,CONSTRAINT fk_security_scans_repo_id FOREIGN KEY (ss_repo_id)
    REFERENCES repositories (repo_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE

,CONSTRAINT fk_security_scans_space_id FOREIGN KEY (ss_space_id)
    REFERENCES spaces (space_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE

,CONSTRAINT fk_security_scans_created_by FOREIGN KEY (ss_created_by)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL
);

CREATE INDEX idx_security_scans_repo_id ON security_scans (ss_repo_id);
CREATE INDEX idx_security_scans_space_id ON security_scans (ss_space_id);
CREATE INDEX idx_security_scans_status ON security_scans (ss_status, ss_created);

CREATE TABLE scan_findings (
 sf_id               BIGSERIAL PRIMARY KEY
,sf_scan_id          BIGINT NOT NULL
,sf_identifier       TEXT NOT NULL
,sf_severity         TEXT NOT NULL
,sf_category         TEXT NOT NULL
,sf_title            TEXT NOT NULL
,sf_description      TEXT
,sf_file_path        TEXT NOT NULL
,sf_line_start       INTEGER NOT NULL
,sf_line_end         INTEGER NOT NULL
,sf_rule             TEXT NOT NULL
,sf_snippet          TEXT
,sf_suggestion       TEXT
,sf_status           TEXT NOT NULL DEFAULT 'open'
,sf_cwe              TEXT
,sf_created          BIGINT NOT NULL
,sf_updated          BIGINT NOT NULL

,UNIQUE(sf_scan_id, sf_identifier)

,CONSTRAINT fk_scan_findings_scan_id FOREIGN KEY (sf_scan_id)
    REFERENCES security_scans (ss_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);

CREATE INDEX idx_scan_findings_scan_id ON scan_findings (sf_scan_id);
CREATE INDEX idx_scan_findings_severity ON scan_findings (sf_severity);
CREATE INDEX idx_scan_findings_status ON scan_findings (sf_status);
CREATE INDEX idx_scan_findings_created ON scan_findings (sf_created DESC);
