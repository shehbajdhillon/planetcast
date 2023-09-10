package dubbing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"os"
	"os/exec"
	"planetcastdev/database"
	"planetcastdev/graph/model"
	"planetcastdev/httpmiddleware"
	"planetcastdev/replicatemiddleware"
	"planetcastdev/storage"
	"planetcastdev/utils"
	"strings"
	"time"

	"github.com/tabbed/pqtype"
)

func getAudioFileName(identifier string, id int64) string {
	audioFileName := fmt.Sprintf("%s_%d_audio_file.mp3", identifier, id)
	return audioFileName
}

func getVideoSegmentName(identifier string, id int64) string {
	videoSegmentName := fmt.Sprintf("%s_%d_video_segment.mp4", identifier, id)
	return videoSegmentName
}

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

func getTranscript(fileName string, file io.ReadSeeker) (*WhisperOutput, error) {

	file.Seek(0, io.SeekStart)

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	MODEL := "whisper-1"
	RESPONSE_FORMAT := "verbose_json"
	URL := "https://api.openai.com/v1/audio/transcriptions"

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", fileName)

	if err != nil {
		return nil, fmt.Errorf("Failed to create form file: %s", err.Error())
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("Failed to copy file data: %s", err.Error())
	}
	writer.WriteField("model", MODEL)
	writer.WriteField("response_format", RESPONSE_FORMAT)
	err = writer.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close writer %s", err.Error())
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
		return nil, fmt.Errorf("Request failed:" + err.Error())
	}

	file.Seek(0, io.SeekStart)
	var whisperOutput WhisperOutput
	json.Unmarshal(responseBody, &whisperOutput)
	log.Println("Whisper request processes successfully for:", fileName)

	return &whisperOutput, nil
}

type CreateTransformationParams struct {
	ProjectID      int64
	TargetLanguage model.SupportedLanguage
	FileName       string
	File           io.ReadSeeker
	IsSource       bool
}

func CreateTransformation(
	ctx context.Context,
	queries *database.Queries,
	args CreateTransformationParams,
) (database.Transformation, error) {

	transcriptPtr, _ := getTranscript(args.FileName, args.File)
	transcriptObj := *transcriptPtr
	jsonBytes, err := json.Marshal(transcriptObj)

	transformation, err := queries.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      args.ProjectID,
		TargetLanguage: args.TargetLanguage.String(),
		TargetMedia:    args.FileName,
		Transcript:     pqtype.NullRawMessage{RawMessage: jsonBytes, Valid: true},
		IsSource:       args.IsSource,
		Status:         "complete",
		Progress:       100,
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
	sourceTransformation database.Transformation,
	targetTransformation database.Transformation,
	identifier string,
) (database.Transformation, error) {

	fileUrl := storage.Connect().GetFileLink(sourceTransformation.TargetMedia)
	//download original media, then save it as identifier.mp4
	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "GET",
		Url:    fileUrl,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "audio/mp4",
		},
	})
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Error downloading original audio file from S3: %s", err.Error())
	}
	err = ioutil.WriteFile(identifier+".mp4", responseBody, 0644)
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Error writing audio file: %s", err.Error())
	}

	var whisperOutput WhisperOutput
	json.Unmarshal(sourceTransformation.Transcript.RawMessage, &whisperOutput)

	// call chatgpt, convert the source text to target text
	sourceSegments := whisperOutput.Segments
	translatedSegmentsPtr, err := fetchAndDub(ctx, sourceSegments, sourceTransformation.ProjectID, identifier, model.SupportedLanguage(targetTransformation.TargetLanguage), targetTransformation.ID, queries)
	translatedSegments := *translatedSegmentsPtr

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not process translated segments " + err.Error())
	}

	newFileName, err := concatSegments(translatedSegments, identifier)
	file, err := os.Open(newFileName)
	if err != nil {
		fmt.Println("Error opening the file:", err)
	}
	defer file.Close()
	storage.Connect().Upload(newFileName, file)
	utils.DeleteFiles([]string{identifier + ".mp4", identifier + "_dubbed.mp4"})

	// get the target text, and parse it
	json.Unmarshal(targetTransformation.Transcript.RawMessage, &whisperOutput)
	whisperOutput.Segments = translatedSegments
	whisperOutput.Language = strings.ToLower(string(targetTransformation.TargetLanguage))

	// store the target text in db
	jsonBytes, err := json.Marshal(whisperOutput)

	targetTransformation, err = queries.UpdateTranscriptById(ctx, database.UpdateTranscriptByIdParams{
		ID:         targetTransformation.ID,
		Transcript: pqtype.NullRawMessage{Valid: true, RawMessage: jsonBytes},
	})

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not update transformation: " + err.Error())
	}

	targetTransformation, err = queries.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     targetTransformation.ID,
		Status: "complete",
	})

	targetTransformation, err = queries.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
		ID:       targetTransformation.ID,
		Progress: 100,
	})

	// return the update transformation
	return targetTransformation, nil
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

