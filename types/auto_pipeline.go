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

// DetectedStack represents the autodetected technology stack of a repository.
type DetectedStack struct {
	Languages  []LanguageInfo `json:"languages"`
	Frameworks []string       `json:"frameworks"`
	HasTests   bool           `json:"has_tests"`
	HasDocker  bool           `json:"has_docker"`
	HasCI      bool           `json:"has_ci"` // existing .harness/ pipeline
	BuildTool  string         `json:"build_tool,omitempty"`
	EntryPoint string         `json:"entry_point,omitempty"`
}

// LanguageInfo stores info about a detected language.
type LanguageInfo struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	Primary    bool    `json:"primary"`
}

// AutoPipelineConfig is a generated pipeline configuration.
type AutoPipelineConfig struct {
	RepoID     int64         `json:"repo_id"`
	Identifier string        `json:"identifier"`
	YAML       string        `json:"yaml"`
	Stack      DetectedStack `json:"detected_stack"`
	Generated  int64         `json:"generated"`
}

// GeneratePipelineInput is the request body for generating an auto-pipeline.
type GeneratePipelineInput struct {
	Branch string `json:"branch,omitempty"` // defaults to repo default branch
}
