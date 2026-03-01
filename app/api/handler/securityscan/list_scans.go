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

	"github.com/harness/gitness/app/api/controller/securityscan"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
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

		page := request.GetQueryParamInt(r, "page", 0)
		size := request.GetQueryParamInt(r, "limit", 20)
		sortStr := request.GetQueryParam(r, "sort", "created")
		orderStr := request.GetQueryParam(r, "order", "desc")
		statusStr := request.GetQueryParam(r, "status", "")
		scanTypeStr := request.GetQueryParam(r, "scan_type", "")
		triggeredByStr := request.GetQueryParam(r, "triggered_by", "")

		filter := &types.ScanResultFilter{
			Page:       page,
			Size:       size,
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
