// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package types

// AgentRole defines the specialization of an agent session.
type AgentRole string

const (
	AgentRoleMonitor   AgentRole = "monitor"   // Watches resources, detects anomalies
	AgentRoleRemediate AgentRole = "remediate" // Generates and applies fixes
	AgentRoleValidate  AgentRole = "validate"  // Validates applied fixes via pipelines
	AgentRoleTriage    AgentRole = "triage"    // Correlates signals and assesses severity
	AgentRoleGeneral   AgentRole = "general"   // Full access, no specialization
)

// AgentSessionStatus tracks session lifecycle.
type AgentSessionStatus string

const (
	AgentSessionActive    AgentSessionStatus = "active"
	AgentSessionWorking   AgentSessionStatus = "working"
	AgentSessionCompleted AgentSessionStatus = "completed"
	AgentSessionClosed    AgentSessionStatus = "closed"
)

// AgentSession represents an active agent connection with scoped capabilities.
type AgentSession struct {
	ID            string             `json:"id"`
	SpaceID       int64              `json:"space_id"`
	Role          AgentRole          `json:"role"`
	Status        AgentSessionStatus `json:"status"`
	AllowedTools  []string           `json:"allowed_tools"`
	ClientInfo    string             `json:"client_info"`
	CurrentTaskID string             `json:"current_task_id,omitempty"`
	Metadata      map[string]string  `json:"metadata,omitempty"`
	Created       int64              `json:"created"`
	Updated       int64              `json:"updated"`
	LastActiveAt  int64              `json:"last_active_at"`
}

// CreateAgentSessionInput is the request body for creating a session.
type CreateAgentSessionInput struct {
	Role         AgentRole         `json:"role"`
	ClientInfo   string            `json:"client_info,omitempty"`
	AllowedTools []string          `json:"allowed_tools,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// AgentHandoff represents a work handoff between two agent sessions.
type AgentHandoff struct {
	FromSessionID string `json:"from_session_id"`
	ToSessionID   string `json:"to_session_id"`
	TaskID        string `json:"task_id"`
	Reason        string `json:"reason"`
	Context       string `json:"context,omitempty"` // Serialized context passed between agents
	Created       int64  `json:"created"`
}

// AgentHandoffInput is the request body for creating a handoff.
type AgentHandoffInput struct {
	ToSessionID string `json:"to_session_id"`
	TaskID      string `json:"task_id"`
	Reason      string `json:"reason"`
	Context     string `json:"context,omitempty"`
}

// AgentWorkflowStatus tracks workflow lifecycle.
type AgentWorkflowStatus string

const (
	AgentWorkflowPending   AgentWorkflowStatus = "pending"
	AgentWorkflowRunning   AgentWorkflowStatus = "running"
	AgentWorkflowCompleted AgentWorkflowStatus = "completed"
	AgentWorkflowFailed    AgentWorkflowStatus = "failed"
)

// AgentWorkflowTemplate defines a reusable multi-agent workflow.
type AgentWorkflowTemplate string

const (
	WorkflowSelfHeal   AgentWorkflowTemplate = "self_heal"  // Monitor → Remediate → Validate
	WorkflowFullAudit  AgentWorkflowTemplate = "full_audit" // Triage → Remediate → Validate
	WorkflowContinuous AgentWorkflowTemplate = "continuous" // Monitor loop
)

// AgentWorkflow is a running multi-agent workflow instance.
type AgentWorkflow struct {
	ID          string                `json:"id"`
	SpaceID     int64                 `json:"space_id"`
	Template    AgentWorkflowTemplate `json:"template"`
	Status      AgentWorkflowStatus   `json:"status"`
	Steps       []AgentWorkflowStep   `json:"steps"`
	CurrentStep int                   `json:"current_step"`
	SessionIDs  []string              `json:"session_ids"`
	Created     int64                 `json:"created"`
	Updated     int64                 `json:"updated"`
}

// AgentWorkflowStep is one step in a workflow.
type AgentWorkflowStep struct {
	Name        string              `json:"name"`
	Role        AgentRole           `json:"role"`
	SessionID   string              `json:"session_id,omitempty"`
	Status      AgentWorkflowStatus `json:"status"`
	Result      string              `json:"result,omitempty"`
	StartedAt   int64               `json:"started_at,omitempty"`
	CompletedAt int64               `json:"completed_at,omitempty"`
}

// StartAgentWorkflowInput is the request body for starting a workflow.
type StartAgentWorkflowInput struct {
	Template AgentWorkflowTemplate `json:"template"`
}
