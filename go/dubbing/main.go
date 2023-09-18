package dubbing

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"planetcastdev/database"
	"planetcastdev/ffmpegmiddleware"
	"planetcastdev/graph/model"
	"planetcastdev/httpmiddleware"
	"planetcastdev/replicatemiddleware"
	"planetcastdev/storage"
	"planetcastdev/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tabbed/pqtype"
	"go.uber.org/zap"
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

type Dubbing struct {
	storage  *storage.Storage
	database *database.Queries
	logger   *zap.Logger
	ffmpeg   *ffmpegmiddleware.Ffmpeg
}

type DubbingConnectProps struct {
	Storage  *storage.Storage
	Database *database.Queries
	Logger   *zap.Logger
	Ffmpeg   *ffmpegmiddleware.Ffmpeg
}

func Connect(args DubbingConnectProps) *Dubbing {
	return &Dubbing{storage: args.Storage, database: args.Database, logger: args.Logger, ffmpeg: args.Ffmpeg}
}

func (d *Dubbing) getTranscript(fileName string) (*WhisperOutput, error) {

	fileUrl := d.storage.GetFileLink(fileName)

	retries := 5

	var output any

	for retries > 0 {

		replicateRequestBody := map[string]interface{}{
			"version": "4a60104c44dd709fc08a03dfeca6c6906257633dd03fd58663ec896a4eeba30e",
			"input": map[string]interface{}{
				"audio":           fileUrl,
				"model":           "large-v2",
				"word_timestamps": true,
			},
		}
		jsonBody, err := json.Marshal(replicateRequestBody)
		output, err = replicatemiddleware.MakeRequest(bytes.NewBuffer(jsonBody))

		if err == nil {
			break
		} else {
			retries -= 1
			d.logger.Error("Whisper request failed, retrying after 5s", zap.Error(err))
			time.Sleep(5 * time.Second)
		}
	}

	if retries <= 0 {
		d.logger.Error("Failed to transcribe whisper request")
		return nil, nil
	}

	outputJson, ok := output.(map[string]interface{})

	if !ok {
		d.logger.Error("Could not parse whisper json output")
		return nil, fmt.Errorf("Could not parse whisper json output")
	}

	responseBody, err := json.Marshal(outputJson)
	if err != nil {
		d.logger.Error("Could not parse whisper output to bytes")
		return nil, fmt.Errorf("Could not parse whisper json body to bytes")
	}

	var whisperOutput WhisperOutput
	err = json.Unmarshal(responseBody, &whisperOutput)
	if err != nil {
		d.logger.Error("Could not parse whisper bytes to struct")
		return nil, fmt.Errorf("Could not parse whisper bytes to struct")
	}
	d.logger.Info("Whisper request processes successfully for:", zap.String("fileName", fileName))

	return &whisperOutput, nil
}

type CreateTransformationParams struct {
	ProjectID      int64
	TargetLanguage model.SupportedLanguage
	FileName       string
	IsSource       bool
}

func (d *Dubbing) CreateTransformation(
	ctx context.Context,
	args CreateTransformationParams,
) (database.Transformation, error) {

	transcriptPtr, err := d.getTranscript(args.FileName)

	if err != nil {
		d.logger.Error("Failed to generate transcript", zap.Error(err))
		return database.Transformation{}, err
	}

	transcriptObj := *transcriptPtr
	jsonBytes, err := json.Marshal(transcriptObj)

	if err != nil {
		d.logger.Error("Failed to parse json ", zap.Error(err))
		return database.Transformation{}, err
	}

	transformation, err := d.database.CreateTransformation(ctx, database.CreateTransformationParams{
		ProjectID:      args.ProjectID,
		TargetLanguage: args.TargetLanguage.String(),
		TargetMedia:    args.FileName,
		Transcript:     pqtype.NullRawMessage{RawMessage: jsonBytes, Valid: true},
		IsSource:       args.IsSource,
		Status:         "complete",
		Progress:       100,
	})

	if err != nil {
		d.logger.Error("Error occured", zap.Error(err))
		return database.Transformation{}, err
	}

	return transformation, nil
}

