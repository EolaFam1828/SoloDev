// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed body unmarshal.
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

package techdebt

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/EolaFam1828/SoloDev/app/api/controller/techdebt"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
	"github.com/EolaFam1828/SoloDev/types"
)

// HandleCreate returns a http.HandlerFunc that creates a new technical debt item.
func HandleCreate(ctrl *techdebt.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			render.BadRequestf(ctx, w, "Invalid Request Body: %s", err)
			return
		}

		in := new(types.TechDebtCreateInput)
		if err := json.Unmarshal(body, in); err != nil {
			render.BadRequestf(ctx, w, "Invalid Request Body: %s", err)
			return
		}

		td, err := ctrl.Create(ctx, session, spaceRef, in)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusCreated, td)
	}
}
