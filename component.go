package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type Component struct {
	Id            *string            `json:"id"`
	AssetId       *string            `json:"assetId"`
	Attributes    map[string] string `json:"attributes"`
	Make          *string            `json:"make"`
	Model         *string            `json:"model"`
	SerialNumber  *string            `json:"serialNumber"`
	SiteId        *string            `json:"siteId"`
	Type          *string            `json:"type"`
}

type ComponentSearchCriteria struct {
	AssetId       *string            `json:"assetId"`
	ComponentId   *string            `json:"componentId"`
	SerialNumber  *string            `json:"serialNumber"`
	SiteId        *string            `json:"siteId"`
	ComponentType *string            `json:"componentType"`
}

const COMPONENT_TYPE_BLADE_A = "BladeA"
const COMPONENT_TYPE_BLADE_B = "BladeB"
const COMPONENT_TYPE_BLADE_C = "BladeC"
const COMPONENT_TYPE_BLADE_D = "BladeD"
const COMPONENT_TYPE_BLADE_E = "BladeE"
const COMPONENT_TYPE_GEARBOX = "Gearbox"

const componentRootURI = "%s/component"
const componentGetURI = componentRootURI + "/%s"
const componentSaveURI = componentRootURI
const componentSearchURI = componentRootURI + "/search"

type ComponentServiceClient struct {
	env *Environment
}

func (client ComponentServiceClient) Create(obj *Component) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(componentSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "PUT", url, data, nil)
	return err
}

func (client ComponentServiceClient) Get(id string) (*Component, error) {
	log.Printf("Loading site for %s", id)
	url := fmt.Sprintf(componentGetURI, client.env.ServiceURI, id)
	result := new(Component)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client ComponentServiceClient) Search(criteria *ComponentSearchCriteria) ([]Component, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]Component, 0)
	url := fmt.Sprintf(componentSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}

func (client ComponentServiceClient) Update(obj *Component) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(componentSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
