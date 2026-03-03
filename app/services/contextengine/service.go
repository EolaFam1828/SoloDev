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

package contextengine

import (
	"bytes"
	"context"
	"encoding/json"
	stdliberrors "errors"
	"fmt"
	"io"
	"strings"

	"github.com/harness/gitness/app/services/vectorstore"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/git"
	"github.com/harness/gitness/git/parser"
	"github.com/harness/gitness/types"

	"github.com/rs/zerolog/log"
)

const (
	maxSourceCodeBytes  int64 = 64 * 1024
	vectorTopK                = 3
	vectorMinSimilarity       = 0.15
)

// Service builds structured context bundles for AI remediation.
type Service struct {
	repoStore   store.RepoStore
	git         git.Interface
	vectorStore *vectorstore.Store
	vectorIdx   *vectorstore.Indexer
}

// NewService creates a new context engine service.
func NewService(
	repoStore store.RepoStore,
	gitClient git.Interface,
	vectorStore *vectorstore.Store,
	vectorIdx *vectorstore.Indexer,
) *Service {
	return &Service{
		repoStore:   repoStore,
		git:         gitClient,
		vectorStore: vectorStore,
		vectorIdx:   vectorIdx,
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

	// Fragment 5: Vector-retrieved related code (Context Engine v2).
	if s.vectorStore != nil && rem.RepoID > 0 {
		s.addVectorContext(ctx, bundle, rem)
	}

	return bundle, nil
}

// addVectorContext queries the vector store for related code chunks and adds them as fragments.
func (s *Service) addVectorContext(ctx context.Context, bundle *ContextBundle, rem *types.Remediation) {
	// Ensure the repo is indexed (lazy indexing on first access).
	if s.vectorStore.RepoChunkCount(rem.RepoID) == 0 && s.vectorIdx != nil {
		if _, err := s.vectorIdx.IndexRepo(ctx, rem.RepoID, rem.Branch); err != nil {
			log.Ctx(ctx).Warn().Err(err).Int64("repo_id", rem.RepoID).Msg("vector indexing failed")
			return
		}
	}

	// Build query from error log + description + file path.
	var queryParts []string
	if rem.ErrorLog != "" {
		queryParts = append(queryParts, rem.ErrorLog)
	}
	if rem.Description != "" {
		queryParts = append(queryParts, rem.Description)
	}
	if rem.FilePath != "" {
		queryParts = append(queryParts, rem.FilePath)
	}
	if len(queryParts) == 0 {
		return
	}

	queryVec := vectorstore.Embed(strings.Join(queryParts, " "))
	results := s.vectorStore.Search(queryVec, rem.RepoID, vectorTopK)

	for _, result := range results {
		// Skip low-similarity results.
		if result.Similarity < vectorMinSimilarity {
			continue
		}
		// Skip the exact file already included as source code.
		if result.Chunk.FilePath == rem.FilePath {
			continue
		}

		label := fmt.Sprintf("Related Code (%s:%d-%d, sim=%.2f)",
			result.Chunk.FilePath, result.Chunk.StartLine, result.Chunk.EndLine, result.Similarity)

		bundle.AddFragment(ContextFragment{
			Label:    label,
			Content:  result.Chunk.Content,
			Source:   SourceVectorSearch,
			FilePath: result.Chunk.FilePath,
		})
	}
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
