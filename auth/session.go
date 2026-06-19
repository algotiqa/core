//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package auth

import (
	"time"

	"github.com/algotiqa/core/auth/role"
)

//=============================================================================

type UserSession struct {
	SessionID  string
	Username   string
	OnBehalfOf string
	Name       string
	Surname    string
	Email      string
	IssuedAt   time.Time
	Expiry     time.Time
	Roles      map[role.Role]any
}

//=============================================================================

func (us *UserSession) IsUserInRole(roles []role.Role) bool {
	for _, r := range roles {
		if _, ok := us.Roles[r]; ok {
			return true
		}
	}

	return false
}

//=============================================================================

func (us *UserSession) IsAdmin() bool {
	_, ok := us.Roles[role.Admin]
	return ok
}

//=============================================================================
