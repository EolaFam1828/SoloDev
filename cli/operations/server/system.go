// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Added MCP accessor methods (Authenticator, SoloDevModules, ErrorBridge).
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

package server

import (
	"context"

	"github.com/harness/gitness/app/auth/authn"
	"github.com/harness/gitness/app/bootstrap"
	"github.com/harness/gitness/app/pipeline/resolver"
	"github.com/harness/gitness/app/router"
	"github.com/harness/gitness/app/server"
	"github.com/harness/gitness/app/services"
	"github.com/harness/gitness/app/services/errorbridge"
	"github.com/harness/gitness/http"
	"github.com/harness/gitness/ssh"

	"github.com/drone/runner-go/poller"
)

// System stores high level System sub-routines.
type System struct {
	bootstrap       bootstrap.Bootstrap
	server          *server.Server
	sshServer       *ssh.Server
	resolverManager *resolver.Manager
	poller          *poller.Poller
	services        services.Services
	metricServer    http.ListenAndServeServer

	// MCP-related fields (optional, set via SetMCPDeps)
	authenticator  authn.Authenticator
	soloDevModules *router.SoloDevModules
	errorBridge    *errorbridge.Bridge
}

func ProvideNoOpMetricServer() http.ListenAndServeServer {
	return http.NoOpListenAndServeServer{}
}

func NewSystem(
	bootstrap bootstrap.Bootstrap,
	server *server.Server,
	sshServer *ssh.Server,
	poller *poller.Poller,
	resolverManager *resolver.Manager,
	services services.Services,
	metricServer http.ListenAndServeServer,
) *System {
	return &System{
		bootstrap:       bootstrap,
		server:          server,
		sshServer:       sshServer,
		poller:          poller,
		resolverManager: resolverManager,
		services:        services,
		metricServer:    metricServer,
	}
}

// SetMCPDeps stores MCP-related dependencies for the MCP CLI subcommand.
func (s *System) SetMCPDeps(
	authenticator authn.Authenticator,
	soloDevModules *router.SoloDevModules,
	errorBridge *errorbridge.Bridge,
) {
	s.authenticator = authenticator
	s.soloDevModules = soloDevModules
	s.errorBridge = errorBridge
}

// Authenticator returns the authenticator.
func (s *System) Authenticator() authn.Authenticator {
	return s.authenticator
}

// SoloDevModules returns the SoloDev controller modules.
func (s *System) SoloDevModules() *router.SoloDevModules {
	return s.soloDevModules
}

// ErrorBridge returns the error bridge.
func (s *System) ErrorBridge() *errorbridge.Bridge {
	return s.errorBridge
}

// Bootstrap runs the system bootstrap.
func (s *System) Bootstrap(ctx context.Context) error {
	return s.bootstrap(ctx)
}
