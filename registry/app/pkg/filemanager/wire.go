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

package filemanager

import (
	"github.com/EolaFam1828/SoloDev/registry/app/events/replication"
	"github.com/EolaFam1828/SoloDev/registry/app/pkg/docker"
	"github.com/EolaFam1828/SoloDev/registry/app/services/hook"
	"github.com/EolaFam1828/SoloDev/registry/app/storage"
	"github.com/EolaFam1828/SoloDev/registry/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	gitnesstypes "github.com/EolaFam1828/SoloDev/types"

	"github.com/google/wire"
)

func Provider(
	registryDao store.RegistryRepository, genericBlobDao store.GenericBlobRepository,
	nodesDao store.NodesRepository,
	tx dbtx.Transactor,
	config *gitnesstypes.Config,
	storageService *storage.Service,
	bucketService docker.BucketService,
	replicationReporter replication.Reporter,
	blobActionHook hook.BlobActionHook,
) FileManager {
	// Pass the BucketService to use the unified implementation
	return NewFileManager(registryDao, genericBlobDao, nodesDao, tx,
		config, storageService, bucketService, replicationReporter, blobActionHook)
}

var Set = wire.NewSet(Provider)

var WireSet = wire.NewSet(Set)
