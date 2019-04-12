package gowindams

import (
	"encoding/json"
	"fmt"
	"log"
)

type WorkOrder struct {
	OrderNumber *string `json:"orderNumber"`
	Description *string `json:"description"`
	RequestDate *string `json:"requestDate"`
	Scope       *string `json:"scope"`
	SiteId      *string `json:"siteId"`
	Status      *string `json:"status"`
	Type        *string `json:"type"`
}

type WorkOrderSearchCriteria struct {
	SiteId      *string  `json:"siteId"`
	Statuses    []string `json:"statuses"`
	Type        *string  `json:"workOrderType"`
}

const WORK_ORDER_SCOPE_ADHOC = "adHoc"
const WORK_ORDER_SCOPE_ALARM = "alarm"
const WORK_ORDER_SCOPE_ANNUAL = "annual"
const WORK_ORDER_SCOPE_END_OF_WARRANTY = "endOfWarranty"
const WORK_ORDER_SCOPE_TURBINE_ALARM = "turbineAlarm"
const WORK_ORDER_SCOPE_FAILURE = "failure"
const WORK_ORDER_SCOPE_INSPECTION = "inspection"
const WORK_ORDER_SCOPE_REPAIR = "repair"
const WORK_ORDER_SCOPE_RETROFIT = "retrofit"

const WORK_ORDER_STATUS_COMPLETED = "Completed"
const WORK_ORDER_STATUS_IMAGES_PROCESSED = "ImagesProcessed"
const WORK_ORDER_STATUS_IMAGES_UPLOADED = "ImagesUploaded"
const WORK_ORDER_STATUS_ONSITE = "Onsite"
const WORK_ORDER_STATUS_REQUESTED = "Requested"

const WORK_ORDER_TYPE = "TurbineBladeInspection"

const woRootURI = "%s/workOrder"
const woGetURI = woRootURI + "/%s"
const woSaveURI = woRootURI
const woSearchURI = woRootURI + "/search"

type WorkOrderServiceClient struct {
	env *Environment
}

func (client WorkOrderServiceClient) Create(obj *WorkOrder) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(woSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "PUT", url, data, nil)
	return err
}

func (client WorkOrderServiceClient) Get(id string) (*WorkOrder, error) {
	log.Printf("Loading work order for %s", id)
	url := fmt.Sprintf(woGetURI, client.env.ServiceURI, id)
	result := new(WorkOrder)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client WorkOrderServiceClient) Search(criteria *WorkOrderSearchCriteria) ([]WorkOrder, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]WorkOrder, 0)
	url := fmt.Sprintf(woSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, &results)
	return results, err
}

func (client WorkOrderServiceClient) Update(obj *WorkOrder) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(woSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
