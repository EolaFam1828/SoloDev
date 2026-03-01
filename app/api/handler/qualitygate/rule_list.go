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

package qualitygate

import (
	"net/http"

	"github.com/harness/gitness/app/api/controller/qualitygate"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// HandleRuleList is an HTTP handler for listing quality rules.
func HandleRuleList(ctrl *qualitygate.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		filter := &types.QualityRuleFilter{
			ListQueryFilter: types.ListQueryFilter{
				Page: request.GetQueryParamAsInt(r, "page", 0),
				Size: request.GetQueryParamAsInt(r, "limit", 20),
			},
		}

		// Optional filters
		if category := r.URL.Query().Get("category"); category != "" {
			cat := enum.QualityRuleCategory(category)
			filter.Category = &cat
		}

		if enforcement := r.URL.Query().Get("enforcement"); enforcement != "" {
			enf := enum.QualityEnforcement(enforcement)
			filter.Enforcement = &enf
		}

		if enabled := r.URL.Query().Get("enabled"); enabled != "" {
			enabledBool := enabled == "true"
			filter.Enabled = &enabledBool
		}

		out, err := ctrl.ListRules(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, out)
	}
}
