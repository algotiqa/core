//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package auth

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/algotiqa/core"
	"github.com/algotiqa/core/req"
	"github.com/coreos/go-oidc/v3/oidc"
)

//=============================================================================

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

//=============================================================================

type RestContext struct {
	sync.RWMutex
	clientId      string
	clientSecret  string
	client        *http.Client
	provider      *oidc.Provider
	tokenResponse *TokenResponse
	tokenDate     time.Time
}

//=============================================================================

var restContext *RestContext

//=============================================================================
//===
//=== Public functions
//===
//=============================================================================

func InitAuthentication(auth *core.Authentication) {
	client   := req.GetDefaultClient()
	ccontext := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ccontext, auth.Authority)
	core.ExitIfError(err)

	restContext = &RestContext{
		clientId:     auth.ClientId,
		clientSecret: auth.ClientSecret,
		client:       client,
		provider:     provider,
	}
}

//=============================================================================

func Token() (string, error) {
	restContext.Lock()
	defer restContext.Unlock()

	if restContext.tokenResponse == nil || isTokenExpired() {
		t, err := getToken()
		if err != nil {
			slog.Error("Cannot get authentication token", "error", err)
			return "", err
		}

		restContext.tokenResponse = t
		restContext.tokenDate = time.Now()
	}

	return restContext.tokenResponse.AccessToken, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getToken() (*TokenResponse, error) {
	params := "grant_type=client_credentials&client_id=" + restContext.clientId + "&client_secret=" + restContext.clientSecret
	resp   := TokenResponse{}
	url    := restContext.provider.Endpoint().TokenURL

	body   := []byte(params)
	reader := bytes.NewReader(body)

	rq, err := http.NewRequest("POST", url, reader)
	if err != nil {
		slog.Error("Error creating a POST request", "error", err.Error())
		return nil, err
	}

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := restContext.client.Do(rq)
	err = req.BuildResponse(res, err, &resp)

	return &resp, err
}

//=============================================================================

func isTokenExpired() bool {
	date := restContext.tokenDate
	now := time.Now()

	curDur := now.Sub(date)
	maxDur := time.Duration(restContext.tokenResponse.ExpiresIn*9/10) * time.Second

	return curDur >= maxDur
}

//=============================================================================
