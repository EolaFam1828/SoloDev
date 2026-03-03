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
