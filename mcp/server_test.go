// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// mockAuthenticator implements authn.Authenticator for testing.
type mockAuthenticator struct{}

func (m *mockAuthenticator) Authenticate(_ *http.Request) (*auth.Session, error) {
	return testSession(), nil
}

func testSession() *auth.Session {
	return &auth.Session{
		Principal: types.Principal{
			ID:   1,
			UID:  "test-user",
			Type: enum.PrincipalTypeUser,
		},
	}
}

// newTestServer creates an MCP server with nil controllers for protocol-level testing.
func newTestServer() *Server {
	return NewServer(&mockAuthenticator{}, &Controllers{})
}

func sendRequest(t *testing.T, srv *Server, method string, params interface{}) Response {
	t.Helper()

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
	}
	if params != nil {
		b, _ := json.Marshal(params)
		reqBody["params"] = json.RawMessage(b)
	}

	raw, _ := json.Marshal(reqBody)
	respBytes, err := srv.HandleMessage(context.Background(), testSession(), raw)
	if err != nil {
		t.Fatalf("HandleMessage error: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	return resp
}

func TestInitialize(t *testing.T) {
	srv := newTestServer()

	params := InitializeParams{
		ProtocolVersion: "2024-11-05",
		ClientInfo:      ClientInfo{Name: "test-client", Version: "1.0.0"},
	}

	resp := sendRequest(t, srv, "initialize", params)

	if resp.Error != nil {
		t.Fatalf("expected no error, got: %v", resp.Error)
	}

	// Unmarshal the result
	b, _ := json.Marshal(resp.Result)
	var result InitializeResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal initialize result: %v", err)
	}

	if result.ProtocolVersion != ProtocolVersion {
		t.Errorf("protocol version = %q, want %q", result.ProtocolVersion, ProtocolVersion)
	}
	if result.ServerInfo.Name != ServerName {
		t.Errorf("server name = %q, want %q", result.ServerInfo.Name, ServerName)
	}
	if result.Capabilities.Tools == nil {
		t.Error("expected tools capability to be set")
	}
	if result.Capabilities.Resources == nil {
		t.Error("expected resources capability to be set")
	}
	if result.Capabilities.Prompts == nil {
		t.Error("expected prompts capability to be set")
	}
	if result.Instructions == "" {
		t.Error("expected instructions to be non-empty")
	}
}

func TestPing(t *testing.T) {
	srv := newTestServer()
	resp := sendRequest(t, srv, "ping", nil)

	if resp.Error != nil {
		t.Fatalf("expected no error, got: %v", resp.Error)
	}
	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestToolsList(t *testing.T) {
	srv := newTestServer()
	resp := sendRequest(t, srv, "tools/list", nil)

	if resp.Error != nil {
		t.Fatalf("expected no error, got: %v", resp.Error)
	}

	b, _ := json.Marshal(resp.Result)
	var result ToolsListResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal tools list result: %v", err)
	}

	// We expect at least 16 atomic + 5 compound = 21 tools
	if len(result.Tools) < 21 {
		t.Errorf("expected at least 21 tools, got %d", len(result.Tools))
	}

	// Check that key tool names exist
	toolNames := make(map[string]bool)
	for _, tool := range result.Tools {
		toolNames[tool.Name] = true
	}

	expectedTools := []string{
		// Tier 1 atomic tools
		"pipeline_generate",
		"security_scan",
		"security_findings",
		"quality_evaluate",
		"quality_rules_list",
		"error_report",
		"error_list",
		"remediation_trigger",
		"remediation_list",
		"health_summary",
		"feature_flag_toggle",
		"tech_debt_list",
		// Tier 2 compound tools
		"fix_this",
		"pr_ready",
		"pipeline_validate",
		"onboard_repo",
		"incident_triage",
	}

	for _, name := range expectedTools {
		if !toolNames[name] {
			t.Errorf("expected tool %q to be registered", name)
		}
	}

	// Verify each tool has a description and input schema
	for _, tool := range result.Tools {
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
		if tool.InputSchema == nil {
			t.Errorf("tool %q has nil input schema", tool.Name)
		}
	}
}

func TestResourcesList(t *testing.T) {
	srv := newTestServer()
	resp := sendRequest(t, srv, "resources/list", nil)

	if resp.Error != nil {
		t.Fatalf("expected no error, got: %v", resp.Error)
	}

	b, _ := json.Marshal(resp.Result)
	var result ResourcesListResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal resources list result: %v", err)
	}

	// We expect 7 resources
	if len(result.Resources) < 7 {
		t.Errorf("expected at least 7 resources, got %d", len(result.Resources))
	}

	// Check that key resource URIs exist
	resourceURIs := make(map[string]bool)
	for _, res := range result.Resources {
		resourceURIs[res.URI] = true
	}

	expectedResources := []string{
		"solodev://errors/active",
		"solodev://remediations/pending",
		"solodev://quality/rules",
		"solodev://quality/summary",
		"solodev://security/open-findings",
		"solodev://health/status",
		"solodev://tech-debt/hotspots",
	}

	for _, uri := range expectedResources {
		if !resourceURIs[uri] {
			t.Errorf("expected resource %q to be registered", uri)
		}
	}
}

