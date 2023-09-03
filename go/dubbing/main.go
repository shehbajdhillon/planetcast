package dubbing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"planetcastdev/database"
	"planetcastdev/httpmiddleware"
	"planetcastdev/utils"
	"strings"
	"time"

	"github.com/tabbed/pqtype"
)

type WhisperOutput struct {
	Language string    `json:"language"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Id    int64   `json:"id"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

func getTranscript(fileName string, file io.ReadSeeker) WhisperOutput {

	file.Seek(0, io.SeekStart)

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	MODEL := "whisper-1"
	RESPONSE_FORMAT := "verbose_json"
	URL := "https://api.openai.com/v1/audio/transcriptions"

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", fileName)

	if err != nil {
		log.Println("Failed to create form file:", err)
		return WhisperOutput{}
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Println("Failed to copy file data:", err)
		return WhisperOutput{}
	}
	writer.WriteField("model", MODEL)
	writer.WriteField("response_format", RESPONSE_FORMAT)
	err = writer.Close()

	if err != nil {
		log.Println("Failed to close writer", err)
		return WhisperOutput{}
	}

	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "POST",
		Url:    URL,
		Body:   requestBody,
		Headers: map[string]string{
			"Authorization": "Bearer " + API_KEY,
			"Content-Type":  writer.FormDataContentType(),
		},
	})

	if err != nil {
		log.Println("Request failed:", err)
		return WhisperOutput{}
	}

	file.Seek(0, io.SeekStart)
	var whisperOutput WhisperOutput
	json.Unmarshal(responseBody, &whisperOutput)
	log.Println("Whisper request processes successfully for:", fileName)

	return whisperOutput
}

type CreateTransformationParams struct {
	ProjectID      int64
	TargetLanguage database.SupportedLanguage
	FileName       string
	File           io.ReadSeeker
	IsSource       bool
}

func CreateTransformation(
	ctx context.Context,
	queries *database.Queries,
	args CreateTransformationParams,
) (database.Transformation, error) {

	transcript := getTranscript(args.FileName, args.File)
	jsonBytes, err := json.Marshal(transcript)

	transformation, err := queries.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      args.ProjectID,
		TargetLanguage: args.TargetLanguage,
		TargetMedia:    args.FileName,
		Transcript:     pqtype.NullRawMessage{RawMessage: jsonBytes, Valid: true},
		IsSource:       args.IsSource,
	})

	if err != nil {
		log.Println("Error occured:", err.Error())
		return database.Transformation{}, err
	}

	return transformation, nil
}

func CreateTranslation(
	ctx context.Context,
	queries *database.Queries,
	sourceTransformationObject database.Transformation,
	targetTransformationObject database.Transformation,
) (database.Transformation, error) {

	// call chatgpt, convert the source text to target text
	translatedSegments, _ := translateResponse(
		ctx,
		sourceTransformationObject,
		targetTransformationObject.TargetLanguage,
	)

	// get the target text, and parse it
	var whisperOutput WhisperOutput
	json.Unmarshal(targetTransformationObject.Transcript.RawMessage, &whisperOutput)
	whisperOutput.Segments = translatedSegments
	whisperOutput.Language = strings.ToLower(string(targetTransformationObject.TargetLanguage))

	// store the target text in db
	jsonBytes, err := json.Marshal(whisperOutput)

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not process translated segments " + err.Error())
	}

	targetTransformationObject, err = queries.UpdateTranscriptById(ctx, database.UpdateTranscriptByIdParams{
		ID:         targetTransformationObject.ID,
		Transcript: pqtype.NullRawMessage{Valid: true, RawMessage: jsonBytes},
	})

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not update transformation: " + err.Error())
	}

	identifier := fmt.Sprintf("%d-%s-%s", sourceTransformationObject.ProjectID, utils.GetCurrentDateTimeString(), targetTransformationObject.TargetLanguage)

	err = fetchDubbedClips(translatedSegments, identifier)

	// return the update transformation
	return targetTransformationObject, nil
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

