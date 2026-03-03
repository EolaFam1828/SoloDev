// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package healthcheck

import (
	"net/http"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/anomalydetector"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/types/enum"
)

// HandleAnomalyDetect returns a handler that analyzes health check trends for anomalies.
func HandleAnomalyDetect(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	detector *anomalydetector.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		space, err := spaceFinder.FindByRef(ctx, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		if err := apiauth.CheckSpace(ctx, authorizer, session, space, enum.PermissionSpaceView); err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		report, err := detector.Analyze(ctx, space.ID)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, report)
	}
}
