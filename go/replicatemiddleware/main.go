package replicatemiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"planetcastdev/httpmiddleware"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type ReplicateConnectProps struct {
	Logger *zap.Logger
}

type Replicate struct {
	logger               *zap.Logger
	postRequestSemaphore *semaphore.Weighted
	getRequestSemaphore  *semaphore.Weighted
}

func Connect(args ReplicateConnectProps) *Replicate {
	maxPostRequestWorkers := 10
	postReqSem := semaphore.NewWeighted(int64(maxPostRequestWorkers))

	maxGetRequestWorkers := 50
	getReqSem := semaphore.NewWeighted(int64(maxGetRequestWorkers))

	return &Replicate{
		logger:               args.Logger,
		postRequestSemaphore: postReqSem,
		getRequestSemaphore:  getReqSem,
	}
}

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

// @PrcTrace
func (r *Replicate) MakeRequest(ctx context.Context, body *bytes.Buffer, url string) (any, error) {

	requestId, err := r.TriggerRequest(ctx, body, url)

	if err != nil {
		return "", err
	}

	for {
		requestOutput, err := r.FetchRequest(ctx, requestId)
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

// @PrcTrace
func (r *Replicate) FetchRequest(ctx context.Context, requestId string) (*ReplicateGetRequestOutput, error) {

	API_KEY := os.Getenv("REPLICATE_KEY")

	if err := r.getRequestSemaphore.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("Failed to acquire semaphore.")
	}
	defer r.getRequestSemaphore.Release(1)

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

// @PrcTrace
func (r *Replicate) TriggerRequest(ctx context.Context, body *bytes.Buffer, url string) (string, error) {

	API_KEY := os.Getenv("REPLICATE_KEY")

	if err := r.postRequestSemaphore.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("Failed to acquire semaphore.")
	}
	defer r.postRequestSemaphore.Release(1)

	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "POST",
		Url:    url,
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
