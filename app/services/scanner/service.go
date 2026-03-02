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

package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/harness/gitness/app/services/securityremediation"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"
	"github.com/harness/gitness/types"
)

// Service manages the security scanner background jobs.
type Service struct {
	config              Config
	scheduler           *job.Scheduler
	executor            *job.Executor
	scanStore           store.SecurityScanStore
	findingStore        store.ScanFindingStore
	scanners            []Scanner
	securityRemediation *securityremediation.Service
}

// NewService creates a new scanner service.
func NewService(
	config Config,
	scheduler *job.Scheduler,
	executor *job.Executor,
	scanStore store.SecurityScanStore,
	findingStore store.ScanFindingStore,
	securityRemediation *securityremediation.Service,
) *Service {
	if !config.Enabled {
		return nil
	}

	var scanners []Scanner
	scanners = append(scanners, NewSemgrepScanner(config.SemgrepPath, config.SemgrepRules))
	scanners = append(scanners, NewGitleaksScanner(config.GitleaksPath))
	scanners = append(scanners, NewTrivyScanner(config.TrivyPath))

	return &Service{
		config:              config,
		scheduler:           scheduler,
		executor:            executor,
		scanStore:           scanStore,
		findingStore:        findingStore,
		scanners:            scanners,
		securityRemediation: securityRemediation,
	}
}

func (s *Service) Capabilities() types.SecurityScannerStatus {
	status := types.SecurityScannerStatus{
		Enabled:      s != nil && s.config.Enabled,
		Capabilities: make([]types.SecurityScannerCapability, 0, len(s.scanners)),
	}
	if s == nil {
		return status
	}

	for _, scanner := range s.scanners {
		capability := types.SecurityScannerCapability{
			Name:      scanner.Name(),
			Available: scanner.Available(),
		}
		if capability.Available {
			status.Ready = true
		}
		status.Capabilities = append(status.Capabilities, capability)
	}

	return status
}

// Register registers scanner job handlers and schedules recurring poller.
func (s *Service) Register(ctx context.Context) error {
	if err := s.registerJobHandlers(); err != nil {
		return fmt.Errorf("failed to register scanner job handlers: %w", err)
	}

	if err := s.scheduleRecurringJobs(ctx); err != nil {
		return fmt.Errorf("failed to schedule scanner jobs: %w", err)
	}

	return nil
}

func (s *Service) registerJobHandlers() error {
	if err := s.executor.Register(jobTypeScanWorker, &scanWorkerHandler{
		scanStore:           s.scanStore,
		findingStore:        s.findingStore,
		scanners:            s.scanners,
		gitRoot:             s.config.GitRoot,
		securityRemediation: s.securityRemediation,
	}); err != nil {
		return fmt.Errorf("failed to register scan worker handler: %w", err)
	}

	if err := s.executor.Register(jobTypeScanPoller, &scanPollerHandler{
		scanStore: s.scanStore,
		scheduler: s.scheduler,
	}); err != nil {
		return fmt.Errorf("failed to register scan poller handler: %w", err)
	}

	return nil
}

func (s *Service) scheduleRecurringJobs(ctx context.Context) error {
	return s.scheduler.AddRecurring(ctx, jobTypeScanPoller, jobTypeScanPoller, jobCronScanPoller, jobMaxDurPoller)
}

// TriggerScanJob submits an immediate scan worker job for the given scan.
func (s *Service) TriggerScanJob(ctx context.Context, scan *types.ScanResult) error {
	if !s.Capabilities().Ready {
		return fmt.Errorf("no security scanners are available")
	}

	data, _ := json.Marshal(scanJobInput{ScanID: scan.ID})
	return s.scheduler.RunJob(ctx, job.Definition{
		UID:     fmt.Sprintf("scan-%s", scan.Identifier),
		Type:    jobTypeScanWorker,
		Timeout: 5 * time.Minute,
		Data:    string(data),
	})
}
