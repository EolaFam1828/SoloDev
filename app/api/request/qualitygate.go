// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package request

import "net/http"

const (
	PathParamQualityGateRuleIdentifier = "rule_identifier"
	PathParamQualityGateEvalIdentifier = "identifier"
)

// GetQualityGateRuleIdentifierFromPath extracts the quality gate rule identifier from the path.
func GetQualityGateRuleIdentifierFromPath(r *http.Request) (string, error) {
	return PathParamOrError(r, PathParamQualityGateRuleIdentifier)
}

// GetQualityGateEvalIdentifierFromPath extracts the quality gate evaluation identifier from the path.
func GetQualityGateEvalIdentifierFromPath(r *http.Request) (string, error) {
	return PathParamOrError(r, PathParamQualityGateEvalIdentifier)
}
