package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type AssetInspection struct {
	Id               *string            `json:"id"`
	AssetId          *string            `json:"assetId"`
	Attributes       map[string] string `json:"attributes"`
	DateOfInspection *string            `json:"dateOfInspection"`
	FailureDate      *string            `json:"failureDate"`
	OrderNumber      *string            `json:"orderNumber"`
	ProcessedBy      *string            `json:"processedBy"`
	Resources        []string           `json:"resources"`
	SiteId           *string            `json:"siteId"`
	Status           *string            `json:"status"`
	Type             *string            `json:"type"`
}

type AssetInspectionSearchCriteria struct {
	AssetId          *string            `json:"assetId"`
	SiteId           *string            `json:"siteId"`
	OrderNumber      *string            `json:"orderNumber"`
	Status           *string            `json:"status"`
}

const ASSET_INSPECTION_STATUS_IN_PROCESS = "In_Process"
const ASSET_INSPECTION_STATUS_PENDING = "Pending"
const ASSET_INSPECTION_STATUS_PROCESSED = "Processed"
const ASSET_INSPECTION_STATUS_RELEASED = "Released"
const ASSET_INSPECTION_STATUS_UPLOADED = "Uploaded"

const ASSET_INSPECTION_TYPE_DRONE = "DroneBladeInspection"
const ASSET_INSPECTION_TYPE_GROUND = "GroundBladeInspection"

const assetInspectionRootURI = "%s/assetInspection"
const assetInspectionGetURI = assetInspectionRootURI + "/%s"
const assetInspectionSaveURI = assetInspectionRootURI
const assetInspectionSearchURI = assetInspectionRootURI + "/search"

type AssetInspectionServiceClient struct {
	env *Environment
}

func (client AssetInspectionServiceClient) Create(obj *AssetInspection) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(assetInspectionSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "PUT", url, data, nil)
	return err
}

func (client AssetInspectionServiceClient) Get(id string) (*AssetInspection, error) {
	log.Printf("Loading site for %s", id)
	url := fmt.Sprintf(assetInspectionGetURI, client.env.ServiceURI, id)
	result := new(AssetInspection)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client AssetInspectionServiceClient) Search(criteria *AssetInspectionSearchCriteria) ([]AssetInspection, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]AssetInspection, 0)
	url := fmt.Sprintf(assetInspectionSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}

func (client AssetInspectionServiceClient) Update(obj *AssetInspection) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(assetInspectionSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