func fetchAndDub(
	ctx context.Context,
	segments []Segment,
	projectId int64,
	identifier string,
	targetLanguage model.SupportedLanguage,
	targetTransformationId int64,
	queries *database.Queries,
) (*[]Segment, error) {

	translatedSegments := []Segment{}

	queries.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     targetTransformationId,
		Status: "processing",
	})

	for idx, segment := range segments {

		translatedSegment, err := translateSegment(ctx, segment, targetLanguage)
		if err != nil {
			return nil, fmt.Errorf("Failed to translated segment %d/%d: %s", idx+1, len(segments), err.Error())
		}
		log.Println("Translation Progress for Project", projectId, "from", targetLanguage, "to", targetLanguage+":", idx+1, "/", len(segments))

		err = fetchDubbedClip(*translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could fetch dubbed clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		log.Println("Fetched dubbed audio file for Project", projectId, ":", idx+1, "/", len(segments))

		err = dubVideoClip(*translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could not process clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		log.Println("Dubbed video clip for Project", projectId, ":", idx+1, "/", len(segments))

		err = lipSyncClip(*translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could not lip sync clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		log.Println("Synced video clip for Project", projectId, ":", idx+1, "/", len(segments))

		translatedSegments = append(translatedSegments, *translatedSegment)

		percentage := 100 * (float64(len(translatedSegments)) / float64(len(segments)))
		percentage = math.Round(percentage*100) / 100

		queries.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
			ID:       targetTransformationId,
			Progress: percentage,
		})

	}

	return &translatedSegments, nil
}

func fetchDubbedClip(segment Segment, identifier string) error {

	url := "https://api.elevenlabs.io/v1/text-to-speech/XMQab44ShF40jzdHBoXu"
	API_KEY := os.Getenv("ELEVEN_LABS_KEY")

	retries := 5

	for retries > 0 {

		id := segment.Id
		audioFileName := getAudioFileName(identifier, id)

		data := VoiceRequest{
			Text:    segment.Text,
			ModelID: "eleven_multilingual_v2",
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

			return nil

		} else {
			retries -= 1
			log.Printf("Request Failed, retrying after 5 seconds, retries left %d: %s\n", retries, err.Error())
			time.Sleep(5 * time.Second)
		}

	}

	return fmt.Errorf("Failed to call elevenlabs")
}

func dubVideoClip(segment Segment, identifier string) error {

	id := segment.Id

	audioFileName := getAudioFileName(identifier, id)
	videoSegmentName := getVideoSegmentName(identifier, id)
	originalVideoSegmentName := "original_" + videoSegmentName
	dubbedVideoSegmentName := "dubbed_" + videoSegmentName

	start := segment.Start
	end := segment.End

	originalLength := end - start

	audioFileDuarion, _ := utils.GetAudioFileDuration(audioFileName)

	ratio := math.Max(audioFileDuarion/originalLength, 1.0)

	generateVideoClip := fmt.Sprintf("ffmpeg -i file:'%s.mp4' -ss %f -to %f file:'%s'", identifier, start, end, originalVideoSegmentName)
	err := exec.Command("sh", "-c", generateVideoClip).Run()

	if err != nil {
		return fmt.Errorf("Clip extraction failed: %s\n%s\n", err.Error(), generateVideoClip)
	}

	stretchVideoClip := fmt.Sprintf("ffmpeg -i file:'%s' -vf 'setpts=%f*PTS' -c:a copy file:'%s'", originalVideoSegmentName, ratio, videoSegmentName)
	err = exec.Command("sh", "-c", stretchVideoClip).Run()

	if err != nil {
		return fmt.Errorf("Clip stretching failed: %s\n%s\n", err.Error(), generateVideoClip)
	}

	dubVideoClip := fmt.Sprintf("ffmpeg -i file:'%s' -i file:'%s' -c:v copy -map 0:v:0 -map 1:a:0 file:'%s'",
		videoSegmentName, audioFileName, dubbedVideoSegmentName)
	err = exec.Command("sh", "-c", dubVideoClip).Run()

	if err != nil {
		return fmt.Errorf("Clip dubbing failed: %s\n%s\n", err.Error(), generateVideoClip)
	}

	utils.DeleteFiles([]string{videoSegmentName, originalVideoSegmentName, audioFileName})

	return nil
}

