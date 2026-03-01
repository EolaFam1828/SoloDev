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

// Package errorbridge connects the Error Tracker to the AI Remediation system.
// When a new error is reported, this bridge automatically creates a remediation
// task with the error's stack trace and context, enabling instant AI-driven fixes.
package errorbridge

import (
	"context"
	"fmt"
	"log"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
)

// Bridge watches for new errors and creates remediation tasks automatically.
type Bridge struct {
	remediationStore store.RemediationStore
	enabled          bool
	autoTriggerFatal bool // only auto-trigger on fatal/error severity, not warnings
}

// NewBridge creates a new Bridge.
func NewBridge(
	remediationStore store.RemediationStore,
	enabled bool,
) *Bridge {
	return &Bridge{
		remediationStore: remediationStore,
		enabled:          enabled,
		autoTriggerFatal: true,
	}
}

// OnErrorReported is called when a new error is reported via the Error Tracker.
// It automatically creates a pending AI remediation task with full context.
func (b *Bridge) OnErrorReported(ctx context.Context, errorGroup *types.ErrorGroup, occurrence *types.ErrorOccurrence) {
	if !b.enabled {
		return
	}

	// Skip warnings unless explicitly configured
	if b.autoTriggerFatal && errorGroup.Severity == types.ErrorSeverityWarning {
		return
	}

	// Skip if already resolved or ignored
	if errorGroup.Status == types.ErrorGroupStatusResolved ||
		errorGroup.Status == types.ErrorGroupStatusIgnored {
		return
	}

	// Build the remediation task
	rem := &types.Remediation{
		SpaceID:       errorGroup.SpaceID,
		RepoID:        errorGroup.RepoID,
		Identifier:    fmt.Sprintf("err-%s-%d", errorGroup.Identifier, types.NowMillis()),
		Title:         fmt.Sprintf("[Auto] Fix: %s", errorGroup.Title),
		Description:   errorGroup.Message,
		Status:        types.RemediationStatusPending,
		TriggerSource: types.RemediationTriggerError,
		TriggerRef:    errorGroup.Identifier,
		ErrorLog:      occurrence.StackTrace,
		FilePath:      errorGroup.FilePath,
		SourceCode:    "", // Would be populated from git if repo context is available
		CreatedBy:     errorGroup.CreatedBy,
		Created:       types.NowMillis(),
		Updated:       types.NowMillis(),
		Version:       1,
	}

	if err := b.remediationStore.Create(ctx, rem); err != nil {
		log.Printf("[errorbridge] failed to auto-create remediation for error %s: %v",
			errorGroup.Identifier, err)
		return
	}

	log.Printf("[errorbridge] auto-created remediation %s for error %s (severity=%s)",
		rem.Identifier, errorGroup.Identifier, errorGroup.Severity)
}

// OnPipelineFailed is called when a pipeline execution fails.
// It creates a remediation task with the build logs as context.
func (b *Bridge) OnPipelineFailed(
	ctx context.Context,
	spaceID int64,
	repoID int64,
	executionNumber int64,
	branch string,
	commitSHA string,
	errorLog string,
	createdBy int64,
) {
	if !b.enabled {
		return
	}

	rem := &types.Remediation{
		SpaceID:       spaceID,
		RepoID:        repoID,
		Identifier:    fmt.Sprintf("pipe-%d-%d", executionNumber, types.NowMillis()),
		Title:         fmt.Sprintf("[Auto] Fix pipeline failure #%d", executionNumber),
		Status:        types.RemediationStatusPending,
		TriggerSource: types.RemediationTriggerPipeline,
		TriggerRef:    fmt.Sprintf("%d", executionNumber),
		Branch:        branch,
		CommitSHA:     commitSHA,
		ErrorLog:      errorLog,
		CreatedBy:     createdBy,
		Created:       types.NowMillis(),
		Updated:       types.NowMillis(),
		Version:       1,
	}

	if err := b.remediationStore.Create(ctx, rem); err != nil {
		log.Printf("[errorbridge] failed to auto-create remediation for pipeline #%d: %v",
			executionNumber, err)
		return
	}

	log.Printf("[errorbridge] auto-created remediation %s for pipeline failure #%d",
		rem.Identifier, executionNumber)
}
