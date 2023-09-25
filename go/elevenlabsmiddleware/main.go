package elevenlabsmiddleware

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

type ElevenLabs struct {
	logger    *zap.Logger
	semaphore *semaphore.Weighted
}

type ElevenLabsConnectProps struct {
	Logger *zap.Logger
}

func Connect(args ElevenLabsConnectProps) *ElevenLabs {
	maxWorkers := 5
	sem := semaphore.NewWeighted(int64(maxWorkers))
	return &ElevenLabs{logger: args.Logger, semaphore: sem}
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
}

type VoiceRequest struct {
	Text         string        `json:"text"`
	ModelID      string        `json:"model_id"`
	VoiceSetting VoiceSettings `json:"voice_settings"`
}

type MakeRequestProps struct {
	Retries      int
	RequestInput VoiceRequest
}

func (e *ElevenLabs) MakeRequest(ctx context.Context, args MakeRequestProps) ([]byte, error) {

	retries := args.Retries

	URL := "https://api.elevenlabs.io/v1/text-to-speech/rU18Fk3uSDhmg5Xh41o4"
	API_KEY := os.Getenv("ELEVEN_LABS_KEY")
	data := args.RequestInput

	for retries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - retries)

		payload, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("Error encoding JSON for ElevenLabs request: %s", err.Error())
		}

		if err := e.semaphore.Acquire(ctx, 1); err != nil {
			defer e.semaphore.Release(1)
			return nil, fmt.Errorf("Failed to acquire semaphore.")
		}
		audioContent, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Method: "POST",
			Url:    URL,
			Body:   bytes.NewBuffer(payload),
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "audio/mpeg",
				"xi-api-key":   API_KEY,
			},
		})
		time.Sleep(1 * time.Second)

		e.semaphore.Release(1)

		if err == nil {
			return audioContent, nil
		} else {
			retries -= 1
			e.logger.Error("Request to Eleven Labs failed, retrying after sleeping", zap.Int("retries_left", retries), zap.Error(err), zap.Int("sleep_time", sleepTime))
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}

	}

	return nil, fmt.Errorf("Failed to generate audio")

}
