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

package errortracker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/api/controller/space"
	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	errortrackerevents "github.com/harness/gitness/app/events/errortracker"
	"github.com/harness/gitness/app/services/errorbridge"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type Controller struct {
	tx                 dbtx.Transactor
	authorizer         authz.Authorizer
	spaceFinder        refcache.SpaceFinder
	repoFinder         refcache.RepoFinder
	errorTrackerStore  store.ErrorTrackerStore
	principalInfoCache store.PrincipalInfoCache
	eventReporter      *errortrackerevents.Reporter
	errorBridge        *errorbridge.Bridge
}

func NewController(
	tx dbtx.Transactor,
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	errorTrackerStore store.ErrorTrackerStore,
	principalInfoCache store.PrincipalInfoCache,
	eventReporter *errortrackerevents.Reporter,
) *Controller {
	return &Controller{
		tx:                 tx,
		authorizer:         authorizer,
		spaceFinder:        spaceFinder,
		repoFinder:         repoFinder,
		errorTrackerStore:  errorTrackerStore,
		principalInfoCache: principalInfoCache,
		eventReporter:      eventReporter,
	}
}

// SetErrorBridge sets the optional error-to-AI remediation bridge.
// When set, new errors will automatically trigger AI remediation tasks.
func (c *Controller) SetErrorBridge(bridge *errorbridge.Bridge) {
	c.errorBridge = bridge
}

// ReportError reports a new error occurrence and creates or updates an error group.
func (c *Controller) ReportError(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	in *types.ReportErrorInput,
) (*types.ErrorGroup, error) {
	if in == nil {
		return nil, usererror.BadRequest("Request body cannot be empty")
	}

	// Get space and check permission
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	// Normalize severity if not provided
	if in.Severity == "" {
		in.Severity = types.ErrorSeverityError
	}

	// Normalize environment
	if in.Environment == "" {
		in.Environment = "production"
	}

	// Generate fingerprint
	fingerprint := store.Fingerprint(in.Title, in.StackTrace)

	// Create tags JSON
	var tagsJSON json.RawMessage
	if len(in.Tags) > 0 {
		tagsJSON, _ = json.Marshal(in.Tags)
	}

	now := types.NowMillis()

	// Create error group
	errorGroup := &types.ErrorGroup{
		SpaceID:         space.ID,
		RepoID:          0, // Could be populated if repo is provided in request
		Identifier:      in.Identifier,
		Title:           in.Title,
		Message:         in.Message,
		Fingerprint:     fingerprint,
		Status:          types.ErrorGroupStatusOpen,
		Severity:        in.Severity,
		FirstSeen:       now,
		LastSeen:        now,
		OccurrenceCount: 1,
		FilePath:        in.FilePath,
		LineNumber:      in.LineNumber,
		FunctionName:    in.FunctionName,
		Language:        in.Language,
		Tags:            tagsJSON,
		CreatedBy:       session.Principal.ID,
		Created:         now,
		Updated:         now,
		Version:         1,
	}

	// Use transaction to create/update error group and add occurrence
	var result *types.ErrorGroup
	err = c.tx.WithTx(ctx, func(tx context.Context) error {
		// Create or update error group
		if err := c.errorTrackerStore.CreateOrUpdateErrorGroup(tx, errorGroup); err != nil {
			return fmt.Errorf("failed to create or update error group: %w", err)
		}

		// Create error occurrence
		occurrence := &types.ErrorOccurrence{
			ErrorGroupID: errorGroup.ID,
			StackTrace:   in.StackTrace,
			Environment:  in.Environment,
			Runtime:      in.Runtime,
			OS:           in.OS,
			Arch:         in.Arch,
			Metadata:     in.Metadata,
			CreatedAt:    now,
		}

		if err := c.errorTrackerStore.CreateErrorOccurrence(tx, occurrence); err != nil {
			return fmt.Errorf("failed to create error occurrence: %w", err)
		}

		// Fetch the updated error group
		updated, err := c.errorTrackerStore.FindByFingerprint(tx, space.ID, fingerprint)
		if err != nil {
			return fmt.Errorf("failed to fetch updated error group: %w", err)
		}

		result = updated
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Report event
	if c.eventReporter != nil {
		c.eventReporter.ErrorReported(ctx, result)
	}

	// Auto-trigger AI remediation via error bridge
	if c.errorBridge != nil {
		// Build a minimal occurrence for the bridge
		bridgeOccurrence := &types.ErrorOccurrence{
			ErrorGroupID: result.ID,
			StackTrace:   in.StackTrace,
			Environment:  in.Environment,
			Runtime:      in.Runtime,
			OS:           in.OS,
			Arch:         in.Arch,
		}
		c.errorBridge.OnErrorReported(ctx, result, bridgeOccurrence)
	}

	return result, nil
}

// ListErrors returns a list of error groups for a space.
func (c *Controller) ListErrors(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	opts types.ErrorTrackerListOptions,
) ([]types.ErrorGroup, error) {
	// Get space and check permission
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	return c.errorTrackerStore.List(ctx, space.ID, opts)
}

// GetError returns details for a specific error group.
func (c *Controller) GetError(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
) (*types.ErrorGroupDetail, error) {
	// Get space and check permission
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	errorGroup, err := c.errorTrackerStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find error group: %w", err)
	}

	if errorGroup == nil {
		return nil, usererror.NotFound("Error group not found")
	}

	// Fetch sample occurrences
	occurrences, err := c.errorTrackerStore.ListOccurrences(ctx, errorGroup.ID, 5, 0)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to fetch occurrences: %v\n", err)
		occurrences = []types.ErrorOccurrence{}
	}

	// Fetch related user information
	detail := &types.ErrorGroupDetail{
		ErrorGroup:        errorGroup,
		OccurrencesSample: occurrences,
	}

	// Fetch user info if assigned
	if errorGroup.AssignedTo != nil {
		if userInfo, err := c.principalInfoCache.Get(ctx, *errorGroup.AssignedTo); err == nil {
			detail.AssignedUser = userInfo
		}
	}

	// Fetch resolved by user info
	if errorGroup.ResolvedBy != nil {
		if userInfo, err := c.principalInfoCache.Get(ctx, *errorGroup.ResolvedBy); err == nil {
			detail.ResolvedByUser = userInfo
		}
	}

	// Fetch created by user info
	if userInfo, err := c.principalInfoCache.Get(ctx, errorGroup.CreatedBy); err == nil {
		detail.CreatedByUser = userInfo
	}

	return detail, nil
}

