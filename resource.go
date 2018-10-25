package gowindams

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ImagePosition struct {
	Side *string  `json:"side"`
	X    *float32 `json:"x"`
	Y    *float32 `json:"y"`
}

const SideLeadingEdge = "LeadingEdge"
const SidePressureSide = "PressureSide"
const SideSuctionSide = "SuctionSide"
const SideTrailingEdge = "TrailingEdge"

type ImageScaleRequest struct {
	ResultType     string   `json:"resultType"`
	ScaleOperation string   `json:"scaleOperation"`
	Depth          *float64 `json:"depth"`
	Height         *float64 `json:"height"`
	Width          *float64 `json:"width"`
}

const ImageScaleResultTypeGIF = "GIF"
const ImageScaleResultTypeJPG = "JPEG"
const ImageScaleResultTypePNG = "PNG"
const ImageScaleOperationScaleToFit = "ScaleToFit"
const ImageScaleOperationScaleToHeight = "ScaleToHeight"
const ImageScaleOperationScaleToSize = "ScaleToSize"
const ImageScaleOperationScaleToWidth = "ScaleToWidth"

type SensorReading struct {
	Type      *string      `json:"type"`
	Value     *float64     `json:"value"`
	StartTime *WindAMSTime `json:"startTime"`
	Endtime   *WindAMSTime `json:"endTime"`
}

const SensorReadingTypeBarometric = "Barometric"
const SensorReadingTypeRelAltitude = "RelativeAltitude"

type ResourceMetadata struct {
	ResourceId            *string         `json:"resourceId"`
	AssetId               *string         `json:"assetId"`
	AssetInspectionId     *string         `json:"assetInspectionId"`
	ComponentId           *string         `json:"componentId"`
	ComponentInspectionId *string         `json:"componentInspectionId"`
	ContentType           *string         `json:"contentType"`
	FormId                *string         `json:"formId"`
	Sequence              *int            `json:"sequence"`
	SourceResourceId      *string         `json:"sourceResourceId"`
	Timestamp             *WindAMSTime    `json:"timestamp"`
	Location              *GeoPoint       `json:"location"`
	Name                  *string         `json:"name"`
	OrderNumber           *string         `json:"orderNumber"`
	Pass                  *int32          `json:"pass"`
	Position              *ImagePosition  `json:"position"`
	ProcessedBy           *string         `json:"processedBy"`
	Readings              []SensorReading `json:"readings"`
	SiteId                *string         `json:"siteId"`
	Size                  *Dimension      `json:"size"`
	Status                *string         `json:"status"`
	SubmissionId          *string         `json:"submissionId"`
	DownloadURL           *string         `json:"downloadURL"`
	SourceURL             *string         `json:"sourceURL"`
	ZoomifyId             *string         `json:"zoomifyId"`
	ZoomifyURL            *string         `json:"zoomifyURL"`
}

type ResourceSearchCriteria struct {
	AssetId               *string `json:"assetId"`
	AssetInspectionId     *string `json:"assetInspectionId"`
	ComponentId           *string `json:"componentId"`
	ComponentInspectionId *string `json:"componentInspectionId"`
	SourceResourceId      *string `json:"sourceResourceId"`
	OrderNumber           *string `json:"orderNumber"`
	Pass                  *int32  `json:"pass"`
	SiteId                *string `json:"siteId"`
	Status                *string `json:"status"`
}

const ResourceStatusArchived = "Archived"
const ResourceStatusNotForDisplay = "NotForDisplay"
const ResourceStatusProcessed = "Processed"
const ResourceStatusQueuedForUpload = "QueuedForUpload"
const ResourceStatusReleased = "Released"

type ResourceServiceClient struct {
	env *Environment
}

const resourceRootURI = "%s/resource"
const resourceScaleURI = resourceRootURI + "/%s/scale"
const resourceGetURI = resourceRootURI + "/%s"
const resourceSaveURI = resourceRootURI
const resourceSearchURI = resourceRootURI + "/search"
const resourceUpDownloadURI = "%s/multimedia/%s"

func (client ResourceServiceClient) Get(resourceId string) (*ResourceMetadata, error) {
	log.Printf("Loading resource metadata for %s", resourceId)
	url := fmt.Sprintf(resourceGetURI, client.env.ServiceURI, resourceId)
	result := new(ResourceMetadata)
	err := executeRestCall(client.env, "GET", url, nil, result)
	return result, err
}

func (client ResourceServiceClient) Save(rmeta *ResourceMetadata) error {
	data, err := json.Marshal(rmeta)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(resourceSaveURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, nil)
	return err
}

func (client ResourceServiceClient) Scale(resourceId string, imageScaleRequest ImageScaleRequest) (*ResourceMetadata, error) {
	data, err := json.Marshal(imageScaleRequest)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(resourceScaleURI, client.env.ServiceURI, resourceId)
	result := new(ResourceMetadata)
	err = executeRestCall(client.env, "POST", url, data, result)
	return result, err
}

func (client ResourceServiceClient) Search(criteria *ResourceSearchCriteria) ([]ResourceMetadata, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, err
	}
	results := make([]ResourceMetadata, 0)
	url := fmt.Sprintf(resourceSearchURI, client.env.ServiceURI)
	err = executeRestCall(client.env, "POST", url, data, results)
	return results, err
}

func (client ResourceServiceClient) Download(resourceId string) (*io.ReadCloser, error) {
	url := fmt.Sprintf(resourceUpDownloadURI, client.env.ServiceURI, resourceId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	} else {
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Got response code %d downloading the image \"%s\" from \"%s\"", resp.StatusCode, resourceId, url)
		} else {
			return &resp.Body, err
		}
	}
}

func (client ResourceServiceClient) Upload(resourceId string, contentType string, body *io.Reader) error {
	url := fmt.Sprintf(resourceUpDownloadURI, client.env.ServiceURI, resourceId)
	resp, err := http.Post(url, contentType, *body)
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			return nil
		} else {
			return fmt.Errorf("Received response %d: %s", resp.StatusCode, resp.Status)
		}
	} else {
		return err
	}
}
