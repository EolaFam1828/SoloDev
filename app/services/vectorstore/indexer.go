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

package vectorstore

import (
	"bytes"
	"context"
	"crypto/sha256"
	stdliberrors "errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	appstore "github.com/harness/gitness/app/store"
	"github.com/harness/gitness/git"
	"github.com/harness/gitness/git/parser"
	"github.com/harness/gitness/types"

	"github.com/rs/zerolog/log"
)

const (
	chunkLines      = 60 // lines per chunk
	chunkOverlap    = 10 // overlap between consecutive chunks
	maxFileBytes    = 128 * 1024
	maxFilesPerRepo = 500
)

// knownCodeExtensions maps file extensions to language names for indexing.
var knownCodeExtensions = map[string]string{
	".go":    "go",
	".js":    "javascript",
	".ts":    "typescript",
	".tsx":   "typescript",
	".jsx":   "javascript",
	".py":    "python",
	".java":  "java",
	".rs":    "rust",
	".rb":    "ruby",
	".c":     "c",
	".cpp":   "cpp",
	".h":     "c",
	".hpp":   "cpp",
	".cs":    "csharp",
	".swift": "swift",
	".kt":    "kotlin",
	".scala": "scala",
	".sh":    "shell",
	".yaml":  "yaml",
	".yml":   "yaml",
	".json":  "json",
	".toml":  "toml",
	".sql":   "sql",
	".proto": "protobuf",
	".tf":    "terraform",
}

// Indexer crawls a repository's tree and indexes code chunks into the vector store.
type Indexer struct {
	repoStore appstore.RepoStore
	git       git.Interface
	store     *Store
}

// NewIndexer creates a new code indexer.
func NewIndexer(
	repoStore appstore.RepoStore,
	gitClient git.Interface,
	store *Store,
) *Indexer {
	return &Indexer{
		repoStore: repoStore,
		git:       gitClient,
		store:     store,
	}
}

// IndexRepo indexes all code files in a repository at the given git ref.
func (idx *Indexer) IndexRepo(ctx context.Context, repoID int64, gitRef string) (int, error) {
	repo, err := idx.repoStore.Find(ctx, repoID)
	if err != nil {
		return 0, fmt.Errorf("find repo: %w", err)
	}

	if gitRef == "" {
		gitRef = repo.DefaultBranch
	}
	if gitRef == "" {
		gitRef = "main"
	}

	// Clear previous index for this repo.
	idx.store.ClearRepo(repoID)

	// List the tree at root.
	entries, err := idx.listTree(ctx, repo.GitUID, gitRef, "")
	if err != nil {
		return 0, fmt.Errorf("list tree: %w", err)
	}

	var totalChunks int
	filesIndexed := 0

	for _, entry := range entries {
		if filesIndexed >= maxFilesPerRepo {
			break
		}
		if entry.Type != git.TreeNodeTypeBlob {
			continue
		}

		lang := detectLanguage(entry.Path)
		if lang == "" {
			continue
		}

		content, err := idx.fetchFileContent(ctx, repo.GitUID, gitRef, entry.Path)
		if err != nil {
			log.Ctx(ctx).Debug().Err(err).Str("path", entry.Path).Msg("skipping file")
			continue
		}
		if content == "" {
			continue
		}

		chunks := chunkContent(repoID, entry.Path, lang, content)
		for i := range chunks {
			chunks[i].Vector = Embed(chunks[i].Content)
		}

		idx.store.IndexBatch(chunks)
		totalChunks += len(chunks)
		filesIndexed++
	}

	log.Ctx(ctx).Info().
		Int64("repo_id", repoID).
		Int("files", filesIndexed).
		Int("chunks", totalChunks).
		Str("ref", gitRef).
		Msg("vector index built")

	return totalChunks, nil
}

// listTree recursively lists all tree entries under a path.
func (idx *Indexer) listTree(ctx context.Context, repoUID, gitRef, path string) ([]git.TreeNode, error) {
	resp, err := idx.git.ListTreeNodes(ctx, &git.ListTreeNodeParams{
		ReadParams: git.ReadParams{RepoUID: repoUID},
		GitREF:     gitRef,
		Path:       path,
	})
	if err != nil {
		return nil, err
	}

	var entries []git.TreeNode
	for _, node := range resp.Nodes {
		if node.Type == git.TreeNodeTypeBlob {
			entries = append(entries, node)
		} else if node.Type == git.TreeNodeTypeTree {
			subEntries, err := idx.listTree(ctx, repoUID, gitRef, node.Path)
			if err != nil {
				continue
			}
			entries = append(entries, subEntries...)
		}
	}
	return entries, nil
}

// fetchFileContent loads a file's text from git, returning empty string for binary or oversized files.
func (idx *Indexer) fetchFileContent(ctx context.Context, repoUID, gitRef, path string) (string, error) {
	node, err := idx.git.GetTreeNode(ctx, &git.GetTreeNodeParams{
		ReadParams: git.ReadParams{RepoUID: repoUID},
		GitREF:     gitRef,
		Path:       path,
	})
	if err != nil {
		return "", err
	}
	if node.Node.Type != git.TreeNodeTypeBlob {
		return "", nil
	}

	blob, err := idx.git.GetBlob(ctx, &git.GetBlobParams{
		ReadParams: git.ReadParams{RepoUID: repoUID},
		SHA:        node.Node.SHA,
		SizeLimit:  int64(maxFileBytes) + 1,
	})
	if err != nil {
		return "", err
	}
	defer blob.Content.Close()

	if blob.Size > int64(maxFileBytes) {
		return "", nil
	}

	raw, err := io.ReadAll(blob.Content)
	if err != nil {
		return "", err
	}
	if _, ok := parser.IsLFSPointer(ctx, raw, blob.Size); ok {
		return "", nil
	}

	scanner, _, err := parser.ReadTextFile(bytes.NewReader(raw), nil)
	if err != nil {
		if stdliberrors.Is(err, parser.ErrBinaryFile) {
			return "", nil
		}
		return "", err
	}

	var text strings.Builder
	for scanner.Scan() {
		text.Write(scanner.Bytes())
	}
	if scanner.Err() != nil && stdliberrors.Is(scanner.Err(), parser.ErrBinaryFile) {
		return "", nil
	}

	return text.String(), nil
}

// chunkContent splits file content into overlapping chunks.
func chunkContent(repoID int64, filePath, lang, content string) []types.CodeChunk {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil
	}

	var chunks []types.CodeChunk
	step := chunkLines - chunkOverlap
	if step <= 0 {
		step = chunkLines
	}

	for start := 0; start < len(lines); start += step {
		end := start + chunkLines
		if end > len(lines) {
			end = len(lines)
		}

		chunkText := strings.Join(lines[start:end], "\n")
		if strings.TrimSpace(chunkText) == "" {
			continue
		}

		id := chunkID(repoID, filePath, start+1)
		chunks = append(chunks, types.CodeChunk{
			ID:        id,
			RepoID:    repoID,
			FilePath:  filePath,
			StartLine: start + 1,
			EndLine:   end,
			Content:   chunkText,
			Language:  lang,
		})

		if end >= len(lines) {
			break
		}
	}

	return chunks
}

func chunkID(repoID int64, filePath string, startLine int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d:%s:%d", repoID, filePath, startLine)))
	return fmt.Sprintf("%x", h[:8])
}

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	return knownCodeExtensions[ext]
}
