ALTER TABLE security_scans
ADD COLUMN ss_failure_reason TEXT NOT NULL DEFAULT '';
