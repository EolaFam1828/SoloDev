// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed request API calls (ParsePage/ParseLimit/ParseQuery).
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
	"net/http"

	"github.com/EolaFam1828/SoloDev/app/api/controller/errortracker"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
	"github.com/EolaFam1828/SoloDev/types"
)

// HandleErrorList is an HTTP handler for listing error groups.
func HandleErrorList(ctrl *errortracker.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		// Parse query parameters
		opts := types.ErrorTrackerListOptions{
			ListQueryFilter: types.ListQueryFilter{
				Pagination: types.Pagination{
					Page: request.ParsePage(r),
					Size: request.ParseLimit(r),
				},
				Query: request.ParseQuery(r),
			},
		}

		// Parse status filter
		if statusStr := r.URL.Query().Get("status"); statusStr != "" {
			status := types.ErrorGroupStatus(statusStr)
			opts.Status = &status
		}

		// Parse severity filter
		if severityStr := r.URL.Query().Get("severity"); severityStr != "" {
			severity := types.ErrorSeverity(severityStr)
			opts.Severity = &severity
		}

		// Parse language filter
		if language := r.URL.Query().Get("language"); language != "" {
			opts.Language = &language
		}

		errorGroups, err := ctrl.ListErrors(ctx, session, spaceRef, opts)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		if errorGroups == nil {
			errorGroups = []types.ErrorGroup{}
		}

		render.JSON(w, http.StatusOK, errorGroups)
	}
}
