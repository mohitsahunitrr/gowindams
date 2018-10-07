package gowindams

import (
	"encoding/json"
	"testing"
	"time"
)

const Example_RMeta =
	`{
		"formId":null,
		"resourceId":"72241b14-c2c8-479f-91b5-df82c82f3333",
		"assetInspectionId":"6c9d4e91-2b5b-4fbc-9fbb-0a518a479465",
		"componentId":"06dd6a73-4b15-4f63-b338-c44a0499f3a0",
		"orderNumber":"1234",
		"pass":null,
		"componentInspectionId":"674c7d59-5a7e-466f-a239-7b8bb36320dc",
		"zoomifyID":"9c259c19-8384-46c7-882f-345ed48ca194",
		"sourceURL":"https://cirblobstoragedev.blob.core.windows.net/cirdevcontainer/85e7bfd2-dba8-4376-944a-907272659617.jpeg",
		"submissionId":null,
		"assetId":"db538d4a-05e0-4e92-8c36-3fccd4d9dc58",
		"name":"511390001.jpg",
		"siteId":"d326dda1-4b60-429e-bdf5-7ac41081613e",
		"position":{"side":"PressureSide","x":0.5,"y":0.16},
		"contentType":"image/jpeg",
		"status":"Released",
		"timestamp":"20151216T224819.000+0000"
	}`

const Example_Timestamp = "20151216T224819.000+0000"

func TestDateParse(testing *testing.T) {
	t, err := time.Parse(DATE_TIME_FORMAT, Example_Timestamp)
	if err != nil {
	    testing.Fatal("error:", err)
	}
	if t.Year() != 2015 {
		testing.Fatalf("Invalid year %d", t.Year())
	}
	if t.Month() != 12 {
		testing.Fatalf("Invalid month %d", t.Month())
	}
	if t.Day() != 16 {
		testing.Fatalf("Invalid year %d", t.Year())
	}
	if t.Hour() != 22 {
		testing.Fatalf("Invalid hour %d", t.Year())
	}
	if t.Minute() != 48 {
		testing.Fatalf("Invalid minutes %d", t.Year())
	}
	if t.Second() != 19 {
		testing.Fatalf("Invalid seconds %d", t.Year())
	}
	if t.Nanosecond() != 0 {
		testing.Fatalf("Invalid nonseconds %d", t.Year())
	}
}

func TestUnMarshalResourceMetadata(testing *testing.T) {
	rmeta := new(ResourceMetadata)
	err := json.Unmarshal([]byte(Example_RMeta), &rmeta)
	if err != nil {
	    testing.Fatal("error:", err)
	}

	if rmeta.ResourceId == nil {
		testing.Fatal("Invalid resourceId: ", rmeta.ResourceId)
	}

	if rmeta.Timestamp == nil {
		testing.Fatal("invalid timestamp: ", rmeta.Timestamp)
	}
}
