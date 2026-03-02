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

package securityscan

import (
	"context"
	"fmt"
	"time"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/controller/space"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/errors"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"

	"github.com/google/uuid"
)

type Controller struct {
	authorizer        authz.Authorizer
	spaceFinder       refcache.SpaceFinder
	repoFinder        refcache.RepoFinder
	scanResultStore   store.SecurityScanStore
	scanFindingStore  store.ScanFindingStore
	scannerService    ScannerService
}

// ScannerService is the interface used by the controller to trigger scan jobs.
type ScannerService interface {
	TriggerScanJob(ctx context.Context, scan *types.ScanResult) error
}

func NewController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	scanResultStore store.SecurityScanStore,
	scanFindingStore store.ScanFindingStore,
	scannerService ScannerService,
) *Controller {
	return &Controller{
		authorizer:       authorizer,
		spaceFinder:      spaceFinder,
		repoFinder:       repoFinder,
		scanResultStore:  scanResultStore,
		scanFindingStore: scanFindingStore,
		scannerService:   scannerService,
	}
}

// getRepoCheckAccess retrieves a repository and verifies access permissions.
func (c *Controller) getRepoCheckAccess(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	reqPermission enum.Permission,
) (*types.RepositoryCore, error) {
	if repoRef == "" {
		return nil, errors.InvalidArgument("A valid repository reference must be provided.")
	}

	repo, err := c.repoFinder.FindByRef(ctx, repoRef)
	if err != nil {
		return nil, fmt.Errorf("failed to find repo: %w", err)
	}

	if err := apiauth.CheckRepoState(ctx, session, repo, reqPermission); err != nil {
		return nil, err
	}

	if err = apiauth.CheckRepo(ctx, c.authorizer, session, repo, reqPermission); err != nil {
		return nil, fmt.Errorf("failed to verify authorization: %w", err)
	}

	return repo, nil
}

// getSpaceCheckAccess retrieves a space and verifies access permissions.
func (c *Controller) getSpaceCheckAccess(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	permission enum.Permission,
) (*types.SpaceCore, error) {
	return space.GetSpaceCheckAuth(ctx, c.spaceFinder, c.authorizer, session, spaceRef, permission)
}

// TriggerScan triggers a new security scan.
func (c *Controller) TriggerScan(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	in *types.ScanResultInput,
) (*types.ScanResult, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoPush)
	if err != nil {
		return nil, err
	}

	if in.ScanType == nil {
		return nil, errors.InvalidArgument("scan_type is required")
	}

	if _, valid := (*in.ScanType).Sanitize(); !valid {
		return nil, errors.InvalidArgument(fmt.Sprintf("scan_type '%s' is invalid", *in.ScanType))
	}

	if in.CommitSHA == nil || *in.CommitSHA == "" {
		return nil, errors.InvalidArgument("commit_sha is required")
	}

	branch := "main"
	if in.Branch != nil && *in.Branch != "" {
		branch = *in.Branch
	}

	triggeredBy := enum.SecurityScanTriggerManual
	if in.TriggeredBy != nil {
		triggeredBy = *in.TriggeredBy
	}

	now := time.Now().UnixMilli()
	scan := &types.ScanResult{
		SpaceID:   repo.ParentID,
		RepoID:    repo.ID,
		Identifier: uuid.New().String(),
		ScanType:  *in.ScanType,
		Status:    enum.SecurityScanStatusPending,
		CommitSHA: *in.CommitSHA,
		Branch:    branch,
		TriggeredBy: triggeredBy,
		CreatedBy: session.Principal.ID,
		Created:   now,
		Updated:   now,
		Version:   0,
	}

	if err := c.scanResultStore.Create(ctx, scan); err != nil {
		return nil, fmt.Errorf("failed to create scan: %w", err)
	}

	// Submit a background scan job if scanner service is available.
	if c.scannerService != nil {
		if err := c.scannerService.TriggerScanJob(ctx, scan); err != nil {
			// Log but don't fail — the scan record was created and the poller will pick it up.
			_ = err
		}
	}

	return scan, nil
}

// FindScan finds a security scan.
func (c *Controller) FindScan(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	scanIdentifier string,
) (*types.ScanResult, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoView)
	if err != nil {
		return nil, err
	}

	scan, err := c.scanResultStore.FindByIdentifier(ctx, repo.ID, scanIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find scan: %w", err)
	}

	return scan, nil
}

