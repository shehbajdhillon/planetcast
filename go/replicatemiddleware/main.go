package replicatemiddleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"planetcastdev/httpmiddleware"
	"time"
)

type ReplicateGetRequestOutput struct {
	ID     string `json:"id"`
	Output *any   `json:"output"`
	Status string `json:"status"`
}

type ReplicateTriggerRequestOutput struct {
	ID     string  `json:"id"`
	Error  *string `json:"error"`
	Status string  `json:"status"`
}

func MakeRequest(body *bytes.Buffer) (any, error) {

	requestId, err := TriggerRequest(body)

	if err != nil {
		return "", err
	}

	for {
		requestOutput, err := FetchRequest(requestId)
		if err != nil {
			return "", err
		}
		if requestOutput.Status == "succeeded" {
			return *requestOutput.Output, nil
		}
		if requestOutput.Status == "failed" {
			return "", fmt.Errorf("Replicate failed to sync video")
		}
		time.Sleep(500 * time.Millisecond)
	}

}

func FetchRequest(requestId string) (*ReplicateGetRequestOutput, error) {

	API_KEY := os.Getenv("REPLICATE_KEY")

	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "GET",
		Url:    fmt.Sprintf("https://api.replicate.com/v1/predictions/%s", requestId),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Token %s", API_KEY),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot make call to replicate: %s", err.Error())
	}

	var replicateOutput ReplicateGetRequestOutput
	json.Unmarshal(responseBody, &replicateOutput)

	return &replicateOutput, nil

}

func TriggerRequest(body *bytes.Buffer) (string, error) {

	API_KEY := os.Getenv("REPLICATE_KEY")

	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "POST",
		Url:    "https://api.replicate.com/v1/predictions",
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Token %s", API_KEY),
		},
		Body: body,
	})

	if err != nil {
		return "", fmt.Errorf("Cannot make call to replicate: %s", err.Error())
	}

	var replicateOutput ReplicateTriggerRequestOutput
	json.Unmarshal(responseBody, &replicateOutput)

	if replicateOutput.Error != nil {
		return "", fmt.Errorf("Replicate request failed: %s", *replicateOutput.Error)
	}

	return replicateOutput.ID, nil
}
