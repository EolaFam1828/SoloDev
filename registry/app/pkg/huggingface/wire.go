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

package huggingface

import (
	urlprovider "github.com/EolaFam1828/SoloDev/app/url"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/base"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/filemanager"
	"github.com/EolaFam1828/SoloDev/registry/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"

	"github.com/google/wire"
)

// WireSet provides a wire set for the huggingface package.
var WireSet = wire.NewSet(
	LocalRegistryProvider,
)

// ProvideLocalRegistry provides a huggingface local registry.
func LocalRegistryProvider(
	localBase base.LocalBase,
	fileManager filemanager.FileManager,
	proxyStore store.UpstreamProxyConfigRepository,
	tx dbtx.Transactor,
	registryDao store.RegistryRepository,
	imageDao store.ImageRepository,
	artifactDao store.ArtifactRepository,
	urlProvider urlprovider.Provider,
) LocalRegistry {
	registry := NewLocalRegistry(localBase, fileManager, proxyStore, tx, registryDao, imageDao, artifactDao,
		urlProvider)
	base.Register(registry)
	return registry
}
