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

//go:build wireinject
// +build wireinject

package main

import (
	"context"

	controllerairemediation "github.com/EolaFam1828/SoloDev/app/api/controller/airemediation"
	controllerautopipeline "github.com/EolaFam1828/SoloDev/app/api/controller/autopipeline"
	checkcontroller "github.com/EolaFam1828/SoloDev/app/api/controller/check"
	"github.com/EolaFam1828/SoloDev/app/api/controller/connector"
	controllererrortracker "github.com/EolaFam1828/SoloDev/app/api/controller/errortracker"
	"github.com/EolaFam1828/SoloDev/app/api/controller/execution"
	controllerfeatureflag "github.com/EolaFam1828/SoloDev/app/api/controller/featureflag"
	githookCtrl "github.com/EolaFam1828/SoloDev/app/api/controller/githook"
	gitspaceCtrl "github.com/EolaFam1828/SoloDev/app/api/controller/gitspace"
	controllerhealthcheck "github.com/EolaFam1828/SoloDev/app/api/controller/healthcheck"
	infraproviderCtrl "github.com/EolaFam1828/SoloDev/app/api/controller/infraprovider"
	controllerkeywordsearch "github.com/EolaFam1828/SoloDev/app/api/controller/keywordsearch"
	"github.com/EolaFam1828/SoloDev/app/api/controller/lfs"
	"github.com/EolaFam1828/SoloDev/app/api/controller/limiter"
	controllerlogs "github.com/EolaFam1828/SoloDev/app/api/controller/logs"
	"github.com/EolaFam1828/SoloDev/app/api/controller/migrate"
	"github.com/EolaFam1828/SoloDev/app/api/controller/pipeline"
	"github.com/EolaFam1828/SoloDev/app/api/controller/plugin"
	"github.com/EolaFam1828/SoloDev/app/api/controller/principal"
	"github.com/EolaFam1828/SoloDev/app/api/controller/pullreq"
	controllerqualitygate "github.com/EolaFam1828/SoloDev/app/api/controller/qualitygate"
	"github.com/EolaFam1828/SoloDev/app/api/controller/repo"
	"github.com/EolaFam1828/SoloDev/app/api/controller/reposettings"
	controllersecurityscan "github.com/EolaFam1828/SoloDev/app/api/controller/securityscan"
	"github.com/EolaFam1828/SoloDev/app/api/controller/secret"
	"github.com/EolaFam1828/SoloDev/app/api/controller/service"
	"github.com/EolaFam1828/SoloDev/app/api/controller/serviceaccount"
	"github.com/EolaFam1828/SoloDev/app/api/controller/space"
	"github.com/EolaFam1828/SoloDev/app/api/controller/system"
	controllertechdebt "github.com/EolaFam1828/SoloDev/app/api/controller/techdebt"
	"github.com/EolaFam1828/SoloDev/app/api/controller/template"
	controllertrigger "github.com/EolaFam1828/SoloDev/app/api/controller/trigger"
	"github.com/EolaFam1828/SoloDev/app/api/controller/upload"
	"github.com/EolaFam1828/SoloDev/app/api/controller/user"
	"github.com/EolaFam1828/SoloDev/app/api/controller/usergroup"
	controllerwebhook "github.com/EolaFam1828/SoloDev/app/api/controller/webhook"
	"github.com/EolaFam1828/SoloDev/app/api/openapi"
	"github.com/EolaFam1828/SoloDev/app/auth/authn"
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/bootstrap"
	connectorservice "github.com/EolaFam1828/SoloDev/app/connector"
	airemediationevents "github.com/EolaFam1828/SoloDev/app/events/airemediation"
	aitaskevent "github.com/EolaFam1828/SoloDev/app/events/aitask"
	checkevents "github.com/EolaFam1828/SoloDev/app/events/check"
	errortrackerevents "github.com/EolaFam1828/SoloDev/app/events/errortracker"
	gitevents "github.com/EolaFam1828/SoloDev/app/events/git"
	gitspaceevents "github.com/EolaFam1828/SoloDev/app/events/gitspace"
	gitspacedeleteevents "github.com/EolaFam1828/SoloDev/app/events/gitspacedelete"
	gitspaceinfraevents "github.com/EolaFam1828/SoloDev/app/events/gitspaceinfra"
	gitspaceoperationsevents "github.com/EolaFam1828/SoloDev/app/events/gitspaceoperations"
	pipelineevents "github.com/EolaFam1828/SoloDev/app/events/pipeline"
	pullreqevents "github.com/EolaFam1828/SoloDev/app/events/pullreq"
	qualitygateevent "github.com/EolaFam1828/SoloDev/app/events/qualitygate"
	repoevents "github.com/EolaFam1828/SoloDev/app/events/repo"
	ruleevents "github.com/EolaFam1828/SoloDev/app/events/rule"
	userevents "github.com/EolaFam1828/SoloDev/app/events/user"
	"github.com/EolaFam1828/SoloDev/app/gitspace/infrastructure"
	"github.com/EolaFam1828/SoloDev/app/gitspace/logutil"
	"github.com/EolaFam1828/SoloDev/app/gitspace/orchestrator"
	containerorchestrator "github.com/EolaFam1828/SoloDev/app/gitspace/orchestrator/container"
	"github.com/EolaFam1828/SoloDev/app/gitspace/orchestrator/ide"
	"github.com/EolaFam1828/SoloDev/app/gitspace/orchestrator/runarg"
	"github.com/EolaFam1828/SoloDev/app/gitspace/platformconnector"
	"github.com/EolaFam1828/SoloDev/app/gitspace/platformsecret"
	"github.com/EolaFam1828/SoloDev/app/gitspace/scm"
	gitspacesecret "github.com/EolaFam1828/SoloDev/app/gitspace/secret"
	"github.com/EolaFam1828/SoloDev/app/pipeline/canceler"
	"github.com/EolaFam1828/SoloDev/app/pipeline/commit"
	"github.com/EolaFam1828/SoloDev/app/pipeline/converter"
	"github.com/EolaFam1828/SoloDev/app/pipeline/file"
	"github.com/EolaFam1828/SoloDev/app/pipeline/manager"
	"github.com/EolaFam1828/SoloDev/app/pipeline/resolver"
	"github.com/EolaFam1828/SoloDev/app/pipeline/runner"
	"github.com/EolaFam1828/SoloDev/app/pipeline/scheduler"
	"github.com/EolaFam1828/SoloDev/app/pipeline/triggerer"
	"github.com/EolaFam1828/SoloDev/app/router"
	"github.com/EolaFam1828/SoloDev/app/server"
	"github.com/EolaFam1828/SoloDev/app/services"
	"github.com/EolaFam1828/SoloDev/app/services/aiworker"
	"github.com/EolaFam1828/SoloDev/app/services/autolink"
	"github.com/EolaFam1828/SoloDev/app/services/branch"
	"github.com/EolaFam1828/SoloDev/app/services/cleanup"
	"github.com/EolaFam1828/SoloDev/app/services/codecomments"
	"github.com/EolaFam1828/SoloDev/app/services/codeowners"
	"github.com/EolaFam1828/SoloDev/app/services/dotrange"
	"github.com/EolaFam1828/SoloDev/app/services/errorbridge"
	"github.com/EolaFam1828/SoloDev/app/services/exporter"
	gitspacedeleteeventservice "github.com/EolaFam1828/SoloDev/app/services/gitspacedeleteevent"
	"github.com/EolaFam1828/SoloDev/app/services/gitspaceevent"
	"github.com/EolaFam1828/SoloDev/app/services/gitspaceservice"
	"github.com/EolaFam1828/SoloDev/app/services/gitspacesettings"
	"github.com/EolaFam1828/SoloDev/app/services/importer"
	"github.com/EolaFam1828/SoloDev/app/services/instrument"
	"github.com/EolaFam1828/SoloDev/app/services/keyfetcher"
	"github.com/EolaFam1828/SoloDev/app/services/keywordsearch"
	svclabel "github.com/EolaFam1828/SoloDev/app/services/label"
	"github.com/EolaFam1828/SoloDev/app/services/languageanalyzer"
	"github.com/EolaFam1828/SoloDev/app/services/locker"
	"github.com/EolaFam1828/SoloDev/app/services/merge"
	"github.com/EolaFam1828/SoloDev/app/services/metric"
	migrateservice "github.com/EolaFam1828/SoloDev/app/services/migrate"
	"github.com/EolaFam1828/SoloDev/app/services/notification"
	"github.com/EolaFam1828/SoloDev/app/services/notification/mailer"
	"github.com/EolaFam1828/SoloDev/app/services/protection"
	"github.com/EolaFam1828/SoloDev/app/services/publicaccess"
	"github.com/EolaFam1828/SoloDev/app/services/publickey"
	pullreqservice "github.com/EolaFam1828/SoloDev/app/services/pullreq"
	"github.com/EolaFam1828/SoloDev/app/services/qualityeval"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/services/remoteauth"
	reposervice "github.com/EolaFam1828/SoloDev/app/services/repo"
	"github.com/EolaFam1828/SoloDev/app/services/rules"
	"github.com/EolaFam1828/SoloDev/app/services/scanner"
	secretservice "github.com/EolaFam1828/SoloDev/app/services/secret"
	"github.com/EolaFam1828/SoloDev/app/services/settings"
	spaceSvc "github.com/EolaFam1828/SoloDev/app/services/space"
	"github.com/EolaFam1828/SoloDev/app/services/spacefinder"
	"github.com/EolaFam1828/SoloDev/app/services/tokengenerator"
	"github.com/EolaFam1828/SoloDev/app/services/trigger"
	"github.com/EolaFam1828/SoloDev/app/services/usage"
	usergroupservice "github.com/EolaFam1828/SoloDev/app/services/usergroup"
	"github.com/EolaFam1828/SoloDev/app/services/webhook"
	"github.com/EolaFam1828/SoloDev/app/sse"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/app/store/cache"
	"github.com/EolaFam1828/SoloDev/app/store/database"
	"github.com/EolaFam1828/SoloDev/app/store/logs"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/audit"
	"github.com/EolaFam1828/SoloDev/blob"
	cliserver "github.com/EolaFam1828/SoloDev/cli/operations/server"
	"github.com/EolaFam1828/SoloDev/encrypt"
	"github.com/EolaFam1828/SoloDev/events"
	"github.com/EolaFam1828/SoloDev/git"
	"github.com/EolaFam1828/SoloDev/git/api"
	"github.com/EolaFam1828/SoloDev/git/storage"
	infraproviderpkg "github.com/EolaFam1828/SoloDev/infraprovider"
	"github.com/EolaFam1828/SoloDev/job"
	"github.com/EolaFam1828/SoloDev/livelog"
	"github.com/EolaFam1828/SoloDev/lock"
	"github.com/EolaFam1828/SoloDev/pubsub"
	registryevents "github.com/EolaFam1828/SoloDev/registry/app/events/artifact"
	registrypostporcessingevents "github.com/EolaFam1828/SoloDev/registry/app/events/asyncprocessing"
	replicationevents "github.com/EolaFam1828/SoloDev/registry/app/events/replication"
	registryhelpers "github.com/EolaFam1828/SoloDev/registry/app/helpers"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/docker"
	cargoutils "github.com/EolaFam1828/SoloDev/registry/app/utils/cargo"
	gopackageutils "github.com/EolaFam1828/SoloDev/registry/app/utils/gopackage"
	registryhandlers "github.com/EolaFam1828/SoloDev/registry/job"
	registryindex "github.com/EolaFam1828/SoloDev/registry/services/asyncprocessing"
	registrywebhooks "github.com/EolaFam1828/SoloDev/registry/services/webhook"
	"github.com/EolaFam1828/SoloDev/ssh"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/check"

	"github.com/google/wire"
)

