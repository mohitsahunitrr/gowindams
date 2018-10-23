package gowindams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const tokenURL = "https://login.microsoftonline.com/%s/oauth2/token"

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type`
	ExpiresIn int64 `json:"expires_in`
	ExpiresOn int64 `json:"expires_on`
	NotBefore int64 `json:"not_before`
	Resource string `json:"resource`
}

type accessTokenProvider struct {
	clientId string
	tenantId string
	clientSecret string
	mutex sync.Mutex
	tokenCache map[string] *AccessTokenResponse
}

func (provider accessTokenProvider) obtainAccessToken(resource string) (string, error) {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()

	// If there is a valid cert in the cache, use it.
	resp, exists := provider.tokenCache[resource]
	if exists {
		// See if the cert has expired
		now := time.Now().Unix()
		if now >= resp.ExpiresOn {
			exists = false
		}
	}
	if exists {
		return resp.AccessToken, nil
	}

	// Query for a new token
	resp, err := provider.queryAccessToken(resource)
	if err != nil {
		return "", err
	} else {
		// Cache the token
		provider.tokenCache[resource] = resp
		return resp.AccessToken, nil
	}
}

func (provider accessTokenProvider) queryAccessToken(resource string) (*AccessTokenResponse, error) {
	params := make(url.Values)
	params["grant_type"] = []string{"grant_type"}
	params["client_id"] = []string{provider.clientId}
	params["client_secret"] = []string{provider.clientSecret}
	params["resource"] = []string{resource}

	resp, err := http.PostForm(
		fmt.Sprintf(tokenURL, provider.tenantId),
		params)
	if err != nil {
		return nil, err
	} else {
		atresp := new(AccessTokenResponse)
		data, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(data, &atresp)
		if err != nil {
			log.Printf("GOWINDAMS: Error obtaining access token for resource %s: %s\n", resource, err)
			return nil, err
		} else {
			log.Printf("GOWINDAMS: Successfully obtained access token for resource %s\n", resource)
			return atresp, nil
		}
	}
}

func NewProvider(envCfg *EnvironmentConfig) *accessTokenProvider {
	provider := accessTokenProvider {
		clientId: envCfg.ClientId,
		tenantId: envCfg.TenantId,
		clientSecret: envCfg.ClientSecret,
		tokenCache: make(map[string] *AccessTokenResponse),
	}
	return &provider
}
