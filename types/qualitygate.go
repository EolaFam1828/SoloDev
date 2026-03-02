// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"encoding/json"

	"github.com/EolaFam1828/SoloDev/types/enum"
)

// QualityRule represents an individual code quality rule/policy.
type QualityRule struct {
	ID              int64                    `json:"id"`
	SpaceID         int64                    `json:"space_id"`
	Identifier      string                   `json:"identifier"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description,omitempty"`
	Category        enum.QualityRuleCategory `json:"category"`
	Enforcement     enum.QualityEnforcement  `json:"enforcement"`
	Condition       string                   `json:"condition"`
	TargetRepoIDs   json.RawMessage          `json:"target_repo_ids"`
	TargetBranches  json.RawMessage          `json:"target_branches"`
	Enabled         bool                     `json:"enabled"`
	Tags            json.RawMessage          `json:"tags,omitempty"`
	CreatedBy       int64                    `json:"created_by"`
	Created         int64                    `json:"created"`
	Updated         int64                    `json:"updated"`
	Version         int64                    `json:"version"`
}

// QualityRuleFilter holds filter options for listing quality rules.
type QualityRuleFilter struct {
	ListQueryFilter
	Category   *enum.QualityRuleCategory
	Enforcement *enum.QualityEnforcement
	Enabled    *bool
}

// QualityEvaluation represents the result of evaluating rules against code.
type QualityEvaluation struct {
	ID              int64                 `json:"id"`
	SpaceID         int64                 `json:"space_id"`
	RepoID          int64                 `json:"repo_id"`
	Identifier      string                `json:"identifier"`
	CommitSHA       string                `json:"commit_sha"`
	Branch          string                `json:"branch"`
	OverallStatus   enum.QualityStatus    `json:"overall_status"`
	RulesEvaluated  int                   `json:"rules_evaluated"`
	RulesPassed     int                   `json:"rules_passed"`
	RulesFailed     int                   `json:"rules_failed"`
	RulesWarned     int                   `json:"rules_warned"`
	Results         json.RawMessage       `json:"results"`
	TriggeredBy     enum.QualityTrigger   `json:"triggered_by"`
	PipelineID      int64                 `json:"pipeline_id,omitempty"`
	Duration        int64                 `json:"duration_ms"`
	CreatedBy       int64                 `json:"created_by"`
	Created         int64                 `json:"created"`
}

// QualityEvaluationResult represents the result of a single rule evaluation.
type QualityEvaluationResult struct {
	RuleID        int64       `json:"rule_id"`
	RuleName      string      `json:"rule_name"`
	Status        string      `json:"status"` // "passed", "failed", "warning"
	ActualValue   interface{} `json:"actual_value,omitempty"`
	ExpectedValue interface{} `json:"expected_value,omitempty"`
	Message       string      `json:"message,omitempty"`
}

// QualityEvaluationFilter holds filter options for listing evaluations.
type QualityEvaluationFilter struct {
	ListQueryFilter
	OverallStatus *enum.QualityStatus
	TriggeredBy   *enum.QualityTrigger
}

// QualitySummary represents aggregate quality statistics for a space.
type QualitySummary struct {
	SpaceID             int64   `json:"space_id"`
	TotalRepositories   int     `json:"total_repositories"`
	RepositoriesPassed  int     `json:"repositories_passed"`
	RepositoriesFailed  int     `json:"repositories_failed"`
	RepositoriesWarned  int     `json:"repositories_warned"`
	AverageCoverage     float64 `json:"average_coverage,omitempty"`
	TotalEvaluations    int     `json:"total_evaluations"`
	FailedEvaluations   int     `json:"failed_evaluations"`
	LastEvaluationTime  int64   `json:"last_evaluation_time,omitempty"`
}
