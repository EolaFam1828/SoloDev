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

package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// GitleaksScanner runs Gitleaks secret detection.
type GitleaksScanner struct {
	path string
}

// NewGitleaksScanner creates a new Gitleaks scanner.
func NewGitleaksScanner(path string) *GitleaksScanner {
	if path == "" {
		path = "gitleaks"
	}
	return &GitleaksScanner{path: path}
}

func (s *GitleaksScanner) Name() string { return "gitleaks" }

func (s *GitleaksScanner) Available() bool {
	_, err := exec.LookPath(s.path)
	return err == nil
}

type gitleaksResult struct {
	Description string `json:"Description"`
	File        string `json:"File"`
	StartLine   int    `json:"StartLine"`
	EndLine     int    `json:"EndLine"`
	RuleID      string `json:"RuleID"`
}

func (s *GitleaksScanner) Scan(ctx context.Context, repoDir string) ([]types.ScanFinding, error) {
	cmd := exec.CommandContext(ctx, s.path, "detect", "--source", repoDir,
		"--report-format", "json", "--report-path", "/dev/stdout", "--no-git")
	output, err := cmd.Output()
	if err != nil {
		// Gitleaks exits with code 1 when leaks are found.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// findings exist — output was captured on stdout
		} else {
			return nil, fmt.Errorf("gitleaks failed: %w", err)
		}
	}

	if len(output) == 0 {
		return nil, nil
	}

	var results []gitleaksResult
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, fmt.Errorf("failed to parse gitleaks output: %w", err)
	}

	findings := make([]types.ScanFinding, 0, len(results))
	for _, r := range results {
		findings = append(findings, types.ScanFinding{
			Identifier:  r.RuleID,
			Title:       r.RuleID,
			Description: r.Description,
			Severity:    enum.SecurityFindingSeverityCritical,
			Category:    enum.SecurityFindingCategorySecret,
			FilePath:    r.File,
			LineStart:   r.StartLine,
			LineEnd:     r.EndLine,
			Status:      enum.SecurityFindingStatusOpen,
		})
	}

	return findings, nil
}