func (d *Dubbing) CreateTranslation(
	ctx context.Context,
	sourceTransformation database.Transformation,
	targetTransformation database.Transformation,
	identifier string,
) (database.Transformation, error) {

	fileUrl := d.storage.GetFileLink(sourceTransformation.TargetMedia)
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
	translatedSegmentsPtr, err := d.fetchAndDub(ctx, sourceSegments, sourceTransformation.ProjectID, identifier, model.SupportedLanguage(targetTransformation.TargetLanguage), targetTransformation.ID)
	if err != nil {
		return database.Transformation{}, fmt.Errorf("Error translating and fetching: %s", err.Error())
	}
	translatedSegments := *translatedSegmentsPtr

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not process translated segments " + err.Error())
	}

	newFileName, err := d.concatSegments(ctx, translatedSegments, identifier)
	if err != nil {
		d.logger.Error("Error concatenating segments", zap.Error(err))
	}

	file, err := os.Open(newFileName)
	if err != nil {
		d.logger.Error("Error opening the file", zap.Error(err))
	}

	defer file.Close()
	d.storage.Upload(newFileName, file)
	utils.DeleteFiles([]string{identifier + ".mp4", identifier + "_dubbed.mp4"})

	// get the target text, and parse it
	json.Unmarshal(targetTransformation.Transcript.RawMessage, &whisperOutput)
	whisperOutput.Segments = translatedSegments
	whisperOutput.Language = strings.ToLower(string(targetTransformation.TargetLanguage))

	// store the target text in db
	jsonBytes, err := json.Marshal(whisperOutput)

	targetTransformation, err = d.database.UpdateTranscriptById(ctx, database.UpdateTranscriptByIdParams{
		ID:         targetTransformation.ID,
		Transcript: pqtype.NullRawMessage{Valid: true, RawMessage: jsonBytes},
	})

	if err != nil {
		return database.Transformation{}, fmt.Errorf("Could not update transformation: " + err.Error())
	}

	targetTransformation, err = d.database.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     targetTransformation.ID,
		Status: "complete",
	})

	targetTransformation, err = d.database.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
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

func (d *Dubbing) fetchAndDub(
	ctx context.Context,
	segments []Segment,
	projectId int64,
	identifier string,
	targetLanguage model.SupportedLanguage,
	targetTransformationId int64,
) (*[]Segment, error) {

	translatedSegments := []Segment{}

	d.database.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     targetTransformationId,
		Status: "processing",
	})

	for idx, segment := range segments {

		logProgress := func(stage string) {
			d.logger.Info(
				stage,
				zap.Int("project_id", int(projectId)),
				zap.Int("transformation_id", int(targetTransformationId)),
				zap.String("target_language", string(targetLanguage)),
				zap.Int("segments_processed", idx+1),
				zap.Int("total_segments", len(segments)),
			)
		}

		translatedSegment, err := d.translateSegment(ctx, segment, targetLanguage)
		if err != nil {
			return nil, fmt.Errorf("Failed to translated segment %d/%d: %s", idx+1, len(segments), err.Error())
		}
		logProgress("Translation Progress")

		err = d.fetchDubbedClip(*translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could fetch dubbed clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		logProgress("Audio Generation Progress")

		err = d.dubVideoClip(ctx, *translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could not process clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		logProgress("Dubbing Progress")

		err = d.lipSyncClip(*translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could not lip sync clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		logProgress("Lip Syncing Progress")

		translatedSegments = append(translatedSegments, *translatedSegment)

		percentage := 100 * (float64(len(translatedSegments)) / float64(len(segments)))
		percentage = math.Round(percentage*100) / 100

		d.database.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
			ID:       targetTransformationId,
			Progress: percentage,
		})

		time.Sleep(500 * time.Millisecond)

	}

	return &translatedSegments, nil
}

func (d *Dubbing) fetchDubbedClip(segment Segment, identifier string) error {

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
			d.logger.Error("Request to Eleven Labs failed, retrying after 5 seconds", zap.Int("retries_left", retries), zap.Error(err))
			time.Sleep(5 * time.Second)
		}

	}

	return fmt.Errorf("Failed to call elevenlabs")
}

