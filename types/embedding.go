// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package types

// CodeChunk represents a chunk of source code stored in the vector index.
type CodeChunk struct {
	ID        string    `json:"id"`
	RepoID    int64     `json:"repo_id"`
	FilePath  string    `json:"file_path"`
	StartLine int       `json:"start_line"`
	EndLine   int       `json:"end_line"`
	Content   string    `json:"content"`
	Language  string    `json:"language,omitempty"`
	Vector    []float64 `json:"-"` // embedding vector, not serialized
}

// VectorSearchResult pairs a code chunk with its similarity score.
type VectorSearchResult struct {
	Chunk      CodeChunk `json:"chunk"`
	Similarity float64   `json:"similarity"`
}
