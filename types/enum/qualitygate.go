// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Updated quality gate enum types.
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

package enum

// QualityRuleCategory defines the category of a quality rule.
type QualityRuleCategory string

func (QualityRuleCategory) Enum() []any { return toInterfaceSlice(qualityRuleCategories) }
func (c QualityRuleCategory) Sanitize() (QualityRuleCategory, bool) {
	return Sanitize(c, GetAllQualityRuleCategories)
}
func GetAllQualityRuleCategories() ([]QualityRuleCategory, QualityRuleCategory) {
	return qualityRuleCategories, ""
}

// QualityRuleCategory enumeration.
const (
	QualityRuleCategoryCoverage      QualityRuleCategory = "coverage"
	QualityRuleCategoryComplexity    QualityRuleCategory = "complexity"
	QualityRuleCategoryDocumentation QualityRuleCategory = "documentation"
	QualityRuleCategoryNaming        QualityRuleCategory = "naming"
	QualityRuleCategoryTesting       QualityRuleCategory = "testing"
	QualityRuleCategorySecurity      QualityRuleCategory = "security"
	QualityRuleCategoryStyle         QualityRuleCategory = "style"
	QualityRuleCategoryCustom        QualityRuleCategory = "custom"
)

var qualityRuleCategories = sortEnum([]QualityRuleCategory{
	QualityRuleCategoryCoverage,
	QualityRuleCategoryComplexity,
	QualityRuleCategoryDocumentation,
	QualityRuleCategoryNaming,
	QualityRuleCategoryTesting,
	QualityRuleCategorySecurity,
	QualityRuleCategoryStyle,
	QualityRuleCategoryCustom,
})

// QualityEnforcement defines how a quality rule is enforced.
type QualityEnforcement string

func (QualityEnforcement) Enum() []any { return toInterfaceSlice(qualityEnforcements) }
func (e QualityEnforcement) Sanitize() (QualityEnforcement, bool) {
	return Sanitize(e, GetAllQualityEnforcements)
}
func GetAllQualityEnforcements() ([]QualityEnforcement, QualityEnforcement) {
	return qualityEnforcements, ""
}

// QualityEnforcement enumeration.
const (
	QualityEnforcementBlock QualityEnforcement = "block" // fails pipeline
	QualityEnforcementWarn  QualityEnforcement = "warn"  // warning only
	QualityEnforcementInfo  QualityEnforcement = "info"  // informational
)

var qualityEnforcements = sortEnum([]QualityEnforcement{
	QualityEnforcementBlock,
	QualityEnforcementWarn,
	QualityEnforcementInfo,
})

// QualityStatus defines the overall status of a quality evaluation.
type QualityStatus string

func (QualityStatus) Enum() []any { return toInterfaceSlice(qualityStatuses) }
func (s QualityStatus) Sanitize() (QualityStatus, bool) {
	return Sanitize(s, GetAllQualityStatuses)
}
func GetAllQualityStatuses() ([]QualityStatus, QualityStatus) {
	return qualityStatuses, ""
}

// QualityStatus enumeration.
const (
	QualityStatusPassed QualityStatus = "passed"
	QualityStatusFailed QualityStatus = "failed"
	QualityStatusWarned QualityStatus = "warning"
)

var qualityStatuses = sortEnum([]QualityStatus{
	QualityStatusPassed,
	QualityStatusFailed,
	QualityStatusWarned,
})

// QualityTrigger defines what triggered a quality evaluation.
type QualityTrigger string

func (QualityTrigger) Enum() []any { return toInterfaceSlice(qualityTriggers) }
func (t QualityTrigger) Sanitize() (QualityTrigger, bool) {
	return Sanitize(t, GetAllQualityTriggers)
}
func GetAllQualityTriggers() ([]QualityTrigger, QualityTrigger) {
	return qualityTriggers, ""
}

// QualityTrigger enumeration.
const (
	QualityTriggerPipeline QualityTrigger = "pipeline"
	QualityTriggerManual   QualityTrigger = "manual"
	QualityTriggerPullReq  QualityTrigger = "pr"
)

var qualityTriggers = sortEnum([]QualityTrigger{
	QualityTriggerPipeline,
	QualityTriggerManual,
	QualityTriggerPullReq,
})

// Helper methods for status determination.
func (s QualityStatus) IsFailed() bool {
	return s == QualityStatusFailed
}

func (s QualityStatus) IsPassed() bool {
	return s == QualityStatusPassed
}

func (s QualityStatus) IsWarned() bool {
	return s == QualityStatusWarned
}

// DetermineOverallStatus determines overall status based on rule results.
func DetermineOverallStatus(failed, warned, total int) QualityStatus {
	if failed > 0 {
		return QualityStatusFailed
	}
	if warned > 0 {
		return QualityStatusWarned
	}
	return QualityStatusPassed
}
