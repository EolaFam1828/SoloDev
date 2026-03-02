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

package generic

import (
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	gitnessstore "github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/audit"
	"github.com/EolaFam1828/SoloDev/registry/app/api/interfaces"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/filemanager"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/generic"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/quarantine"
	"github.com/EolaFam1828/SoloDev/registry/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"

	"github.com/google/wire"
)

func DBStoreProvider(
	imageDao store.ImageRepository,
	artifactDao store.ArtifactRepository,
	bandwidthStatDao store.BandwidthStatRepository,
	downloadStatDao store.DownloadStatRepository,
	registryDao store.RegistryRepository,
) *DBStore {
	return NewDBStore(registryDao, imageDao, artifactDao, bandwidthStatDao, downloadStatDao)
}

func ControllerProvider(
	spaceStore gitnessstore.SpaceStore,
	authorizer authz.Authorizer,
	fileManager filemanager.FileManager,
	dBStore *DBStore,
	tx dbtx.Transactor,
	spaceFinder refcache.SpaceFinder,
	local generic.LocalRegistry,
	proxy generic.Proxy,
	quarantineFinder quarantine.Finder,
	dependencyFirewallChecker interfaces.DependencyFirewallChecker,
	auditService audit.Service,
	packageWrapper interfaces.PackageWrapper,
) *Controller {
	return NewController(
		spaceStore,
		authorizer,
		fileManager,
		dBStore,
		tx,
		spaceFinder,
		local,
		proxy,
		quarantineFinder,
		dependencyFirewallChecker,
		auditService,
		packageWrapper,
	)
}

var DBStoreSet = wire.NewSet(DBStoreProvider)
var ControllerSet = wire.NewSet(ControllerProvider)

var WireSet = wire.NewSet(ControllerSet, DBStoreSet)
