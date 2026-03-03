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

package agentorchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harness/gitness/app/api/controller/agentorchestrator"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
)

// HandleCreateSession returns a handler that creates an agent session.
func HandleCreateSession(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		var input types.CreateAgentSessionInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.BadRequestf(ctx, w, "invalid request body: %s", err)
			return
		}

		result, err := ctrl.CreateSession(ctx, session, spaceRef, &input)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusCreated, result)
	}
}

// HandleListSessions returns a handler that lists agent sessions.
func HandleListSessions(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		sessions, err := ctrl.ListSessions(ctx, session, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, sessions)
	}
}

// HandleFindSession returns a handler that finds an agent session.
func HandleFindSession(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		sessionID := chi.URLParam(r, "session_id")
		result, err := ctrl.FindSession(ctx, session, spaceRef, sessionID)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, result)
	}
}

// HandleCloseSession returns a handler that closes an agent session.
func HandleCloseSession(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		sessionID := chi.URLParam(r, "session_id")
		if err := ctrl.CloseSession(ctx, session, spaceRef, sessionID); err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// HandleHandoff returns a handler that performs a work handoff between sessions.
func HandleHandoff(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		fromSessionID := chi.URLParam(r, "session_id")
		var input types.AgentHandoffInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.BadRequestf(ctx, w, "invalid request body: %s", err)
			return
		}

		result, err := ctrl.Handoff(ctx, session, spaceRef, fromSessionID, &input)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, result)
	}
}

// HandleStartWorkflow returns a handler that starts a multi-agent workflow.
func HandleStartWorkflow(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		var input types.StartAgentWorkflowInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.BadRequestf(ctx, w, "invalid request body: %s", err)
			return
		}

		result, err := ctrl.StartWorkflow(ctx, session, spaceRef, &input)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusCreated, result)
	}
}

// HandleListWorkflows returns a handler that lists workflows.
func HandleListWorkflows(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		workflows, err := ctrl.ListWorkflows(ctx, session, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, workflows)
	}
}

// HandleFindWorkflow returns a handler that finds a workflow.
func HandleFindWorkflow(ctrl *agentorchestrator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		workflowID := chi.URLParam(r, "workflow_id")
		result, err := ctrl.FindWorkflow(ctx, session, spaceRef, workflowID)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, result)
	}
}
