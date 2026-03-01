CREATE TABLE feature_flags (
 ff_id INTEGER PRIMARY KEY AUTOINCREMENT
,ff_space_id INTEGER NOT NULL
,ff_identifier TEXT NOT NULL
,ff_name TEXT NOT NULL
,ff_description TEXT
,ff_kind TEXT NOT NULL
,ff_default_on_variation TEXT NOT NULL
,ff_default_off_variation TEXT NOT NULL
,ff_enabled BOOLEAN NOT NULL DEFAULT 0
,ff_variations TEXT NOT NULL DEFAULT '[]'
,ff_tags TEXT NOT NULL DEFAULT '[]'
,ff_permanent BOOLEAN NOT NULL DEFAULT 0
,ff_created_by INTEGER NOT NULL
,ff_created BIGINT NOT NULL
,ff_updated BIGINT NOT NULL
,ff_version BIGINT NOT NULL DEFAULT 0
,UNIQUE(ff_space_id, ff_identifier)
,CONSTRAINT fk_feature_flags_space_id FOREIGN KEY (ff_space_id)
    REFERENCES spaces (space_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_feature_flags_created_by FOREIGN KEY (ff_created_by)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);

CREATE INDEX idx_feature_flags_space_id ON feature_flags(ff_space_id);
CREATE INDEX idx_feature_flags_identifier ON feature_flags(ff_identifier);
