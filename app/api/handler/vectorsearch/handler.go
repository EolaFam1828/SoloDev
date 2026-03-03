// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package vectorsearch

import (
	"encoding/json"
	"net/http"

	"github.com/harness/gitness/app/api/controller/vectorsearch"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
)

// HandleIndex returns a handler that indexes a repository for vector search.
func HandleIndex(ctrl *vectorsearch.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		var input vectorsearch.IndexInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.BadRequestf(ctx, w, "invalid request body: %s", err)
			return
		}

		out, err := ctrl.IndexRepo(ctx, session, spaceRef, &input)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, out)
	}
}

// HandleSearch returns a handler that performs vector similarity search.
func HandleSearch(ctrl *vectorsearch.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		var input vectorsearch.SearchInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.BadRequestf(ctx, w, "invalid request body: %s", err)
			return
		}

		results, err := ctrl.Search(ctx, session, spaceRef, &input)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, results)
	}
}

// HandleStats returns a handler that reports vector index statistics.
func HandleStats(ctrl *vectorsearch.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		stats, err := ctrl.Stats(ctx, session, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, stats)
	}
}
