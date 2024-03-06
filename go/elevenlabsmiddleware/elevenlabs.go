package elevenlabsmiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"planetcastdev/httpmiddleware"
	"planetcastdev/utils"
	"time"

	"go.uber.org/zap"
)

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
}

type VoiceRequest struct {
	Text         string        `json:"text"`
	ModelID      string        `json:"model_id"`
	VoiceSetting VoiceSettings `json:"voice_settings"`
}

type ElevenLabsRequestArgs struct {
	AudioFileName string
	Text          string
	Gender        string
}

func (e *ElevenLabs) ElevenLabsMakeRequest(ctx context.Context, args ElevenLabsRequestArgs) ([]byte, error) {

	originalVoiceId := "XMQab44ShF40jzdHBoXu" //Fallback voice id
	if args.Gender == "female" {
		originalVoiceId = "21m00Tcm4TlvDq8ikWAM"
	}

	voiceId := originalVoiceId

	/**
		fileName := args.AudioFileName
		voiceId, err := e.cloneVoice(ctx, fileName)
		if err != nil {
			voiceId = originalVoiceId
		}
	  **/

	data := VoiceRequest{
		Text:    args.Text,
		ModelID: "eleven_multilingual_v2",
		VoiceSetting: VoiceSettings{
			Stability:       1.0,
			SimilarityBoost: 1.0,
		},
	}

	audioContent, err := e.elevenLabsPerformTextToSpeech(ctx, data, voiceId)
	/**
		if voiceId != originalVoiceId {
			e.elevenLabsDeleteClonedVoice(voiceId)
		}
	  **/

	if err != nil {
		return nil, err
	}
	return audioContent, nil
}

func (e *ElevenLabs) elevenLabsPerformTextToSpeech(ctx context.Context, args VoiceRequest, voiceId string) ([]byte, error) {

	URL := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceId)
	payload, err := json.Marshal(args)
	if err != nil {
		e.logger.Error("Could not read ElevenLabs request body", zap.Error(err), zap.Any("request_body", args))
		return nil, fmt.Errorf("Could not read ElevenLabs request body: %s", err.Error())
	}
	textToSpeechRetries := 5

	if err := e.semaphore.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("Failed to acquire semaphore.")
	}
	defer e.semaphore.Release(1)

	for textToSpeechRetries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - textToSpeechRetries)

		audioContent, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Method: "POST",
			Url:    URL,
			Body:   bytes.NewBuffer(payload),
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "audio/mpeg",
				"xi-api-key":   e.apiKey,
			},
		})

		if err == nil {
			return audioContent, nil
		} else {
			textToSpeechRetries -= 1
			e.logger.Error(
				"Request to Eleven Labs for text to speech failed, retrying after sleeping",
				zap.Int("retries_left", textToSpeechRetries),
				zap.Error(err),
				zap.Int("sleep_time", sleepTime),
			)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}

	}

	return nil, fmt.Errorf("Failed to generate audio")
}

type ElevenLabsVoiceCloneResponse struct {
	VoiceId string `json:"voice_id"`
}

func (e *ElevenLabs) cloneVoice(ctx context.Context, fileName string) (string, error) {

	voiceCloningRetries := 2

	file, err := os.Open(fileName)
	URL := "https://api.elevenlabs.io/v1/voices/add"

	if err != nil {
		return "", fmt.Errorf("Could not open file for voice cloning: %s, %s", fileName, err.Error())
	}
	defer file.Close()

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	filePart, err := writer.CreateFormFile("files", fileName)
	if err != nil {
		return "", fmt.Errorf("Could not create form field for voice cloning: %s, %s", fileName, err.Error())
	}

	_, err = io.Copy(filePart, file)

	if err != nil {
		return "", fmt.Errorf("Could not copy file into request body for voice cloning: %s, %s", fileName, err.Error())
	}

	writer.WriteField("name", fileName)

	err = writer.Close()

	if err != nil {
		return "", fmt.Errorf("Could write request for voice cloning: %s, %s", fileName, err.Error())
	}

	var responseBody []byte

	for voiceCloningRetries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(2 - voiceCloningRetries)

		responseBody, err = httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Method: "POST",
			Url:    URL,
			Body:   requestBody,
			Headers: map[string]string{
				"Content-Type": writer.FormDataContentType(),
				"Accept":       "application/json",
				"xi-api-key":   e.apiKey,
			},
		})

		if err != nil {
			voiceCloningRetries -= 1
			e.logger.Error(
				"Request to Eleven Labs for voice cloning failed, retrying after sleeping",
				zap.Int("retries_left", voiceCloningRetries),
				zap.Error(err),
				zap.Int("sleep_time", sleepTime),
			)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			break
		}
	}

	if voiceCloningRetries <= 0 {
		return "", fmt.Errorf("Could not complete voice cloning network request: %s, %s", fileName, err.Error())
	}

	var response ElevenLabsVoiceCloneResponse
	json.Unmarshal(responseBody, &response)

	return response.VoiceId, nil
}

func (e *ElevenLabs) elevenLabsDeleteClonedVoice(voiceId string) error {

	voiceIdUrl := fmt.Sprintf("https://api.elevenlabs.io/v1/voices/%s", voiceId)
	deleteVoiceRetries := 2

	for deleteVoiceRetries > 0 {
		sleepTime := utils.GetExponentialDelaySeconds(2 - deleteVoiceRetries)
		_, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Url:    voiceIdUrl,
			Method: "DELETE",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"xi-api-key":   e.apiKey,
			},
		})
		if err == nil {
			return nil
		} else {
			deleteVoiceRetries -= 1
			e.logger.Error(
				"Request to Eleven Labs for voice id deletion failed, retrying after sleeping",
				zap.Int("retries_left", deleteVoiceRetries),
				zap.Error(err),
				zap.Int("sleep_time", sleepTime),
				zap.String("voice_id", voiceId),
			)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

	return fmt.Errorf("Could not delete voice id from eleven labs %s", voiceId)
}
