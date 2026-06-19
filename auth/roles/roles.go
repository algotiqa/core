//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package roles

import "github.com/algotiqa/core/auth/role"

//=============================================================================

var Admin = []role.Role{role.Admin}
var User = []role.Role{role.User}
var Service = []role.Role{role.Service}

var Admin_User = []role.Role{role.Admin, role.User}
var Admin_Service = []role.Role{role.Admin, role.Service}
var Admin_User_Service = []role.Role{role.Admin, role.User, role.Service}

//=============================================================================
