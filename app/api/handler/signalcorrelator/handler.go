// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package signalcorrelator

import (
	"net/http"
	"strconv"

	"github.com/harness/gitness/app/api/controller/signalcorrelator"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
)

// HandleCorrelate returns a http.HandlerFunc that runs signal correlation for a space.
func HandleCorrelate(ctrl *signalcorrelator.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		filter := types.CorrelatedIncidentFilter{}
		if v := r.URL.Query().Get("window_minutes"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				filter.WindowMinutes = parsed
			}
		}
		if v := r.URL.Query().Get("min_signals"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				filter.MinSignals = parsed
			}
		}
		if v := r.URL.Query().Get("repo_id"); v != "" {
			if parsed, err := strconv.ParseInt(v, 10, 64); err == nil && parsed > 0 {
				filter.RepoID = &parsed
			}
		}

		incidents, err := ctrl.Correlate(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		if incidents == nil {
			incidents = []types.CorrelatedIncident{}
		}

		render.JSON(w, http.StatusOK, incidents)
	}
}
