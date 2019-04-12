package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type StatusEvent struct {
	Status               *string            `json:"status"`
	Timestamp            *WindAMSTime       `json:"timestamp"`
	UserId               *string            `json:"userId"`
}

type ComponentInspection struct {
	Id                   *string            `json:"id"`
	AssetId              *string            `json:"assetId"`
	AssetInspectionId    *string            `json:"assetInspectionId"`
	Attributes           map[string] string `json:"attributes"`
	ComponentId          *string            `json:"componentId"`
	Description          *string            `json:"description"`
	InspectionPosition   *string            `json:"inspectionPosition"`
	InspectionUserId     *string            `json:"inspectionUserId"`
	OrderNumber          *string            `json:"orderNumber"`
	ReasonForService     *string            `json:"reasonForService"`
	Resources            []string           `json:"resources"`
	PlateImageResourceId *string            `json:"plateImageResourceId"`
	SiteId               *string            `json:"siteId"`
	Source               *string            `json:"source"`
	StatusHistory        []StatusEvent      `json:"statusHistory"`
	Type                 *string            `json:"type"`
	VendorId             *string            `json:"vendorId"`
}

type ComponentInspectionSearchCriteria struct {
	ComponentId          *string            `json:"componentId"`
	SiteId               *string            `json:"siteId"`
	OrderNumber          *string            `json:"orderNumber"`
	Status               *string            `json:"status"`
}

const COMP_INSPECTION_SOURCE_PH          = "InspecTools"
const COMP_INSPECTION_VESTAS             = "Vestas"

const COMP_INSPECTION_STATUS_APPROVED    = "InspectionApproved"
const COMP_INSPECTION_STATUS_COMPLETED   = "InspectionCompleted"
const COMP_INSPECTION_STATUS_LOCKED      = "InspectionLocked"
const COMP_INSPECTION_STATUS_STARTED     = "InspectionStarted"
const COMP_INSPECTION_STATUS_SUBMITTED   = "InspectionSubmitted"
const COMP_INSPECTION_STATUS_TRANSMITTED = "InspectionTransmitted"

const COMP_INSPECTION_TYPE_DRONE = "DroneBladeInspection"
const COMP_INSPECTION_TYPE_GROUND = "GroundBladeInspection"

const componentInspectionRootURI = "%s/componentInspection"
const componentInspectionGetURI = componentInspectionRootURI + "/%s"
const componentInspectionSaveURI = componentInspectionRootURI
const componentInspectionSearchURI = componentInspectionRootURI + "/search"

type ComponentInspectionServiceClient struct {
	env *Environment
}

func (client ComponentInspectionServiceClient) Create(obj *ComponentInspection) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(componentInspectionSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "PUT", url, data, nil)
	return err
}

func (client ComponentInspectionServiceClient) Get(id string) (*ComponentInspection, error) {
	log.Printf("Loading site for %s", id)
	url := fmt.Sprintf(componentInspectionGetURI, client.env.ServiceURI, id)
	result := new(ComponentInspection)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client ComponentInspectionServiceClient) Search(criteria *ComponentInspectionSearchCriteria) ([]ComponentInspection, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]ComponentInspection, 0)
	url := fmt.Sprintf(componentInspectionSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}

func (client ComponentInspectionServiceClient) Update(obj *ComponentInspection) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(componentInspectionSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
