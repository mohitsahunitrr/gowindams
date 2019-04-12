package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type Site struct {
	Id             *string `json:"id"`
	Name           *string `json:"name"`
	OrganizationId *string `json:"organizationId"`
}

type SiteSearchCriteria struct {
	OrganizationId *string `json:"organizationId"`
}

const siteRootURI = "%s/site"
const siteGetURI = siteRootURI + "/%s"
const siteSaveURI = siteRootURI
const siteSearchURI = siteRootURI + "/search"

type SiteServiceClient struct {
	env *Environment
}

func (client SiteServiceClient) Create(obj *Site) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(siteSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "PUT", url, data, nil)
	return err
}

func (client SiteServiceClient) Get(id string) (*Site, error) {
	log.Printf("Loading site for %s", id)
	url := fmt.Sprintf(siteGetURI, client.env.ServiceURI, id)
	result := new(Site)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client SiteServiceClient) Search(criteria *SiteSearchCriteria) ([]Site, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]Site, 0)
	url := fmt.Sprintf(siteSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}

func (client SiteServiceClient) Update(obj *Site) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(siteSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
