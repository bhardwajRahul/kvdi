/*
Copyright 2020,2021 Avi Zimmerman

This file is part of kvdi.

kvdi is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

kvdi is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with kvdi.  If not, see <https://www.gnu.org/licenses/>.
*/

package vault

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault/api"

	appv1 "github.com/kvdi/kvdi/apis/app/v1"
)

// DefaultTokenPath is where the k8s serviceaccount token is mounted inside the
// container.
const DefaultTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"

// AuthRequest represents a request for a vault token using the k8s JWT.
// There is probably a struct defined in the libary for this somewhere.
type AuthRequest struct {
	JWT  string `json:"jwt"`
	Role string `json:"role"`
}

// getClientToken will read the k8s serviceaccount token and use it to request
// a vault login token.
func getK8sAuth(crConfig *appv1.VaultConfig, vaultConfig *api.Config) (*api.Secret, error) {
	tokenBytes, err := os.ReadFile(DefaultTokenPath)
	if err != nil {
		return nil, err
	}
	authURLStr := fmt.Sprintf("%s/v1/auth/kubernetes/login", vaultConfig.Address)
	body, err := json.Marshal(&AuthRequest{JWT: string(tokenBytes), Role: crConfig.GetAuthRole()})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, authURLStr, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	res, err := vaultConfig.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(resBody))
	}
	authResponse := &api.Secret{}
	return authResponse, json.Unmarshal(resBody, authResponse)
}

// runTokenRefreshLoop waits for 60 seconds before token expiry and either renews
// or requests a new token.
func (p *Provider) runTokenRefreshLoop(authInfo *api.Secret) {
	var err error
	ticker := newAuthTicker(authInfo.Auth)
	for {
		select {
		case <-p.stopCh:
			vaultLogger.Info("Stopping token refresh loop")
			return
		case <-ticker.C:
			vaultLogger.Info("Refreshing client token")
			if authInfo != nil && authInfo.Auth.Renewable {
				authInfo, err = p.client.Auth().Token().RenewSelf(authInfo.Auth.LeaseDuration)
				if err == nil {
					p.client.SetToken(authInfo.Auth.ClientToken)
					ticker = newAuthTicker(authInfo.Auth)
					continue
				}
				vaultLogger.Error(err, "Failed to renew token, requesting a new one")
				// If there was an error we can try a full login
			}
			var err error
			authInfo, err = p.getAuth(p.crConfig, p.vaultConfig)
			if err != nil {
				vaultLogger.Error(err, "Failed to acquire a new vault token, retrying in 10 seconds")
				ticker = time.NewTicker(time.Duration(10) * time.Second)
				continue
			}
			p.client.SetToken(authInfo.Auth.ClientToken)
			ticker = newAuthTicker(authInfo.Auth)
			continue
		}
	}
}

// newAuthTicker returns a ticker for 60 seconds before the expiry of the given
// token information.
func newAuthTicker(auth *api.SecretAuth) *time.Ticker {
	return time.NewTicker(time.Duration(auth.LeaseDuration-60) * time.Second)
}
