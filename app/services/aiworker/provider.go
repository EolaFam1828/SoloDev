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
	"fmt"
)

// LLMRequest represents a request to an LLM provider.
type LLMRequest struct {
	SystemPrompt string
	UserPrompt   string
	MaxTokens    int
	Temperature  float64
}

// LLMResponse represents a response from an LLM provider.
type LLMResponse struct {
	Content    string
	TokensUsed int
}

// LLMProvider is the interface for all LLM provider adapters.
type LLMProvider interface {
	Name() string
	Complete(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
}

// ProvideProvider creates an LLM provider based on configuration.
func ProvideProvider(config Config) (LLMProvider, error) {
	switch config.Provider {
	case "anthropic":
		return NewAnthropicProvider(config.APIKey, config.Model), nil
	case "openai":
		return NewOpenAIProvider(config.APIKey, config.Model), nil
	case "gemini":
		return NewGeminiProvider(config.APIKey, config.Model), nil
	case "ollama":
		return NewOllamaProvider(config.Model), nil
	case "":
		return nil, nil // No provider configured — AI remediation disabled.
	default:
		return nil, fmt.Errorf("unknown AI remediation provider: %s", config.Provider)
	}
}
