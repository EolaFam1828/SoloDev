// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Added SoloDevModules struct, SetupSoloDevModules, and AI remediation/auto-pipeline routes.
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

	"github.com/EolaFam1828/SoloDev/app/api/controller/airemediation"
	"github.com/EolaFam1828/SoloDev/app/api/controller/autopipeline"
	"github.com/EolaFam1828/SoloDev/app/api/controller/errortracker"
	"github.com/EolaFam1828/SoloDev/app/api/controller/featureflag"
	"github.com/EolaFam1828/SoloDev/app/api/controller/healthcheck"
	"github.com/EolaFam1828/SoloDev/app/api/controller/qualitygate"
	"github.com/EolaFam1828/SoloDev/app/api/controller/securityscan"
	"github.com/EolaFam1828/SoloDev/app/api/controller/techdebt"
	handlerairemediation "github.com/EolaFam1828/SoloDev/app/api/handler/airemediation"
	handlerautopipeline "github.com/EolaFam1828/SoloDev/app/api/handler/autopipeline"
	handlererrortracker "github.com/EolaFam1828/SoloDev/app/api/handler/errortracker"
	handlerfeatureflag "github.com/EolaFam1828/SoloDev/app/api/handler/featureflag"
	handlerhealthcheck "github.com/EolaFam1828/SoloDev/app/api/handler/healthcheck"
	handlerqualitygate "github.com/EolaFam1828/SoloDev/app/api/handler/qualitygate"
	handlersecurityscan "github.com/EolaFam1828/SoloDev/app/api/handler/securityscan"
	handlersolodev "github.com/EolaFam1828/SoloDev/app/api/handler/solodev"
	handlertechdebt "github.com/EolaFam1828/SoloDev/app/api/handler/techdebt"
	"github.com/go-chi/chi/v5"
)

const (
	// Path parameter names for the modules
	PathParamFeatureFlagIdentifier = "identifier"
	PathParamTechDebtID            = "identifier"
	PathParamSecurityScanID        = "scan_identifier"
	PathParamHealthCheckID         = "identifier"
	PathParamErrorID               = "error_identifier"
	PathParamQualityGateEvalID     = "identifier"
	PathParamQualityGateRuleID     = "rule_identifier"
	PathParamRemediationID         = "remediation_identifier"
)

// SoloDevModules holds controllers for all SoloDev-specific modules.
type SoloDevModules struct {
	FeatureFlagCtrl  *featureflag.Controller
	TechDebtCtrl     *techdebt.Controller
	SecurityScanCtrl *securityscan.Controller
	HealthCheckCtrl  *healthcheck.Controller
	ErrorTrackerCtrl *errortracker.Controller
	QualityGateCtrl  *qualitygate.Controller
	RemediationCtrl  *airemediation.Controller
	AutoPipelineCtrl *autopipeline.Controller
}

// SetupSoloDevModules registers all SoloDev module routes under the current space router.
// This is the single entry point called from setupSpaces() in api.go.
func SetupSoloDevModules(r chi.Router, m *SoloDevModules) {
	if m == nil {
		return
	}

	r.Get("/solodev/overview", handlersolodev.HandleOverview(m.SecurityScanCtrl, m.ErrorTrackerCtrl, m.RemediationCtrl))

	// Existing modules
	if m.FeatureFlagCtrl != nil {
		setupFeatureFlags(r, m.FeatureFlagCtrl)
	}
	if m.TechDebtCtrl != nil {
		setupTechDebt(r, m.TechDebtCtrl)
	}
	if m.SecurityScanCtrl != nil {
		setupSecurityScans(r, m.SecurityScanCtrl)
	}
	if m.HealthCheckCtrl != nil {
		setupHealthCheckMonitor(r, m.HealthCheckCtrl)
	}
	if m.ErrorTrackerCtrl != nil {
		setupErrorTracker(r, m.ErrorTrackerCtrl)
	}
	if m.QualityGateCtrl != nil {
		setupQualityGates(r, m.QualityGateCtrl)
	}

	// New Solo-AI modules
	if m.RemediationCtrl != nil {
		setupRemediations(r, m.RemediationCtrl)
	}
	if m.AutoPipelineCtrl != nil {
		setupAutoPipeline(r, m.AutoPipelineCtrl)
	}
}

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
		r.Get("/summary", handlersecurityscan.HandleGetSecuritySummary(securityScanCtrl))

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
		r.Post("/", handlererrortracker.HandleErrorReport(errorTrackerCtrl))
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

// setupRemediations registers AI remediation routes.
func setupRemediations(r chi.Router, remediationCtrl *airemediation.Controller) {
	r.Route("/remediations", func(r chi.Router) {
		r.Post("/", handlerairemediation.HandleTrigger(remediationCtrl))
		r.Post("/from-security-finding", handlerairemediation.HandleTriggerFromSecurityFinding(remediationCtrl))
		r.Get("/", handlerairemediation.HandleList(remediationCtrl))
		r.Get("/summary", handlerairemediation.HandleSummary(remediationCtrl))

		r.Route(fmt.Sprintf("/{%s}", PathParamRemediationID), func(r chi.Router) {
			r.Get("/", handlerairemediation.HandleGet(remediationCtrl))
			r.Patch("/", handlerairemediation.HandleUpdate(remediationCtrl))
		})
	})
}

// setupAutoPipeline registers auto-pipeline generation routes.
func setupAutoPipeline(r chi.Router, autoPipelineCtrl *autopipeline.Controller) {
	r.Route("/auto-pipeline", func(r chi.Router) {
		r.Post("/generate", handlerautopipeline.HandleGenerate(autoPipelineCtrl))
	})
}
