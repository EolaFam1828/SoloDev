// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed path parameter extraction.
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

package techdebt

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/harness/gitness/app/api/controller/techdebt"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
)

// HandleList returns a http.HandlerFunc that lists technical debt items.
func HandleList(ctrl *techdebt.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		// Parse query parameters
		filter := &types.TechDebtFilter{}

		// Parse severity filter
		if severityStr := r.URL.Query().Get("severity"); severityStr != "" {
			filter.Severity = strings.Split(severityStr, ",")
		}

		// Parse status filter
		if statusStr := r.URL.Query().Get("status"); statusStr != "" {
			filter.Status = strings.Split(statusStr, ",")
		}

		// Parse category filter
		if categoryStr := r.URL.Query().Get("category"); categoryStr != "" {
			filter.Category = strings.Split(categoryStr, ",")
		}

		// Parse repo ID filter
		if repoIDStr := r.URL.Query().Get("repo"); repoIDStr != "" {
			if repoID, err := strconv.ParseInt(repoIDStr, 10, 64); err == nil {
				filter.RepoID = repoID
			}
		}

		// Parse pagination
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
				filter.Page = page
			} else {
				filter.Page = 1
			}
		} else {
			filter.Page = 1
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				filter.Limit = limit
			} else {
				filter.Limit = 20
			}
		} else {
			filter.Limit = 20
		}

		resp, err := ctrl.List(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, resp)
	}
}

// HandleSummary returns a http.HandlerFunc that returns aggregated statistics.
func HandleSummary(ctrl *techdebt.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		// Parse query parameters
		filter := &types.TechDebtFilter{}

		if repoIDStr := r.URL.Query().Get("repo"); repoIDStr != "" {
			if repoID, err := strconv.ParseInt(repoIDStr, 10, 64); err == nil {
				filter.RepoID = repoID
			}
		}

		summary, err := ctrl.Summary(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, summary)
	}
}

// HandleResolve returns a http.HandlerFunc that quickly resolves a technical debt item.
func HandleResolve(ctrl *techdebt.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		identifier, err := request.PathParamOrError(r, "identifier")
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		td, err := ctrl.Resolve(ctx, session, spaceRef, identifier)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, td)
	}
}
