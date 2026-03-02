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
	"strings"

	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

// SemgrepScanner runs Semgrep SAST analysis.
type SemgrepScanner struct {
	path  string
	rules string
}

// NewSemgrepScanner creates a new Semgrep scanner.
func NewSemgrepScanner(path, rules string) *SemgrepScanner {
	if path == "" {
		path = "semgrep"
	}
	if rules == "" {
		rules = "auto"
	}
	return &SemgrepScanner{path: path, rules: rules}
}

func (s *SemgrepScanner) Name() string { return "semgrep" }

func (s *SemgrepScanner) Available() bool {
	_, err := exec.LookPath(s.path)
	return err == nil
}

type semgrepOutput struct {
	Results []semgrepResult `json:"results"`
}

type semgrepResult struct {
	CheckID string `json:"check_id"`
	Path    string `json:"path"`
	Start   struct {
		Line int `json:"line"`
	} `json:"start"`
	End struct {
		Line int `json:"line"`
	} `json:"end"`
	Extra struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	} `json:"extra"`
}

func (s *SemgrepScanner) Scan(ctx context.Context, repoDir string) ([]types.ScanFinding, error) {
	cmd := exec.CommandContext(ctx, s.path, "--config", s.rules, "--json", repoDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Semgrep exits with code 1 when findings exist — that is not an error.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// output is valid json with findings
		} else if len(output) == 0 {
			return nil, fmt.Errorf("semgrep failed: %w", err)
		}
	}

	var parsed semgrepOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse semgrep output: %w", err)
	}

	findings := make([]types.ScanFinding, 0, len(parsed.Results))
	for _, r := range parsed.Results {
		findings = append(findings, types.ScanFinding{
			Identifier:  r.CheckID,
			Title:       r.CheckID,
			Description: r.Extra.Message,
			Severity:    mapSemgrepSeverity(r.Extra.Severity),
			Category:    enum.SecurityFindingCategorySAST,
			FilePath:    r.Path,
			LineStart:   r.Start.Line,
			LineEnd:     r.End.Line,
			Status:      enum.SecurityFindingStatusOpen,
		})
	}

	return findings, nil
}

func mapSemgrepSeverity(s string) enum.SecurityFindingSeverity {
	switch strings.ToUpper(s) {
	case "ERROR":
		return enum.SecurityFindingSeverityHigh
	case "WARNING":
		return enum.SecurityFindingSeverityMedium
	case "INFO":
		return enum.SecurityFindingSeverityLow
	default:
		return enum.SecurityFindingSeverityInfo
	}
}
