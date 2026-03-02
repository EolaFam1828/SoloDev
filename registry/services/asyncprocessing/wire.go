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

package asyncprocessing

import (
	"context"

	"github.com/EolaFam1828/SoloDev/app/services/locker"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/events"
	"github.com/EolaFam1828/SoloDev/registry/app/api/interfaces"
	"github.com/EolaFam1828/SoloDev/registry/app/events/asyncprocessing"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/filemanager"
	"github.com/EolaFam1828/SoloDev/registry/app/store"
	"github.com/EolaFam1828/SoloDev/registry/app/utils/cargo"
	"github.com/EolaFam1828/SoloDev/registry/app/utils/gopackage"
	"github.com/EolaFam1828/SoloDev/secret"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideRegistryPostProcessingConfig,
	ProvideService,
	ProvideRpmHelper,
)

func ProvideService(
	ctx context.Context,
	tx dbtx.Transactor,
	rpmRegistryHelper RpmHelper,
	cargoRegistryHelper cargo.RegistryHelper,
	gopackageRegistryHelper gopackage.RegistryHelper,
	locker *locker.Locker,
	artifactsReaderFactory *events.ReaderFactory[*asyncprocessing.Reader],
	config Config,
	registryDao store.RegistryRepository,
	taskRepository store.TaskRepository,
	taskSourceRepository store.TaskSourceRepository,
	taskEventRepository store.TaskEventRepository,
	eventsSystem *events.System,
	postProcessingReporter *asyncprocessing.Reporter,
	packageWrapper interfaces.PackageWrapper,
) (*Service, error) {
	return NewService(
		ctx,
		tx,
		rpmRegistryHelper,
		cargoRegistryHelper,
		gopackageRegistryHelper,
		locker,
		artifactsReaderFactory,
		config,
		registryDao,
		taskRepository,
		taskSourceRepository,
		taskEventRepository,
		eventsSystem,
		postProcessingReporter,
		packageWrapper,
	)
}

func ProvideRegistryPostProcessingConfig(config *types.Config) Config {
	return Config{
		EventReaderName: config.InstanceID,
		Concurrency:     config.Registry.PostProcessing.Concurrency,
		MaxRetries:      config.Registry.PostProcessing.MaxRetries,
		AllowLoopback:   config.Registry.PostProcessing.AllowLoopback,
	}
}

func ProvideRpmHelper(
	fileManager filemanager.FileManager,
	artifactDao store.ArtifactRepository,
	upstreamProxyStore store.UpstreamProxyConfigRepository,
	spaceFinder refcache.SpaceFinder,
	secretService secret.Service,
	registryDao store.RegistryRepository,
) RpmHelper {
	return NewRpmHelper(
		fileManager,
		artifactDao,
		upstreamProxyStore,
		spaceFinder,
		secretService,
		registryDao,
	)
}
