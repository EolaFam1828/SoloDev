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

package qualitygate

import (
	"encoding/json"
	"fmt"

	"github.com/EolaFam1828/SoloDev/types"
)

// SanitizeRepoIDs sanitizes and validates a list of repository IDs.
func SanitizeRepoIDs(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Validate all IDs are positive
	for _, id := range ids {
		if id <= 0 {
			return nil, fmt.Errorf("invalid repository ID: %d", id)
		}
	}

	// Remove duplicates
	seen := make(map[int64]bool)
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	return result, nil
}

// SanitizeBranches sanitizes and validates a list of branch patterns.
func SanitizeBranches(branches []string) ([]string, error) {
	if len(branches) == 0 {
		return nil, nil
	}

	// Remove duplicates and validate
	seen := make(map[string]bool)
	result := make([]string, 0, len(branches))
	for _, branch := range branches {
		if branch == "" {
			return nil, fmt.Errorf("branch pattern cannot be empty")
		}
		if !seen[branch] {
			seen[branch] = true
			result = append(result, branch)
		}
	}

	return result, nil
}

// SanitizeTags sanitizes and validates a list of tags.
func SanitizeTags(tags []string) ([]string, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	// Remove duplicates and validate
	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag == "" {
			return nil, fmt.Errorf("tag cannot be empty")
		}
		if !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}

	return result, nil
}

// EncodeJSON encodes any value to JSON bytes.
func EncodeJSON(data interface{}) (json.RawMessage, error) {
	if data == nil {
		return nil, nil
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}

	if len(bytes) == 0 || string(bytes) == "null" {
		return nil, nil
	}

	return json.RawMessage(bytes), nil
}

// DecodeJSON decodes JSON bytes to the target value.
func DecodeJSON(data json.RawMessage, target interface{}) error {
	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, target)
}

// DecodeJSONArray decodes JSON array bytes to a slice of values.
func DecodeJSONArray(data json.RawMessage) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var result []interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON array: %w", err)
	}

	return result, nil
}

// DecodeRepoIDs decodes JSON array to repository IDs.
func DecodeRepoIDs(data json.RawMessage) ([]int64, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var ids []int64
	err := json.Unmarshal(data, &ids)
	if err != nil {
		return nil, fmt.Errorf("failed to decode repository IDs: %w", err)
	}

	return ids, nil
}

// DecodeBranches decodes JSON array to branches.
func DecodeBranches(data json.RawMessage) ([]string, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var branches []string
	err := json.Unmarshal(data, &branches)
	if err != nil {
		return nil, fmt.Errorf("failed to decode branches: %w", err)
	}

	return branches, nil
}

// DecodeTags decodes JSON array to tags.
func DecodeTags(data json.RawMessage) ([]string, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var tags []string
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tags: %w", err)
	}

	return tags, nil
}

// DecodeResults decodes JSON array to evaluation results.
func DecodeResults(data json.RawMessage) ([]types.QualityEvaluationResult, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var results []types.QualityEvaluationResult
	err := json.Unmarshal(data, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to decode evaluation results: %w", err)
	}

	return results, nil
}
