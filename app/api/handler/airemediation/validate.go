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

package airemediation

import (
	"encoding/json"
	"net/http"

	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
)

type validateRequest struct {
	PipelineIdentifier string `json:"pipeline_identifier"`
}

// HandleValidate is an HTTP handler for triggering validation on a remediation.
func HandleValidate(ctrl *airemediation.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		identifier, err := request.GetRemediationIdentifierFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		var in validateRequest
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				render.BadRequestf(ctx, w, "Invalid Request Body: %s.", err)
				return
			}
		}

		rem, err := ctrl.ValidateRemediation(ctx, session, spaceRef, identifier, in.PipelineIdentifier)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, rem)
	}
}
