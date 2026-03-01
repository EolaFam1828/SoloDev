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

package featureflag

import (
	"github.com/harness/gitness/types"
)

type UpdateInput struct {
	Name                  *string          `json:"name,omitempty"`
	Description           *string          `json:"description,omitempty"`
	Kind                  *string          `json:"kind,omitempty"`
	DefaultOnVariation    *string          `json:"default_on_variation,omitempty"`
	DefaultOffVariation   *string          `json:"default_off_variation,omitempty"`
	Enabled               *bool            `json:"enabled,omitempty"`
	Variations            *[]types.Variation `json:"variations,omitempty"`
	Tags                  *[]string        `json:"tags,omitempty"`
	Permanent             *bool            `json:"permanent,omitempty"`
}

type ToggleInput struct {
	Enabled bool `json:"enabled"`
}
