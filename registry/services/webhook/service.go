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

package webhook

import (
	"context"
	"fmt"
	"time"

	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	gitnesswebhook "github.com/EolaFam1828/SoloDev/app/services/webhook"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/encrypt"
	"github.com/EolaFam1828/SoloDev/events"
	events2 "github.com/EolaFam1828/SoloDev/registry/app/events/artifact"
	registrystore "github.com/EolaFam1828/SoloDev/registry/app/store"
	"github.com/EolaFam1828/SoloDev/secret"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/stream"
)

const (
	eventsReaderGroupName = "gitness:webhook"
)

// Verify Service implements ServiceInterface.
var _ ServiceInterface = (*Service)(nil)

// Service is responsible for processing webhook events.
type Service struct {
	WebhookExecutor    *gitnesswebhook.WebhookExecutor
	tx                 dbtx.Transactor
	urlProvider        url.Provider
	spaceStore         store.SpaceStore
	principalStore     store.PrincipalStore
	config             gitnesswebhook.Config
	spacePathStore     store.SpacePathStore
	registryRepository registrystore.RegistryRepository
	spaceFinder        refcache.SpaceFinder
}

func NewService(
	ctx context.Context,
	config gitnesswebhook.Config,
	tx dbtx.Transactor,
	artifactsReaderFactory *events.ReaderFactory[*events2.Reader],
	webhookStore registrystore.WebhooksRepository,
	webhookExecutionStore registrystore.WebhooksExecutionRepository,
	spaceStore store.SpaceStore,
	urlProvider url.Provider,
	principalStore store.PrincipalStore,
	webhookURLProvider gitnesswebhook.URLProvider,
	spacePathStore store.SpacePathStore,
	secretService secret.Service,
	registryRepository registrystore.RegistryRepository,
	encrypter encrypt.Encrypter,
	spaceFinder refcache.SpaceFinder,
) (*Service, error) {
	if err := config.Prepare(); err != nil {
		return nil, fmt.Errorf("provided webhook service config is invalid: %w", err)
	}
	webhookExecutorStore := &RegistryWebhookExecutorStore{
		webhookStore:          webhookStore,
		webhookExecutionStore: webhookExecutionStore,
	}
	executor := gitnesswebhook.NewWebhookExecutor(config, webhookURLProvider, encrypter, spacePathStore,
		secretService, principalStore, webhookExecutorStore, gitnesswebhook.ArtifactRegistryTrigger)

	service := &Service{
		WebhookExecutor:    executor,
		tx:                 tx,
		spaceStore:         spaceStore,
		urlProvider:        urlProvider,
		principalStore:     principalStore,
		config:             config,
		spacePathStore:     spacePathStore,
		registryRepository: registryRepository,
		spaceFinder:        spaceFinder,
	}

	_, err := artifactsReaderFactory.Launch(ctx, eventsReaderGroupName, config.EventReaderName,
		func(r *events2.Reader) error {
			const idleTimeout = 1 * time.Minute
			r.Configure(
				stream.WithConcurrency(config.Concurrency),
				stream.WithHandlerOptions(
					stream.WithIdleTimeout(idleTimeout),
					stream.WithMaxRetries(config.MaxRetries),
				))

			// register events
			_ = r.RegisterArtifactCreated(service.handleEventArtifactCreated)
			_ = r.RegisterArtifactDeleted(service.handleEventArtifactDeleted)

			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to launch registry event reader for webhooks: %w", err)
	}

	return service, nil
}

func (s *Service) ReTriggerWebhookExecution(
	ctx context.Context,
	webhookExecutionID int64,
) (*gitnesswebhook.TriggerResult, error) {
	return s.WebhookExecutor.RetriggerWebhookExecution(ctx, webhookExecutionID)
}