func (d *Dubbing) dubVideoClip(ctx context.Context, segment Segment, identifier string) error {

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

	generateVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s.mp4' -ss %f -to %f file:'%s'", identifier, start, end, originalVideoSegmentName)
	_, err := d.ffmpeg.Run(ctx, generateVideoClip)

	if err != nil {
		return fmt.Errorf("Clip extraction failed: %s\n%s\n", err.Error(), generateVideoClip)
	}

	stretchVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -vf 'setpts=%f*PTS' -c:a copy file:'%s'", originalVideoSegmentName, ratio, videoSegmentName)
	_, err = d.ffmpeg.Run(ctx, stretchVideoClip)

	if err != nil {
		return fmt.Errorf("Clip stretching failed: %s\n%s\n", err.Error(), generateVideoClip)
	}

	dubVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -i file:'%s' -c:v copy -map 0:v:0 -map 1:a:0 file:'%s'",
		videoSegmentName, audioFileName, dubbedVideoSegmentName)
	_, err = d.ffmpeg.Run(ctx, dubVideoClip)

	if err != nil {
		return fmt.Errorf("Clip dubbing failed: %s\n%s\n", err.Error(), generateVideoClip)
	}

	utils.DeleteFiles([]string{videoSegmentName, originalVideoSegmentName, audioFileName})

	return nil
}

