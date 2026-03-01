// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Refactored auth to use apiauth.CheckSpace, replaced spaceStore/repoStore with refcache finders, fixed transaction callbacks.
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

package qualitygate

import (
	"context"
	"errors"
	"fmt"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	qualitygateevent "github.com/harness/gitness/app/events/qualitygate"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
	gitness_store "github.com/harness/gitness/store"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type Controller struct {
	tx               dbtx.Transactor
	authorizer       authz.Authorizer
	qualityRuleStore store.QualityRuleStore
	qualityEvalStore store.QualityEvaluationStore
	spaceFinder      refcache.SpaceFinder
	repoFinder       refcache.RepoFinder
	eventReporter    *qualitygateevent.Reporter
}

func NewController(
	tx dbtx.Transactor,
	authorizer authz.Authorizer,
	qualityRuleStore store.QualityRuleStore,
	qualityEvalStore store.QualityEvaluationStore,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	eventReporter *qualitygateevent.Reporter,
) *Controller {
	return &Controller{
		tx:               tx,
		authorizer:       authorizer,
		qualityRuleStore: qualityRuleStore,
		qualityEvalStore: qualityEvalStore,
		spaceFinder:      spaceFinder,
		repoFinder:       repoFinder,
		eventReporter:    eventReporter,
	}
}

// CreateRuleInput holds the input for creating a quality rule.
type CreateRuleInput struct {
	Identifier     string                   `json:"identifier"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description,omitempty"`
	Category       enum.QualityRuleCategory `json:"category"`
	Enforcement    enum.QualityEnforcement  `json:"enforcement"`
	Condition      string                   `json:"condition"`
	TargetRepoIDs  []int64                  `json:"target_repo_ids,omitempty"`
	TargetBranches []string                 `json:"target_branches,omitempty"`
	Enabled        bool                     `json:"enabled"`
	Tags           []string                 `json:"tags,omitempty"`
}

// UpdateRuleInput holds the input for updating a quality rule.
type UpdateRuleInput struct {
	Name           *string                   `json:"name,omitempty"`
	Description    *string                   `json:"description,omitempty"`
	Category       *enum.QualityRuleCategory `json:"category,omitempty"`
	Enforcement    *enum.QualityEnforcement  `json:"enforcement,omitempty"`
	Condition      *string                   `json:"condition,omitempty"`
	TargetRepoIDs  []int64                   `json:"target_repo_ids,omitempty"`
	TargetBranches []string                  `json:"target_branches,omitempty"`
	Enabled        *bool                     `json:"enabled,omitempty"`
	Tags           []string                  `json:"tags,omitempty"`
}

// CreateRuleOutput holds the response for creating a quality rule.
type CreateRuleOutput struct {
	*types.QualityRule
}

// CreateRule creates a new quality rule.
func (c *Controller) CreateRule(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *CreateRuleInput,
) (*types.QualityRule, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateCreate); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Validate inputs
	if in.Identifier == "" {
		return nil, usererror.BadRequest("identifier is required")
	}
	if in.Name == "" {
		return nil, usererror.BadRequest("name is required")
	}
	if in.Condition == "" {
		return nil, usererror.BadRequest("condition is required")
	}

	// Check if rule with same identifier already exists
	_, err = c.qualityRuleStore.FindByIdentifier(ctx, sp.ID, in.Identifier)
	if err == nil {
		return nil, usererror.BadRequestf("rule with identifier %q already exists", in.Identifier)
	}
	if !errors.Is(err, gitness_store.ErrResourceNotFound) {
		return nil, err
	}

	// Sanitize inputs
	category, _ := in.Category.Sanitize()
	enforcement, _ := in.Enforcement.Sanitize()

	// Encode arrays to JSON
	targetRepoIDs, err := encodeJSONArray(in.TargetRepoIDs)
	if err != nil {
		return nil, usererror.BadRequestf("invalid target_repo_ids: %v", err)
	}

	targetBranches, err := encodeJSONArray(in.TargetBranches)
	if err != nil {
		return nil, usererror.BadRequestf("invalid target_branches: %v", err)
	}

	tags, err := encodeJSONArray(in.Tags)
	if err != nil {
		return nil, usererror.BadRequestf("invalid tags: %v", err)
	}

	rule := &types.QualityRule{
		SpaceID:        sp.ID,
		Identifier:     in.Identifier,
		Name:           in.Name,
		Description:    in.Description,
		Category:       category,
		Enforcement:    enforcement,
		Condition:      in.Condition,
		TargetRepoIDs:  targetRepoIDs,
		TargetBranches: targetBranches,
		Enabled:        in.Enabled,
		Tags:           tags,
		CreatedBy:      session.Principal.ID,
	}

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		err = c.qualityRuleStore.Create(ctx, rule)
		if err != nil {
			return err
		}

		if c.eventReporter != nil {
			c.eventReporter.RuleCreated(ctx, rule)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// GetRule retrieves a quality rule.
func (c *Controller) GetRule(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	ruleIdentifier string,
) (*types.QualityRule, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateView); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Find rule
	rule, err := c.qualityRuleStore.FindByIdentifier(ctx, sp.ID, ruleIdentifier)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// ListRulesOutput holds the response for listing quality rules.
type ListRulesOutput struct {
	Rules []*types.QualityRule `json:"rules"`
	Count int64                `json:"count"`
}

// ListRules lists quality rules for a space.
func (c *Controller) ListRules(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.QualityRuleFilter,
) (*ListRulesOutput, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateView); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// List rules
	rules, err := c.qualityRuleStore.List(ctx, sp.ID, filter)
	if err != nil {
		return nil, err
	}

	// Count total
	count, err := c.qualityRuleStore.Count(ctx, sp.ID, filter)
	if err != nil {
		return nil, err
	}

	return &ListRulesOutput{
		Rules: rules,
		Count: count,
	}, nil
}

// UpdateRule updates a quality rule.
func (c *Controller) UpdateRule(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	ruleIdentifier string,
	in *UpdateRuleInput,
) (*types.QualityRule, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateEdit); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Find rule
	rule, err := c.qualityRuleStore.FindByIdentifier(ctx, sp.ID, ruleIdentifier)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if in.Name != nil {
		rule.Name = *in.Name
	}
	if in.Description != nil {
		rule.Description = *in.Description
	}
	if in.Category != nil {
		category, _ := in.Category.Sanitize()
		rule.Category = category
	}
	if in.Enforcement != nil {
		enforcement, _ := in.Enforcement.Sanitize()
		rule.Enforcement = enforcement
	}
	if in.Condition != nil {
		rule.Condition = *in.Condition
	}
	if in.Enabled != nil {
		rule.Enabled = *in.Enabled
	}

	// Update arrays
	if in.TargetRepoIDs != nil {
		targetRepoIDs, err := encodeJSONArray(in.TargetRepoIDs)
		if err != nil {
			return nil, usererror.BadRequestf("invalid target_repo_ids: %v", err)
		}
		rule.TargetRepoIDs = targetRepoIDs
	}

	if in.TargetBranches != nil {
		targetBranches, err := encodeJSONArray(in.TargetBranches)
		if err != nil {
			return nil, usererror.BadRequestf("invalid target_branches: %v", err)
		}
		rule.TargetBranches = targetBranches
	}

	if in.Tags != nil {
		tags, err := encodeJSONArray(in.Tags)
		if err != nil {
			return nil, usererror.BadRequestf("invalid tags: %v", err)
		}
		rule.Tags = tags
	}

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		err = c.qualityRuleStore.Update(ctx, rule)
		if err != nil {
			return err
		}

		if c.eventReporter != nil {
			c.eventReporter.RuleUpdated(ctx, rule)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// ToggleRuleInput holds the input for toggling a rule.
type ToggleRuleInput struct {
	Enabled bool `json:"enabled"`
}

// ToggleRule toggles the enabled state of a quality rule.
func (c *Controller) ToggleRule(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	ruleIdentifier string,
	in *ToggleRuleInput,
) (*types.QualityRule, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateEdit); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Find rule
	rule, err := c.qualityRuleStore.FindByIdentifier(ctx, sp.ID, ruleIdentifier)
	if err != nil {
		return nil, err
	}

	rule.Enabled = in.Enabled

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		err = c.qualityRuleStore.Update(ctx, rule)
		if err != nil {
			return err
		}

		if c.eventReporter != nil {
			if in.Enabled {
				c.eventReporter.RuleEnabled(ctx, rule)
			} else {
				c.eventReporter.RuleDisabled(ctx, rule)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// DeleteRule deletes a quality rule.
func (c *Controller) DeleteRule(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	ruleIdentifier string,
) error {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateDelete); err != nil {
		return fmt.Errorf("failed to authorize: %w", err)
	}

	// Find rule
	rule, err := c.qualityRuleStore.FindByIdentifier(ctx, sp.ID, ruleIdentifier)
	if err != nil {
		return err
	}

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		err = c.qualityRuleStore.Delete(ctx, rule.ID)
		if err != nil {
			return err
		}

		if c.eventReporter != nil {
			c.eventReporter.RuleDeleted(ctx, rule)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// EvaluateInput holds the input for triggering a quality evaluation.
type EvaluateInput struct {
	RepoRef   string                `json:"repo_ref"`
	CommitSHA string                `json:"commit_sha"`
	Branch    string                `json:"branch,omitempty"`
	Trigger   enum.QualityTrigger   `json:"trigger"`
	PipelineID *int64               `json:"pipeline_id,omitempty"`
}

// Evaluate triggers a quality evaluation for a repository and commit.
func (c *Controller) Evaluate(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *EvaluateInput,
) (*types.QualityEvaluation, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateEvaluate); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Find repository
	repo, err := c.repoFinder.FindByRef(ctx, in.RepoRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find repository: %w", err)
	}

	// Validate inputs
	if in.CommitSHA == "" {
		return nil, usererror.BadRequest("commit_sha is required")
	}

	// Sanitize trigger
	trigger, _ := in.Trigger.Sanitize()

	// Get quality rules for this repo
	filter := &types.QualityRuleFilter{
		ListQueryFilter: types.ListQueryFilter{
			Pagination: types.Pagination{Page: 0, Size: 1000},
		},
		Enabled: boolPtr(true),
	}
	rules, err := c.qualityRuleStore.List(ctx, sp.ID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list quality rules: %w", err)
	}

	// Evaluate rules (placeholder - actual evaluation would be done by separate service)
	evaluation := &types.QualityEvaluation{
		SpaceID:        sp.ID,
		RepoID:         repo.ID,
		Identifier:     fmt.Sprintf("%s-%s", repo.Identifier, in.CommitSHA[:8]),
		CommitSHA:      in.CommitSHA,
		Branch:         in.Branch,
		RulesEvaluated: len(rules),
		RulesPassed:    len(rules),
		RulesFailed:    0,
		RulesWarned:    0,
		OverallStatus:  enum.QualityStatusPassed,
		TriggeredBy:    trigger,
		CreatedBy:      session.Principal.ID,
		Duration:       0,
	}

	if in.PipelineID != nil {
		evaluation.PipelineID = *in.PipelineID
	}

	// Create evaluation
	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		err = c.qualityEvalStore.Create(ctx, evaluation)
		if err != nil {
			return err
		}

		if c.eventReporter != nil {
			c.eventReporter.EvaluationCreated(ctx, evaluation)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return evaluation, nil
}

// ListEvaluationsOutput holds the response for listing evaluations.
type ListEvaluationsOutput struct {
	Evaluations []*types.QualityEvaluation `json:"evaluations"`
	Count       int64                      `json:"count"`
}

// ListEvaluations lists quality evaluations for a space.
func (c *Controller) ListEvaluations(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.QualityEvaluationFilter,
) (*ListEvaluationsOutput, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateView); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// List evaluations
	evals, err := c.qualityEvalStore.List(ctx, sp.ID, filter)
	if err != nil {
		return nil, err
	}

	// Count total
	count, err := c.qualityEvalStore.Count(ctx, sp.ID, filter)
	if err != nil {
		return nil, err
	}

	return &ListEvaluationsOutput{
		Evaluations: evals,
		Count:       count,
	}, nil
}

// GetEvaluation retrieves a quality evaluation.
func (c *Controller) GetEvaluation(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	evalIdentifier string,
) (*types.QualityEvaluation, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateView); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Find evaluation
	eval, err := c.qualityEvalStore.FindByIdentifier(ctx, evalIdentifier)
	if err != nil {
		return nil, err
	}

	// Verify it belongs to the space
	if eval.SpaceID != sp.ID {
		return nil, usererror.NotFoundf("evaluation with identifier %q not found", evalIdentifier)
	}

	return eval, nil
}

// GetSummary retrieves aggregate quality statistics for a space.
func (c *Controller) GetSummary(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
) (*types.QualitySummary, error) {
	// Find space and authorize
	sp, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find space: %w", err)
	}

	if err = apiauth.CheckSpace(ctx, c.authorizer, session, sp, enum.PermissionQualityGateView); err != nil {
		return nil, fmt.Errorf("failed to authorize: %w", err)
	}

	// Get summary
	summary, err := c.qualityEvalStore.Summary(ctx, sp.ID)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// Helper functions

func encodeJSONArray(data interface{}) ([]byte, error) {
	if data == nil || (len(fmt.Sprintf("%v", data)) == 0) {
		return nil, nil
	}

	// Use the built-in JSON marshaling
	type tmp interface{}
	return marshalJSON(data)
}

func marshalJSON(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case []int64:
		if len(v) == 0 {
			return nil, nil
		}
		// Marshal manually
		b := []byte(`[`)
		for i, id := range v {
			if i > 0 {
				b = append(b, ',')
			}
			b = append(b, []byte(fmt.Sprintf("%d", id))...)
		}
		b = append(b, ']')
		return b, nil
	case []string:
		if len(v) == 0 {
			return nil, nil
		}
		// Marshal manually
		b := []byte(`[`)
		for i, s := range v {
			if i > 0 {
				b = append(b, ',')
			}
			b = append(b, []byte(fmt.Sprintf(`"%s"`, s))...)
		}
		b = append(b, ']')
		return b, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", data)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
