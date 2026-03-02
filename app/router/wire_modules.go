// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/controller/autopipeline"
	"github.com/harness/gitness/app/api/controller/errortracker"
	"github.com/harness/gitness/app/api/controller/featureflag"
	"github.com/harness/gitness/app/api/controller/healthcheck"
	"github.com/harness/gitness/app/api/controller/qualitygate"
	"github.com/harness/gitness/app/api/controller/securityscan"
	"github.com/harness/gitness/app/api/controller/techdebt"
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
) *SoloDevModules {
	return &SoloDevModules{
		FeatureFlagCtrl:  featureFlagCtrl,
		TechDebtCtrl:     techDebtCtrl,
		SecurityScanCtrl: securityScanCtrl,
		HealthCheckCtrl:  healthCheckCtrl,
		ErrorTrackerCtrl: errorTrackerCtrl,
		QualityGateCtrl:  qualityGateCtrl,
		RemediationCtrl:  remediationCtrl,
		AutoPipelineCtrl: autoPipelineCtrl,
	}
}
