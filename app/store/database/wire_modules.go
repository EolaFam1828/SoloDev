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

package database

import (
	"github.com/harness/gitness/app/store"
	"github.com/jmoiron/sqlx"
)

// ProvideFeatureFlagStore provides a feature flag store.
func ProvideFeatureFlagStore(db *sqlx.DB) store.FeatureFlagStore {
	return NewFeatureFlagStore(db)
}

// ProvideTechDebtStore provides a technical debt store.
func ProvideTechDebtStore(db *sqlx.DB) store.TechDebtStore {
	return NewTechDebtStore(db)
}

// ProvideSecurityScanStore provides a security scan store.
func ProvideSecurityScanStore(db *sqlx.DB) store.SecurityScanStore {
	return NewSecurityScanStore(db)
}

// ProvideHealthCheckStore provides a health check store.
func ProvideHealthCheckStore(db *sqlx.DB) store.HealthCheckStore {
	return NewHealthCheckStore(db)
}

// ProvideErrorTrackerStore provides an error tracker store.
func ProvideErrorTrackerStore(db *sqlx.DB) store.ErrorTrackerStore {
	return NewErrorTrackerStore(db)
}

// ProvideQualityGateStore provides a quality gate store.
func ProvideQualityGateStore(db *sqlx.DB) store.QualityGateStore {
	return NewQualityGateStore(db)
}

// ProvideRemediationStore provides an AI remediation store.
func ProvideRemediationStore(db *sqlx.DB) store.RemediationStore {
	return NewRemediationStore(db)
}
