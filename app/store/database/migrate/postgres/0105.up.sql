-- Copyright 2023 Harness, Inc.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

-- Create quality_rules table
CREATE TABLE IF NOT EXISTS quality_rules (
    qr_id                   BIGSERIAL PRIMARY KEY,
    qr_space_id             BIGINT NOT NULL,
    qr_identifier           VARCHAR(255) NOT NULL,
    qr_name                 VARCHAR(255) NOT NULL,
    qr_description          TEXT,
    qr_category             VARCHAR(50) NOT NULL,
    qr_enforcement          VARCHAR(50) NOT NULL,
    qr_condition            TEXT NOT NULL,
    qr_target_repo_ids      JSONB,
    qr_target_branches      JSONB,
    qr_enabled              BOOLEAN NOT NULL DEFAULT TRUE,
    qr_tags                 JSONB,
    qr_created_by           BIGINT NOT NULL,
    qr_created              BIGINT NOT NULL,
    qr_updated              BIGINT NOT NULL,
    qr_version              BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_quality_rules_space_id FOREIGN KEY (qr_space_id) REFERENCES spaces(space_id) ON DELETE CASCADE,
    CONSTRAINT uk_quality_rules_space_identifier UNIQUE(qr_space_id, qr_identifier)
);

CREATE INDEX idx_quality_rules_space_id ON quality_rules(qr_space_id);
CREATE INDEX idx_quality_rules_category ON quality_rules(qr_category);
CREATE INDEX idx_quality_rules_enforcement ON quality_rules(qr_enforcement);
CREATE INDEX idx_quality_rules_enabled ON quality_rules(qr_enabled);

-- Create quality_evaluations table
CREATE TABLE IF NOT EXISTS quality_evaluations (
    qe_id                   BIGSERIAL PRIMARY KEY,
    qe_space_id             BIGINT NOT NULL,
    qe_repo_id              BIGINT NOT NULL,
    qe_identifier           VARCHAR(255) NOT NULL,
    qe_commit_sha           VARCHAR(40) NOT NULL,
    qe_branch               VARCHAR(255),
    qe_overall_status       VARCHAR(50) NOT NULL,
    qe_rules_evaluated      INTEGER NOT NULL,
    qe_rules_passed         INTEGER NOT NULL,
    qe_rules_failed         INTEGER NOT NULL,
    qe_rules_warned         INTEGER NOT NULL,
    qe_results              JSONB,
    qe_triggered_by         VARCHAR(50) NOT NULL,
    qe_pipeline_id          BIGINT,
    qe_duration_ms          BIGINT,
    qe_created_by           BIGINT NOT NULL,
    qe_created              BIGINT NOT NULL,
    CONSTRAINT fk_quality_evaluations_space_id FOREIGN KEY (qe_space_id) REFERENCES spaces(space_id) ON DELETE CASCADE,
    CONSTRAINT fk_quality_evaluations_repo_id FOREIGN KEY (qe_repo_id) REFERENCES repositories(repo_id) ON DELETE CASCADE
);

CREATE INDEX idx_quality_evaluations_space_id ON quality_evaluations(qe_space_id);
CREATE INDEX idx_quality_evaluations_repo_id ON quality_evaluations(qe_repo_id);
CREATE INDEX idx_quality_evaluations_identifier ON quality_evaluations(qe_identifier);
CREATE INDEX idx_quality_evaluations_overall_status ON quality_evaluations(qe_overall_status);
CREATE INDEX idx_quality_evaluations_triggered_by ON quality_evaluations(qe_triggered_by);
CREATE INDEX idx_quality_evaluations_created ON quality_evaluations(qe_created);