func fetchDubbedClips(segments []Segment, identifier string) error {

	url := "https://api.elevenlabs.io/v1/text-to-speech/2EiwWnXFnvU5JabPnv8n"
	API_KEY := os.Getenv("ELEVEN_LABS_KEY")

	for idx, s := range segments {

		for {

			id := s.Id
			audioFileName := fmt.Sprintf("%s_%d_audio_file.mp3", identifier, id)

			data := VoiceRequest{
				Text:    s.Text,
				ModelID: "eleven_multilingual_v1",
				VoiceSetting: VoiceSettings{
					Stability:       1.0,
					SimilarityBoost: 1.0,
				},
			}

			payload, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("Error encoding JSON for ElevenLabs request: %s", err.Error())
			}

			responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
				Method: "POST",
				Url:    url,
				Body:   bytes.NewBuffer(payload),
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Accept":       "audio/mpeg",
					"xi-api-key":   API_KEY,
				},
			})

			if err == nil {

				audioContent := responseBody
				if err != nil {
					return fmt.Errorf("Error reading response from ElevenLabs: %s", err.Error())
				}

				err = ioutil.WriteFile(audioFileName, audioContent, 0644)
				if err != nil {
					return fmt.Errorf("Error writing audio file: %s", err.Error())
				}

				log.Println("Audio file saved successfully:", idx+1, "out of", len(segments))
				time.Sleep(500 * time.Millisecond)
				break

			} else {
				log.Println("Request Failed: " + err.Error())
				time.Sleep(5 * time.Second)
			}

		}
	}

	return nil
}

func dubAudioClips(segments []Segment, identifier string) error {

	for _, s := range segments {
		id := s.Id

		audioFileName := fmt.Sprintf("%s-%d-audio-file.mp3", identifier, id)
		videoSegmentName := fmt.Sprintf("%s-%d-video-segment.mp4", identifier, id)
		dubbedVideoSegmentName := fmt.Sprintf("%s-%d-video-segment-dubbed.mp4", identifier, id)

		generateVideoClipCmd := exec.Command("ffmpeg", "-i", "iio_video_2.mp4", "-ss", fmt.Sprintf("%.2f", s.Start), "-to", fmt.Sprintf("%.2f", s.End), videoSegmentName)
		if err := generateVideoClipCmd.Run(); err != nil {
			return fmt.Errorf("Could not clip source video: %s", err.Error())
		}

		dubVideoClipCmd := exec.Command("ffmpeg", "-i", videoSegmentName, "-i", audioFileName, "-c:v", "copy", "-map", "0:v:0", "-map", "1:a:0", dubbedVideoSegmentName)
		if err := dubVideoClipCmd.Run(); err != nil {
			return fmt.Errorf("Could not dub the clip: %s", err.Error())
		}

	}

	return nil

}

func concatSegments(segments []Segment, identifier string) {
}

func cleanUp(segments []Segment, identifier string) {
}

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

func translateResponse(
	ctx context.Context,
	sourceTransformationObject database.Transformation,
	targetLanguage database.SupportedLanguage,
) ([]Segment, error) {

	var whisperOutput WhisperOutput
	json.Unmarshal(sourceTransformationObject.Transcript.RawMessage, &whisperOutput)

	segments := whisperOutput.Segments

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	URL := "https://api.openai.com/v1/chat/completions"

	log.Println("Translating project", sourceTransformationObject.ProjectID, "from", sourceTransformationObject.TargetLanguage, "to", targetLanguage)

	for idx, segment := range segments {

		retries := 5

		for retries > 0 {

			systemPrompt := fmt.Sprintf("Translate the following text to street spoken, informal %s. Provide the output in %s alphabet. Just give the output.", targetLanguage, targetLanguage)
			chatGptInput := ChatRequestInput{
				Model: "gpt-4",
				Messages: []ChatCompletionMessage{
					{Role: "system", Content: systemPrompt},
					{Role: "user", Content: segment.Text},
				},
			}

			jsonData, err := json.Marshal(chatGptInput)
			if err != nil {
				return []Segment{}, fmt.Errorf("Could not generate request body: " + err.Error())
			}

			respBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
				Method: "POST",
				Url:    URL,
				Body:   bytes.NewBuffer(jsonData),
				Headers: map[string]string{
					"Authorization": "Bearer " + API_KEY,
					"Content-Type":  "application/json",
				},
			})

			if err != nil {
				return []Segment{}, fmt.Errorf("Could not parse request body successfully: " + err.Error())
			}

			var chatResponse ChatCompletionResponse
			json.Unmarshal(respBody, &chatResponse)

			time.Sleep(500 * time.Millisecond)

			if len(chatResponse.Choices) == 0 {
				retries -= 1
				log.Println("Open AI request failed, sleeping for 5s, retries left:", retries)
				time.Sleep(5 * time.Second)
			} else {
				segment.Text = chatResponse.Choices[0].Message.Content
				segments[idx] = segment
				log.Println("Translation Progress for Project", sourceTransformationObject.ProjectID, "from", sourceTransformationObject.TargetLanguage, "to", targetLanguage+":", idx+1, "/", len(segments))
				break
			}
		}
	}
	log.Println(segments)

	return segments, nil

}
