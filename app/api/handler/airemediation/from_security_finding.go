// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package airemediation

import (
	"encoding/json"
	"net/http"

	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
)

// HandleTriggerFromSecurityFinding creates or reuses a remediation for a security finding.
func HandleTriggerFromSecurityFinding(ctrl *airemediation.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		in := new(types.CreateRemediationFromSecurityFindingInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			render.BadRequestf(ctx, w, "Invalid Request Body: %s.", err)
			return
		}

		rem, created, err := ctrl.TriggerRemediationFromSecurityFinding(ctx, session, spaceRef, in)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		status := http.StatusCreated
		if !created {
			status = http.StatusOK
		}
		render.JSON(w, status, rem)
	}
}
