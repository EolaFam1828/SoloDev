// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
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

package airemediation

import "github.com/EolaFam1828/SoloDev/types"

// Base contains basic AI remediation event data.
type Base struct {
	RemediationID int64  `json:"remediation_id"`
	Identifier    string `json:"identifier"`
	SpaceID       int64  `json:"space_id"`
}

// RemediationTriggered is emitted when a new remediation is triggered.
type RemediationTriggered struct {
	Base
	Remediation *types.Remediation `json:"remediation"`
}

// RemediationCompleted is emitted when a remediation finishes (success or failure).
type RemediationCompleted struct {
	Base
	Status types.RemediationStatus `json:"status"`
}

// RemediationApplied is emitted when a user accepts the AI fix.
type RemediationApplied struct {
	Base
	FixBranch string `json:"fix_branch"`
	PRLink    string `json:"pr_link,omitempty"`
}
