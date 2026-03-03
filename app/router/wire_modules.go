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
