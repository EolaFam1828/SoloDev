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
	"context"
	"fmt"

	apiauth "github.com/EolaFam1828/SoloDev/app/api/auth"
	"github.com/EolaFam1828/SoloDev/app/api/usererror"
	"github.com/EolaFam1828/SoloDev/app/auth"
	"github.com/EolaFam1828/SoloDev/types"
	"github.com/EolaFam1828/SoloDev/types/enum"
)

// Delete deletes a user.
func (c *Controller) Delete(ctx context.Context, session *auth.Session,
	userUID string) error {
	user, err := findUserFromUID(ctx, c.principalStore, userUID)
	if err != nil {
		return err
	}

	// Fail if the user being deleted is the only admin in DB
	if user.Admin {
		admUsrCount, err := c.principalStore.CountUsers(ctx, &types.UserFilter{Admin: true})
		if err != nil {
			return fmt.Errorf("failed to check admin user count: %w", err)
		}

		if admUsrCount == 1 {
			return usererror.BadRequest("Cannot delete the only admin user")
		}
	}

	// Ensure principal has required permissions on parent
	if err = apiauth.CheckUser(ctx, c.authorizer, session, user, enum.PermissionUserDelete); err != nil {
		return err
	}

	return c.principalStore.DeleteUser(ctx, user.ID)
}
