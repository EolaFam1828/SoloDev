// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package airemediation

import (
	"net/http"
	"strconv"

	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
)

// HandleMetrics returns time-windowed remediation metrics.
func HandleMetrics(ctrl *airemediation.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		windowDays := 30
		if v := r.URL.Query().Get("window_days"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				windowDays = parsed
			}
		}

		metrics, err := ctrl.GetMetrics(ctx, session, spaceRef, windowDays)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, metrics)
	}
}
