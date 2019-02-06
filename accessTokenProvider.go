package gowindams

import (
	"crypto/rsa"
	"github.com/lestrrat/go-jwx/jwk"
	"strings"
	"sync"
	"time"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type`
	ExpiresIn int64 `json:"expires_in`
	ExpiresOn int64 `json:"expires_on`
	NotBefore int64 `json:"not_before`
	Resource string `json:"resource`
}

type AccessTokenErrorResponse struct {
	Error string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes []int `json:"error_codes"`
	Timestamp string `json:"timestamp"`
	TraceId string `json:"trace_id"`
	CorrelationId string `json:"correlation_id"`
}

type accessTokenProvider interface {
	getMutex() *sync.Mutex
	getTokenCache() *map[string] *AccessTokenResponse
	getWellKnown() ([]byte, error)
	queryAccessToken(resource string) (*AccessTokenResponse, error)
}

func NewProvider(envCfg *EnvironmentConfig) accessTokenProvider {
	s := strings.ToLower(envCfg.AccessTokenProvider)

	if strings.Contains(s, "auth0") {
		provider := auth0AccessTokenProvider{
			clientId:     envCfg.ClientId,
			tenantId:     envCfg.TenantId,
			clientSecret: envCfg.ClientSecret,
			tokenCache:   make(map[string]*AccessTokenResponse),
		}
		return &provider
	}
	if strings.Contains(s, "aad") {
		provider := aadAccessTokenProvider{
			clientId:     envCfg.ClientId,
			tenantId:     envCfg.TenantId,
			clientSecret: envCfg.ClientSecret,
			tokenCache:   make(map[string]*AccessTokenResponse),
		}
		return &provider
	}
	return nil
}


func obtainAccessToken(provider accessTokenProvider, resource string) (string, error) {
	provider.getMutex().Lock()
	defer provider.getMutex().Unlock()

	tokenCache := *provider.getTokenCache()
	// If there is a valid cert in the cache, use it.
	resp, exists := tokenCache[resource]
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
		tokenCache[resource] = resp
		return resp.AccessToken, nil
	}
}


func obtainSigningKeys(provider accessTokenProvider) (map[string]interface{}, error) {
	body, err := provider.getWellKnown()
	if err != nil {
		return nil, err
	}
	keys := make(map[string]interface{})
	set, _ := jwk.Parse(body)
	for _, key := range set.Keys {
		publicKey, _ := key.Materialize()
		if publicKey, ok := publicKey.(*rsa.PublicKey); ok {
			keys[key.KeyID()] = publicKey
		} // else, not an rsa key
	}
	return keys, nil
}
