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
	"planetcastdev/httpmiddleware"
	"planetcastdev/storage"
	"strconv"
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

	transcriptPtr, _ := getTranscript(args.FileName, args.File)
	transcriptObj := *transcriptPtr
	jsonBytes, err := json.Marshal(transcriptObj)

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

	audioContent := responseBody
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Error downloading original audio file from S3: %s", err.Error())
	}

	err = ioutil.WriteFile(identifier+".mp4", audioContent, 0644)
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Error writing audio file: %s", err.Error())
	}

	log.Println("Source Audio File Downloaded")

	// call chatgpt, convert the source text to target text
	translatedSegments, _ := translateResponse(ctx, sourceTransformation, targetTransformation.TargetLanguage)

	// get the target text, and parse it
	var whisperOutput WhisperOutput
	json.Unmarshal(targetTransformation.Transcript.RawMessage, &whisperOutput)
	whisperOutput.Segments = translatedSegments
	whisperOutput.Language = strings.ToLower(string(targetTransformation.TargetLanguage))

	// store the target text in db
	jsonBytes, err := json.Marshal(whisperOutput)

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not process translated segments " + err.Error())
	}

	targetTransformation, err = queries.UpdateTranscriptById(ctx, database.UpdateTranscriptByIdParams{
		ID:         targetTransformation.ID,
		Transcript: pqtype.NullRawMessage{Valid: true, RawMessage: jsonBytes},
	})

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not update transformation: " + err.Error())
	}

	newFileName := identifier + "_dubbed.mp4"

	err = fetchDubbedClips(translatedSegments, identifier)
	err = dubVideoClips(translatedSegments, identifier)
	err = concatSegments(translatedSegments, identifier)

	file, err := os.Open(newFileName)
	if err != nil {
		fmt.Println("Error opening the file:", err)
	}
	defer file.Close()

	storage.Connect().Upload(newFileName, file)

	cleanUp(translatedSegments, identifier)

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

func fetchDubbedClips(segments []Segment, identifier string) error {

	url := "https://api.elevenlabs.io/v1/text-to-speech/2EiwWnXFnvU5JabPnv8n"
	API_KEY := os.Getenv("ELEVEN_LABS_KEY")

	for idx, s := range segments {

		for {

			id := s.Id
			audioFileName := getAudioFileName(identifier, id)

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

func dubVideoClips(segments []Segment, identifier string) error {

	for idx, s := range segments {
		id := s.Id

		audioFileName := getAudioFileName(identifier, id)
		videoSegmentName := getVideoSegmentName(identifier, id)
		dubbedVideoSegmentName := "dubbed_" + videoSegmentName
		originalVideoSegmentName := "original_" + videoSegmentName

		start := s.Start
		end := s.End

		originalLength := end - start

		audioFileDuarion, _ := getAudioFileDuration(audioFileName)

		ratio := math.Max(audioFileDuarion/originalLength, 1.0)

		generateVideoClip := fmt.Sprintf("ffmpeg -i file:'%s.mp4' -ss %f -to %f file:'%s'", identifier, start, end, originalVideoSegmentName)
		err := exec.Command("sh", "-c", generateVideoClip).Run()

		if err != nil {
			log.Printf("Could not extract video clip %d/%d: %s\n%s\n", idx+1, len(segments), err.Error(), generateVideoClip)
		}

		stretchVideoClip := fmt.Sprintf("ffmpeg -i file:'%s' -vf 'setpts=%f*PTS' -c:a copy file:'%s'", originalVideoSegmentName, ratio, videoSegmentName)
		err = exec.Command("sh", "-c", stretchVideoClip).Run()

		if err != nil {
			log.Printf("Could not stretch video clip %d/%d: %s\n%s\n", idx+1, len(segments), err.Error(), stretchVideoClip)
		}

		dubVideoClip := fmt.Sprintf("ffmpeg -i file:'%s' -i file:'%s' -c:v copy -map 0:v:0 -map 1:a:0 file:'%s'",
			videoSegmentName, audioFileName, dubbedVideoSegmentName)
		err = exec.Command("sh", "-c", dubVideoClip).Run()

		if err != nil {
			log.Printf("Could not dub video clip %d/%d: %s\n%s", idx+1, len(segments), err.Error(), dubVideoClip)
		}

		log.Printf("Dubbed video clip %d/%d\n", idx+1, len(segments))

	}

	return nil

}

func concatSegments(segments []Segment, identifier string) error {

	inputList := []string{}
	filterList := []string{}

	for idx, s := range segments {
		id := s.Id

		videoSegmentName := getVideoSegmentName(identifier, id)
		syncedSegmentName := "dubbed_" + videoSegmentName
		inputList = append(inputList, fmt.Sprintf("-i file:'%s'", syncedSegmentName))
		filterList = append(filterList, fmt.Sprintf("[%d:v][%d:a]", idx, idx))
	}

	filterList = append(filterList, fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", len(segments)))

	inputArgs := strings.Join(inputList, " ")
	filterComplex := strings.Join(filterList, "")

	ffmpegCmd := fmt.Sprintf("ffmpeg %s -filter_complex '%s' -map '[v]' -map '[a]' file:'%s_dubbed.mp4'",
		inputArgs, filterComplex, identifier)

	log.Println("Concating segments:", ffmpegCmd)

	err := exec.Command("sh", "-c", ffmpegCmd).Run()

	if err != nil {
		log.Println("Could not concat segments:", ffmpegCmd)
	}

	return err
}

func cleanUp(segments []Segment, identifier string) {
	for _, s := range segments {
		id := s.Id

		audioFileName := getAudioFileName(identifier, id)
		videoSegmentName := getVideoSegmentName(identifier, id)
		dubbedSegmentName := "dubbed_" + videoSegmentName
		syncedSegmentName := "synced_" + videoSegmentName
		originalSegmentName := "original_" + videoSegmentName

		removeCmd :=
			fmt.Sprintf(
				"rm -rf %s %s %s %s %s",
				audioFileName,
				videoSegmentName,
				dubbedSegmentName,
				syncedSegmentName,
				originalSegmentName,
			)
		exec.Command("sh", "-c", removeCmd).Run()
	}

	removeCmd := fmt.Sprintf("rm -rf %s %s", identifier+".mp4", identifier+"_dubbed.mp4")
	exec.Command("sh", "-c", removeCmd).Run()
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

func getAudioFileDuration(fileName string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"format=duration", "-of", "default=noprint_wrappers=1:nokey=1", fileName)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil

}
