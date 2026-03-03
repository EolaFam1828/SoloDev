// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package openapi

import (
	"net/http"

	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/types"

	"github.com/gotidy/ptr"
	"github.com/swaggest/openapi-go/openapi3"
)

type securityScanRequest struct {
	spaceRequest
	RepoRef string `query:"repo_ref" required:"true"`
}

type triggerSecurityScanRequest struct {
	securityScanRequest
	types.ScanResultInput
}

type securityScanIdentifierRequest struct {
	securityScanRequest
	ScanIdentifier string `path:"scan_identifier" required:"true"`
}

type securityScanSummaryRequest struct {
	spaceRequest
	RepoRef string `query:"repo_ref"`
}

type createRemediationFromSecurityFindingRequest struct {
	spaceRequest
	types.CreateRemediationFromSecurityFindingInput
}

type remediationIdentifierRequest struct {
	spaceRequest
	Identifier string `path:"remediation_identifier" required:"true"`
}

type scanResultListResponse struct {
	Data  []*types.ScanResult `json:"data"`
	Count int64               `json:"count"`
}

type scanFindingListResponse struct {
	Data  []*types.ScanFinding `json:"data"`
	Count int64                `json:"count"`
}

var queryParameterRepoRef = openapi3.ParameterOrRef{
	Parameter: &openapi3.Parameter{
		Name:        request.QueryParamRepoRef,
		In:          openapi3.ParameterInQuery,
		Description: ptr.String("Repository reference inside the selected space."),
		Required:    ptr.Bool(false),
		Schema: &openapi3.SchemaOrRef{
			Schema: &openapi3.Schema{
				Type: ptrSchemaType(openapi3.SchemaTypeString),
			},
		},
	},
}

