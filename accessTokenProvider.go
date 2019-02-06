package gowindams

import "strings"

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

type jwt struct {
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	X5t string `json:"x5t"`
	N string `json:"n"`
	E string `json:"e"`
	X5c []string `json:"x5c"`
}

type jwts struct {
	Keys []jwt `json:"Keys"`
}

type accessTokenProvider interface {
	obtainAccessToken(string) (string, error)
	obtainSigningKeys() (map[string][]byte, error)
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
