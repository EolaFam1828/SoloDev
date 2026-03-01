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

package enum

import "strings"

// SecurityScanType defines the types of security scans available.
type SecurityScanType string

func (SecurityScanType) Enum() []any { return toInterfaceSlice(securityScanTypes) }
func (s SecurityScanType) Sanitize() (SecurityScanType, bool) {
	return Sanitize(s, GetAllSecurityScanTypes)
}

func GetAllSecurityScanTypes() ([]SecurityScanType, SecurityScanType) {
	return securityScanTypes, ""
}

const (
	// SecurityScanTypeSAST describes a Static Application Security Testing scan.
	SecurityScanTypeSAST SecurityScanType = "sast"

	// SecurityScanTypeSCA describes a Software Composition Analysis scan.
	SecurityScanTypeSCA SecurityScanType = "sca"

	// SecurityScanTypeSecretDetection describes a secret detection scan.
	SecurityScanTypeSecretDetection SecurityScanType = "secret_detection"
)

var securityScanTypes = sortEnum([]SecurityScanType{
	SecurityScanTypeSAST,
	SecurityScanTypeSCA,
	SecurityScanTypeSecretDetection,
})

// SecurityScanStatus defines the status of a security scan.
type SecurityScanStatus string

func (SecurityScanStatus) Enum() []any { return toInterfaceSlice(securityScanStatuses) }
func (s SecurityScanStatus) Sanitize() (SecurityScanStatus, bool) {
	return Sanitize(s, GetAllSecurityScanStatuses)
}

func GetAllSecurityScanStatuses() ([]SecurityScanStatus, SecurityScanStatus) {
	return securityScanStatuses, ""
}

const (
	// SecurityScanStatusPending describes a scan that is pending execution.
	SecurityScanStatusPending SecurityScanStatus = "pending"

	// SecurityScanStatusRunning describes a scan that is currently running.
	SecurityScanStatusRunning SecurityScanStatus = "running"

	// SecurityScanStatusCompleted describes a scan that has completed successfully.
	SecurityScanStatusCompleted SecurityScanStatus = "completed"

	// SecurityScanStatusFailed describes a scan that has failed.
	SecurityScanStatusFailed SecurityScanStatus = "failed"
)

var securityScanStatuses = sortEnum([]SecurityScanStatus{
	SecurityScanStatusPending,
	SecurityScanStatusRunning,
	SecurityScanStatusCompleted,
	SecurityScanStatusFailed,
})

// SecurityScanTrigger defines what triggered a security scan.
type SecurityScanTrigger string

func (SecurityScanTrigger) Enum() []any { return toInterfaceSlice(securityScanTriggers) }
func (s SecurityScanTrigger) Sanitize() (SecurityScanTrigger, bool) {
	return Sanitize(s, GetAllSecurityScanTriggers)
}

func GetAllSecurityScanTriggers() ([]SecurityScanTrigger, SecurityScanTrigger) {
	return securityScanTriggers, ""
}

const (
	// SecurityScanTriggerManual describes a scan triggered manually.
	SecurityScanTriggerManual SecurityScanTrigger = "manual"

	// SecurityScanTriggerPipeline describes a scan triggered by a pipeline.
	SecurityScanTriggerPipeline SecurityScanTrigger = "pipeline"

	// SecurityScanTriggerWebhook describes a scan triggered by a webhook.
	SecurityScanTriggerWebhook SecurityScanTrigger = "webhook"
)

var securityScanTriggers = sortEnum([]SecurityScanTrigger{
	SecurityScanTriggerManual,
	SecurityScanTriggerPipeline,
	SecurityScanTriggerWebhook,
})

// SecurityScanAttr defines scan attributes that can be used for sorting and filtering.
type SecurityScanAttr int

const (
	SecurityScanAttrNone SecurityScanAttr = iota
	SecurityScanAttrID
	SecurityScanAttrIdentifier
	SecurityScanAttrScanType
	SecurityScanAttrStatus
	SecurityScanAttrCreated
	SecurityScanAttrUpdated
)

// ParseSecurityScanAttr parses the security scan attribute string
// and returns the equivalent enumeration.
func ParseSecurityScanAttr(s string) SecurityScanAttr {
	switch strings.ToLower(s) {
	case id:
		return SecurityScanAttrID
	case identifier:
		return SecurityScanAttrIdentifier
	case "scan_type", "scantype":
		return SecurityScanAttrScanType
	case status:
		return SecurityScanAttrStatus
	case created, createdAt:
		return SecurityScanAttrCreated
	case updated, updatedAt:
		return SecurityScanAttrUpdated
	default:
		return SecurityScanAttrNone
	}
}

// String returns the string representation of the attribute.
func (a SecurityScanAttr) String() string {
	switch a {
	case SecurityScanAttrID:
		return id
	case SecurityScanAttrIdentifier:
		return identifier
	case SecurityScanAttrScanType:
		return "scan_type"
	case SecurityScanAttrStatus:
		return status
	case SecurityScanAttrCreated:
		return created
	case SecurityScanAttrUpdated:
		return updated
	case SecurityScanAttrNone:
		return ""
	default:
		return undefined
	}
}

// SecurityFindingSeverity defines the severity levels of security findings.
type SecurityFindingSeverity string

func (SecurityFindingSeverity) Enum() []any { return toInterfaceSlice(securityFindingSeverities) }
func (s SecurityFindingSeverity) Sanitize() (SecurityFindingSeverity, bool) {
	return Sanitize(s, GetAllSecurityFindingSeverities)
}

