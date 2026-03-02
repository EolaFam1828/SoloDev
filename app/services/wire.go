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

package services

import (
	"github.com/EolaFam1828/SoloDev/app/services/aitaskevent"
	"github.com/EolaFam1828/SoloDev/app/services/aiworker"
	"github.com/EolaFam1828/SoloDev/app/services/branch"
	"github.com/EolaFam1828/SoloDev/app/services/cleanup"
	"github.com/EolaFam1828/SoloDev/app/services/gitspace"
	"github.com/EolaFam1828/SoloDev/app/services/gitspacedeleteevent"
	"github.com/EolaFam1828/SoloDev/app/services/gitspaceevent"
	"github.com/EolaFam1828/SoloDev/app/services/gitspaceinfraevent"
	"github.com/EolaFam1828/SoloDev/app/services/gitspaceoperationsevent"
	"github.com/EolaFam1828/SoloDev/app/services/infraprovider"
	"github.com/EolaFam1828/SoloDev/app/services/instrument"
	"github.com/EolaFam1828/SoloDev/app/services/keywordsearch"
	"github.com/EolaFam1828/SoloDev/app/services/languageanalyzer"
	"github.com/EolaFam1828/SoloDev/app/services/metric"
	"github.com/EolaFam1828/SoloDev/app/services/notification"
	"github.com/EolaFam1828/SoloDev/app/services/pullreq"
	"github.com/EolaFam1828/SoloDev/app/services/repo"
	"github.com/EolaFam1828/SoloDev/app/services/scanner"
	"github.com/EolaFam1828/SoloDev/app/services/trigger"
	"github.com/EolaFam1828/SoloDev/app/services/webhook"
	"github.com/EolaFam1828/SoloDev/job"
	"github.com/EolaFam1828/SoloDev/registry/job/handler"
	registryasyncprocessing "github.com/EolaFam1828/SoloDev/registry/services/asyncprocessing"
	registrywebhooks "github.com/EolaFam1828/SoloDev/registry/services/webhook"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideServices,
)

type Services struct {
	Webhook                        *webhook.Service
	PullReq                        *pullreq.Service
	Trigger                        *trigger.Service
	JobScheduler                   *job.Scheduler
	MetricCollector                *metric.CollectorJob
	RepoSizeCalculator             *repo.SizeCalculator
	Repo                           *repo.Service
	Cleanup                        *cleanup.Service
	Notification                   *notification.Service
	Keywordsearch                  *keywordsearch.Service
	GitspaceService                *GitspaceServices
	Instrumentation                instrument.Service
	instrumentConsumer             instrument.Consumer
	instrumentRepoCounter          *instrument.RepositoryCount
	registryWebhooksService        *registrywebhooks.Service
	Branch                         *branch.Service
	registryAsyncProcessingService *registryasyncprocessing.Service
	languageAnalyzer               languageanalyzer.LanguageAnalyzer
	Scanner                        *scanner.Service
	AIWorker                       *aiworker.Service
}

type GitspaceServices struct {
	GitspaceEvent              *gitspaceevent.Service
	infraProvider              *infraprovider.Service
	gitspace                   *gitspace.Service
	gitspaceInfraEventSvc      *gitspaceinfraevent.Service
	gitspaceOperationsEventSvc *gitspaceoperationsevent.Service
	gitspaceDeleteEventSvc     *gitspacedeleteevent.Service
	aiTaskEventSvc             *aitaskevent.Service
}

func ProvideGitspaceServices(
	gitspaceEventSvc *gitspaceevent.Service,
	gitspaceDeleteEventSvc *gitspacedeleteevent.Service,
	infraProviderSvc *infraprovider.Service,
	gitspaceSvc *gitspace.Service,
	gitspaceInfraEventSvc *gitspaceinfraevent.Service,
	gitspaceOperationsEventSvc *gitspaceoperationsevent.Service,
	aiTaskEventSvc *aitaskevent.Service,
) *GitspaceServices {
	return &GitspaceServices{
		GitspaceEvent:              gitspaceEventSvc,
		infraProvider:              infraProviderSvc,
		gitspace:                   gitspaceSvc,
		gitspaceInfraEventSvc:      gitspaceInfraEventSvc,
		gitspaceOperationsEventSvc: gitspaceOperationsEventSvc,
		gitspaceDeleteEventSvc:     gitspaceDeleteEventSvc,
		aiTaskEventSvc:             aiTaskEventSvc,
	}
}

func ProvideServices(
	webhooksSvc *webhook.Service,
	pullReqSvc *pullreq.Service,
	triggerSvc *trigger.Service,
	jobScheduler *job.Scheduler,
	metricCollector *metric.CollectorJob,
	repoSizeCalculator *repo.SizeCalculator,
	repo *repo.Service,
	cleanupSvc *cleanup.Service,
	notificationSvc *notification.Service,
	keywordsearchSvc *keywordsearch.Service,
	gitspaceSvc *GitspaceServices,
	instrumentation instrument.Service,
	instrumentConsumer instrument.Consumer,
	instrumentRepoCounter *instrument.RepositoryCount,
	registryWebhooksService *registrywebhooks.Service,
	branchSvc *branch.Service,
	registryAsyncProcessingService *registryasyncprocessing.Service,
	registryJobRpmRegistryIndex *handler.JobRpmRegistryIndex,
	languageAnalyzer languageanalyzer.LanguageAnalyzer,
	scannerSvc *scanner.Service,
	aiWorkerSvc *aiworker.Service,
) Services {
	return Services{
		Webhook:                        webhooksSvc,
		PullReq:                        pullReqSvc,
		Trigger:                        triggerSvc,
		JobScheduler:                   jobScheduler,
		MetricCollector:                metricCollector,
		RepoSizeCalculator:             repoSizeCalculator,
		Repo:                           repo,
		Cleanup:                        cleanupSvc,
		Notification:                   notificationSvc,
		Keywordsearch:                  keywordsearchSvc,
		GitspaceService:                gitspaceSvc,
		Instrumentation:                instrumentation,
		instrumentConsumer:             instrumentConsumer,
		instrumentRepoCounter:          instrumentRepoCounter,
		registryWebhooksService:        registryWebhooksService,
		Branch:                         branchSvc,
		registryAsyncProcessingService: registryAsyncProcessingService,
		languageAnalyzer:               languageAnalyzer,
		Scanner:                        scannerSvc,
		AIWorker:                       aiWorkerSvc,
	}
}