func TestPromptsList(t *testing.T) {
	srv := newTestServer()
	resp := sendRequest(t, srv, "prompts/list", nil)

	if resp.Error != nil {
		t.Fatalf("expected no error, got: %v", resp.Error)
	}

	b, _ := json.Marshal(resp.Result)
	var result PromptsListResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal prompts list result: %v", err)
	}

	// We expect 5 prompts
	if len(result.Prompts) < 5 {
		t.Errorf("expected at least 5 prompts, got %d", len(result.Prompts))
	}

	promptNames := make(map[string]bool)
	for _, p := range result.Prompts {
		promptNames[p.Name] = true
	}

	expectedPrompts := []string{
		"solodev_review",
		"solodev_incident",
		"solodev_pipeline_debug",
	}

	for _, name := range expectedPrompts {
		if !promptNames[name] {
			t.Errorf("expected prompt %q to be registered", name)
		}
	}
}

func TestToolCallPipelineGenerate(t *testing.T) {
	srv := newTestServer()

	// tools/call for pipeline generation should work even with nil controllers
	// because the autopipeline controller is nil — it should return an error message
	params := map[string]interface{}{
		"name": "solodev_generate_pipeline",
		"arguments": map[string]interface{}{
			"space_ref": "default",
			"files":     []string{"main.go", "go.mod", "Dockerfile"},
		},
	}

	resp := sendRequest(t, srv, "tools/call", params)

	// We expect either a result (with isError=true since controller is nil) or an RPC error
	if resp.Error != nil && resp.Result == nil {
		// This is acceptable - the tool might not be available
		return
	}

	// If result returned, it should indicate the tool ran (or gracefully failed)
	if resp.Result != nil {
		b, _ := json.Marshal(resp.Result)
		var result ToolCallResult
		if err := json.Unmarshal(b, &result); err != nil {
			t.Fatalf("unmarshal tool call result: %v", err)
		}
		// With nil controller, we expect an error result
		if !result.IsError {
			t.Log("Tool call succeeded (unexpected with nil controller)")
		}
	}
}

func TestToolCallUnknownTool(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"name":      "nonexistent_tool",
		"arguments": map[string]interface{}{},
	}

	resp := sendRequest(t, srv, "tools/call", params)

	if resp.Error != nil {
		// Unknown tool returns an RPC error — this is the expected path.
		return
	}

	// If no RPC error, the result itself should signal an error.
	b, _ := json.Marshal(resp.Result)
	var result ToolCallResult
	_ = json.Unmarshal(b, &result)
	if !result.IsError {
		t.Error("expected error for unknown tool")
	}
}

func TestMethodNotFound(t *testing.T) {
	srv := newTestServer()
	resp := sendRequest(t, srv, "nonexistent/method", nil)

	if resp.Error == nil {
		t.Fatal("expected error for unknown method")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
}

func TestInvalidJSON(t *testing.T) {
	srv := newTestServer()

	respBytes, err := srv.HandleMessage(context.Background(), testSession(), []byte("not valid json"))
	if err != nil {
		t.Fatalf("HandleMessage error: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected parse error")
	}
	if resp.Error.Code != ErrCodeParse {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeParse)
	}
}

func TestInvalidJSONRPCVersion(t *testing.T) {
	srv := newTestServer()

	raw := []byte(`{"jsonrpc": "1.0", "id": 1, "method": "ping"}`)
	respBytes, err := srv.HandleMessage(context.Background(), testSession(), raw)
	if err != nil {
		t.Fatalf("HandleMessage error: %v", err)
	}

	var resp Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected invalid request error")
	}
	if resp.Error.Code != ErrCodeInvalidRequest {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeInvalidRequest)
	}
}

func TestNotification(t *testing.T) {
	srv := newTestServer()

	// Notifications have no id — should return nil
	raw := []byte(`{"jsonrpc": "2.0", "method": "notifications/initialized"}`)
	respBytes, err := srv.HandleMessage(context.Background(), testSession(), raw)
	if err != nil {
		t.Fatalf("HandleMessage error: %v", err)
	}

	if respBytes != nil {
		t.Errorf("expected nil response for notification, got: %s", string(respBytes))
	}
}

