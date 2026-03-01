// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed pagination and query parsing.
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

package featureflag

import (
	"net/http"

	"github.com/harness/gitness/app/api/controller/featureflag"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
)

func HandleList(ffCtrl *featureflag.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		filter := &types.FeatureFlagFilter{}
		filter.Page = request.ParsePage(r)
		filter.Size = request.ParseLimit(r)
		filter.Query = request.ParseQuery(r)

		ffs, err := ffCtrl.List(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		count, err := ffCtrl.Count(ctx, session, spaceRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.Pagination(r, w, filter.Page, filter.Size, int(count))
		render.JSON(w, http.StatusOK, ffs)
	}
}
