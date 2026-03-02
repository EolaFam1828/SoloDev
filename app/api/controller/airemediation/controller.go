// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
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

package airemediation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/controller/space"
	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	airemediationevents "github.com/harness/gitness/app/events/airemediation"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/services/securityremediation"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// Controller implements the business logic for AI remediations.
type Controller struct {
	authorizer          authz.Authorizer
	spaceFinder         refcache.SpaceFinder
	repoFinder          refcache.RepoFinder
	remediationStore    store.RemediationStore
	scanResultStore     store.SecurityScanStore
	scanFindingStore    store.ScanFindingStore
	eventReporter       *airemediationevents.Reporter
	securityRemediation *securityremediation.Service
	aiAvailable         bool
}

// NewController returns a new Controller.
func NewController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	remediationStore store.RemediationStore,
	scanResultStore store.SecurityScanStore,
	scanFindingStore store.ScanFindingStore,
	eventReporter *airemediationevents.Reporter,
	securityRemediation *securityremediation.Service,
	aiAvailable bool,
) *Controller {
	return &Controller{
		authorizer:          authorizer,
		spaceFinder:         spaceFinder,
		repoFinder:          repoFinder,
		remediationStore:    remediationStore,
		scanResultStore:     scanResultStore,
		scanFindingStore:    scanFindingStore,
		eventReporter:       eventReporter,
		securityRemediation: securityRemediation,
		aiAvailable:         aiAvailable,
	}
}

// generateIdentifier creates a short unique identifier for a remediation.
func generateIdentifier() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "rem-" + hex.EncodeToString(b)
}

// TriggerRemediation creates a new AI remediation task.
func (c *Controller) TriggerRemediation(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *types.TriggerRemediationInput,
) (*types.Remediation, error) {
	if !c.aiAvailable {
		return nil, usererror.New(http.StatusServiceUnavailable, "AI remediation is not configured")
	}

	if in == nil {
		return nil, usererror.BadRequest("Request body cannot be empty")
	}

	sp, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	now := types.NowMillis()

	rem := &types.Remediation{
		SpaceID:       sp.ID,
		Identifier:    generateIdentifier(),
		Title:         in.Title,
		Description:   in.Description,
		Status:        types.RemediationStatusPending,
		TriggerSource: in.TriggerSource,
		TriggerRef:    in.TriggerRef,
		Branch:        in.Branch,
		CommitSHA:     in.CommitSHA,
		ErrorLog:      in.ErrorLog,
		SourceCode:    in.SourceCode,
		FilePath:      in.FilePath,
		AIModel:       in.AIModel,
		Metadata:      in.Metadata,
		CreatedBy:     session.Principal.ID,
		Created:       now,
		Updated:       now,
		Version:       1,
	}

	if err := c.remediationStore.Create(ctx, rem); err != nil {
		return nil, fmt.Errorf("failed to create remediation: %w", err)
	}

	// Report event
	if c.eventReporter != nil {
		c.eventReporter.RemediationTriggered(ctx, rem)
	}

	return rem, nil
}

func (c *Controller) TriggerRemediationFromSecurityFinding(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *types.CreateRemediationFromSecurityFindingInput,
) (*types.Remediation, bool, error) {
	if !c.aiAvailable || c.securityRemediation == nil {
		return nil, false, usererror.New(http.StatusServiceUnavailable, "AI remediation is not configured")
	}
	if in == nil {
		return nil, false, usererror.BadRequest("Request body cannot be empty")
	}
	if in.RepoRef == "" {
		return nil, false, usererror.BadRequest("repo_ref is required")
	}
	if in.ScanIdentifier == "" {
		return nil, false, usererror.BadRequest("scan_identifier is required")
	}
	if in.FindingID <= 0 {
		return nil, false, usererror.BadRequest("finding_id is required")
	}

	_, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, false, err
	}

	repo, err := c.getRepoCheckAccess(ctx, session, in.RepoRef, enum.PermissionRepoEdit)
	if err != nil {
		return nil, false, err
	}

	scan, err := c.scanResultStore.FindByIdentifier(ctx, repo.ID, in.ScanIdentifier)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find security scan: %w", err)
	}

	finding, err := c.scanFindingStore.Find(ctx, in.FindingID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find security finding: %w", err)
	}
	if finding.ScanID != scan.ID {
		return nil, false, usererror.BadRequest("finding does not belong to the specified scan")
	}

	rem, created, err := c.securityRemediation.CreateFromFinding(ctx, scan, finding, session.Principal.ID, false)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create security remediation: %w", err)
	}

	if created && c.eventReporter != nil {
		c.eventReporter.RemediationTriggered(ctx, rem)
	}

	return rem, created, nil
}