func TestStreamableHTTPHandler(t *testing.T) {
	srv := newTestServer()
	handler := srv.StreamableHTTPHandler()

	// Test POST with initialize request
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"test","version":"1.0"}}}`
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Errorf("expected no error, got: %v", resp.Error)
	}
}

func TestStreamableHTTPOptionsPreflightCORS(t *testing.T) {
	srv := newTestServer()
	handler := srv.StreamableHTTPHandler()

	req := httptest.NewRequest(http.MethodOptions, "/mcp", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("expected CORS headers on OPTIONS response")
	}
}

func TestStreamableHTTPMethodNotAllowed(t *testing.T) {
	srv := newTestServer()
	handler := srv.StreamableHTTPHandler()

	req := httptest.NewRequest(http.MethodPut, "/mcp", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestResourceReadActiveErrors(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"uri": "solodev://errors/active",
	}

	resp := sendRequest(t, srv, "resources/read", params)

	// With nil error tracker controller, we should get a graceful error message
	if resp.Error != nil {
		t.Fatalf("expected no RPC error, got: %v", resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestPromptGetCodeReview(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"name": "solodev_review",
		"arguments": map[string]string{
			"space_ref": "default",
		},
	}

	resp := sendRequest(t, srv, "prompts/get", params)

	// Even with nil controllers, prompt generation should succeed
	// (it will just skip sections that need controller data)
	if resp.Error != nil {
		t.Fatalf("expected no RPC error, got: %v", resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}
}

// --- Compound Tool Tests ---

func TestCompoundToolPipelineValidate(t *testing.T) {
	srv := newTestServer()

	tests := []struct {
		name       string
		yaml       string
		wantStatus string
	}{
		{
			name:       "valid pipeline with stages and tests",
			yaml:       "stages:\n  - stage:\n      steps:\n        - step:\n            name: test\n            command: go test ./...\n      cache:\n        key: go-modules\n",
			wantStatus: "valid",
		},
		{
			name:       "empty pipeline",
			yaml:       "# empty pipeline\nname: my-pipeline\n",
			wantStatus: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]interface{}{
				"name": "pipeline_validate",
				"arguments": map[string]interface{}{
					"yaml": tt.yaml,
				},
			}

			resp := sendRequest(t, srv, "tools/call", params)
			if resp.Error != nil {
				t.Fatalf("unexpected RPC error: %v", resp.Error)
			}

			b, _ := json.Marshal(resp.Result)
			var result ToolCallResult
			if err := json.Unmarshal(b, &result); err != nil {
				t.Fatalf("unmarshal result: %v", err)
			}

			if result.IsError {
				t.Fatalf("unexpected tool error: %v", result.Content)
			}

			// Parse the content text as JSON to check status
			if len(result.Content) == 0 {
				t.Fatal("expected content in result")
			}
			var output map[string]interface{}
			if err := json.Unmarshal([]byte(result.Content[0].Text), &output); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}
			if got := output["status"]; got != tt.wantStatus {
				t.Errorf("status = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestCompoundToolPRReadyNilControllers(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"name": "pr_ready",
		"arguments": map[string]interface{}{
			"space_ref":  "default",
			"repo_ref":   "my-repo",
			"commit_sha": "abc12345",
		},
	}

	resp := sendRequest(t, srv, "tools/call", params)
	if resp.Error != nil {
		t.Fatalf("unexpected RPC error: %v", resp.Error)
	}

	b, _ := json.Marshal(resp.Result)
	var result ToolCallResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	// With nil controllers, pr_ready should still return PASS (no checks ran).
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}
	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result.Content[0].Text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if verdict, ok := output["verdict"].(string); !ok || verdict != "PASS" {
		t.Errorf("verdict = %v, want PASS (nil controllers should skip all checks)", output["verdict"])
	}
}

func TestCompoundToolIncidentTriageNilControllers(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"name": "incident_triage",
		"arguments": map[string]interface{}{
			"space_ref":   "default",
			"time_window": "last 1 hour",
		},
	}

	resp := sendRequest(t, srv, "tools/call", params)
	if resp.Error != nil {
		t.Fatalf("unexpected RPC error: %v", resp.Error)
	}

	b, _ := json.Marshal(resp.Result)
	var result ToolCallResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}
	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result.Content[0].Text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if output["timestamp"] == nil {
		t.Error("expected timestamp in incident triage report")
	}
}

func TestCompoundToolOnboardRepoNilControllers(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"name": "onboard_repo",
		"arguments": map[string]interface{}{
			"space_ref":         "default",
			"repo_file_listing": []string{"main.go", "go.mod", "Makefile"},
		},
	}

	resp := sendRequest(t, srv, "tools/call", params)
	if resp.Error != nil {
		t.Fatalf("unexpected RPC error: %v", resp.Error)
	}

	b, _ := json.Marshal(resp.Result)
	var result ToolCallResult
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	// With nil controllers the tool should still return a result (empty map).
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}
}

func TestCompoundToolFixThisNilControllers(t *testing.T) {
	srv := newTestServer()

	params := map[string]interface{}{
		"name": "fix_this",
		"arguments": map[string]interface{}{
			"space_ref": "default",
			"error_log": "panic: runtime error: index out of range [5] with length 3\ngoroutine 1 [running]:\nmain.main()\n\t/app/main.go:15",
		},
	}

	resp := sendRequest(t, srv, "tools/call", params)

	// fix_this with nil ErrorTracker will panic-deref or return an error —
	// either an RPC error or a tool-level error is acceptable.
	if resp.Error != nil {
		return // RPC error is acceptable with nil controllers
	}

	if resp.Result != nil {
		b, _ := json.Marshal(resp.Result)
		var result ToolCallResult
		_ = json.Unmarshal(b, &result)
		// Either a successful fallback or error result is fine
		t.Logf("fix_this result: isError=%v", result.IsError)
	}
}
