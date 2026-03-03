// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package agentorchestrator

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/harness/gitness/types"

	"github.com/rs/zerolog/log"
)

// roleToolScope defines the default allowed tools for each agent role.
var roleToolScope = map[types.AgentRole][]string{
	types.AgentRoleMonitor: {
		"error_list", "health_status", "security_scan_list",
		"remediation_list", "pipeline_list",
	},
	types.AgentRoleRemediate: {
		"error_list", "error_detail", "remediation_trigger", "remediation_status",
		"remediation_apply", "fix_this", "security_scan_trigger",
	},
	types.AgentRoleValidate: {
		"remediation_status", "remediation_validate", "pipeline_list",
		"pipeline_trigger", "quality_gate_evaluate",
	},
	types.AgentRoleTriage: {
		"error_list", "error_detail", "health_status", "security_scan_list",
		"incident_triage", "signal_correlate",
	},
	types.AgentRoleGeneral: nil, // nil = all tools allowed
}

// Service manages agent sessions and multi-agent workflow coordination.
type Service struct {
	mu        sync.RWMutex
	sessions  map[string]*types.AgentSession
	handoffs  []types.AgentHandoff
	workflows map[string]*types.AgentWorkflow
}

// NewService creates a new agent orchestrator.
func NewService() *Service {
	return &Service{
		sessions:  make(map[string]*types.AgentSession),
		workflows: make(map[string]*types.AgentWorkflow),
	}
}

// CreateSession creates a new agent session with role-based tool scoping.
func (s *Service) CreateSession(_ context.Context, spaceID int64, input *types.CreateAgentSessionInput) (*types.AgentSession, error) {
	if input.Role == "" {
		input.Role = types.AgentRoleGeneral
	}

	allowedTools := input.AllowedTools
	if len(allowedTools) == 0 {
		allowedTools = roleToolScope[input.Role]
	}

	now := types.NowMillis()
	session := &types.AgentSession{
		ID:           generateID("ags"),
		SpaceID:      spaceID,
		Role:         input.Role,
		Status:       types.AgentSessionActive,
		AllowedTools: allowedTools,
		ClientInfo:   input.ClientInfo,
		Metadata:     input.Metadata,
		Created:      now,
		Updated:      now,
		LastActiveAt: now,
	}

	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()

	return session, nil
}

// FindSession returns a session by ID.
func (s *Service) FindSession(_ context.Context, sessionID string) (*types.AgentSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("agent session not found: %s", sessionID)
	}
	return session, nil
}

// ListSessions returns all sessions for a space.
func (s *Service) ListSessions(_ context.Context, spaceID int64) []*types.AgentSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.AgentSession
	for _, session := range s.sessions {
		if session.SpaceID == spaceID {
			result = append(result, session)
		}
	}
	return result
}

// CloseSession closes an agent session.
func (s *Service) CloseSession(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("agent session not found: %s", sessionID)
	}
	session.Status = types.AgentSessionClosed
	session.Updated = types.NowMillis()
	return nil
}

// IsToolAllowed checks whether a session is permitted to call a given tool.
func (s *Service) IsToolAllowed(_ context.Context, sessionID, toolName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return false
	}
	// nil allowedTools = all tools allowed (general role).
	if session.AllowedTools == nil {
		return true
	}
	for _, t := range session.AllowedTools {
		if t == toolName {
			return true
		}
	}
	return false
}

// Handoff transfers work from one agent session to another.
func (s *Service) Handoff(_ context.Context, fromSessionID string, input *types.AgentHandoffInput) (*types.AgentHandoff, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	from, ok := s.sessions[fromSessionID]
	if !ok {
		return nil, fmt.Errorf("source session not found: %s", fromSessionID)
	}
	to, ok := s.sessions[input.ToSessionID]
	if !ok {
		return nil, fmt.Errorf("target session not found: %s", input.ToSessionID)
	}

	now := types.NowMillis()

	// Clear task from source, assign to target.
	from.CurrentTaskID = ""
	from.Status = types.AgentSessionActive
	from.Updated = now

	to.CurrentTaskID = input.TaskID
	to.Status = types.AgentSessionWorking
	to.Updated = now
	to.LastActiveAt = now

	handoff := types.AgentHandoff{
		FromSessionID: fromSessionID,
		ToSessionID:   input.ToSessionID,
		TaskID:        input.TaskID,
		Reason:        input.Reason,
		Context:       input.Context,
		Created:       now,
	}

	s.handoffs = append(s.handoffs, handoff)
	return &handoff, nil
}

