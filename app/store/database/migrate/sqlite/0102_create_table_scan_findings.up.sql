CREATE TABLE scan_findings (
 sf_id               INTEGER PRIMARY KEY AUTOINCREMENT
,sf_scan_id          INTEGER NOT NULL
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
,sf_created          INTEGER NOT NULL
,sf_updated          INTEGER NOT NULL

,UNIQUE(sf_scan_id, sf_identifier)

,CONSTRAINT fk_scan_findings_scan_id FOREIGN KEY (sf_scan_id)
    REFERENCES security_scans (ss_id) ON DELETE CASCADE
);

CREATE INDEX idx_scan_findings_scan_id ON scan_findings (sf_scan_id);
CREATE INDEX idx_scan_findings_severity ON scan_findings (sf_severity);
CREATE INDEX idx_scan_findings_status ON scan_findings (sf_status);
CREATE INDEX idx_scan_findings_created ON scan_findings (sf_created DESC);
