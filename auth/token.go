//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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
	"bytes"
	"context"
	"github.com/algotiqa/core"
	"github.com/algotiqa/core/req"
	"github.com/coreos/go-oidc/v3/oidc"
	"log/slog"
	"net/http"
	"sync"
	"time"
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
	client := req.GetClient("bf")
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
	resp := TokenResponse{}
	url := restContext.provider.Endpoint().TokenURL

	body := []byte(params)
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
