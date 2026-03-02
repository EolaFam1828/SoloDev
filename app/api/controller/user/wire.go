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

package user

import (
	"github.com/EolaFam1828/SoloDev/app/auth/authz"
	userevents "github.com/EolaFam1828/SoloDev/app/events/user"
	"github.com/EolaFam1828/SoloDev/app/services/refcache"
	"github.com/EolaFam1828/SoloDev/app/store"
	"github.com/EolaFam1828/SoloDev/store/database/dbtx"
	"github.com/EolaFam1828/SoloDev/types/check"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
	tx dbtx.Transactor,
	principalUIDCheck check.PrincipalUID,
	authorizer authz.Authorizer,
	principalStore store.PrincipalStore,
	tokenStore store.TokenStore,
	membershipStore store.MembershipStore,
	publicKeyStore store.PublicKeyStore,
	publicKeySubKeyStore store.PublicKeySubKeyStore,
	gitSignatureResultStore store.GitSignatureResultStore,
	eventReporter *userevents.Reporter,
	repoFinder refcache.RepoFinder,
	favoriteStore store.FavoriteStore,
) *Controller {
	return NewController(
		tx,
		principalUIDCheck,
		authorizer,
		principalStore,
		tokenStore,
		membershipStore,
		publicKeyStore,
		publicKeySubKeyStore,
		gitSignatureResultStore,
		eventReporter,
		repoFinder,
		favoriteStore)
}
