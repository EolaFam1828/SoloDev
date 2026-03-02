// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package router

import (
	controllerairemediation "github.com/EolaFam1828/SoloDev/app/api/controller/airemediation"
	controllerautopipeline "github.com/EolaFam1828/SoloDev/app/api/controller/autopipeline"
	controllererrortracker "github.com/EolaFam1828/SoloDev/app/api/controller/errortracker"
	controllerfeatureflag "github.com/EolaFam1828/SoloDev/app/api/controller/featureflag"
	controllerhealthcheck "github.com/EolaFam1828/SoloDev/app/api/controller/healthcheck"
	controllerqualitygate "github.com/EolaFam1828/SoloDev/app/api/controller/qualitygate"
	controllersecurityscan "github.com/EolaFam1828/SoloDev/app/api/controller/securityscan"
	controllertechdebt "github.com/EolaFam1828/SoloDev/app/api/controller/techdebt"
	"github.com/EolaFam1828/SoloDev/app/services/errorbridge"
)

// ProvideSoloDevModules assembles the SoloDev controller bundle and attaches
// the error bridge before the router and MCP layers consume it.
func ProvideSoloDevModules(
	featureFlagCtrl *controllerfeatureflag.Controller,
	techDebtCtrl *controllertechdebt.Controller,
	securityScanCtrl *controllersecurityscan.Controller,
	healthCheckCtrl *controllerhealthcheck.Controller,
	errorTrackerCtrl *controllererrortracker.Controller,
	qualityGateCtrl *controllerqualitygate.Controller,
	remediationCtrl *controllerairemediation.Controller,
	autoPipelineCtrl *controllerautopipeline.Controller,
	errorBridge *errorbridge.Bridge,
) *SoloDevModules {
	if errorTrackerCtrl != nil {
		errorTrackerCtrl.SetErrorBridge(errorBridge)
	}

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
