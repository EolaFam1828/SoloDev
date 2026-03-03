// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package agentorchestrator

import (
	"context"
	"fmt"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/agentorchestrator"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// Controller provides REST operations for agent session management.
type Controller struct {
	authorizer  authz.Authorizer
	spaceFinder refcache.SpaceFinder
	service     *agentorchestrator.Service
}

// ProvideController creates a new agent orchestrator controller.
func ProvideController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	service *agentorchestrator.Service,
) *Controller {
	return &Controller{
		authorizer:  authorizer,
		spaceFinder: spaceFinder,
		service:     service,
	}
}

// CreateSession creates a new agent session.
func (c *Controller) CreateSession(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	input *types.CreateAgentSessionInput,
) (*types.AgentSession, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceEdit); err != nil {
		return nil, err
	}
	return c.service.CreateSession(ctx, space.ID, input)
}

// FindSession finds an agent session by ID.
func (c *Controller) FindSession(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	sessionID string,
) (*types.AgentSession, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}
	return c.service.FindSession(ctx, sessionID)
}

// ListSessions lists all agent sessions for a space.
func (c *Controller) ListSessions(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
) ([]*types.AgentSession, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}
	return c.service.ListSessions(ctx, space.ID), nil
}

// CloseSession closes an agent session.
func (c *Controller) CloseSession(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	sessionID string,
) error {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceEdit); err != nil {
		return err
	}
	return c.service.CloseSession(ctx, sessionID)
}

// Handoff transfers work between sessions.
func (c *Controller) Handoff(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	fromSessionID string,
	input *types.AgentHandoffInput,
) (*types.AgentHandoff, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceEdit); err != nil {
		return nil, err
	}
	return c.service.Handoff(ctx, fromSessionID, input)
}

// StartWorkflow starts a multi-agent workflow.
func (c *Controller) StartWorkflow(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	input *types.StartAgentWorkflowInput,
) (*types.AgentWorkflow, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceEdit); err != nil {
		return nil, err
	}
	return c.service.StartWorkflow(ctx, space.ID, input)
}

// ListWorkflows lists all workflows for a space.
func (c *Controller) ListWorkflows(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
) ([]*types.AgentWorkflow, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}
	return c.service.ListWorkflows(ctx, space.ID), nil
}

// FindWorkflow finds a workflow by ID.
func (c *Controller) FindWorkflow(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	workflowID string,
) (*types.AgentWorkflow, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}
	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}
	return c.service.FindWorkflow(ctx, workflowID)
}