func soloDevOperations(reflector *openapi3.Reflector) {
	opTriggerSecurityScan := openapi3.Operation{}
	opTriggerSecurityScan.WithTags("solodev")
	opTriggerSecurityScan.WithMapOfAnything(map[string]any{"operationId": "triggerSecurityScan"})
	_ = reflector.SetRequest(&opTriggerSecurityScan, new(triggerSecurityScanRequest), http.MethodPost)
	_ = reflector.SetJSONResponse(&opTriggerSecurityScan, new(types.ScanResult), http.StatusCreated)
	_ = reflector.SetJSONResponse(&opTriggerSecurityScan, new(usererror.Error), http.StatusBadRequest)
	_ = reflector.SetJSONResponse(&opTriggerSecurityScan, new(usererror.Error), http.StatusServiceUnavailable)
	_ = reflector.Spec.AddOperation(http.MethodPost, "/spaces/{space_ref}/security-scans", opTriggerSecurityScan)

	opListSecurityScans := openapi3.Operation{}
	opListSecurityScans.WithTags("solodev")
	opListSecurityScans.WithMapOfAnything(map[string]any{"operationId": "listSecurityScans"})
	opListSecurityScans.WithParameters(QueryParameterPage, QueryParameterLimit, queryParameterRepoRef)
	_ = reflector.SetRequest(&opListSecurityScans, new(securityScanRequest), http.MethodGet)
	_ = reflector.SetJSONResponse(&opListSecurityScans, new(scanResultListResponse), http.StatusOK)
	_ = reflector.Spec.AddOperation(http.MethodGet, "/spaces/{space_ref}/security-scans", opListSecurityScans)

	opFindSecurityScan := openapi3.Operation{}
	opFindSecurityScan.WithTags("solodev")
	opFindSecurityScan.WithMapOfAnything(map[string]any{"operationId": "findSecurityScan"})
	opFindSecurityScan.WithParameters(queryParameterRepoRef)
	_ = reflector.SetRequest(&opFindSecurityScan, new(securityScanIdentifierRequest), http.MethodGet)
	_ = reflector.SetJSONResponse(&opFindSecurityScan, new(types.ScanResult), http.StatusOK)
	_ = reflector.Spec.AddOperation(http.MethodGet, "/spaces/{space_ref}/security-scans/{scan_identifier}", opFindSecurityScan)

	opListSecurityFindings := openapi3.Operation{}
	opListSecurityFindings.WithTags("solodev")
	opListSecurityFindings.WithMapOfAnything(map[string]any{"operationId": "listSecurityFindings"})
	opListSecurityFindings.WithParameters(QueryParameterPage, QueryParameterLimit, queryParameterRepoRef)
	_ = reflector.SetRequest(&opListSecurityFindings, new(securityScanIdentifierRequest), http.MethodGet)
	_ = reflector.SetJSONResponse(&opListSecurityFindings, new(scanFindingListResponse), http.StatusOK)
	_ = reflector.Spec.AddOperation(http.MethodGet, "/spaces/{space_ref}/security-scans/{scan_identifier}/findings", opListSecurityFindings)

	opGetSecuritySummary := openapi3.Operation{}
	opGetSecuritySummary.WithTags("solodev")
	opGetSecuritySummary.WithMapOfAnything(map[string]any{"operationId": "getSecuritySummary"})
	opGetSecuritySummary.WithParameters(queryParameterRepoRef)
	_ = reflector.SetRequest(&opGetSecuritySummary, new(securityScanSummaryRequest), http.MethodGet)
	_ = reflector.SetJSONResponse(&opGetSecuritySummary, new(types.SecuritySummary), http.StatusOK)
	_ = reflector.Spec.AddOperation(http.MethodGet, "/spaces/{space_ref}/security-scans/summary", opGetSecuritySummary)

	opCreateRemediationFromFinding := openapi3.Operation{}
	opCreateRemediationFromFinding.WithTags("solodev")
	opCreateRemediationFromFinding.WithMapOfAnything(map[string]any{"operationId": "createRemediationFromSecurityFinding"})
	_ = reflector.SetRequest(&opCreateRemediationFromFinding, new(createRemediationFromSecurityFindingRequest), http.MethodPost)
	_ = reflector.SetJSONResponse(&opCreateRemediationFromFinding, new(types.Remediation), http.StatusCreated)
	_ = reflector.SetJSONResponse(&opCreateRemediationFromFinding, new(types.Remediation), http.StatusOK)
	_ = reflector.SetJSONResponse(&opCreateRemediationFromFinding, new(usererror.Error), http.StatusBadRequest)
	_ = reflector.SetJSONResponse(&opCreateRemediationFromFinding, new(usererror.Error), http.StatusServiceUnavailable)
	_ = reflector.Spec.AddOperation(http.MethodPost, "/spaces/{space_ref}/remediations/from-security-finding", opCreateRemediationFromFinding)

	opApplyRemediation := openapi3.Operation{}
	opApplyRemediation.WithTags("solodev")
	opApplyRemediation.WithMapOfAnything(map[string]any{"operationId": "applyRemediation"})
	_ = reflector.SetRequest(&opApplyRemediation, new(remediationIdentifierRequest), http.MethodPost)
	_ = reflector.SetJSONResponse(&opApplyRemediation, new(types.Remediation), http.StatusOK)
	_ = reflector.SetJSONResponse(&opApplyRemediation, new(usererror.Error), http.StatusConflict)
	_ = reflector.SetJSONResponse(&opApplyRemediation, new(usererror.Error), http.StatusServiceUnavailable)
	_ = reflector.Spec.AddOperation(http.MethodPost, "/spaces/{space_ref}/remediations/{remediation_identifier}/apply", opApplyRemediation)

	opGetMcpCatalog := openapi3.Operation{}
	opGetMcpCatalog.WithTags("system")
	opGetMcpCatalog.WithMapOfAnything(map[string]any{"operationId": "getMCPCatalog"})
	_ = reflector.SetRequest(&opGetMcpCatalog, nil, http.MethodGet)
	_ = reflector.SetJSONResponse(&opGetMcpCatalog, new(types.MCPCatalog), http.StatusOK)
	_ = reflector.Spec.AddOperation(http.MethodGet, "/system/mcp/catalog", opGetMcpCatalog)

	opGetSoloDevOverview := openapi3.Operation{}
	opGetSoloDevOverview.WithTags("solodev")
	opGetSoloDevOverview.WithMapOfAnything(map[string]any{"operationId": "getSoloDevOverview"})
	_ = reflector.SetRequest(&opGetSoloDevOverview, new(spaceRequest), http.MethodGet)
	_ = reflector.SetJSONResponse(&opGetSoloDevOverview, new(types.SoloDevOverview), http.StatusOK)
	_ = reflector.Spec.AddOperation(http.MethodGet, "/spaces/{space_ref}/solodev/overview", opGetSoloDevOverview)
}
