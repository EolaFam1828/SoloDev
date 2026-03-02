// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
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

package autopipeline

import (
	"encoding/json"
	"net/http"

	"github.com/EolaFam1828/SoloDev/app/api/controller/autopipeline"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
)

type generateRequest struct {
	Files []string `json:"files"`
}

// HandleGenerate is an HTTP handler for generating an auto-pipeline.
func HandleGenerate(ctrl *autopipeline.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		in := new(generateRequest)
		err = json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			render.BadRequestf(ctx, w, "Invalid Request Body: %s.", err)
			return
		}

		config, err := ctrl.GenerateAutoConfig(ctx, session, spaceRef, in.Files)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, config)
	}
}
