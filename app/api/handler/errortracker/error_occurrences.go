// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed request API calls.
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

// HandleErrorOccurrences is an HTTP handler for listing error occurrences.
func HandleErrorOccurrences(ctrl *errortracker.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		identifier, err := request.GetErrorIdentifierFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		limit := request.ParseLimit(r)
		offset := request.ParsePage(r) * limit

		occurrences, err := ctrl.ListOccurrences(ctx, session, spaceRef, identifier, limit, offset)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		if occurrences == nil {
			occurrences = []types.ErrorOccurrence{}
		}

		render.JSON(w, http.StatusOK, occurrences)
	}
}
