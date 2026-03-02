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

	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// TrivyScanner runs Trivy SCA/vulnerability scanning.
type TrivyScanner struct {
	path string
}

// NewTrivyScanner creates a new Trivy scanner.
func NewTrivyScanner(path string) *TrivyScanner {
	if path == "" {
		path = "trivy"
	}
	return &TrivyScanner{path: path}
}

func (s *TrivyScanner) Name() string { return "trivy" }

func (s *TrivyScanner) Available() bool {
	_, err := exec.LookPath(s.path)
	return err == nil
}

type trivyOutput struct {
	Results []trivyResult `json:"Results"`
}

type trivyResult struct {
	Target          string               `json:"Target"`
	Vulnerabilities []trivyVulnerability `json:"Vulnerabilities"`
}

type trivyVulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	Title            string `json:"Title"`
	Description      string `json:"Description"`
	Severity         string `json:"Severity"`
}

func (s *TrivyScanner) Scan(ctx context.Context, repoDir string) ([]types.ScanFinding, error) {
	cmd := exec.CommandContext(ctx, s.path, "fs", "--format", "json", "--scanners", "vuln", repoDir)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("trivy failed: %w", err)
	}

	var parsed trivyOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse trivy output: %w", err)
	}

	var findings []types.ScanFinding
	for _, result := range parsed.Results {
		for _, vuln := range result.Vulnerabilities {
			desc := vuln.Description
			if desc == "" {
				desc = vuln.Title
			}
			findings = append(findings, types.ScanFinding{
				Identifier:  vuln.VulnerabilityID,
				Title:       fmt.Sprintf("%s in %s@%s", vuln.VulnerabilityID, vuln.PkgName, vuln.InstalledVersion),
				Description: desc,
				Severity:    mapTrivySeverity(vuln.Severity),
				Category:    enum.SecurityFindingCategorySCA,
				FilePath:    result.Target,
				Status:      enum.SecurityFindingStatusOpen,
			})
		}
	}

	return findings, nil
}

func mapTrivySeverity(s string) enum.SecurityFindingSeverity {
	switch strings.ToUpper(s) {
	case "CRITICAL":
		return enum.SecurityFindingSeverityCritical
	case "HIGH":
		return enum.SecurityFindingSeverityHigh
	case "MEDIUM":
		return enum.SecurityFindingSeverityMedium
	case "LOW":
		return enum.SecurityFindingSeverityLow
	default:
		return enum.SecurityFindingSeverityInfo
	}
}
