// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package system

import (
	"net/http"

	"github.com/harness/gitness/app/api/render"
	"github.com/harness/gitness/mcp"
)

func HandleGetMCPCatalog(controllers *mcp.Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, http.StatusOK, mcp.BuildCatalog(controllers))
	}
}
