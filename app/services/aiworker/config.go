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

// Config holds configuration for the AI remediation worker.
type Config struct {
	Enabled         bool
	Provider        string // anthropic, openai, gemini, ollama
	APIKey          string
	Model           string
	MaxTokens       int
	Temperature     float64
	CreateFixBranch bool

	// AutoApplyMinConfidence, when > 0, auto-applies (creates fix branch + draft PR)
	// completed remediations whose confidence score meets or exceeds this threshold.
	// Overrides CreateFixBranch for qualifying remediations.
	AutoApplyMinConfidence float64

	// AutoValidateAfterApply, when true, automatically triggers pipeline validation
	// after a successful auto-apply.
	AutoValidateAfterApply bool

	// AutoMergeAfterValidation, when true, auto-merges the draft PR after
	// pipeline validation passes. Completes the fully autonomous self-healing loop.
	AutoMergeAfterValidation bool
}
