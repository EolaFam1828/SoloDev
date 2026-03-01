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

package securityscan

import (
	"net/http"
	"strconv"

	"github.com/harness/gitness/app/api/controller/securityscan"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"

	"github.com/go-chi/chi/v5"
)

// HandleUpdateFindingStatus returns a http.HandlerFunc that updates a finding status.
func HandleUpdateFindingStatus(scanCtrl *securityscan.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		findingIDStr := chi.URLParam(r, "finding_id")
		findingID, err := strconv.ParseInt(findingIDStr, 10, 64)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		in := new(types.ScanFindingStatusUpdate)
		if err := render.DecodeJSON(r, in); err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		finding, err := scanCtrl.UpdateFindingStatus(ctx, session, spaceRef, findingID, in)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, finding)
	}
}
