//  Copyright 2023 Harness, Inc.
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

package base

import (
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/audit"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/filemanager"
	registryrefcache "github.com/EolaFam1828/SoloDev/registry/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/registry/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"

	"github.com/google/wire"
)

func LocalBaseProvider(
	registryDao store.RegistryRepository,
	registryFinder registryrefcache.RegistryFinder,
	fileManager filemanager.FileManager,
	tx dbtx.Transactor,
	imageDao store.ImageRepository,
	artifactDao store.ArtifactRepository,
	nodesDao store.NodesRepository,
	tagsDao store.PackageTagRepository,
	authorizer authz.Authorizer,
	spaceFinder refcache.SpaceFinder,
	auditService audit.Service,
) LocalBase {
	return NewLocalBase(
		registryDao, registryFinder, fileManager, tx, imageDao, artifactDao, nodesDao,
		tagsDao, authorizer, spaceFinder, auditService,
	)
}

var WireSet = wire.NewSet(LocalBaseProvider)
