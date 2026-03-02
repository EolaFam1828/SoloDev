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

package cargo

import (
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/filemanager"
	"github.com/EolaFam1828/SoloDev/registry/app/store"

	"github.com/google/wire"
)

func LocalRegistryHelperProvider(
	fileManager filemanager.FileManager,
	artifactDao store.ArtifactRepository,
	spaceFinder refcache.SpaceFinder,
) RegistryHelper {
	return NewRegistryHelper(fileManager, artifactDao, spaceFinder)
}

var WireSet = wire.NewSet(LocalRegistryHelperProvider)
