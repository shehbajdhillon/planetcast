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

	// return the update transformation
	return targetTransformationObject, nil
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
