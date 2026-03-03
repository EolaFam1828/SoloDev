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

package sologate

import (
	"encoding/json"
	"net/http"

	apiauth "github.com/harness/gitness/app/api/auth"
	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

// HandleGetConfig returns a http.HandlerFunc that gets the Solo Gate config for a space.
func HandleGetConfig(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	configStore store.SoloGateConfigStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		space, err := spaceFinder.FindByRef(ctx, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		if err := apiauth.CheckSpace(ctx, authorizer, session, space, enum.PermissionSpaceView); err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		config, err := configStore.FindBySpaceID(ctx, space.ID)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		// Return defaults if no config exists.
		if config == nil {
			config = &types.SoloGateConfig{
				SpaceID:         space.ID,
				EnforcementMode: types.EnforcementModeStrict,
			}
		}

		render.JSON(w, http.StatusOK, config)
	}
}

// HandleUpdateConfig returns a http.HandlerFunc that updates the Solo Gate config for a space.
func HandleUpdateConfig(
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	configStore store.SoloGateConfigStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		spaceRef, err := request.GetSpaceRefFromPath(r)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		space, err := spaceFinder.FindByRef(ctx, spaceRef)
		if err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		if err := apiauth.CheckSpace(ctx, authorizer, session, space, enum.PermissionSpaceEdit); err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		var input types.UpdateSoloGateConfigInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			render.BadRequestf(ctx, w, "invalid request body: %s", err)
			return
		}

		// Fetch existing or create defaults.
		config, _ := configStore.FindBySpaceID(ctx, space.ID)
		if config == nil {
			config = &types.SoloGateConfig{
				SpaceID:         space.ID,
				EnforcementMode: types.EnforcementModeStrict,
			}
		}

		// Apply updates.
		if input.EnforcementMode != nil {
			config.EnforcementMode = *input.EnforcementMode
		}
		if input.AutoIgnoreLow != nil {
			config.AutoIgnoreLow = *input.AutoIgnoreLow
		}
		if input.AutoTriageKnown != nil {
			config.AutoTriageKnown = *input.AutoTriageKnown
		}
		if input.AIAutoFix != nil {
			config.AIAutoFix = *input.AIAutoFix
		}
		if input.LogTechDebt != nil {
			config.LogTechDebt = *input.LogTechDebt
		}

		if err := configStore.Upsert(ctx, config); err != nil {
			render.TranslatedUserError(ctx, w, err)
			return
		}

		render.JSON(w, http.StatusOK, config)
	}
}
