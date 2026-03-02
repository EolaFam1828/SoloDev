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
	"encoding/json"
	"fmt"
	"time"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/job"
)

// remPollerHandler finds pending remediations and spawns worker jobs.
type remPollerHandler struct {
	remStore  store.RemediationStore
	scheduler *job.Scheduler
}

func (h *remPollerHandler) Handle(ctx context.Context, _ string, _ job.ProgressReporter) (string, error) {
	pending, err := h.remStore.ListPendingGlobal(ctx, 10)
	if err != nil {
		return "", fmt.Errorf("failed to list pending remediations: %w", err)
	}

	spawned := 0
	for _, rem := range pending {
		data, _ := json.Marshal(remJobInput{RemediationID: rem.ID})
		err := h.scheduler.RunJob(ctx, job.Definition{
			UID:     fmt.Sprintf("ai-rem-%s", rem.Identifier),
			Type:    jobTypeRemWorker,
			Timeout: 2 * time.Minute,
			Data:    string(data),
		})
		if err != nil {
			continue
		}
		spawned++
	}

	return fmt.Sprintf("spawned %d remediation jobs", spawned), nil
}
