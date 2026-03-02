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
	"path/filepath"
	"time"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

const (
	jobTypeScanWorker = "scanner-worker"
	jobTypeScanPoller = "scanner-poller"
	jobCronScanPoller = "*/15 * * * * *" // every 15 seconds
	jobMaxDurPoller   = 30 * time.Second
)

// scanWorkerHandler processes a single security scan.
type scanWorkerHandler struct {
	scanStore    store.SecurityScanStore
	findingStore store.ScanFindingStore
	scanners     []Scanner
	gitRoot      string
}

type scanJobInput struct {
	ScanID int64 `json:"scan_id"`
}

func (h *scanWorkerHandler) Handle(ctx context.Context, data string, _ job.ProgressReporter) (string, error) {
	var input scanJobInput
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("failed to parse scan job input: %w", err)
	}

	scan, err := h.scanStore.Find(ctx, input.ScanID)
	if err != nil {
		return "", fmt.Errorf("failed to find scan %d: %w", input.ScanID, err)
	}

	// Mark as running.
	scan.Status = enum.SecurityScanStatusRunning
	scan.Updated = time.Now().UnixMilli()
	if err := h.scanStore.Update(ctx, scan); err != nil {
		return "", fmt.Errorf("failed to update scan status to running: %w", err)
	}

	startTime := time.Now()
	repoDir := filepath.Join(h.gitRoot, fmt.Sprintf("%d", scan.RepoID))

	var allFindings []types.ScanFinding
	var scanErrors []string

	for _, s := range h.scanners {
		if !s.Available() {
			continue
		}

		findings, err := s.Scan(ctx, repoDir)
		if err != nil {
			scanErrors = append(scanErrors, fmt.Sprintf("%s: %v", s.Name(), err))
			continue
		}

		allFindings = append(allFindings, findings...)
	}

	// Insert findings.
	var critical, high, medium, low int
	for i := range allFindings {
		f := &allFindings[i]
		f.ScanID = scan.ID
		f.Created = time.Now().UnixMilli()
		f.Updated = f.Created

		if err := h.findingStore.Create(ctx, f); err != nil {
			scanErrors = append(scanErrors, fmt.Sprintf("insert finding: %v", err))
			continue
		}

		switch f.Severity {
		case enum.SecurityFindingSeverityCritical:
			critical++
		case enum.SecurityFindingSeverityHigh:
			high++
		case enum.SecurityFindingSeverityMedium:
			medium++
		case enum.SecurityFindingSeverityLow:
			low++
		}
	}

	// Update scan with results.
	scan.TotalIssues = len(allFindings)
	scan.CriticalCount = critical
	scan.HighCount = high
	scan.MediumCount = medium
	scan.LowCount = low
	scan.Duration = time.Since(startTime).Milliseconds()
	scan.Updated = time.Now().UnixMilli()

	if len(scanErrors) > 0 && len(allFindings) == 0 {
		scan.Status = enum.SecurityScanStatusFailed
	} else {
		scan.Status = enum.SecurityScanStatusCompleted
	}

	if err := h.scanStore.Update(ctx, scan); err != nil {
		return "", fmt.Errorf("failed to update scan with results: %w", err)
	}

	return fmt.Sprintf("completed: %d findings", len(allFindings)), nil
}

// scanPollerHandler finds pending scans and spawns worker jobs.
type scanPollerHandler struct {
	scanStore store.SecurityScanStore
	scheduler *job.Scheduler
}

func (h *scanPollerHandler) Handle(ctx context.Context, _ string, _ job.ProgressReporter) (string, error) {
	pending, err := h.scanStore.ListByStatus(ctx, enum.SecurityScanStatusPending, 10)
	if err != nil {
		return "", fmt.Errorf("failed to list pending scans: %w", err)
	}

	spawned := 0
	for _, scan := range pending {
		data, _ := json.Marshal(scanJobInput{ScanID: scan.ID})
		err := h.scheduler.RunJob(ctx, job.Definition{
			UID:     fmt.Sprintf("scan-%s", scan.Identifier),
			Type:    jobTypeScanWorker,
			Timeout: 5 * time.Minute,
			Data:    string(data),
		})
		if err != nil {
			continue
		}
		spawned++
	}

	return fmt.Sprintf("spawned %d scan jobs", spawned), nil
}