func initSystem(ctx context.Context, config *types.Config) (*cliserver.System, error) {
	wire.Build(
		cliserver.ProvideSystem,
		cliserver.ProvideRedis,
		bootstrap.WireSet,
		cliserver.ProvideDatabaseConfig,
		database.WireSet,
		cliserver.ProvideBlobStoreConfig,
		mailer.WireSet,
		notification.WireSet,
		blob.WireSet,
		dbtx.WireSet,
		cache.WireSetSpace,
		cache.WireSetRepo,
		refcache.WireSet,
		spacefinder.WireSet,
		router.WireSet,
		pullreqservice.WireSet,
		services.WireSet,
		services.ProvideGitspaceServices,
		cliserver.ProvideScannerConfig,
		scanner.WireSet,
		cliserver.ProvideAIWorkerConfig,
		aiworker.WireSet,
		qualityeval.WireSet,
		errorbridge.WireSet,
		server.WireSet,
		cliserver.ProvideNoOpMetricServer,
		url.WireSet,
		spaceSvc.ProvideNoopResourceMover,
		spaceSvc.WireSet,
		space.WireSet,
		limiter.WireSet,
		publicaccess.WireSet,
		repo.WireSet,
		reposettings.WireSet,
		pullreq.WireSet,
		merge.WireSet,
		controllerwebhook.WireSet,
		controllerwebhook.ProvidePreprocessor,
		svclabel.WireSet,
		serviceaccount.WireSet,
		user.WireSet,
		upload.WireSet,
		service.WireSet,
		principal.WireSet,
		usergroupservice.WireSet,
		system.WireSet,
		authn.WireSet,
		authz.WireSet,
		infrastructure.WireSet,
		infraproviderpkg.WireSet,
		gitspaceevents.WireSet,
		pipelineevents.WireSet,
		infraproviderCtrl.WireSet,
		gitspaceCtrl.WireSet,
		gitevents.WireSet,
		pullreqevents.WireSet,
		repoevents.WireSet,
		ruleevents.WireSet,
		userevents.WireSet,
		storage.WireSet,
		api.WireSet,
		cliserver.ProvideGitConfig,
		git.WireSet,
		store.WireSet,
		check.WireSet,
		encrypt.WireSet,
		cliserver.ProvideEventsConfig,
		events.WireSet,
		cliserver.ProvideWebhookConfig,
		cliserver.ProvideNotificationConfig,
		webhook.WireSet,
		languageanalyzer.WireSet,
		cliserver.ProvideTriggerConfig,
		trigger.WireSet,
		tokengenerator.WireSet,
		githookCtrl.ExtenderWireSet,
		githookCtrl.WireSet,
		cliserver.ProvideLockConfig,
		lock.WireSet,
		locker.WireSet,
		cliserver.ProvidePubsubConfig,
		pubsub.WireSet,
		cliserver.ProvideJobsConfig,
		job.WireSet,
		cliserver.ProvideCleanupConfig,
		cleanup.WireSet,
		codecomments.WireSet,
		protection.WireSet,
		checkcontroller.WireSet,
		controllerfeatureflag.WireSet,
		controllertechdebt.WireSet,
		controllersecurityscan.WireSet,
		controllerhealthcheck.WireSet,
		controllererrortracker.WireSet,
		controllerqualitygate.WireSet,
		controllerairemediation.WireSet,
		controllerautopipeline.WireSet,
		execution.WireSet,
		pipeline.WireSet,
		logs.WireSet,
		livelog.WireSet,
		controllerlogs.WireSet,
		secret.WireSet,
		connector.WireSet,
		connectorservice.WireSet,
		template.WireSet,
		manager.WireSet,
		triggerer.WireSet,
		file.WireSet,
		converter.WireSet,
		runner.WireSet,
		sse.WireSet,
		scheduler.WireSet,
		commit.WireSet,
		controllertrigger.WireSet,
		plugin.WireSet,
		resolver.WireSet,
		importer.WireSet,
		importer.ProvideConnectorService,
		migrateservice.WireSet,
		canceler.WireSet,
		exporter.WireSet,
		metric.WireSet,
		reposervice.WireSet,
		cliserver.ProvideCodeOwnerConfig,
		codeowners.WireSet,
		gitspaceevent.WireSet,
		cliserver.ProvideKeywordSearchConfig,
		keywordsearch.WireSet,
		rules.WireSet,
		rules.ProvideValidator,
		controllerkeywordsearch.WireSet,
		settings.WireSet,
		usergroup.WireSet,
		openapi.WireSet,
		repo.ProvideRepoCheck,
		audit.WireSet,
		ssh.WireSet,
		publickey.WireSet,
		keyfetcher.ProvideService,
		remoteauth.WireSet,
		migrate.WireSet,
		scm.WireSet,
		platformconnector.WireSet,
		platformsecret.WireSet,
		gitspacesecret.WireSet,
		orchestrator.WireSet,
		containerorchestrator.WireSet,
		cliserver.ProvideIDEVSCodeWebConfig,
		cliserver.ProvideDockerConfig,
		cliserver.ProvideGitspaceEventConfig,
		cliserver.ProvideGitspaceDeleteEventConfig,
		logutil.WireSet,
		cliserver.ProvideGitspaceOrchestratorConfig,
		ide.WireSet,
		gitspaceinfraevents.WireSet,
		aitaskevent.WireSet,
		gitspaceservice.WireSet,
		gitspacesettings.WireSet,
		gitspaceoperationsevents.WireSet,
		cliserver.ProvideGitspaceInfraProvisionerConfig,
		cliserver.ProvideIDEVSCodeConfig,
		cliserver.ProvideIDECursorConfig,
		cliserver.ProvideIDEWindsurfConfig,
		cliserver.ProvideIDEJetBrainsConfig,
		instrument.WireSet,
		airemediationevents.WireSet,
		docker.ProvideReporter,
		errortrackerevents.WireSet,
		qualitygateevent.WireSet,
		secretservice.WireSet,
		runarg.WireSet,
		lfs.WireSet,
		usage.WireSet,
		registryevents.WireSet,
		registrywebhooks.WireSet,
		gitspacedeleteevents.WireSet,
		gitspacedeleteeventservice.WireSet,
		registryindex.WireSet,
		cliserver.ProvideBranchConfig,
		branch.WireSet,
		autolink.WireSet,
		dotrange.WireSet,
		cargoutils.WireSet,
		gopackageutils.WireSet,
		registrypostporcessingevents.ProvideAsyncProcessingReporter,
		registrypostporcessingevents.ProvideReaderFactory,
		checkevents.WireSet,
		registryhelpers.WireSet,
		replicationevents.ProvideNoOpReplicationReporter,
		registryhandlers.WireSet,
	)
	return &cliserver.System{}, nil
}
