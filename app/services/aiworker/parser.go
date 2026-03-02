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
	"regexp"
	"strconv"
	"strings"
)

// ParsedResponse holds the extracted diff and confidence from an AI response.
type ParsedResponse struct {
	Diff       string
	Confidence float64
	RawText    string
}

var (
	diffBlockRe   = regexp.MustCompile("(?s)```(?:diff)?\\s*\\n(.*?)```")
	confidenceRe  = regexp.MustCompile(`(?i)CONFIDENCE:\s*([\d.]+)`)
)

// ParseAIResponse extracts the diff block and confidence score from raw AI output.
func ParseAIResponse(raw string) ParsedResponse {
	result := ParsedResponse{
		RawText: raw,
	}

	// Extract diff block.
	if matches := diffBlockRe.FindStringSubmatch(raw); len(matches) > 1 {
		result.Diff = strings.TrimSpace(matches[1])
	}

	// Extract confidence score.
	if matches := confidenceRe.FindStringSubmatch(raw); len(matches) > 1 {
		if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
			if val >= 0 && val <= 1.0 {
				result.Confidence = val
			}
		}
	}

	return result
}
