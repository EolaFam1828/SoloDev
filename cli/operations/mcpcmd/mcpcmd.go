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

// Package mcpcmd implements the `gitness mcp` CLI subcommands for starting
// the MCP server in stdio or SSE mode.
package mcpcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/harness/gitness/app/auth/authn"
	"github.com/harness/gitness/cli/operations/server"
	"github.com/harness/gitness/mcp"
	"github.com/harness/gitness/types"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// MCPInitializer is a function that initializes the MCP server dependencies.
// It reuses the main system initializer to get controller access.
type MCPInitializer func(context.Context, *types.Config) (*server.System, error)

type stdioCommand struct {
	envfile     string
	initializer MCPInitializer
}

type sseCommand struct {
	envfile     string
	port        int
	initializer MCPInitializer
}

// Register adds `gitness mcp` subcommands to the CLI application.
func Register(app *kingpin.Application, initializer MCPInitializer) {
	mcpCmd := app.Command("mcp", "Run the SoloDev MCP server")

	// gitness mcp --stdio
	stdioCmd := &stdioCommand{initializer: initializer}
	stdio := mcpCmd.Command("stdio", "Start MCP server in stdio mode (for Claude Desktop)").
		Default().
		Action(stdioCmd.run)
	stdio.Arg("envfile", "load the environment variable file").
		Default("").
		StringVar(&stdioCmd.envfile)

	// gitness mcp sse --port 3001
	sseCmd := &sseCommand{initializer: initializer}
	sse := mcpCmd.Command("sse", "Start MCP server in HTTP/SSE mode (for remote clients)").
		Action(sseCmd.run)
	sse.Arg("envfile", "load the environment variable file").
		Default("").
		StringVar(&sseCmd.envfile)
	sse.Flag("port", "Port to listen on for SSE transport").
		Default("3001").
		IntVar(&sseCmd.port)
}

func (c *stdioCommand) run(_ *kingpin.ParseContext) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = godotenv.Load(c.envfile)

	config, err := server.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	server.SetupLogger(config)
	ctx = log.Logger.WithContext(ctx)

	// Initialize the full system to get controller access
	system, err := c.initializer(ctx, config)
	if err != nil {
		return fmt.Errorf("init system: %w", err)
	}

	// Extract the authenticator and controllers from the system
	authenticator, controllers := extractMCPDeps(system)

	// Build the MCP server
	mcpServer := mcp.NewServer(authenticator, controllers)

	log.Info().Msg("Starting MCP server in stdio mode")
	return mcpServer.ServeStdio(ctx)
}

func (c *sseCommand) run(_ *kingpin.ParseContext) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = godotenv.Load(c.envfile)

	config, err := server.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	server.SetupLogger(config)
	ctx = log.Logger.WithContext(ctx)

	// Initialize the full system to get controller access
	system, err := c.initializer(ctx, config)
	if err != nil {
		return fmt.Errorf("init system: %w", err)
	}

	// Bootstrap the system for DB access
	if err := system.Bootstrap(ctx); err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}

	// Extract the authenticator and controllers from the system
	authenticator, controllers := extractMCPDeps(system)

	// Build the MCP server
	mcpServer := mcp.NewServer(authenticator, controllers)

	addr := fmt.Sprintf(":%d", c.port)
	log.Info().Str("addr", addr).Msg("Starting MCP server in SSE mode")
	return mcpServer.ServeSSE(ctx, addr)
}

// extractMCPDeps pulls the authenticator and SoloDev controllers from the initialized System.
// This is implemented in wire_mcp.go to access the system's internal fields.
func extractMCPDeps(system *server.System) (authn.Authenticator, *mcp.Controllers) {
	// Get the authenticator from the system
	authenticator := system.Authenticator()

	// Get the SoloDev modules from the system
	modules := system.SoloDevModules()

	controllers := &mcp.Controllers{}
	if modules != nil {
		controllers.AutoPipeline = modules.AutoPipelineCtrl
		controllers.SecurityScan = modules.SecurityScanCtrl
		controllers.QualityGate = modules.QualityGateCtrl
		controllers.ErrorTracker = modules.ErrorTrackerCtrl
		controllers.Remediation = modules.RemediationCtrl
		controllers.HealthCheck = modules.HealthCheckCtrl
		controllers.FeatureFlag = modules.FeatureFlagCtrl
		controllers.TechDebt = modules.TechDebtCtrl
	}

	// Get the error bridge if available
	controllers.ErrorBridge = system.ErrorBridge()

	return authenticator, controllers
}

// init registers the SOLODEV_SPACE env var default
func init() {
	if os.Getenv("SOLODEV_SPACE") == "" {
		_ = os.Setenv("SOLODEV_SPACE", "default")
	}
}
