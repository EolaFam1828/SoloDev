// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package contextengine

import "encoding/json"

// ContextSource identifies where a fragment originated.
type ContextSource string

const (
	SourceErrorLog        ContextSource = "error_log"
	SourceUserInput       ContextSource = "user_input"
	SourceGitFetch        ContextSource = "git_fetch"
	SourceSecurityFinding ContextSource = "security_finding"
	SourceScanMetadata    ContextSource = "scan_metadata"
	SourceVectorSearch    ContextSource = "vector_search"
)

// ContextFragment is a single piece of context with provenance.
type ContextFragment struct {
	Label        string        `json:"label"`
	Content      string        `json:"content"`
	Source       ContextSource `json:"source"`
	FilePath     string        `json:"file_path,omitempty"`
	TrimmedBytes int64         `json:"trimmed_bytes,omitempty"`
}

// ContextBundle is the structured context assembled for one remediation.
type ContextBundle struct {
	Fragments     []ContextFragment `json:"fragments"`
	TotalCharsEst int               `json:"total_chars_est"`
	TriggerSource string            `json:"trigger_source"`
	TriggerRef    string            `json:"trigger_ref,omitempty"`
	Branch        string            `json:"branch,omitempty"`
	CommitSHA     string            `json:"commit_sha,omitempty"`
	FilePath      string            `json:"file_path,omitempty"`
}

// AddFragment appends a fragment and updates the character estimate.
func (b *ContextBundle) AddFragment(f ContextFragment) {
	b.Fragments = append(b.Fragments, f)
	b.TotalCharsEst += len(f.Content)
}

// JSON returns the bundle as JSON bytes for metadata storage.
func (b *ContextBundle) JSON() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}