// StartWorkflow creates and starts a multi-agent workflow from a template.
func (s *Service) StartWorkflow(ctx context.Context, spaceID int64, input *types.StartAgentWorkflowInput) (*types.AgentWorkflow, error) {
	steps := templateSteps(input.Template)
	if len(steps) == 0 {
		return nil, fmt.Errorf("unknown workflow template: %s", input.Template)
	}

	now := types.NowMillis()
	workflow := &types.AgentWorkflow{
		ID:          generateID("wf"),
		SpaceID:     spaceID,
		Template:    input.Template,
		Status:      types.AgentWorkflowPending,
		Steps:       steps,
		CurrentStep: 0,
		Created:     now,
		Updated:     now,
	}

	// Create a session for each step and assign it.
	var sessionIDs []string
	for i, step := range workflow.Steps {
		session, err := s.CreateSession(ctx, spaceID, &types.CreateAgentSessionInput{
			Role:       step.Role,
			ClientInfo: fmt.Sprintf("workflow-%s-step-%d", workflow.ID, i),
		})
		if err != nil {
			return nil, fmt.Errorf("create session for step %d: %w", i, err)
		}
		workflow.Steps[i].SessionID = session.ID
		sessionIDs = append(sessionIDs, session.ID)
	}
	workflow.SessionIDs = sessionIDs

	s.mu.Lock()
	s.workflows[workflow.ID] = workflow
	s.mu.Unlock()

	log.Ctx(ctx).Info().
		Str("workflow_id", workflow.ID).
		Str("template", string(input.Template)).
		Int("steps", len(steps)).
		Msg("multi-agent workflow started")

	return workflow, nil
}

// FindWorkflow returns a workflow by ID.
func (s *Service) FindWorkflow(_ context.Context, workflowID string) (*types.AgentWorkflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wf, ok := s.workflows[workflowID]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}
	return wf, nil
}

// ListWorkflows returns all workflows for a space.
func (s *Service) ListWorkflows(_ context.Context, spaceID int64) []*types.AgentWorkflow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.AgentWorkflow
	for _, wf := range s.workflows {
		if wf.SpaceID == spaceID {
			result = append(result, wf)
		}
	}
	return result
}

// AdvanceWorkflow advances a workflow to the next step.
func (s *Service) AdvanceWorkflow(_ context.Context, workflowID, result string) (*types.AgentWorkflow, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	wf, ok := s.workflows[workflowID]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	now := types.NowMillis()

	// Complete current step.
	if wf.CurrentStep < len(wf.Steps) {
		wf.Steps[wf.CurrentStep].Status = types.AgentWorkflowCompleted
		wf.Steps[wf.CurrentStep].Result = result
		wf.Steps[wf.CurrentStep].CompletedAt = now
	}

	// Advance to next step.
	wf.CurrentStep++
	wf.Updated = now

	if wf.CurrentStep >= len(wf.Steps) {
		wf.Status = types.AgentWorkflowCompleted
	} else {
		wf.Status = types.AgentWorkflowRunning
		wf.Steps[wf.CurrentStep].Status = types.AgentWorkflowRunning
		wf.Steps[wf.CurrentStep].StartedAt = now
	}

	return wf, nil
}

// templateSteps returns the step definitions for a workflow template.
func templateSteps(template types.AgentWorkflowTemplate) []types.AgentWorkflowStep {
	switch template {
	case types.WorkflowSelfHeal:
		return []types.AgentWorkflowStep{
			{Name: "detect", Role: types.AgentRoleMonitor, Status: types.AgentWorkflowPending},
			{Name: "fix", Role: types.AgentRoleRemediate, Status: types.AgentWorkflowPending},
			{Name: "verify", Role: types.AgentRoleValidate, Status: types.AgentWorkflowPending},
		}
	case types.WorkflowFullAudit:
		return []types.AgentWorkflowStep{
			{Name: "triage", Role: types.AgentRoleTriage, Status: types.AgentWorkflowPending},
			{Name: "fix", Role: types.AgentRoleRemediate, Status: types.AgentWorkflowPending},
			{Name: "verify", Role: types.AgentRoleValidate, Status: types.AgentWorkflowPending},
		}
	case types.WorkflowContinuous:
		return []types.AgentWorkflowStep{
			{Name: "monitor", Role: types.AgentRoleMonitor, Status: types.AgentWorkflowPending},
		}
	default:
		return nil
	}
}

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(b))
}
