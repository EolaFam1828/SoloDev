CREATE INDEX IF NOT EXISTS idx_security_scans_status ON security_scans (ss_status, ss_created);