// UpdateError updates an error group status or assignment.
func (c *Controller) UpdateError(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
	in *types.UpdateErrorGroupInput,
) (*types.ErrorGroup, error) {
	// Get space and check permission
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceEdit)
	if err != nil {
		return nil, err
	}

	errorGroup, err := c.errorTrackerStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find error group: %w", err)
	}

	if errorGroup == nil {
		return nil, usererror.NotFound("Error group not found")
	}

	// Update status if provided
	if in.Status != nil && *in.Status != errorGroup.Status {
		err = c.errorTrackerStore.UpdateStatus(ctx, errorGroup.ID, *in.Status, session.Principal.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update error group status: %w", err)
		}

		// Report status change event
		if c.eventReporter != nil {
			c.eventReporter.ErrorStatusChanged(ctx, errorGroup, *in.Status)
		}
	}

	// Update assignment if provided
	if in.AssignedTo != errorGroup.AssignedTo {
		err = c.errorTrackerStore.UpdateAssignment(ctx, errorGroup.ID, in.AssignedTo)
		if err != nil {
			return nil, fmt.Errorf("failed to update error group assignment: %w", err)
		}

		// Report assignment change event
		if c.eventReporter != nil {
			c.eventReporter.ErrorAssigned(ctx, errorGroup, in.AssignedTo)
		}
	}

	// Fetch updated error group
	updated, err := c.errorTrackerStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated error group: %w", err)
	}

	return updated, nil
}

// ListOccurrences returns a list of error occurrences for an error group.
func (c *Controller) ListOccurrences(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	identifier string,
	limit int,
	offset int,
) ([]types.ErrorOccurrence, error) {
	// Get space and check permission
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	errorGroup, err := c.errorTrackerStore.FindByIdentifier(ctx, space.ID, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find error group: %w", err)
	}

	if errorGroup == nil {
		return nil, usererror.NotFound("Error group not found")
	}

	// Default limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	return c.errorTrackerStore.ListOccurrences(ctx, errorGroup.ID, limit, offset)
}

// GetSummary returns summary statistics for error groups in a space.
func (c *Controller) GetSummary(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
) (*types.ErrorTrackerSummary, error) {
	// Get space and check permission
	space, err := c.getSpaceCheckAccess(ctx, session, spaceRef, enum.PermissionSpaceView)
	if err != nil {
		return nil, err
	}

	return c.errorTrackerStore.GetSummary(ctx, space.ID)
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
