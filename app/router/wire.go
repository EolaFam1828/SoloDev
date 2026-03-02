// Copyright 2023 Harness, Inc.
// Modified by EolaFam1828 (2026) — Added soloDevModules parameter to ProvideRouter for SoloDev modules.
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

package router

import (
	"context"
	"strings"

	"github.com/EolaFam1828/SoloDev/app/api/controller/check"
	"github.com/EolaFam1828/SoloDev/app/api/controller/connector"
	"github.com/EolaFam1828/SoloDev/app/api/controller/execution"
	"github.com/EolaFam1828/SoloDev/app/api/controller/githook"
	"github.com/EolaFam1828/SoloDev/app/api/controller/gitspace"
	"github.com/EolaFam1828/SoloDev/app/api/controller/infraprovider"
	"github.com/EolaFam1828/SoloDev/app/api/controller/keywordsearch"
	"github.com/EolaFam1828/SoloDev/app/api/controller/lfs"
	"github.com/EolaFam1828/SoloDev/app/api/controller/logs"
	"github.com/EolaFam1828/SoloDev/app/api/controller/migrate"
	"github.com/EolaFam1828/SoloDev/app/api/controller/pipeline"
	"github.com/EolaFam1828/SoloDev/app/api/controller/plugin"
	"github.com/EolaFam1828/SoloDev/app/api/controller/principal"
	"github.com/EolaFam1828/SoloDev/app/api/controller/pullreq"
	"github.com/EolaFam1828/SoloDev/app/api/controller/repo"
	"github.com/EolaFam1828/SoloDev/app/api/controller/reposettings"
	"github.com/EolaFam1828/SoloDev/app/api/controller/secret"
	"github.com/EolaFam1828/SoloDev/app/api/controller/serviceaccount"
	"github.com/EolaFam1828/SoloDev/app/api/controller/space"
	"github.com/EolaFam1828/SoloDev/app/api/controller/system"
	"github.com/EolaFam1828/SoloDev/app/api/controller/template"
	"github.com/EolaFam1828/SoloDev/app/api/controller/trigger"
	"github.com/EolaFam1828/SoloDev/app/api/controller/upload"
	"github.com/EolaFam1828/SoloDev/app/api/controller/user"
	"github.com/EolaFam1828/SoloDev/app/api/controller/usergroup"
	"github.com/EolaFam1828/SoloDev/app/api/controller/webhook"
	"github.com/EolaFam1828/SoloDev/app/api/openapi"
	"github.com/EolaFam1828/SoloDev/app/auth/authn"
	"github.com/EolaFam1828/SoloDev/app/services/usage"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/git"
	"github.com/EolaFam1828/SoloDev/registry/app/api"
	"github.com/EolaFam1828/SoloDev/registry/app/api/router"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideRouter,
	ProvideSoloDevModules,
	api.WireSet,
)

func GetGitRoutingHost(ctx context.Context, urlProvider url.Provider) string {
	// use url provider as it has the latest data.
	gitHostname := urlProvider.GetGITHostname(ctx)
	apiHostname := urlProvider.GetAPIHostname(ctx)

	// only use host name to identify git traffic if it differs from api hostname.
	// TODO: Can we make this even more flexible - aka use the full base urls to route traffic?
	gitRoutingHost := ""
	if !strings.EqualFold(gitHostname, apiHostname) {
		gitRoutingHost = gitHostname
	}
	return gitRoutingHost
}

// ProvideRouter provides ordered list of routers.
func ProvideRouter(
	appCtx context.Context,
	config *types.Config,
	authenticator authn.Authenticator,
	repoCtrl *repo.Controller,
	repoSettingsCtrl *reposettings.Controller,
	executionCtrl *execution.Controller,
	logCtrl *logs.Controller,
	spaceCtrl *space.Controller,
	pipelineCtrl *pipeline.Controller,
	secretCtrl *secret.Controller,
	triggerCtrl *trigger.Controller,
	connectorCtrl *connector.Controller,
	templateCtrl *template.Controller,
	pluginCtrl *plugin.Controller,
	pullreqCtrl *pullreq.Controller,
	webhookCtrl *webhook.Controller,
	githookCtrl *githook.Controller,
	git git.Interface,
	saCtrl *serviceaccount.Controller,
	userCtrl *user.Controller,
	principalCtrl principal.Controller,
	userGroupCtrl *usergroup.Controller,
	checkCtrl *check.Controller,
	sysCtrl *system.Controller,
	blobCtrl *upload.Controller,
	searchCtrl *keywordsearch.Controller,
	infraProviderCtrl *infraprovider.Controller,
	gitspaceCtrl *gitspace.Controller,
	migrateCtrl *migrate.Controller,
	urlProvider url.Provider,
	openapi openapi.Service,
	registryRouter router.AppRouter,
	usageSender usage.Sender,
	lfsCtrl *lfs.Controller,
	soloDevModules *SoloDevModules,
) *Router {
	routers := make([]Interface, 4)

	gitRoutingHost := GetGitRoutingHost(appCtx, urlProvider)
	gitHandler := NewGitHandler(
		config,
		urlProvider,
		authenticator,
		repoCtrl,
		usageSender,
		lfsCtrl,
	)
	routers[0] = NewGitRouter(gitHandler, gitRoutingHost)
	routers[1] = router.NewRegistryRouter(registryRouter)

	apiHandler := NewAPIHandler(
		appCtx, config,
		authenticator, repoCtrl, repoSettingsCtrl, executionCtrl, logCtrl, spaceCtrl, pipelineCtrl,
		secretCtrl, triggerCtrl, connectorCtrl, templateCtrl, pluginCtrl, pullreqCtrl, webhookCtrl,
		githookCtrl, git, saCtrl, userCtrl, principalCtrl, userGroupCtrl, checkCtrl, sysCtrl, blobCtrl, searchCtrl,
		infraProviderCtrl, migrateCtrl, gitspaceCtrl, usageSender, soloDevModules)
	routers[2] = NewAPIRouter(apiHandler)

	sec := NewSecure(config)
	webHandler := NewWebHandler(
		authenticator, openapi, sec,
		config.PublicResourceCreationEnabled,
		config.Development.UISourceOverride,
	)
	routers[3] = NewWebRouter(webHandler)

	return NewRouter(routers)
}