// ListRemediations lists remediations for a space.
func (c *Controller) ListRemediations(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	filter *types.RemediationListFilter,
) ([]*types.Remediation, error) {
	sp, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	return c.remediationStore.List(ctx, sp.ID, filter)
}

// GetRemediation retrieves a single remediation by identifier.
func (c *Controller) GetRemediation(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
) (*types.Remediation, error) {
	sp, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	rem, err := c.remediationStore.FindByIdentifier(ctx, sp.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find remediation: %w", err)
	}

	return rem, nil
}

// UpdateRemediation updates a remediation's status and AI results.
func (c *Controller) UpdateRemediation(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
	in *types.UpdateRemediationInput,
) (*types.Remediation, error) {
	sp, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	rem, err := c.remediationStore.FindByIdentifier(ctx, sp.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find remediation: %w", err)
	}

	// Apply updates
	if in.Status != nil {
		oldStatus := rem.Status
		rem.Status = *in.Status

		// If moving to applied, report the event
		if *in.Status == types.RemediationStatusApplied && oldStatus != types.RemediationStatusApplied {
			if c.eventReporter != nil {
				c.eventReporter.RemediationApplied(ctx, rem)
			}
		}

		// If moving to completed or failed, report completion
		if (*in.Status == types.RemediationStatusCompleted || *in.Status == types.RemediationStatusFailed) &&
			(oldStatus != types.RemediationStatusCompleted && oldStatus != types.RemediationStatusFailed) {
			if c.eventReporter != nil {
				c.eventReporter.RemediationCompleted(ctx, rem, *in.Status)
			}
		}
	}
	if in.PatchDiff != "" {
		rem.PatchDiff = in.PatchDiff
	}
	if in.AIResponse != "" {
		rem.AIResponse = in.AIResponse
	}
	if in.FixBranch != "" {
		rem.FixBranch = in.FixBranch
	}
	if in.PRLink != "" {
		rem.PRLink = in.PRLink
	}
	if in.Confidence != nil {
		rem.Confidence = *in.Confidence
	}

	if err := c.remediationStore.Update(ctx, rem); err != nil {
		return nil, fmt.Errorf("failed to update remediation: %w", err)
	}

	return rem, nil
}

// GetSummary returns aggregate remediation statistics for a space.
func (c *Controller) GetSummary(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
) (*types.RemediationSummary, error) {
	sp, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	return c.remediationStore.Summary(ctx, sp.ID)
}

func (c *Controller) AIAvailable() bool {
	return c.aiAvailable
}

// Helper function to get space and check access.
func (c *Controller) getSpaceCheckAccess(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	permission enum.Permission,
) (*types.SpaceCore, error) {
	return space.GetSpaceCheckAuth(ctx, c.spaceFinder, c.authorizer, session, spaceRef, permission)
}

func (c *Controller) getRepoCheckAccess(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	permission enum.Permission,
) (*types.RepositoryCore, error) {
	repo, err := c.repoFinder.FindByRef(ctx, repoRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find repo: %w", err)
	}

	if err := apiauth.CheckRepoState(ctx, session, repo, permission); err != nil {
		return nil, err
	}

	if err = apiauth.CheckRepo(ctx, c.authorizer, session, repo, permission); err != nil {
		return nil, fmt.Errorf("failed to verify authorization: %w", err)
	}

	return repo, nil
}

// Ensure apiauth import is used (for GetSpaceCheckAuth).
var _ = apiauth.CheckSpace
