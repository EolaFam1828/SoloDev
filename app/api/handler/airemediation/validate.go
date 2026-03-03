// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
			_ = json.NewDecoder(r.Body).Decode(&in)
		}

		rem, err := ctrl.ValidateRemediation(ctx, session, spaceRef, identifier, in.PipelineIdentifier)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, rem)
	}
}
