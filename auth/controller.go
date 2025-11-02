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
	"context"
	"encoding/json"
	"github.com/tradalia/core"
	"github.com/tradalia/core/auth/role"
	"github.com/tradalia/core/req"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

//=============================================================================

type OidcController struct {
	authority string
	client    *http.Client
	context   *context.Context
	verifier  *oidc.IDTokenVerifier
	logger    *slog.Logger
	config    any
}

//=============================================================================

type userToken struct {
	JTI      string `json:"jti,omitempty"`
	SID      string `json:"sid,omitempty"`
	Name     string `json:"given_name,omitempty"`
	Surname  string `json:"family_name,omitempty"`
	Username string `json:"preferred_username,omitempty"`
	Email    string `json:"email,omitempty"`

	RealmAccess map[string]json.RawMessage `json:"realm_access,omitempty"`
}

//=============================================================================

func NewOidcController(authority string, client *http.Client, logger *slog.Logger, config any) *OidcController {
	ccontext      := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ccontext, authority)
	core.ExitIfError(err)

	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
	}
	verifier := provider.Verifier(oidcConfig)

	return &OidcController{
		authority: authority,
		client   : client,
		context  : &ccontext,
		verifier : verifier,
		logger   : logger,
		config   : config,
	}
}

//=============================================================================

func (oc *OidcController) Secure(h RestService, roles []role.Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		rawAccessToken := c.Request.Header.Get("Authorization")
		onBehalfOf     := c.Request.Header.Get(req.OnBehalfOf)

		tokens := strings.Split(rawAccessToken, " ")
		if len(tokens) != 2 {
			req.ReturnUnauthorizedError(c, "Authorisation failed due to a bad header")
			return
		}

		idToken, err := oc.verifier.Verify(*oc.context, tokens[1])
		if err != nil {
			req.ReturnUnauthorizedError(c, "Authorisation failed while verifying the token: "+ err.Error())
			return
		}

		var ut userToken
		if err := idToken.Claims(&ut); err != nil {
			req.ReturnUnauthorizedError(c, "Authorization failed while getting claims: "+ err.Error())
			return
		}

		us := buildUserSession(&ut, idToken, onBehalfOf)

		if ! us.IsUserInRole(roles) {
			req.ReturnForbiddenError(c, "User not allowed to access this API: "+ us.Username)
			return
		}

		ctx := &Context{
			Gin    : c,
			Session: us,
			Log    : oc.createLogger(us, c),
			Config : oc.config,
			Token  : tokens[1],
		}

		h(ctx)
	}
}

//=============================================================================

func (oc *OidcController) createLogger(us *UserSession, c *gin.Context) *slog.Logger {
	return oc.logger.With(
		slog.String("client",   c.ClientIP()),
		slog.String("username", us.Username),
	).WithGroup("data")
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func buildUserSession(ut *userToken, it *oidc.IDToken, onBehalfOf string) *UserSession {
	if onBehalfOf == "" {
		onBehalfOf = ut.Username
	}

	return &UserSession{
		SessionID : ut.SID,
		Username  : ut.Username,
		OnBehalfOf: onBehalfOf,
		Name      : ut.Name,
		Surname   : ut.Surname,
		Email     : ut.Email,
		IssuedAt  : it.IssuedAt,
		Expiry    : it.Expiry,
		Roles     : buildRoleMap(ut),
	}
}

//=============================================================================

func buildRoleMap(ut *userToken) map[role.Role]any {
	userRoles := map[role.Role]any{}

	for k,v := range ut.RealmAccess {
		if k == "roles" {
			var realmRoles []string
			err := json.Unmarshal(v, &realmRoles)

			if err == nil {
				for _, r := range realmRoles {
					userRoles[role.Role(r)] = nil
				}
			}
		}
	}

	return userRoles
}

//=============================================================================
