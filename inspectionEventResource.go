package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type InspectionEventPolygon struct {
	Center *GeoPoint					`json:"center"'`
	Geometry []GeoPoint					`json:"geometry"`
	Id *string							`json:"id"`
	Name *string						`json:"name"`
	Severity *int8						`json:"severity"`
	Text *string						`json:"text"`
}

type InspectionEventResource struct {
	AssetId *string						`json:"assetId"`
	Id *string							`json:"id"`
	InspectionEventId *string			`json:"inspectionEventId"`
	OrderNumber *string					`json:"orderNumber"`
	Polygons []InspectionEventPolygon	`json:"polygons"`
	ResourceId *string					`json:"resourceId"`
	SiteId *string						`json:"siteId"`
}

type InspectionEventResourceSearchCriteria struct {
	InspectionEventId *string			`json:"inspectionEventId"`
	ResourceId *string					`json:"resourceId"`
}

type InspectionEventResourceServiceClient struct {
	env *Environment
}

const ieRootURI = "%s/inspectionEventResource"
const ieGetURI = ieRootURI + "/%s"
const ieSaveURI = ieRootURI
const ieSearchURI = ieRootURI + "/search"

func (client InspectionEventResourceServiceClient) Get(id string) (*ResourceMetadata, error) {
	log.Printf("Loading inspection event resource for %s", id)
	url := fmt.Sprintf(ieGetURI, client.env.ServiceURI, id)
	result := new(ResourceMetadata)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client InspectionEventResourceServiceClient) Save(obj *InspectionEventResource) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(ieSaveURI, client.env.ServiceURI)
	if obj.Id == nil {
		err = executeRestCall(client.env, "PUT", url, data, obj)
	} else {
		err = executeRestCall(client.env, "POST", url, data, nil)
	}
	return err
}

func (client InspectionEventResourceServiceClient) Search(criteria *InspectionEventResourceSearchCriteria) ([]InspectionEventResource, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]InspectionEventResource, 0)
	url := fmt.Sprintf(ieSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}
