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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/version"

	"github.com/rs/zerolog/log"
)

const (
	ProtocolVersion = "2024-11-05"
	ServerName      = "solodev"
)

// ToolHandler processes a tool call.
type ToolHandler func(ctx context.Context, session *auth.Session, args json.RawMessage) (*ToolCallResult, error)

// ResourceHandler reads a resource.
type ResourceHandler func(ctx context.Context, session *auth.Session) (*ResourceReadResult, error)

// PromptHandler generates a prompt.
type PromptHandler func(ctx context.Context, session *auth.Session, args map[string]string) (*PromptGetResult, error)

// Server is the MCP protocol router.
type Server struct {
	auth        *MCPAuthenticator
	controllers *Controllers

	tools     map[string]ToolHandler
	resources map[string]ResourceHandler
	prompts   map[string]PromptHandler

	toolDefs   []ToolDefinition
	resDefs    []ResourceDefinition
	promptDefs []PromptDefinition
	catalog    *types.MCPCatalog
}

// RegisterTool adds a tool to the server.
func (s *Server) RegisterTool(def ToolDefinition, handler ToolHandler) {
	if s.catalog != nil && !catalogHasActiveTool(s.catalog, def.Name) {
		return
	}
	s.tools[def.Name] = handler
	s.toolDefs = append(s.toolDefs, def)
}

// RegisterResource adds a resource to the server.
func (s *Server) RegisterResource(def ResourceDefinition, handler ResourceHandler) {
	if s.catalog != nil && !catalogHasActiveResource(s.catalog, def.URI) {
		return
	}
	s.resources[def.URI] = handler
	s.resDefs = append(s.resDefs, def)
}

// RegisterPrompt adds a prompt to the server.
func (s *Server) RegisterPrompt(def PromptDefinition, handler PromptHandler) {
	if s.catalog != nil && !catalogHasActivePrompt(s.catalog, def.Name) {
		return
	}
	s.prompts[def.Name] = handler
	s.promptDefs = append(s.promptDefs, def)
}

// HandleMessage processes a single JSON-RPC request and returns a response.
func (s *Server) HandleMessage(ctx context.Context, session *auth.Session, raw []byte) ([]byte, error) {
	var req Request
	if err := json.Unmarshal(raw, &req); err != nil {
		return s.errorResponse(nil, ErrCodeParse, "parse error")
	}

	if req.JSONRPC != "2.0" {
		return s.errorResponse(req.ID, ErrCodeInvalidRequest, "invalid JSON-RPC version")
	}

	// Handle notifications (no id) — they don't expect a response
	if req.ID == nil || string(req.ID) == "null" {
		s.handleNotification(ctx, req)
		return nil, nil
	}

	result, rpcErr := s.dispatch(ctx, session, req)
	if rpcErr != nil {
		return s.marshalResponse(Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   rpcErr,
		})
	}

	return s.marshalResponse(Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	})
}

func (s *Server) dispatch(ctx context.Context, session *auth.Session, req Request) (interface{}, *ResponseError) {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req.Params)
	case "ping":
		return map[string]interface{}{}, nil
	case "tools/list":
		return s.handleToolsList()
	case "tools/call":
		return s.handleToolCall(ctx, session, req.Params)
	case "resources/list":
		return s.handleResourcesList()
	case "resources/read":
		return s.handleResourceRead(ctx, session, req.Params)
	case "prompts/list":
		return s.handlePromptsList()
	case "prompts/get":
		return s.handlePromptGet(ctx, session, req.Params)
	default:
		return nil, &ResponseError{Code: ErrCodeMethodNotFound, Message: fmt.Sprintf("method not found: %s", req.Method)}
	}
}

func (s *Server) handleNotification(_ context.Context, req Request) {
	switch req.Method {
	case "notifications/initialized":
		log.Info().Msg("MCP client initialized")
	case "notifications/cancelled":
		log.Debug().Msg("MCP client cancelled request")
	default:
		log.Debug().Str("method", req.Method).Msg("unknown notification")
	}
}

