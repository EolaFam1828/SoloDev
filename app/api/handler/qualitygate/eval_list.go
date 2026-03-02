// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed list query filter struct literal.
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

package qualitygate

import (
	"net/http"

	"github.com/EolaFam1828/SoloDev/app/api/controller/qualitygate"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

// HandleEvaluationList is an HTTP handler for listing quality evaluations.
func HandleEvaluationList(ctrl *qualitygate.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		filter := &types.QualityEvaluationFilter{
			ListQueryFilter: request.ParseListQueryFilterFromRequest(r),
		}

		// Optional filters
		if status := r.URL.Query().Get("status"); status != "" {
			s := enum.QualityStatus(status)
			filter.OverallStatus = &s
		}

		if trigger := r.URL.Query().Get("triggered_by"); trigger != "" {
			t := enum.QualityTrigger(trigger)
			filter.TriggeredBy = &t
		}

		out, err := ctrl.ListEvaluations(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, out)
	}
}
