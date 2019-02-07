package gowindams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// In both of the below URLs, %s should be tenant ID such as precisionhawk.auth0.com
const auth0TokenURL = "https://%s/oauth2/token"
const auth0KeysURL = "https://%s/.well-known/jwks.json"

type auth0AccessTokenProvider struct {
	clientId string
	tenantId string
	clientSecret string
}

func (provider auth0AccessTokenProvider) queryAccessToken(resource string) (*AccessTokenResponse, error) {
	// The below code is the same for aadAccessTokenProvider and auth0TokenProvider except for URL building, but may be
	// different for others, so keeping it duplicated for now.
	params := make(url.Values)
	params["grant_type"] = []string{"client_credentials"}
	params["client_id"] = []string{provider.clientId}
	params["client_secret"] = []string{provider.clientSecret}
	params["resource"] = []string{resource}

	resp, err := http.PostForm(
		fmt.Sprintf(auth0TokenURL, provider.tenantId),
		params)
	if err != nil {
		return nil, err
	} else {
		atresp := new(AccessTokenResponse)
		data, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			eresp := new(AccessTokenErrorResponse)
			json.Unmarshal(data, &eresp)
			log.Printf("GOWINDAMS: Error from token request:\t%s", eresp.ErrorDescription)
			return nil, fmt.Errorf("%s", eresp.Error)
		} else {
			err = json.Unmarshal(data, &atresp)
			if err != nil {
				log.Printf("GOWINDAMS: Error obtaining access token for resource %s: %s\n", resource, err)
				return nil, err
			} else {
				//				log.Printf("GOWINDAMS: Successfully obtained access token for resource %s: %+v\n", resource, atresp)
				return atresp, nil
			}
		}
	}
}

func (provider auth0AccessTokenProvider) getWellKnown() ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(auth0KeysURL, provider.tenantId))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to obtain signing Keys, got response code %d", resp.StatusCode)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
}

/* If there is a client secret, we assume it's server to server.  Otherwise, it must be user authenticated. */

func (provider auth0AccessTokenProvider) isServerToServer() bool {
	return provider.clientSecret != ""
}

func (provider auth0AccessTokenProvider) isUserAuthenticated() bool {
	return ! provider.isServerToServer()
}

func (provider auth0AccessTokenProvider) getAuthenticationProviderType() AuthenticationProviderType {
	return AP_Auth0
}
