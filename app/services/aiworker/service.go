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

package aiworker

import (
	"context"
	"fmt"

	"github.com/harness/gitness/app/services/contextengine"
	"github.com/harness/gitness/app/services/remediationdelivery"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/git"
	"github.com/harness/gitness/job"
)

// Service manages the AI remediation background jobs.
type Service struct {
	config        Config
	scheduler     *job.Scheduler
	executor      *job.Executor
	remStore      store.RemediationStore
	repoStore     store.RepoStore
	git           git.Interface
	provider      LLMProvider
	delivery      *remediationdelivery.Service
	contextEngine *contextengine.Service
}

// NewService creates a new AI remediation worker service.
func NewService(
	config Config,
	scheduler *job.Scheduler,
	executor *job.Executor,
	remStore store.RemediationStore,
	repoStore store.RepoStore,
	gitClient git.Interface,
	delivery *remediationdelivery.Service,
	ctxEngine *contextengine.Service,
) (*Service, error) {
	if !config.Enabled {
		return nil, nil
	}

	provider, err := ProvideProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM provider: %w", err)
	}
	if provider == nil {
		return nil, nil
	}

	return &Service{
		config:        config,
		scheduler:     scheduler,
		executor:      executor,
		remStore:      remStore,
		repoStore:     repoStore,
		git:           gitClient,
		provider:      provider,
		delivery:      delivery,
		contextEngine: ctxEngine,
	}, nil
}

func (s *Service) Available() bool {
	return s != nil && s.config.Enabled && s.provider != nil
}

// Register registers AI remediation job handlers and schedules the recurring poller.
func (s *Service) Register(ctx context.Context) error {
	if err := s.registerJobHandlers(); err != nil {
		return fmt.Errorf("failed to register AI remediation job handlers: %w", err)
	}

	if err := s.scheduleRecurringJobs(ctx); err != nil {
		return fmt.Errorf("failed to schedule AI remediation jobs: %w", err)
	}

	return nil
}

func (s *Service) registerJobHandlers() error {
	if err := s.executor.Register(jobTypeRemWorker, &remWorkerHandler{
		remStore:        s.remStore,
		repoStore:       s.repoStore,
		git:             s.git,
		provider:        s.provider,
		config:          s.config,
		deliveryService: s.delivery,
		contextEngine:   s.contextEngine,
	}); err != nil {
		return fmt.Errorf("failed to register remediation worker handler: %w", err)
	}

	if err := s.executor.Register(jobTypeRemPoller, &remPollerHandler{
		remStore:  s.remStore,
		scheduler: s.scheduler,
	}); err != nil {
		return fmt.Errorf("failed to register remediation poller handler: %w", err)
	}

	return nil
}

func (s *Service) scheduleRecurringJobs(ctx context.Context) error {
	return s.scheduler.AddRecurring(ctx, jobTypeRemPoller, jobTypeRemPoller, jobCronRemPoller, jobMaxDurPoller)
}
