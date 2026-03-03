// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package solodev

import (
	"net/http"

	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/controller/errortracker"
	"github.com/harness/gitness/app/api/controller/securityscan"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/mcp"
	"github.com/harness/gitness/types"
)

func HandleOverview(
	securityScanCtrl *securityscan.Controller,
	errorTrackerCtrl *errortracker.Controller,
	remediationCtrl *airemediation.Controller,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		catalog := mcp.BuildCatalog(&mcp.Controllers{
			SecurityScan: securityScanCtrl,
			ErrorTracker: errorTrackerCtrl,
			Remediation:  remediationCtrl,
		})

		overview := types.SoloDevOverview{
			SpaceRef:  spaceRef,
			UpdatedAt: types.NowMillis(),
			Security: types.SoloDevSecurityOverview{
				Availability:     "blocked",
				LatestScanStatus: "not_started",
			},
			Remediation: types.SoloDevRemediationOverview{
				Availability: "blocked",
			},
			Errors: types.SoloDevErrorsOverview{
				Availability: "blocked",
			},
			MCP: types.SoloDevMCPOverview{
				Tools:     catalog.Counts.ActiveTools,
				Resources: catalog.Counts.ActiveResources,
				Prompts:   catalog.Counts.ActivePrompts,
			},
			DeferredDomains: []string{"quality", "health", "pipelines", "feature_flags", "tech_debt"},
		}

		if securityScanCtrl != nil {
			scannerStatus := securityScanCtrl.ScannerStatus()
			if scannerStatus.Ready {
				overview.Security.Availability = "ready"
			} else if scannerStatus.Enabled {
				overview.Security.Availability = "blocked"
			}

			summary, err := securityScanCtrl.GetSecuritySummary(ctx, session, spaceRef, nil)
			if err != nil {
				render.TranslatedUserError(ctx, w, err)
				return
			}
			if summary != nil {
				if summary.LastScanID > 0 {
					overview.Security.LatestScanStatus = "completed"
				}
				overview.Security.OpenFindings = summary.TotalFindings
				overview.Security.Critical = summary.CriticalIssues
				overview.Security.High = summary.HighIssues
				overview.Security.Medium = summary.MediumIssues
				overview.Security.Low = summary.LowIssues
				overview.Security.LastScanTime = summary.LastScanTime
			}
		}

		if remediationCtrl != nil {
			if remediationCtrl.AIAvailable() {
				overview.Remediation.Availability = "ready"
			}
			summary, err := remediationCtrl.GetSummary(ctx, session, spaceRef)
			if err != nil {
				render.TranslatedUserError(ctx, w, err)
				return
			}
			if summary != nil {
				overview.Remediation.Pending = summary.Pending
				overview.Remediation.Processing = summary.Processing
				overview.Remediation.Completed = summary.Completed
				overview.Remediation.Applied = summary.Applied
				overview.Remediation.Failed = summary.Failed
			}
		}

		if errorTrackerCtrl != nil {
			overview.Errors.Availability = "ready"
			summary, err := errorTrackerCtrl.GetSummary(ctx, session, spaceRef)
			if err != nil {
				render.TranslatedUserError(ctx, w, err)
				return
			}
			if summary != nil {
				overview.Errors.Open = summary.OpenErrors
				overview.Errors.Fatal = summary.FatalCount
				overview.Errors.Warning = summary.WarningCount
				overview.Errors.LastSeen = summary.LastUpdated
			}
		}

		render.JSON(w, http.StatusOK, overview)
	}
}
