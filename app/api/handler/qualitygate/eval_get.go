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

package qualitygate

import (
	"net/http"

	"github.com/EolaFam1828/SoloDev/app/api/controller/qualitygate"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
)

// HandleEvaluationGet is an HTTP handler for retrieving a quality evaluation.
func HandleEvaluationGet(ctrl *qualitygate.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		evalIdentifier, err := request.GetQualityGateEvalIdentifierFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		eval, err := ctrl.GetEvaluation(ctx, session, spaceRef, evalIdentifier)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, eval)
	}
}
