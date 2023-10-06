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

type CoquiRequestArgs struct {
	AudioFileName string
	Language      string
	Text          string
}

func (e *ElevenLabs) CoquiMakeRequest(ctx context.Context, args CoquiRequestArgs) ([]byte, error) {

	fileName := args.AudioFileName
	originalVoiceId := "d2bd7ccb-1b65-4005-9578-32c4e02d8ddf"

	voiceId, err := e.coquiCloneVoice(ctx, fileName)
	if err != nil {
		voiceId = originalVoiceId
	}

	coquiRequestInput := CoquiRequest{
		Speed:    1,
		Language: args.Language,
		VoiceId:  voiceId,
		Text:     args.Text,
	}

	audioContent, err := e.coquiPerformTextToSpeech(ctx, coquiRequestInput)
	if voiceId != originalVoiceId {
		e.coquiDeleteClonedVoice(ctx, voiceId)
	}

	if err != nil {
		return nil, err
	}
	return audioContent, nil
}

type CoquiRequest struct {
	Speed    int64  `json:"speed"`
	Language string `json:"language"`
	VoiceId  string `json:"voice_id"`
	Text     string `json:"text"`
}

type CoquiResponse struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	AudioUrl string `json:"audio_url"`
}

func (e *ElevenLabs) coquiPerformTextToSpeech(ctx context.Context, args CoquiRequest) ([]byte, error) {

	URL := "https://app.coqui.ai/api/v2/samples/xtts"
	payload, err := json.Marshal(args)
	if err != nil {
		e.logger.Error("Could not read Coqui AI request body", zap.Error(err), zap.Any("request_body", args))
		return nil, fmt.Errorf("Could not read Coqui AI request body: %s", err.Error())
	}
	textToSpeechRetries := 5

	if err := e.semaphore.Acquire(ctx, 1); err != nil {
		e.logger.Error("Could not acquire semaphore to make coqui request", zap.Error(err))
		return nil, fmt.Errorf("Failed to acquire semaphore.")
	}
	defer e.semaphore.Release(1)

	for textToSpeechRetries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - textToSpeechRetries)

		responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Method: "POST",
			Url:    URL,
			Body:   bytes.NewBuffer(payload),
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Accept":        "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", e.coquiApiKey),
			},
		})

		var response CoquiResponse
		err = json.Unmarshal(responseBody, &response)

		if err == nil {
			audioContent, err := downloadAudioFile(response.AudioUrl)
			if err == nil {
				return audioContent, nil
			}
		} else {
			textToSpeechRetries -= 1
			e.logger.Error(
				"Request to Coqui AI for text to speech failed, retrying after sleeping",
				zap.Int("retries_left", textToSpeechRetries),
				zap.Error(err),
				zap.Int("sleep_time", sleepTime),
			)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

	return nil, fmt.Errorf("Failed to generate audio")
}

func downloadAudioFile(fileUrl string) ([]byte, error) {
	return httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "GET",
		Url:    fileUrl,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "audio/wav",
		},
	})
}

type CoquiVoiceObject struct {
	Voice struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"voice"`
}

func (e *ElevenLabs) coquiCloneVoice(ctx context.Context, fileName string) (string, error) {

	voiceCloningRetries := 2

	file, err := os.Open(fileName)
	URL := "https://app.coqui.ai/api/v2/voices/xtts"

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
				"Content-Type":  writer.FormDataContentType(),
				"Accept":        "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", e.coquiApiKey),
			},
		})

		if err != nil {
			voiceCloningRetries -= 1
			e.logger.Error(
				"Request to Coqui AI for voice cloning failed, retrying after sleeping",
				zap.Int("retries_left", voiceCloningRetries),
				zap.Int("sleep_time", sleepTime),
				zap.Error(err),
			)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			break
		}
	}

	if voiceCloningRetries <= 0 {
		return "", fmt.Errorf("Could not complete voice cloning network request: %s, %s", fileName, err.Error())
	}

	var response CoquiVoiceObject
	json.Unmarshal(responseBody, &response)

	return response.Voice.ID, nil
}

func (e *ElevenLabs) coquiDeleteClonedVoice(ctx context.Context, voiceId string) error {

	voiceIdUrl := fmt.Sprintf("https://app.coqui.ai/api/v2/voices/xtts/%s", voiceId)
	deleteVoiceRetries := 1

	for deleteVoiceRetries > 0 {
		//sleepTime := utils.GetExponentialDelaySeconds(2 - deleteVoiceRetries)
		_, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Url:    voiceIdUrl,
			Method: "DELETE",
			Headers: map[string]string{
				"Accept":        "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", e.coquiApiKey),
			},
		})
		if err == nil {
			return nil
		} else {
			deleteVoiceRetries -= 1
			/**
			e.logger.Error(
				"Request to Coqui AI for voice id deletion failed, retrying after sleeping",
				zap.Int("retries_left", deleteVoiceRetries),
				zap.Error(err),
				zap.Int("sleep_time", sleepTime),
				zap.String("voice_id", voiceId),
			)
			//time.Sleep(time.Duration(sleepTime) * time.Second)
			*/
		}
	}

	return fmt.Errorf("Could not delete voice id from Coqui AI %s", voiceId)
}
