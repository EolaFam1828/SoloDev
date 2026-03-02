// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed request parameter extraction and list parsing.
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

package healthcheck

import (
	"net/http"
	"strconv"

	"github.com/EolaFam1828/SoloDev/app/api/controller/healthcheck"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
)

// HandleFind returns a http.HandlerFunc that finds a health check.
func HandleFind(healthCheckCtrl *healthcheck.Controller) http.HandlerFunc {
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

		hc, err := healthCheckCtrl.Find(ctx, session, spaceRef, identifier)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, hc)
	}
}

// HandleList returns a http.HandlerFunc that lists health checks in a space.
func HandleList(healthCheckCtrl *healthcheck.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		filter := request.ParseListQueryFilterFromRequest(r)

		hcs, err := healthCheckCtrl.List(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, hcs)
	}
}

// HandleListResults returns a http.HandlerFunc that lists results for a health check.
func HandleListResults(healthCheckCtrl *healthcheck.Controller) http.HandlerFunc {
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

		limit := 100
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		results, err := healthCheckCtrl.GetResults(ctx, session, spaceRef, identifier, limit)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, results)
	}
}

// HandleSummary returns a http.HandlerFunc that gets overall uptime stats.
func HandleSummary(healthCheckCtrl *healthcheck.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		summaries, err := healthCheckCtrl.GetSummary(ctx, session, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, summaries)
	}
}
