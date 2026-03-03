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
	"sort"
	"sync"

	"github.com/harness/gitness/types"
)

// Store is an in-memory vector store for code chunks.
// It supports per-repo indexing and approximate nearest-neighbor search via cosine similarity.
type Store struct {
	mu     sync.RWMutex
	chunks map[int64][]types.CodeChunk // repoID → chunks
}

// NewStore creates a new in-memory vector store.
func NewStore() *Store {
	return &Store{
		chunks: make(map[int64][]types.CodeChunk),
	}
}

// Index stores a code chunk with its precomputed embedding vector.
func (s *Store) Index(chunk types.CodeChunk) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chunks[chunk.RepoID] = append(s.chunks[chunk.RepoID], chunk)
}

// IndexBatch stores multiple code chunks at once.
func (s *Store) IndexBatch(chunks []types.CodeChunk) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, c := range chunks {
		s.chunks[c.RepoID] = append(s.chunks[c.RepoID], c)
	}
}

// ClearRepo removes all indexed chunks for a repository.
func (s *Store) ClearRepo(repoID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.chunks, repoID)
}

// Search finds the top-k most similar chunks to the query vector.
// If repoID > 0, only chunks from that repo are searched.
// If repoID == 0, all repos are searched.
func (s *Store) Search(queryVec []float64, repoID int64, topK int) []types.VectorSearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if topK <= 0 {
		topK = 5
	}

	var candidates []types.VectorSearchResult

	if repoID > 0 {
		chunks := s.chunks[repoID]
		for _, c := range chunks {
			sim := CosineSimilarity(queryVec, c.Vector)
			candidates = append(candidates, types.VectorSearchResult{
				Chunk:      c,
				Similarity: sim,
			})
		}
	} else {
		for _, chunks := range s.chunks {
			for _, c := range chunks {
				sim := CosineSimilarity(queryVec, c.Vector)
				candidates = append(candidates, types.VectorSearchResult{
					Chunk:      c,
					Similarity: sim,
				})
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Similarity > candidates[j].Similarity
	})

	if len(candidates) > topK {
		candidates = candidates[:topK]
	}

	return candidates
}

// RepoChunkCount returns the number of indexed chunks for a repository.
func (s *Store) RepoChunkCount(repoID int64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.chunks[repoID])
}

// TotalChunkCount returns the total number of indexed chunks across all repos.
func (s *Store) TotalChunkCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var total int
	for _, chunks := range s.chunks {
		total += len(chunks)
	}
	return total
}