func (d *Dubbing) lipSyncClip(segment Segment, identifier string) error {

	videoSegmentName := getVideoSegmentName(identifier, segment.Id)
	dubbedVideoSegmentName := "dubbed_" + videoSegmentName
	syncedVideoSegmentName := "synced_" + videoSegmentName

	file, err := os.Open(dubbedVideoSegmentName)
	if err != nil {
		d.logger.Error("Error opening file", zap.Error(err))
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		d.logger.Error("Error reading file", zap.Error(err))
	}

	file.Seek(0, io.SeekStart)

	base64Content := base64.StdEncoding.EncodeToString(content)
	dataURL := "data:video/mp4;base64," + base64Content

	replicateRequestBody := map[string]interface{}{
		"version": "8d65e3f4f4298520e079198b493c25adfc43c058ffec924f2aefc8010ed25eef",
		"input": map[string]string{
			"face":  dataURL,
			"audio": dataURL,
		},
	}

	jsonBody, err := json.Marshal(replicateRequestBody)
	outputUrl, err := replicatemiddleware.MakeRequest(bytes.NewBuffer(jsonBody))

	if err != nil {
		d.logger.Error("Replicate Request Failed", zap.Error(err))
		copyFileCmd := fmt.Sprintf("cp %s %s", dubbedVideoSegmentName, syncedVideoSegmentName)
		utils.ExecCommand(copyFileCmd)
		utils.DeleteFiles([]string{dubbedVideoSegmentName})
		return nil
	}

	//download original media, then save it as identifier.mp4
	responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "GET",
		Url:    outputUrl.(string),
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

func (d *Dubbing) concatSegments(ctx context.Context, segments []Segment, identifier string) (string, error) {

	batchSize := 5
	batchFiles := []string{}

	for i := 0; i < len(segments); i += batchSize {
		end := i + batchSize
		if end > len(segments) {
			end = len(segments)
		}

		batch := segments[i:end]
		batchIdentifier := fmt.Sprintf("%s_batch%d_%s", identifier, i/batchSize, uuid.NewString())

		err := d.concatBatchSegments(ctx, batch, batchIdentifier, identifier)
		if err != nil {
			return "", err
		}

		d.logger.Info("Segments concatonated successfully")
		batchFiles = append(batchFiles, batchIdentifier+"_dubbed.mp4")
	}

	// Now we need to concatenate the batch files
	finalOutput, err := d.concatBatchFiles(ctx, batchFiles, identifier, batchSize)
	if err != nil {
		return "", err
	}

	return finalOutput, nil
}

func (d *Dubbing) concatBatchSegments(ctx context.Context, batch []Segment, batchIdentifier string, identifier string) error {
	inputList := []string{}
	filterList := []string{}

	for idx, s := range batch {
		id := s.Id

		videoSegmentName := getVideoSegmentName(identifier, id)
		syncedSegmentName := "synced_" + videoSegmentName
		inputList = append(inputList, fmt.Sprintf("-i file:'%s'", syncedSegmentName))
		filterList = append(filterList, fmt.Sprintf("[%d:v][%d:a]", idx, idx))
	}

	filterList = append(filterList, fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", len(batch)))

	inputArgs := strings.Join(inputList, " ")
	filterComplex := strings.Join(filterList, "")

	ffmpegCmd := fmt.Sprintf("ffmpeg -threads 1 %s -filter_complex '%s' -map '[v]' -map '[a]' -vsync 2 file:'%s_dubbed.mp4'",
		inputArgs, filterComplex, batchIdentifier)

	d.logger.Info("Concatenating segments", zap.String("ffmpeg_command", ffmpegCmd))

	_, err := d.ffmpeg.Run(ctx, ffmpegCmd)

	fileList := []string{}
	for _, s := range batch {
		fileName := "synced_" + getVideoSegmentName(identifier, s.Id)
		fileList = append(fileList, fileName)
	}
	utils.DeleteFiles(fileList)

	if err != nil {
		return fmt.Errorf("Could not concat segments: %s\n%s", err.Error(), ffmpegCmd)
	}

	return nil
}

func (d *Dubbing) concatBatchFiles(ctx context.Context, batchFiles []string, identifier string, batchSize int) (string, error) {
	for len(batchFiles) > 1 {
		newBatchFiles := []string{}

		for i := 0; i < len(batchFiles); i += batchSize {
			end := i + batchSize
			if end > len(batchFiles) {
				end = len(batchFiles)
			}

			batch := batchFiles[i:end]
			batchIdentifier := fmt.Sprintf("%s_finalbatch%d_%s", identifier, i/batchSize, uuid.NewString())

			err := d.concatBatch(ctx, batch, batchIdentifier)
			if err != nil {
				return "", err
			}

			d.logger.Info("Batch files concatenated successfully")
			utils.DeleteFiles(batch)
			newBatchFiles = append(newBatchFiles, batchIdentifier+"_dubbed.mp4")
		}

		batchFiles = newBatchFiles
	}

	// Rename the final batch file to the final output file
	finalOutput := identifier + "_dubbed.mp4"
	err := os.Rename(batchFiles[0], finalOutput)
	if err != nil {
		return "", fmt.Errorf("Could not rename final output file: %s", err.Error())
	}

	return finalOutput, nil
}

func (d *Dubbing) concatBatch(ctx context.Context, batch []string, batchIdentifier string) error {
	inputList := []string{}
	filterList := []string{}

	for idx, fileName := range batch {
		inputList = append(inputList, fmt.Sprintf("-i file:'%s'", fileName))
		filterList = append(filterList, fmt.Sprintf("[%d:v][%d:a]", idx, idx))
	}

	filterList = append(filterList, fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", len(batch)))

	inputArgs := strings.Join(inputList, " ")
	filterComplex := strings.Join(filterList, "")

	ffmpegCmd := fmt.Sprintf("ffmpeg -threads 1 %s -filter_complex '%s' -map '[v]' -map '[a]' -vsync 2 file:'%s_dubbed.mp4'",
		inputArgs, filterComplex, batchIdentifier)

	d.logger.Info("Concatenating batch", zap.String("batch_identifier", batchIdentifier), zap.String("ffmpeg_command", ffmpegCmd))

	_, err := d.ffmpeg.Run(ctx, ffmpegCmd)

	if err != nil {
		return fmt.Errorf("Could not concat batch files: %s\n%s", err.Error(), ffmpegCmd)
	}

	return nil
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

func (d *Dubbing) translateSegment(ctx context.Context, segment Segment, targetLang model.SupportedLanguage) (*Segment, error) {

	retries := 5

	API_KEY := os.Getenv("OPEN_AI_SECRET_KEY")
	URL := "https://api.openai.com/v1/chat/completions"

	for retries > 0 {

		systemPrompt := fmt.Sprintf("Translate the following text to colloquial, everyday spoken %s. Provide the output in %s alphabet. Just give the output.", targetLang, targetLang)
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

		var chatResponse ChatCompletionResponse
		err2 := json.Unmarshal(respBody, &chatResponse)

		if err != nil || err2 != nil || len(chatResponse.Choices) == 0 {

			retries -= 1
			d.logger.Error("Open AI request failed, sleeping for 5s.", zap.Int("retries_left", retries))
			time.Sleep(5 * time.Second)

		} else {

			segment.Text = chatResponse.Choices[0].Message.Content
			return &segment, nil

		}
	}

	return nil, fmt.Errorf("Open AI Requests Failed")
}
