package gowindams

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type ProcessQueueEntry struct {
	Id *int64           `json:"id"`
	ObjectId *string    `json:"objectId"`
	ProcessType *string `json:"processType"`
}

const ProcessTypeCIRScale = "CIR_Scale"
const ProcessTypeZoomifyImage = "Zoomify_Img"

type ProcessQueueServiceClient struct {
	env *Environment
}

const paramError = "error"
const paramProcessor = "processor"
const paramProcessType = "processType"

const pqRootURI = "%s/processQueue"
const pqClaimURI = pqRootURI + "/claim?%s"
const pqEnqueueURI = pqRootURI + "/enqueue"
const pqErroredURI = pqRootURI + "/%d/errored?%s"
const pqProcessedURI = pqRootURI + "/processed"


func (client ProcessQueueServiceClient) Claim(processorId string, processType string) ([]ProcessQueueEntry, error) {
	v := url.Values{}
	v.Add(paramProcessor, processorId)
	v.Add(paramProcessType, processType)
	url := fmt.Sprintf(pqClaimURI, client.env.ServiceURI, v.Encode())
	results := make([]ProcessQueueEntry, 0)
	err := executeRestCall(client.env, "GET", url, nil, &results)
	return results, err
}

func (client ProcessQueueServiceClient) Enqueue(entries []ProcessQueueEntry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(pqEnqueueURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}

func (client ProcessQueueServiceClient) MarkErrored(entryId int64, error string) error {
	v := url.Values{}
	v.Add(paramError, error)
	url := fmt.Sprintf(pqErroredURI, client.env.ServiceURI, entryId, v.Encode())
	err := executeRestCall(client.env, "POST", url, nil, nil)
	return err
}

func (client ProcessQueueServiceClient) MarkProcessed(entries []ProcessQueueEntry) error {
	data, err := json.Marshal(entries);
	if err != nil {
		return err
	}
	url := fmt.Sprintf(pqProcessedURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}
