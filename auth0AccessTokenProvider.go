package gowindams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// In both of the below URLs, %s should be tenant ID such as precisionhawk.auth0.com
const auth0TokenURL = "https://%s/oauth/token"
const auth0KeysURL = "https://%s/.well-known/jwks.json"

type auth0AccessTokenProvider struct {
	clientId string
	tenantId string
	clientSecret string
}

func (provider auth0AccessTokenProvider) queryAccessToken(resource string) (*AccessTokenResponse, error) {

	url := fmt.Sprintf(auth0TokenURL, provider.tenantId)
	payload := strings.NewReader(
		fmt.Sprintf(
			"{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"client_credentials\"}",
			provider.clientId, provider.clientSecret, resource))

	req, _ := http.NewRequest(http.MethodPost, url, payload)
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
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
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to obtain signing Keys, got response code %d", resp.StatusCode)
	} else {
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
