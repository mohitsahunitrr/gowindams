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

// In both of the below URLs, %s should be tenant ID such as precisionhawk.auth0.com
const auth0TokenURL = "https://%s/oauth2/token"
const auth0KeysURL = "https://%s/.well-known/jwks.json"

type auth0AccessTokenProvider struct {
	clientId string
	tenantId string
	clientSecret string
	mutex sync.Mutex
	tokenCache map[string] *AccessTokenResponse
}

func (provider auth0AccessTokenProvider) obtainAccessToken(resource string) (string, error) {
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

func (provider auth0AccessTokenProvider) queryAccessToken(resource string) (*AccessTokenResponse, error) {
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

func (provider auth0AccessTokenProvider) obtainSigningKeys() (map[string][]byte, error) {
	resp, err := http.Get(fmt.Sprintf(auth0KeysURL, provider.tenantId))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to obtain signing Keys, got response code %d", resp.StatusCode)
	} else {
		myjwts := new(jwts)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, myjwts)
		if err != nil {
			return nil, err
		}
		keys := make(map[string][]byte)
		for _, jwt := range myjwts.Keys {
			keys[jwt.Kid] = []byte(jwt.X5c[0])
		}
		return keys, nil
	}
}