func GetAllSecurityFindingSeverities() ([]SecurityFindingSeverity, SecurityFindingSeverity) {
	return securityFindingSeverities, ""
}

const (
	// SecurityFindingSeverityCritical describes a critical severity finding.
	SecurityFindingSeverityCritical SecurityFindingSeverity = "critical"

	// SecurityFindingSeverityHigh describes a high severity finding.
	SecurityFindingSeverityHigh SecurityFindingSeverity = "high"

	// SecurityFindingSeverityMedium describes a medium severity finding.
	SecurityFindingSeverityMedium SecurityFindingSeverity = "medium"

	// SecurityFindingSeverityLow describes a low severity finding.
	SecurityFindingSeverityLow SecurityFindingSeverity = "low"

	// SecurityFindingSeverityInfo describes an info severity finding.
	SecurityFindingSeverityInfo SecurityFindingSeverity = "info"
)

var securityFindingSeverities = sortEnum([]SecurityFindingSeverity{
	SecurityFindingSeverityCritical,
	SecurityFindingSeverityHigh,
	SecurityFindingSeverityMedium,
	SecurityFindingSeverityLow,
	SecurityFindingSeverityInfo,
})

// SecurityFindingCategory defines the categories of security findings.
type SecurityFindingCategory string

func (SecurityFindingCategory) Enum() []any { return toInterfaceSlice(securityFindingCategories) }
func (s SecurityFindingCategory) Sanitize() (SecurityFindingCategory, bool) {
	return Sanitize(s, GetAllSecurityFindingCategories)
}

func GetAllSecurityFindingCategories() ([]SecurityFindingCategory, SecurityFindingCategory) {
	return securityFindingCategories, ""
}

const (
	// SecurityFindingCategoryVulnerability describes a vulnerability finding.
	SecurityFindingCategoryVulnerability SecurityFindingCategory = "vulnerability"

	// SecurityFindingCategorySecret describes a secret finding.
	SecurityFindingCategorySecret SecurityFindingCategory = "secret"

	// SecurityFindingCategoryCodeSmell describes a code smell finding.
	SecurityFindingCategoryCodeSmell SecurityFindingCategory = "code_smell"

	// SecurityFindingCategoryBug describes a bug finding.
	SecurityFindingCategoryBug SecurityFindingCategory = "bug"
)

var securityFindingCategories = sortEnum([]SecurityFindingCategory{
	SecurityFindingCategoryVulnerability,
	SecurityFindingCategorySecret,
	SecurityFindingCategoryCodeSmell,
	SecurityFindingCategoryBug,
})

// SecurityFindingStatus defines the status of a security finding.
type SecurityFindingStatus string

func (SecurityFindingStatus) Enum() []any { return toInterfaceSlice(securityFindingStatuses) }
func (s SecurityFindingStatus) Sanitize() (SecurityFindingStatus, bool) {
	return Sanitize(s, GetAllSecurityFindingStatuses)
}

func GetAllSecurityFindingStatuses() ([]SecurityFindingStatus, SecurityFindingStatus) {
	return securityFindingStatuses, ""
}

const (
	// SecurityFindingStatusOpen describes an open security finding.
	SecurityFindingStatusOpen SecurityFindingStatus = "open"

	// SecurityFindingStatusResolved describes a resolved security finding.
	SecurityFindingStatusResolved SecurityFindingStatus = "resolved"

	// SecurityFindingStatusIgnored describes an ignored security finding.
	SecurityFindingStatusIgnored SecurityFindingStatus = "ignored"

	// SecurityFindingStatusFalsePositive describes a false positive security finding.
	SecurityFindingStatusFalsePositive SecurityFindingStatus = "false_positive"
)

var securityFindingStatuses = sortEnum([]SecurityFindingStatus{
	SecurityFindingStatusOpen,
	SecurityFindingStatusResolved,
	SecurityFindingStatusIgnored,
	SecurityFindingStatusFalsePositive,
})

// SecurityFindingAttr defines finding attributes that can be used for sorting and filtering.
type SecurityFindingAttr int

const (
	SecurityFindingAttrNone SecurityFindingAttr = iota
	SecurityFindingAttrID
	SecurityFindingAttrIdentifier
	SecurityFindingAttrSeverity
	SecurityFindingAttrCategory
	SecurityFindingAttrStatus
	SecurityFindingAttrCreated
	SecurityFindingAttrUpdated
)

// ParseSecurityFindingAttr parses the security finding attribute string
// and returns the equivalent enumeration.
func ParseSecurityFindingAttr(s string) SecurityFindingAttr {
	switch strings.ToLower(s) {
	case id:
		return SecurityFindingAttrID
	case identifier:
		return SecurityFindingAttrIdentifier
	case severity:
		return SecurityFindingAttrSeverity
	case category:
		return SecurityFindingAttrCategory
	case status:
		return SecurityFindingAttrStatus
	case created, createdAt:
		return SecurityFindingAttrCreated
	case updated, updatedAt:
		return SecurityFindingAttrUpdated
	default:
		return SecurityFindingAttrNone
	}
}

// String returns the string representation of the attribute.
func (a SecurityFindingAttr) String() string {
	switch a {
	case SecurityFindingAttrID:
		return id
	case SecurityFindingAttrIdentifier:
		return identifier
	case SecurityFindingAttrSeverity:
		return severity
	case SecurityFindingAttrCategory:
		return category
	case SecurityFindingAttrStatus:
		return status
	case SecurityFindingAttrCreated:
		return created
	case SecurityFindingAttrUpdated:
		return updated
	case SecurityFindingAttrNone:
		return ""
	default:
		return undefined
	}
}