func (s *Server) handleInitialize(params json.RawMessage) (interface{}, *ResponseError) {
	var p InitializeParams
	if params != nil {
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, &ResponseError{Code: ErrCodeInvalidParams, Message: "invalid initialize params"}
		}
	}

	log.Info().
		Str("client", p.ClientInfo.Name).
		Str("client_version", p.ClientInfo.Version).
		Str("protocol", p.ProtocolVersion).
		Msg("MCP initialize")

	return &InitializeResult{
		ProtocolVersion: ProtocolVersion,
		Capabilities: ServerCapability{
			Tools:     &ToolsCapability{},
			Resources: &ResourcesCapability{Subscribe: false},
			Prompts:   &PromptsCapability{},
			Logging:   &LoggingCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    ServerName,
			Version: version.Version.String(),
		},
		Instructions: "SoloDev MCP Server — AI-powered code remediation, security scanning, quality gates, " +
			"error tracking, pipeline generation, health monitoring, feature flags, and tech debt management. " +
			"Use tools to trigger actions, resources for live data, and prompts for expert reasoning chains.",
	}, nil
}

func (s *Server) handleToolsList() (interface{}, *ResponseError) {
	return &ToolsListResult{Tools: s.toolDefs}, nil
}

func (s *Server) handleToolCall(ctx context.Context, session *auth.Session, params json.RawMessage) (interface{}, *ResponseError) {
	var p ToolCallParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &ResponseError{Code: ErrCodeInvalidParams, Message: "invalid tool call params"}
	}

	handler, ok := s.tools[p.Name]
	if !ok {
		return nil, &ResponseError{Code: ErrCodeMethodNotFound, Message: fmt.Sprintf("unknown tool: %s", p.Name)}
	}

	result, err := handler(ctx, session, p.Arguments)
	if err != nil {
		log.Error().Err(err).Str("tool", p.Name).Msg("tool call failed")
		return &ToolCallResult{
			Content: []ContentBlock{TextContent(fmt.Sprintf("Error: %s", err.Error()))},
			IsError: true,
		}, nil
	}

	return result, nil
}

func (s *Server) handleResourcesList() (interface{}, *ResponseError) {
	return &ResourcesListResult{Resources: s.resDefs}, nil
}

func (s *Server) handleResourceRead(ctx context.Context, session *auth.Session, params json.RawMessage) (interface{}, *ResponseError) {
	var p ResourceReadParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &ResponseError{Code: ErrCodeInvalidParams, Message: "invalid resource read params"}
	}

	handler, ok := s.resources[p.URI]
	if !ok {
		return nil, &ResponseError{Code: ErrCodeInvalidParams, Message: fmt.Sprintf("unknown resource: %s", p.URI)}
	}

	result, err := handler(ctx, session)
	if err != nil {
		log.Error().Err(err).Str("uri", p.URI).Msg("resource read failed")
		return nil, &ResponseError{Code: ErrCodeInternal, Message: err.Error()}
	}

	return result, nil
}

func (s *Server) handlePromptsList() (interface{}, *ResponseError) {
	return &PromptsListResult{Prompts: s.promptDefs}, nil
}

func (s *Server) handlePromptGet(ctx context.Context, session *auth.Session, params json.RawMessage) (interface{}, *ResponseError) {
	var p PromptGetParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, &ResponseError{Code: ErrCodeInvalidParams, Message: "invalid prompt get params"}
	}

	handler, ok := s.prompts[p.Name]
	if !ok {
		return nil, &ResponseError{Code: ErrCodeMethodNotFound, Message: fmt.Sprintf("unknown prompt: %s", p.Name)}
	}

	result, err := handler(ctx, session, p.Arguments)
	if err != nil {
		log.Error().Err(err).Str("prompt", p.Name).Msg("prompt get failed")
		return nil, &ResponseError{Code: ErrCodeInternal, Message: err.Error()}
	}

	return result, nil
}

func (s *Server) errorResponse(id json.RawMessage, code int, message string) ([]byte, error) {
	return s.marshalResponse(Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &ResponseError{Code: code, Message: message},
	})
}

func (s *Server) marshalResponse(resp Response) ([]byte, error) {
	return json.Marshal(resp)
}

// SuccessResult creates a ToolCallResult with text content.
func SuccessResult(data interface{}) (*ToolCallResult, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	return &ToolCallResult{
		Content: []ContentBlock{TextContent(string(b))},
	}, nil
}

// ErrorResult creates an error ToolCallResult.
func ErrorResult(msg string) *ToolCallResult {
	return &ToolCallResult{
		Content: []ContentBlock{TextContent(msg)},
		IsError: true,
	}
}
