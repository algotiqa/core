//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package auth

import (
	"github.com/algotiqa/core/auth/role"
	"time"
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
