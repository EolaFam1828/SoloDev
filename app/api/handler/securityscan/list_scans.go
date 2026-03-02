// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Fixed query parameter extraction.
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

	"github.com/EolaFam1828/SoloDev/app/api/controller/securityscan"
	"github.com/EolaFam1828/SoloDev/app/api/render"
	"github.com/EolaFam1828/SoloDev/app/api/request"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

// HandleListScans returns a http.HandlerFunc that lists security scans.
func HandleListScans(scanCtrl *securityscan.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		repoRef, err := request.GetRepoRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		page, err := request.QueryParamAsPositiveInt64OrDefault(r, "page", 0)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}
		size, err := request.QueryParamAsPositiveInt64OrDefault(r, "limit", 20)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}
		sortStr := request.QueryParamOrDefault(r, "sort", "created")
		orderStr := request.QueryParamOrDefault(r, "order", "desc")
		statusStr := request.QueryParamOrDefault(r, "status", "")
		scanTypeStr := request.QueryParamOrDefault(r, "scan_type", "")
		triggeredByStr := request.QueryParamOrDefault(r, "triggered_by", "")

		filter := &types.ScanResultFilter{
			Page:       int(page),
			Size:       int(size),
			Sort:       enum.ParseSecurityScanAttr(sortStr),
			Order:      enum.ParseOrder(orderStr),
			Status:     enum.SecurityScanStatus(statusStr),
			ScanType:   enum.SecurityScanType(scanTypeStr),
			TriggeredBy: enum.SecurityScanTrigger(triggeredByStr),
		}

		scans, count, err := scanCtrl.ListScans(ctx, session, repoRef, filter)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, map[string]interface{}{
			"data":  scans,
			"count": count,
		})
	}
}
