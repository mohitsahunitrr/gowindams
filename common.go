package gowindams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const DATE_FORMAT = "yyyyMMdd"
                        //20060102T150405
const DATE_TIME_FORMAT = "20060102T150405.000-0700" //"yyyyMMddTHHmmss.SSSZ"
const TIME_FORMAT = "'T'HHmmss.SSSZ"

type WindAMSTime time.Time

func (t WindAMSTime) Format(layout string) string {
	return time.Time(t).Format(layout)
}

func (t *WindAMSTime) MarshalJSON() ([]byte, error) {
	s := "\"" + t.Format(DATE_TIME_FORMAT) + "\""
	return []byte(s), nil
}

func (t *WindAMSTime) UnmarshalJSON(data []byte) error {
	s := string(data[:])
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		s := string([]rune(s)[1:len(s) - 1]) // Strip the enclosing quotes
		u, err := time.Parse(DATE_TIME_FORMAT, s)
		v := WindAMSTime(u)
		t = &v
		return err
	} else {
		return errors.New(fmt.Sprintf("Timestamp value (%s) must be enclosed with double-quotes in the JSON.", s))
	}
}

type Dimension struct {
	Depth *float64  `json:"depth"`
	Height *float64 `json:"height"`
	Width *float64  `json:"width"`
}

type GeoPoint struct {
	Accuracy *float64 `json:"accuracy"`
	Altitude *float64 `json:"altitude"`
	Latitude *float64 `json:"latitude"`
}

func executeRestCall(env *Environment, action string, url string, data []byte, results interface{}) error {
	client := &http.Client{
	}

	req, err := http.NewRequest(action, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("GOWINDAMS: Error building http request for %s against %s: %s\n", action, url, err)
		return err
	}
	token, err := env.obtainAccessToken()
	if err != nil {
		return err
	}

	log.Printf("GOWINDAMS: Executing %s against endpoint %s", action, url)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	log.Printf("GOWINDAMS: Response with status code %d for %s against endpoint %s", resp.StatusCode, action, url)

	if err != nil {
		log.Printf("GOWINDAMS: Got error for %s against %s: %s\n", action, url, err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GOWINDAMS: Got error getting response body for %s against %s: %s\n", action, url, err)
		return err
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 204 && results == nil {
			// This is fine.  Processed ok, no content, but we don't expect any.
		} else {
			s := string(body)
			log.Printf("GOWINDAMS: Got status code %d for %s against %s: %s\n", resp.StatusCode, action, url, s)
			return errors.New(s)
		}
	}
	if results != nil {
		err = json.Unmarshal(body, results)
	}
	return err
}
