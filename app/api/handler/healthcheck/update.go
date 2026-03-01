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

package healthcheck

import (
	"encoding/json"
	"net/http"

	"github.com/harness/gitness/app/api/controller/healthcheck"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
)

// HandleUpdate returns a http.HandlerFunc that updates a health check.
func HandleUpdate(healthCheckCtrl *healthcheck.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef := request.GetParameter(r, "space_ref")
		identifier := request.GetParameter(r, "identifier")

		in := new(healthcheck.UpdateInput)
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			render.BadRequestf(ctx, w, "Invalid Request Body: %s.", err)
			return
		}

		hc, err := healthCheckCtrl.Update(ctx, session, spaceRef, identifier, in)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, hc)
	}
}

// HandleToggle returns a http.HandlerFunc that toggles a health check's enabled status.
func HandleToggle(healthCheckCtrl *healthcheck.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef := request.GetParameter(r, "space_ref")
		identifier := request.GetParameter(r, "identifier")

		hc, err := healthCheckCtrl.Toggle(ctx, session, spaceRef, identifier)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, hc)
	}
}