func lipSyncClip(segment Segment, identifier string) error {

	storageClient := storage.Connect()

	videoSegmentName := getVideoSegmentName(identifier, segment.Id)
	dubbedVideoSegmentName := "dubbed_" + videoSegmentName
	syncedVideoSegmentName := "synced_" + videoSegmentName

	file, err := os.Open(dubbedVideoSegmentName)
	if err != nil {
		fmt.Println("Error opening the file:", err)
	}
	defer file.Close()

	storageClient.Upload(dubbedVideoSegmentName, file)
	fileLink := storageClient.GetFileLink(dubbedVideoSegmentName)

	replicateRequestBody := map[string]interface{}{
		"version": "8d65e3f4f4298520e079198b493c25adfc43c058ffec924f2aefc8010ed25eef",
		"input": map[string]string{
			"face":  fileLink,
			"audio": fileLink,
		},
	}

	jsonBody, err := json.Marshal(replicateRequestBody)
	outputUrl, err := replicatemiddleware.MakeRequest(bytes.NewBuffer(jsonBody))
	storageClient.DeleteFile(dubbedVideoSegmentName)

	if err != nil {
		copyFileCmd := fmt.Sprintf("cp %s %s", dubbedVideoSegmentName, syncedVideoSegmentName)
		exec.Command("sh", "-c", copyFileCmd).Run()
		utils.DeleteFiles([]string{dubbedVideoSegmentName})
		return nil
	}

	//download original media, then save it as identifier.mp4
	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "GET",
		Url:    outputUrl,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "audio/mp4",
		},
	})

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(syncedVideoSegmentName, responseBody, 0644)

	if err != nil {
		return err
	}

	utils.DeleteFiles([]string{dubbedVideoSegmentName})

	return nil

}

func concatSegments(segments []Segment, identifier string) (string, error) {

	inputList := []string{}
	filterList := []string{}

	for idx, s := range segments {
		id := s.Id

		videoSegmentName := getVideoSegmentName(identifier, id)
		syncedSegmentName := "synced_" + videoSegmentName
		inputList = append(inputList, fmt.Sprintf("-i file:'%s'", syncedSegmentName))
		filterList = append(filterList, fmt.Sprintf("[%d:v][%d:a]", idx, idx))
	}

	filterList = append(filterList, fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", len(segments)))

	inputArgs := strings.Join(inputList, " ")
	filterComplex := strings.Join(filterList, "")

	ffmpegCmd := fmt.Sprintf("ffmpeg %s -filter_complex '%s' -map '[v]' -map '[a]' -vsync 2 file:'%s_dubbed.mp4'",
		inputArgs, filterComplex, identifier)

	log.Println("Concating segments:", ffmpegCmd)

	err := exec.Command("sh", "-c", ffmpegCmd).Run()

	fileList := []string{}
	for _, s := range segments {
		fileName := "synced_" + getVideoSegmentName(identifier, s.Id)
		fileList = append(fileList, fileName)
	}
	utils.DeleteFiles(fileList)

	if err != nil {
		return "", fmt.Errorf("Could not concat segments: %s", ffmpegCmd)
	}

	return identifier + "_dubbed.mp4", nil
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

func translateSegment(ctx context.Context, segment Segment, targetLang model.SupportedLanguage) (*Segment, error) {

	retries := 5

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	URL := "https://api.openai.com/v1/chat/completions"

	for retries > 0 {

		systemPrompt := fmt.Sprintf("Translate the following text to street spoken, informal %s. Provide the output in %s alphabet. Just give the output.", targetLang, targetLang)
		chatGptInput := ChatRequestInput{
			Model: "gpt-4",
			Messages: []ChatCompletionMessage{
				{Role: "system", Content: systemPrompt},
				{Role: "user", Content: segment.Text},
			},
		}

		jsonData, err := json.Marshal(chatGptInput)
		if err != nil {
			return nil, fmt.Errorf("Could not generate request body: " + err.Error())
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
			return nil, fmt.Errorf("Could not parse request body successfully: " + err.Error())
		}

		var chatResponse ChatCompletionResponse
		json.Unmarshal(respBody, &chatResponse)

		if len(chatResponse.Choices) == 0 {

			retries -= 1
			log.Println("Open AI request failed, sleeping for 5s, retries left:", retries)
			time.Sleep(5 * time.Second)

		} else {
			segment.Text = chatResponse.Choices[0].Message.Content
			return &segment, nil
		}
	}

	return nil, fmt.Errorf("Open AI Requests Failed")
}
