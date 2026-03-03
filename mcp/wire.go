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
	"github.com/harness/gitness/app/api/controller/airemediation"
	"github.com/harness/gitness/app/api/controller/autopipeline"
	"github.com/harness/gitness/app/api/controller/errortracker"
	"github.com/harness/gitness/app/api/controller/featureflag"
	"github.com/harness/gitness/app/api/controller/healthcheck"
	"github.com/harness/gitness/app/api/controller/qualitygate"
	"github.com/harness/gitness/app/api/controller/securityscan"
	"github.com/harness/gitness/app/api/controller/techdebt"
	"github.com/harness/gitness/app/auth/authn"
	"github.com/harness/gitness/app/services/errorbridge"
)

// Controllers holds all SoloDev controllers needed by the MCP server.
type Controllers struct {
	AutoPipeline *autopipeline.Controller
	SecurityScan *securityscan.Controller
	QualityGate  *qualitygate.Controller
	ErrorTracker *errortracker.Controller
	Remediation  *airemediation.Controller
	HealthCheck  *healthcheck.Controller
	FeatureFlag  *featureflag.Controller
	TechDebt     *techdebt.Controller
	ErrorBridge  *errorbridge.Bridge
}

// NewServer creates a fully wired MCP server with all controllers and tools registered.
func NewServer(
	authenticator authn.Authenticator,
	controllers *Controllers,
) *Server {
	auth := NewMCPAuthenticator(authenticator)
	srv := &Server{
		auth:        auth,
		controllers: controllers,
		tools:       make(map[string]ToolHandler),
		resources:   make(map[string]ResourceHandler),
		prompts:     make(map[string]PromptHandler),
		toolDefs:    nil,
		resDefs:     nil,
		promptDefs:  nil,
		catalog:     BuildCatalog(controllers),
	}

	// Register all tiers
	registerAtomicTools(srv)
	registerCompoundTools(srv)
	registerResources(srv)
	registerPrompts(srv)

	return srv
}
