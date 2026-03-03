// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package contextengine

import (
	"bytes"
	"context"
	"encoding/json"
	stdliberrors "errors"
	"fmt"
	"io"
	"strings"

	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/git"
	"github.com/harness/gitness/git/parser"
	"github.com/harness/gitness/types"
)

const maxSourceCodeBytes int64 = 64 * 1024

// Service builds structured context bundles for AI remediation.
type Service struct {
	repoStore store.RepoStore
	git       git.Interface
}

// NewService creates a new context engine service.
func NewService(
	repoStore store.RepoStore,
	gitClient git.Interface,
) *Service {
	return &Service{
		repoStore: repoStore,
		git:       gitClient,
	}
}

// BuildContext assembles a ContextBundle from a remediation record,
// enriching source code from git when possible.
func (s *Service) BuildContext(ctx context.Context, rem *types.Remediation) (*ContextBundle, error) {
	bundle := &ContextBundle{
		TriggerSource: string(rem.TriggerSource),
		TriggerRef:    rem.TriggerRef,
		Branch:        rem.Branch,
		CommitSHA:     rem.CommitSHA,
		FilePath:      rem.FilePath,
	}

	// Fragment 1: Error log (always present, required field).
	if rem.ErrorLog != "" {
		bundle.AddFragment(ContextFragment{
			Label:   "Error Log",
			Content: rem.ErrorLog,
			Source:  SourceErrorLog,
		})
	}

	// Fragment 2: Source code — use existing or enrich from git.
	sourceCode := rem.SourceCode
	sourceOrigin := SourceUserInput
	if sourceCode == "" && rem.FilePath != "" && rem.RepoID > 0 {
		fetched, err := s.fetchSourceFromGit(ctx, rem)
		if err == nil && fetched != "" {
			sourceCode = fetched
			sourceOrigin = SourceGitFetch
			rem.SourceCode = sourceCode
		}
	}
	if sourceCode != "" {
		bundle.AddFragment(ContextFragment{
			Label:    "Source Code",
			Content:  sourceCode,
			Source:   sourceOrigin,
			FilePath: rem.FilePath,
		})
	}

	// Fragment 3: Description (if provided and distinct from error log).
	if rem.Description != "" && rem.Description != rem.ErrorLog {
		bundle.AddFragment(ContextFragment{
			Label:   "Description",
			Content: rem.Description,
			Source:  SourceUserInput,
		})
	}

	// Fragment 4: Security metadata (if trigger is security_scan).
	if rem.TriggerSource == types.RemediationTriggerSecurity && len(rem.Metadata) > 0 {
		secMeta := extractSecurityMeta(rem.Metadata)
		if secMeta != "" {
			bundle.AddFragment(ContextFragment{
				Label:   "Security Context",
				Content: secMeta,
				Source:  SourceSecurityFinding,
			})
		}
	}

	return bundle, nil
}

// fetchSourceFromGit retrieves file content from the git repository.
func (s *Service) fetchSourceFromGit(ctx context.Context, rem *types.Remediation) (string, error) {
	gitRef := strings.TrimSpace(rem.CommitSHA)
	if gitRef == "" {
		gitRef = strings.TrimSpace(rem.Branch)
	}
	if gitRef == "" || s.repoStore == nil || s.git == nil {
		return "", nil
	}

	repo, err := s.repoStore.Find(ctx, rem.RepoID)
	if err != nil {
		return "", fmt.Errorf("find repo: %w", err)
	}

	node, err := s.git.GetTreeNode(ctx, &git.GetTreeNodeParams{
		ReadParams: git.ReadParams{RepoUID: repo.GitUID},
		GitREF:     gitRef,
		Path:       rem.FilePath,
	})
	if err != nil {
		return "", fmt.Errorf("resolve file: %w", err)
	}
	if node.Node.Type != git.TreeNodeTypeBlob {
		return "", nil
	}

	blob, err := s.git.GetBlob(ctx, &git.GetBlobParams{
		ReadParams: git.ReadParams{RepoUID: repo.GitUID},
		SHA:        node.Node.SHA,
		SizeLimit:  maxSourceCodeBytes + 1,
	})
	if err != nil {
		return "", fmt.Errorf("load blob: %w", err)
	}
	defer blob.Content.Close()

	if blob.Size > maxSourceCodeBytes {
		return "", nil
	}

	content, err := io.ReadAll(blob.Content)
	if err != nil {
		return "", fmt.Errorf("read blob: %w", err)
	}
	if _, ok := parser.IsLFSPointer(ctx, content, blob.Size); ok {
		return "", nil
	}

	scanner, _, err := parser.ReadTextFile(bytes.NewReader(content), nil)
	if err != nil {
		if stdliberrors.Is(err, parser.ErrBinaryFile) {
			return "", nil
		}
		return "", fmt.Errorf("inspect source: %w", err)
	}

	var text strings.Builder
	for scanner.Scan() {
		text.Write(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		if stdliberrors.Is(err, parser.ErrBinaryFile) {
			return "", nil
		}
		return "", fmt.Errorf("scan source: %w", err)
	}

	return text.String(), nil
}

// extractSecurityMeta formats security-related metadata fields for the prompt.
func extractSecurityMeta(metadata json.RawMessage) string {
	var meta map[string]interface{}
	if err := json.Unmarshal(metadata, &meta); err != nil {
		return ""
	}

	var b strings.Builder
	fields := []struct {
		key   string
		label string
	}{
		{"severity", "Severity"},
		{"category", "Category"},
		{"rule", "Rule"},
		{"cwe", "CWE"},
	}

	for _, f := range fields {
		if v, ok := meta[f.key]; ok && v != nil {
			if s, ok := v.(string); ok && s != "" {
				fmt.Fprintf(&b, "- **%s:** %s\n", f.label, s)
			}
		}
	}

	return strings.TrimSpace(b.String())
}
