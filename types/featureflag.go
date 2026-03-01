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

package types

import "database/sql/driver"

type Variation struct {
	Identifier string      `json:"identifier"`
	Name       string      `json:"name"`
	Value      interface{} `json:"value"`
}

type FeatureFlag struct {
	ID                    int64        `json:"id" db:"ff_id"`
	SpaceID               int64        `json:"space_id" db:"ff_space_id"`
	Identifier            string       `json:"identifier" db:"ff_identifier"`
	Name                  string       `json:"name" db:"ff_name"`
	Description           string       `json:"description" db:"ff_description"`
	Kind                  string       `json:"kind" db:"ff_kind"`
	DefaultOnVariation    string       `json:"default_on_variation" db:"ff_default_on_variation"`
	DefaultOffVariation   string       `json:"default_off_variation" db:"ff_default_off_variation"`
	Enabled               bool         `json:"enabled" db:"ff_enabled"`
	Variations            []Variation  `json:"variations" db:"ff_variations"`
	Tags                  []string     `json:"tags" db:"ff_tags"`
	Permanent             bool         `json:"permanent" db:"ff_permanent"`
	CreatedBy             int64        `json:"created_by" db:"ff_created_by"`
	Created               int64        `json:"created" db:"ff_created"`
	Updated               int64        `json:"updated" db:"ff_updated"`
	Version               int64        `json:"version" db:"ff_version"`
}

type FeatureFlagFilter struct {
	ListQueryFilter
	Query string `json:"query,omitempty"`
}

// Value implements the driver.Valuer interface for storing in the database.
func (f *FeatureFlag) Value() (driver.Value, error) {
	return f, nil
}
