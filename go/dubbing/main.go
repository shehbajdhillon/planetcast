package dubbing

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"planetcastdev/database"

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

	log.Println("Making request:", req.URL.String())
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
	log.Println(whisperOutput.Segments)

	return whisperOutput
}

func CreateTransformation(ctx context.Context, projectId int64, targetLanguage database.SupportedLanguage, fileNameIdentifier string, file io.ReadSeeker, queries *database.Queries, isSource bool) (database.Transformation, error) {

	transcript := getTranscript(fileNameIdentifier, file)
	jsonBytes, err := json.Marshal(transcript)

	transformation, err := queries.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      projectId,
		TargetLanguage: targetLanguage,
		TargetMedia:    fileNameIdentifier + ".mp4",
		Transcript:     pqtype.NullRawMessage{RawMessage: jsonBytes, Valid: true},
		IsSource:       isSource,
	})

	if err != nil {
		log.Println("Error occured:", err.Error())
		return database.Transformation{}, err
	}

	return transformation, nil
}
