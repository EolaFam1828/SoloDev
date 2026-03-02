// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Registered MCP CLI subcommand.
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

package main

import (
	"github.com/EolaFam1828/SoloDev/app/api/openapi"
	"github.com/EolaFam1828/SoloDev/cli"
	"github.com/EolaFam1828/SoloDev/cli/operations/account"
	"github.com/EolaFam1828/SoloDev/cli/operations/hooks"
	"github.com/EolaFam1828/SoloDev/cli/operations/mcpcmd"
	"github.com/EolaFam1828/SoloDev/cli/operations/migrate"
	"github.com/EolaFam1828/SoloDev/cli/operations/server"
	"github.com/EolaFam1828/SoloDev/cli/operations/swagger"
	"github.com/EolaFam1828/SoloDev/cli/operations/user"
	"github.com/EolaFam1828/SoloDev/cli/operations/users"
	"github.com/EolaFam1828/SoloDev/version"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	application = "solodev"
	description = "SoloDev AI-native DevOps platform"
)

func main() {
	args := cli.GetArguments()

	app := kingpin.New(application, description)

	migrate.Register(app)
	server.Register(app, initSystem)
	mcpcmd.Register(app, initSystem)

	user.Register(app)
	users.Register(app)

	account.RegisterLogin(app)
	account.RegisterRegister(app)
	account.RegisterLogout(app)

	hooks.Register(app)

	swagger.Register(app, openapi.NewOpenAPIService())

	kingpin.Version(version.Version.String())
	kingpin.MustParse(app.Parse(args))
}
