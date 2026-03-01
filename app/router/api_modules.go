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

package router

import (
	"fmt"

	"github.com/harness/gitness/app/api/controller/errortracker"
	"github.com/harness/gitness/app/api/controller/featureflag"
	"github.com/harness/gitness/app/api/controller/healthcheck"
	"github.com/harness/gitness/app/api/controller/qualitygate"
	"github.com/harness/gitness/app/api/controller/securityscan"
	"github.com/harness/gitness/app/api/controller/techdebt"
	handlererrortracker "github.com/harness/gitness/app/api/handler/errortracker"
	handlerfeatureflag "github.com/harness/gitness/app/api/handler/featureflag"
	handlerhealthcheck "github.com/harness/gitness/app/api/handler/healthcheck"
	handlerqualitygate "github.com/harness/gitness/app/api/handler/qualitygate"
	handlersecurityscan "github.com/harness/gitness/app/api/handler/securityscan"
	handlertechdebt "github.com/harness/gitness/app/api/handler/techdebt"
	"github.com/harness/gitness/app/api/request"
	"github.com/go-chi/chi/v5"
)

const (
	// Path parameter names for the new modules
	PathParamFeatureFlagIdentifier = "identifier"
	PathParamTechDebtID            = "identifier"
	PathParamSecurityScanID        = "scan_identifier"
	PathParamHealthCheckID         = "identifier"
	PathParamErrorID               = "error_identifier"
	PathParamQualityGateEvalID     = "identifier"
	PathParamQualityGateRuleID     = "rule_identifier"
)

// setupFeatureFlags registers feature flag routes.
func setupFeatureFlags(r chi.Router, featureFlagCtrl *featureflag.Controller) {
	r.Route("/feature-flags", func(r chi.Router) {
		r.Post("/", handlerfeatureflag.HandleCreate(featureFlagCtrl))
		r.Get("/", handlerfeatureflag.HandleList(featureFlagCtrl))

		r.Route(fmt.Sprintf("/{%s}", PathParamFeatureFlagIdentifier), func(r chi.Router) {
			r.Get("/", handlerfeatureflag.HandleFind(featureFlagCtrl))
			r.Patch("/", handlerfeatureflag.HandleUpdate(featureFlagCtrl))
			r.Delete("/", handlerfeatureflag.HandleDelete(featureFlagCtrl))
			r.Post("/toggle", handlerfeatureflag.HandleToggle(featureFlagCtrl))
		})
	})
}

// setupTechDebt registers technical debt routes.
func setupTechDebt(r chi.Router, techDebtCtrl *techdebt.Controller) {
	r.Route("/tech-debt", func(r chi.Router) {
		r.Post("/", handlertechdebt.HandleCreate(techDebtCtrl))
		r.Get("/", handlertechdebt.HandleList(techDebtCtrl))

		r.Route(fmt.Sprintf("/{%s}", PathParamTechDebtID), func(r chi.Router) {
			r.Get("/", handlertechdebt.HandleFind(techDebtCtrl))
			r.Patch("/", handlertechdebt.HandleUpdate(techDebtCtrl))
			r.Delete("/", handlertechdebt.HandleDelete(techDebtCtrl))
			r.Get("/summary", handlertechdebt.HandleSummary(techDebtCtrl))
		})
	})
}

// setupSecurityScans registers security scan routes.
func setupSecurityScans(r chi.Router, securityScanCtrl *securityscan.Controller) {
	r.Route("/security-scans", func(r chi.Router) {
		r.Post("/", handlersecurityscan.HandleTriggerScan(securityScanCtrl))
		r.Get("/", handlersecurityscan.HandleListScans(securityScanCtrl))

		r.Route(fmt.Sprintf("/{%s}", PathParamSecurityScanID), func(r chi.Router) {
			r.Get("/", handlersecurityscan.HandleFindScan(securityScanCtrl))
			r.Get("/findings", handlersecurityscan.HandleListFindings(securityScanCtrl))
			r.Get("/summary", handlersecurityscan.HandleGetSecuritySummary(securityScanCtrl))
		})
	})
}

// setupHealthCheckMonitor registers health check monitoring routes.
func setupHealthCheckMonitor(r chi.Router, healthCheckCtrl *healthcheck.Controller) {
	r.Route("/health-checks", func(r chi.Router) {
		r.Post("/", handlerhealthcheck.HandleCreate(healthCheckCtrl))
		r.Get("/", handlerhealthcheck.HandleList(healthCheckCtrl))

		r.Route(fmt.Sprintf("/{%s}", PathParamHealthCheckID), func(r chi.Router) {
			r.Get("/", handlerhealthcheck.HandleFind(healthCheckCtrl))
			r.Patch("/", handlerhealthcheck.HandleUpdate(healthCheckCtrl))
			r.Delete("/", handlerhealthcheck.HandleDelete(healthCheckCtrl))
			r.Get("/results", handlerhealthcheck.HandleListResults(healthCheckCtrl))
			r.Get("/summary", handlerhealthcheck.HandleSummary(healthCheckCtrl))
		})
	})
}

// setupErrorTracker registers error tracking routes.
func setupErrorTracker(r chi.Router, errorTrackerCtrl *errortracker.Controller) {
	r.Route("/errors", func(r chi.Router) {
		r.Get("/", handlererrortracker.HandleErrorList(errorTrackerCtrl))
		r.Get("/summary", handlererrortracker.HandleErrorSummary(errorTrackerCtrl))

		r.Route(fmt.Sprintf("/{%s}", PathParamErrorID), func(r chi.Router) {
			r.Get("/", handlererrortracker.HandleErrorDetail(errorTrackerCtrl))
			r.Patch("/", handlererrortracker.HandleErrorUpdate(errorTrackerCtrl))
			r.Get("/occurrences", handlererrortracker.HandleErrorOccurrences(errorTrackerCtrl))
		})
	})
}

// setupQualityGates registers quality gate routes.
func setupQualityGates(r chi.Router, qualityGateCtrl *qualitygate.Controller) {
	r.Route("/quality-gates", func(r chi.Router) {
		r.Post("/evaluate", handlerqualitygate.HandleEvaluate(qualityGateCtrl))
		r.Get("/summary", handlerqualitygate.HandleSummaryGet(qualityGateCtrl))

		// Evaluation endpoints
		r.Route("/evaluations", func(r chi.Router) {
			r.Get("/", handlerqualitygate.HandleEvaluationList(qualityGateCtrl))
			r.Route(fmt.Sprintf("/{%s}", PathParamQualityGateEvalID), func(r chi.Router) {
				r.Get("/", handlerqualitygate.HandleEvaluationGet(qualityGateCtrl))
			})
		})

		// Rule endpoints
		r.Route("/rules", func(r chi.Router) {
			r.Post("/", handlerqualitygate.HandleRuleCreate(qualityGateCtrl))
			r.Get("/", handlerqualitygate.HandleRuleList(qualityGateCtrl))

			r.Route(fmt.Sprintf("/{%s}", PathParamQualityGateRuleID), func(r chi.Router) {
				r.Get("/", handlerqualitygate.HandleRuleGet(qualityGateCtrl))
				r.Patch("/", handlerqualitygate.HandleRuleUpdate(qualityGateCtrl))
				r.Delete("/", handlerqualitygate.HandleRuleDelete(qualityGateCtrl))
				r.Post("/toggle", handlerqualitygate.HandleRuleToggle(qualityGateCtrl))
			})
		})
	})
}
