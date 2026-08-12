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
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/algotiqa/core"
	"github.com/algotiqa/core/auth/role"
	"github.com/algotiqa/core/req"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
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
	ccontext, provider := createContextAndProvider(client, authority)

	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
	}
	verifier := provider.Verifier(oidcConfig)

	return &OidcController{
		authority: authority,
		client:    client,
		context:   &ccontext,
		verifier:  verifier,
		logger:    logger,
		config:    config,
	}
}

//=============================================================================

func (oc *OidcController) Secure(h RestService, roles []role.Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		rawAccessToken := c.Request.Header.Get("Authorization")
		onBehalfOf := c.Request.Header.Get(req.OnBehalfOf)

		tokens := strings.Split(rawAccessToken, " ")
		if len(tokens) != 2 {
			req.ReturnUnauthorizedError(c, "Authorisation failed due to a bad header")
			return
		}

		idToken, err := oc.verifier.Verify(*oc.context, tokens[1])
		if err != nil {
			req.ReturnUnauthorizedError(c, "Authorisation failed while verifying the token: "+err.Error())
			return
		}

		var ut userToken
		if err := idToken.Claims(&ut); err != nil {
			req.ReturnUnauthorizedError(c, "Authorization failed while getting claims: "+err.Error())
			return
		}

		us := buildUserSession(&ut, idToken, onBehalfOf)

		if !us.IsUserInRole(roles) {
			req.ReturnForbiddenError(c, "User not allowed to access this API: "+us.Username)
			return
		}

		ctx := &Context{
			Gin:     c,
			Session: us,
			Log:     oc.createLogger(us, c),
			Config:  oc.config,
			Token:   tokens[1],
		}

		h(ctx)
	}
}

//=============================================================================

func (oc *OidcController) createLogger(us *UserSession, c *gin.Context) *slog.Logger {
	return oc.logger.With(
		slog.String("client", c.ClientIP()),
		slog.String("username", us.Username),
	).WithGroup("data")
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func createContextAndProvider(client *http.Client, authority string) (context.Context, *oidc.Provider){
	var ccontext context.Context
	var provider *oidc.Provider
	var err error

	slog.Info("Connecting to OIDC provider...")

	//--- Retry up to 50 secs to allow the identity provider to start
	//--- Issue: if this container fail fast, it is not restarted by Podman

	for i:=0; i<10; i++ {
		ccontext = oidc.ClientContext(context.Background(), client)
		provider, err = oidc.NewProvider(ccontext, authority)

		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
		slog.Info("Retrying to connect to OIDC provider...")
	}

	core.ExitIfError(err)

	return ccontext, provider
}

//=============================================================================

func buildUserSession(ut *userToken, it *oidc.IDToken, onBehalfOf string) *UserSession {
	if onBehalfOf == "" {
		onBehalfOf = ut.Username
	}

	return &UserSession{
		SessionID:  ut.SID,
		Username:   ut.Username,
		OnBehalfOf: onBehalfOf,
		Name:       ut.Name,
		Surname:    ut.Surname,
		Email:      ut.Email,
		IssuedAt:   it.IssuedAt,
		Expiry:     it.Expiry,
		Roles:      buildRoleMap(ut),
	}
}

//=============================================================================

func buildRoleMap(ut *userToken) map[role.Role]any {
	userRoles := map[role.Role]any{}

	for k, v := range ut.RealmAccess {
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
