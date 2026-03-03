CREATE TABLE health_check_results (
    hcr_id INTEGER PRIMARY KEY AUTOINCREMENT,
    hcr_health_check_id INTEGER NOT NULL,
    hcr_status TEXT NOT NULL,
    hcr_response_time INTEGER NOT NULL DEFAULT 0,
    hcr_status_code INTEGER NOT NULL DEFAULT 0,
    hcr_error_message TEXT NOT NULL DEFAULT '',
    hcr_created_at INTEGER NOT NULL,
    FOREIGN KEY (hcr_health_check_id) REFERENCES health_checks(hc_id) ON DELETE CASCADE
);

CREATE INDEX idx_health_check_results_health_check_id ON health_check_results(hcr_health_check_id);
CREATE INDEX idx_health_check_results_created_at ON health_check_results(hcr_created_at);
CREATE INDEX idx_health_check_results_status ON health_check_results(hcr_health_check_id, hcr_status);
