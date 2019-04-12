package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type Asset struct {
	Id            *string            `json:"id"`
	Attributes    map[string] string `json:"attributes"`
	DateOfInstall *string            `json:"dateOfInstall"`
	Location      *GeoPoint          `json:"location"`
	Make          *string            `json:"make"`
	Model         *string            `json:"model"`
	Name          *string            `json:"name"`
	SerialNumber  *string            `json:"serialNumber"`
	SiteId        *string            `json:"siteId"`
	Type          *string            `json:"type"`
}

type AssetSearchCriteria struct {
	AssetId       *string            `json:"assetId"`
	SiteId        *string            `json:"siteId"`
	AssetType     *string            `json:"assetType"`
	Name          *string            `json:"name"`
	SerialNumber  *string            `json:"serialNumber"`
}

const ASSET_TYPE_SOLAR_PANEL = "Solar_Panel"
const ASSET_TYPE_WIND_TURBINE = "Wind_Turbine_Tower"

const assetRootURI = "%s/asset"
const assetGetURI = assetRootURI + "/%s"
const assetSaveURI = assetRootURI
const assetSearchURI = assetRootURI + "/search"

type AssetServiceClient struct {
	env *Environment
}

func (client AssetServiceClient) Create(obj *Asset) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(assetSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "PUT", url, data, nil)
	return err
}

func (client AssetServiceClient) Get(id string) (*Asset, error) {
	log.Printf("Loading site for %s", id)
	url := fmt.Sprintf(assetGetURI, client.env.ServiceURI, id)
	result := new(Asset)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client AssetServiceClient) Search(criteria *AssetSearchCriteria) ([]Asset, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]Asset, 0)
	url := fmt.Sprintf(assetSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}

func (client AssetServiceClient) Update(obj *Asset) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(assetSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
