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
	"net/http"
	"os"
	"os/exec"
	"planetcastdev/database"
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

func getTranscript(fileNameIdentifier string, file io.ReadSeeker) WhisperOutput {

	file.Seek(0, io.SeekStart)

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	MODEL := "whisper-1"
	RESPONSE_FORMAT := "verbose_json"
	URL := "https://api.openai.com/v1/audio/transcriptions"

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", fileNameIdentifier+".mp4")

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

	req, err := http.NewRequest("POST", URL, requestBody)
	if err != nil {
		log.Println("Failed to create request:", err)
		return WhisperOutput{}
	}

	req.Header.Add("Authorization", "Bearer "+API_KEY)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request failed:", err)
		return WhisperOutput{}
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return WhisperOutput{}
	}

	file.Seek(0, io.SeekStart)
	var whisperOutput WhisperOutput
	json.Unmarshal(responseBody, &whisperOutput)
	log.Println("Whisper request processes successfully for:", fileNameIdentifier)

	return whisperOutput
}

type CreateTransformationParams struct {
	ProjectID          int64
	TargetLanguage     database.SupportedLanguage
	FileNameIdentifier string
	File               io.ReadSeeker
	IsSource           bool
}

func CreateTransformation(
	ctx context.Context,
	queries *database.Queries,
	args CreateTransformationParams,
) (database.Transformation, error) {

	transcript := getTranscript(args.FileNameIdentifier, args.File)
	jsonBytes, err := json.Marshal(transcript)

	transformation, err := queries.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      args.ProjectID,
		TargetLanguage: args.TargetLanguage,
		TargetMedia:    args.FileNameIdentifier + ".mp4",
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

	currentTime := time.Now()
	timeString := strings.ReplaceAll(currentTime.Format("2006-01-02 15:04:05"), " ", "-")
	identifier := fmt.Sprintf("%d-%s-%s", sourceTransformationObject.ProjectID, sourceTransformationObject.TargetLanguage, timeString)

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

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "audio/mpeg")
	headers.Set("xi-api-key", "")
	url := "https://api.elevenlabs.io/v1/text-to-speech/DgzXv5iB8NJnHCwRTzL8"

	for idx, s := range segments {

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

		response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("Error sending request to ElevenLabs: %s", err.Error())
		}

		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			audioContent, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("Error reading response from ElevenLabs: %s", err.Error())
			}

			err = ioutil.WriteFile(fmt.Sprintf("%s-%d-audio-file.mp3", identifier, s.Id), audioContent, 0644)
			if err != nil {
				return fmt.Errorf("Error writing audio file: %s", err.Error())
			}

			log.Println("Audio file saved successfully:", idx+1, "out of", len(segments))
		} else {
			return fmt.Errorf("Request Failed: %d", response.StatusCode)
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

			systemPrompt := fmt.Sprintf("Translate the following text from %s to %s. Just give me the output.", sourceTransformationObject.TargetLanguage, targetLanguage)
			chatGptInput := ChatRequestInput{
				Model: "gpt-3.5-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "system", Content: systemPrompt},
					{Role: "user", Content: segment.Text},
				},
			}

			jsonData, err := json.Marshal(chatGptInput)
			if err != nil {
				return []Segment{}, fmt.Errorf("Could not generate request body: " + err.Error())
			}

			req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
			if err != nil {
				return []Segment{}, fmt.Errorf("Could not generate POST request: " + err.Error())
			}

			req.Header.Add("Authorization", "Bearer "+API_KEY)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return []Segment{}, fmt.Errorf("Could not process request successfully: " + err.Error())
			}

			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return []Segment{}, fmt.Errorf("Could not parse request body successfully: " + err.Error())
			}

			var chatResponse ChatCompletionResponse
			json.Unmarshal(respBody, &chatResponse)

			time.Sleep(2 * time.Second)

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
