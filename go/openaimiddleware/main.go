package openaimiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"planetcastdev/httpmiddleware"
	"planetcastdev/utils"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequestInput struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionChoice struct {
	Index   int                   `json:"index"`
	Message ChatCompletionMessage `json:"message"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
}

type OpenAIConnectProps struct {
	Logger *zap.Logger
}

type OpenAI struct {
	logger    *zap.Logger
	semaphore *semaphore.Weighted
}

func Connect(args OpenAIConnectProps) *OpenAI {
	maxWorkers := 2
	sem := semaphore.NewWeighted(int64(maxWorkers))
	return &OpenAI{logger: args.Logger, semaphore: sem}
}

type MakeAPIRequestProps struct {
	Retries      int
	RequestInput ChatRequestInput
}

func (o *OpenAI) MakeAPIRequest(ctx context.Context, args MakeAPIRequestProps) (*ChatCompletionResponse, error) {

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	URL := "https://api.openai.com/v1/chat/completions"

	chatGptInput := args.RequestInput
	retries := args.Retries

	jsonData, err := json.Marshal(chatGptInput)
	if err != nil {
		return nil, fmt.Errorf("Could not generate request body: " + err.Error())
	}

	for retries >= 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - retries)

		if err := o.semaphore.Acquire(ctx, 1); err != nil {
			return nil, fmt.Errorf("Failed to acquire semaphore.")
		}
		defer o.semaphore.Release(1)

		respBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Method: "POST",
			Url:    URL,
			Body:   bytes.NewBuffer(jsonData),
			Headers: map[string]string{
				"Authorization": "Bearer " + API_KEY,
				"Content-Type":  "application/json",
			},
		})
		time.Sleep(1 * time.Second)

		if err != nil {
			o.logger.Error(
				"Could not make request to OpenAI. Retrying after sleeping.",
				zap.Error(err),
				zap.Int("retries_left", retries),
				zap.Int("sleep_time", sleepTime),
				zap.Any("request_input", chatGptInput),
			)
			retries -= 1
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			var chatResponse ChatCompletionResponse
			err = json.Unmarshal(respBody, &chatResponse)
			if err != nil || len(chatResponse.Choices) == 0 {
				retries -= 1
				o.logger.Error(
					"Could not parse OpenAI Request. Retying after sleeping.",
					zap.Int("retries_left", retries),
					zap.Int("sleep_time", sleepTime),
					zap.Error(err),
					zap.String("response_body", string(respBody)),
					zap.Any("request_input", chatGptInput),
					zap.Int("chat_choices", len(chatResponse.Choices)),
				)
				time.Sleep(time.Duration(sleepTime) * time.Second)
			} else {
				return &chatResponse, nil
			}
		}
	}

	return nil, fmt.Errorf("Open AI Requests Failed")

}
