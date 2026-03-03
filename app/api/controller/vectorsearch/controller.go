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

package vectorsearch

import (
	"context"
	"fmt"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/services/vectorstore"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// Controller exposes vector search operations.
type Controller struct {
	authorizer  authz.Authorizer
	spaceFinder refcache.SpaceFinder
	repoFinder  refcache.RepoFinder
	store       *vectorstore.Store
	indexer     *vectorstore.Indexer
}

// ProvideController creates a new vector search controller.
func ProvideController(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	repoFinder refcache.RepoFinder,
	store *vectorstore.Store,
	indexer *vectorstore.Indexer,
) *Controller {
	return &Controller{
		authorizer:  authorizer,
		spaceFinder: spaceFinder,
		repoFinder:  repoFinder,
		store:       store,
		indexer:     indexer,
	}
}

// IndexInput is the request body for indexing a repository.
type IndexInput struct {
	RepoRef string `json:"repo_ref"`
	GitRef  string `json:"git_ref,omitempty"`
}

// IndexOutput is the response from indexing a repository.
type IndexOutput struct {
	RepoID     int64 `json:"repo_id"`
	ChunkCount int   `json:"chunk_count"`
}

// IndexRepo indexes a repository's code for vector search.
func (c *Controller) IndexRepo(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	input *IndexInput,
) (*IndexOutput, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}

	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}

	repo, err := c.repoFinder.FindByRef(ctx, input.RepoRef)
	if err != nil {
		return nil, fmt.Errorf("find repo: %w", err)
	}

	count, err := c.indexer.IndexRepo(ctx, repo.ID, input.GitRef)
	if err != nil {
		return nil, fmt.Errorf("index repo: %w", err)
	}

	return &IndexOutput{
		RepoID:     repo.ID,
		ChunkCount: count,
	}, nil
}

// SearchInput is the request for vector search.
type SearchInput struct {
	Query  string `json:"query"`
	RepoID int64  `json:"repo_id,omitempty"`
	TopK   int    `json:"top_k,omitempty"`
}

// Search performs a vector similarity search over indexed code.
func (c *Controller) Search(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
	input *SearchInput,
) ([]types.VectorSearchResult, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}

	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}

	topK := input.TopK
	if topK <= 0 {
		topK = 5
	}

	queryVec := vectorstore.Embed(input.Query)
	results := c.store.Search(queryVec, input.RepoID, topK)

	return results, nil
}

// StatsOutput reports index statistics.
type StatsOutput struct {
	TotalChunks int `json:"total_chunks"`
}

// Stats returns vector store statistics.
func (c *Controller) Stats(
	ctx context.Context,
	session *auth.Session,
	spaceRef string,
) (*StatsOutput, error) {
	space, err := c.spaceFinder.FindByRef(ctx, spaceRef)
	if err != nil {
		return nil, fmt.Errorf("find space: %w", err)
	}

	if err := apiauth.CheckSpace(ctx, c.authorizer, session, space, enum.PermissionSpaceView); err != nil {
		return nil, err
	}

	return &StatsOutput{
		TotalChunks: c.store.TotalChunkCount(),
	}, nil
}
