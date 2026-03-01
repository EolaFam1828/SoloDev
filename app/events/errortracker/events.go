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

package errortracker

import "github.com/harness/gitness/types"

// Base contains basic error tracker event data.
type Base struct {
	ErrorGroupID int64  `json:"error_group_id"`
	Identifier   string `json:"identifier"`
	SpaceID      int64  `json:"space_id"`
}

// ErrorReported is emitted when an error is reported.
type ErrorReported struct {
	Base
	ErrorGroup *types.ErrorGroup `json:"error_group"`
}

// ErrorStatusChanged is emitted when an error group status changes.
type ErrorStatusChanged struct {
	Base
	OldStatus types.ErrorGroupStatus `json:"old_status"`
	NewStatus types.ErrorGroupStatus `json:"new_status"`
}

// ErrorAssigned is emitted when an error group is assigned to a user.
type ErrorAssigned struct {
	Base
	AssignedTo *int64 `json:"assigned_to"`
}
