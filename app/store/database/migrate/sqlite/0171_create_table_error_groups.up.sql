CREATE TABLE error_groups (
 eg_id INTEGER PRIMARY KEY AUTOINCREMENT
,eg_space_id INTEGER NOT NULL
,eg_repo_id INTEGER NOT NULL
,eg_identifier TEXT NOT NULL
,eg_title TEXT NOT NULL
,eg_message TEXT NOT NULL
,eg_fingerprint TEXT NOT NULL UNIQUE
,eg_status TEXT NOT NULL
,eg_severity TEXT NOT NULL
,eg_first_seen INTEGER NOT NULL
,eg_last_seen INTEGER NOT NULL
,eg_occurrence_count INTEGER NOT NULL DEFAULT 1
,eg_file_path TEXT
,eg_line_number INTEGER
,eg_function_name TEXT
,eg_language TEXT
,eg_tags JSON
,eg_assigned_to INTEGER
,eg_resolved_at INTEGER
,eg_resolved_by INTEGER
,eg_created_by INTEGER NOT NULL
,eg_created INTEGER NOT NULL
,eg_updated INTEGER NOT NULL
,eg_version INTEGER NOT NULL DEFAULT 1
,CONSTRAINT fk_error_group_space_id FOREIGN KEY (eg_space_id)
    REFERENCES spaces (space_id)
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_error_group_repo_id FOREIGN KEY (eg_repo_id)
    REFERENCES repositories (repo_id)
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_error_group_created_by FOREIGN KEY (eg_created_by)
    REFERENCES principals (principal_id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
,CONSTRAINT fk_error_group_assigned_to FOREIGN KEY (eg_assigned_to)
    REFERENCES principals (principal_id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
,CONSTRAINT fk_error_group_resolved_by FOREIGN KEY (eg_resolved_by)
    REFERENCES principals (principal_id)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);

CREATE INDEX error_groups_space_id_status
    ON error_groups(eg_space_id, eg_status);

CREATE INDEX error_groups_space_id_severity
    ON error_groups(eg_space_id, eg_severity);

CREATE INDEX error_groups_space_id_last_seen
    ON error_groups(eg_space_id, eg_last_seen DESC);

CREATE INDEX error_groups_space_id_identifier
    ON error_groups(eg_space_id, eg_identifier);

CREATE TABLE error_occurrences (
 eo_id INTEGER PRIMARY KEY AUTOINCREMENT
,eo_error_group_id INTEGER NOT NULL
,eo_stack_trace TEXT NOT NULL
,eo_environment TEXT NOT NULL
,eo_runtime TEXT
,eo_os TEXT
,eo_arch TEXT
,eo_metadata JSON
,eo_created_at INTEGER NOT NULL
,CONSTRAINT fk_error_occurrence_group_id FOREIGN KEY (eo_error_group_id)
    REFERENCES error_groups (eg_id)
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);

CREATE INDEX error_occurrences_group_id_created
    ON error_occurrences(eo_error_group_id, eo_created_at DESC);

CREATE INDEX error_occurrences_environment
    ON error_occurrences(eo_environment);