// ListScans lists security scans for a repository.
func (c *Controller) ListScans(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	filter *types.ScanResultFilter,
) ([]*types.ScanResult, int64, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoView)
	if err != nil {
		return nil, 0, err
	}

	if filter == nil {
		filter = &types.ScanResultFilter{}
	}

	scans, count, err := c.scanResultStore.List(ctx, repo.ID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list scans: %w", err)
	}

	return scans, count, nil
}

// UpdateScan updates a security scan.
func (c *Controller) UpdateScan(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	scanIdentifier string,
	in *types.ScanResult,
) (*types.ScanResult, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoPush)
	if err != nil {
		return nil, err
	}

	scan, err := c.scanResultStore.FindByIdentifier(ctx, repo.ID, scanIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find scan: %w", err)
	}

	if in.Status != "" {
		scan.Status = in.Status
	}

	if in.TotalIssues > 0 {
		scan.TotalIssues = in.TotalIssues
	}

	if in.CriticalCount >= 0 {
		scan.CriticalCount = in.CriticalCount
	}

	if in.HighCount >= 0 {
		scan.HighCount = in.HighCount
	}

	if in.MediumCount >= 0 {
		scan.MediumCount = in.MediumCount
	}

	if in.LowCount >= 0 {
		scan.LowCount = in.LowCount
	}

	if in.Duration > 0 {
		scan.Duration = in.Duration
	}

	if err := c.scanResultStore.Update(ctx, scan); err != nil {
		return nil, fmt.Errorf("failed to update scan: %w", err)
	}

	return scan, nil
}

// FindFinding finds a specific finding.
func (c *Controller) FindFinding(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	findingID int64,
) (*types.ScanFinding, error) {
	_, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	finding, err := c.scanFindingStore.Find(ctx, findingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find finding: %w", err)
	}

	return finding, nil
}

// ListFindings lists findings for a scan.
func (c *Controller) ListFindings(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	scanIdentifier string,
	filter *types.ScanFindingFilter,
) ([]*types.ScanFinding, int64, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoView)
	if err != nil {
		return nil, 0, err
	}

	scan, err := c.scanResultStore.FindByIdentifier(ctx, repo.ID, scanIdentifier)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find scan: %w", err)
	}

	if filter == nil {
		filter = &types.ScanFindingFilter{}
	}

	findings, count, err := c.scanFindingStore.ListByScan(ctx, scan.ID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list findings: %w", err)
	}

	return findings, count, nil
}

// UpdateFindingStatus updates the status of a finding.
func (c *Controller) UpdateFindingStatus(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	findingID int64,
	in *types.ScanFindingStatusUpdate,
) (*types.ScanFinding, error) {
	_, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	if in.Status == nil {
		return nil, errors.InvalidArgument("status is required")
	}

	if _, valid := (*in.Status).Sanitize(); !valid {
		return nil, errors.InvalidArgument(fmt.Sprintf("status '%s' is invalid", *in.Status))
	}

	finding, err := c.scanFindingStore.Find(ctx, findingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find finding: %w", err)
	}

	finding.Status = *in.Status

	if err := c.scanFindingStore.Update(ctx, finding); err != nil {
		return nil, fmt.Errorf("failed to update finding: %w", err)
	}

	return finding, nil
}

// GetSecuritySummary retrieves the security posture summary for a space or repository.
func (c *Controller) GetSecuritySummary(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	repoRef *string,
) (*types.SecuritySummary, error) {
	spaceCore, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	summary := &types.SecuritySummary{
		SpaceID: spaceCore.ID,
	}

	if repoRef != nil && *repoRef != "" {
		repo, err := c.getRepoCheckAccess(ctx, session, *repoRef, enum.PermissionRepoView)
		if err != nil {
			return nil, err
		}

		summary.RepoID = repo.ID

		// Get latest scan for this repo
		filter := &types.ScanResultFilter{
			Page: 1,
			Size: 1,
			Sort: enum.SecurityScanAttrCreated,
			Order: enum.OrderDesc,
		}

		scans, _, err := c.scanResultStore.List(ctx, repo.ID, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest scan: %w", err)
		}

		if len(scans) > 0 {
			lastScan := scans[0]
			summary.LastScanID = lastScan.ID
			summary.LastScanTime = lastScan.Created
			summary.TotalFindings = lastScan.TotalIssues
			summary.CriticalIssues = lastScan.CriticalCount
			summary.HighIssues = lastScan.HighCount
			summary.MediumIssues = lastScan.MediumCount
			summary.LowIssues = lastScan.LowCount
		}
	} else {
		// For space-level summary, aggregate from all repos in the space
		// This would require additional store methods for aggregation
		// For now, returning empty summary for space-level view
		summary.LastScanTime = 0
	}

	return summary, nil
}
