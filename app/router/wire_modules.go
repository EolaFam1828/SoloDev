// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package router

import (
	controlleragentorch "github.com/harness/gitness/app/api/controller/agentorchestrator"
	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/controller/autopipeline"
	"github.com/harness/gitness/app/api/controller/errortracker"
	"github.com/harness/gitness/app/api/controller/featureflag"
	"github.com/harness/gitness/app/api/controller/healthcheck"
	"github.com/harness/gitness/app/api/controller/qualitygate"
	"github.com/harness/gitness/app/api/controller/securityscan"
	controllersignalcorrelator "github.com/harness/gitness/app/api/controller/signalcorrelator"
	"github.com/harness/gitness/app/api/controller/techdebt"
	"github.com/harness/gitness/app/api/controller/vectorsearch"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/anomalydetector"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
)

func ProvideSoloDevModules(
	featureFlagCtrl *featureflag.Controller,
	techDebtCtrl *techdebt.Controller,
	securityScanCtrl *securityscan.Controller,
	healthCheckCtrl *healthcheck.Controller,
	errorTrackerCtrl *errortracker.Controller,
	qualityGateCtrl *qualitygate.Controller,
	remediationCtrl *airemediation.Controller,
	autoPipelineCtrl *autopipeline.Controller,
	signalCorrelatorCtrl *controllersignalcorrelator.Controller,
	vectorSearchCtrl *vectorsearch.Controller,
	agentOrchCtrl *controlleragentorch.Controller,
	anomalyDetector *anomalydetector.Service,
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	gateConfigStore store.SoloGateConfigStore,
) *SoloDevModules {
	return &SoloDevModules{
		FeatureFlagCtrl:      featureFlagCtrl,
		TechDebtCtrl:         techDebtCtrl,
		SecurityScanCtrl:     securityScanCtrl,
		HealthCheckCtrl:      healthCheckCtrl,
		ErrorTrackerCtrl:     errorTrackerCtrl,
		QualityGateCtrl:      qualityGateCtrl,
		RemediationCtrl:      remediationCtrl,
		AutoPipelineCtrl:     autoPipelineCtrl,
		SignalCorrelatorCtrl: signalCorrelatorCtrl,
		VectorSearchCtrl:     vectorSearchCtrl,
		AgentOrchCtrl:        agentOrchCtrl,
		AnomalyDetector:      anomalyDetector,
		Authorizer:           authorizer,
		SpaceFinder:          spaceFinder,
		GateConfigStore:      gateConfigStore,
	}
}
