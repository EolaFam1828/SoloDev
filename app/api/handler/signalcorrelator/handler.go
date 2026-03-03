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
