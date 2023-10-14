package dubbing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"planetcastdev/auth"
	"planetcastdev/database"
	"planetcastdev/elevenlabsmiddleware"
	"planetcastdev/email"
	"planetcastdev/ffmpegmiddleware"
	"planetcastdev/httpmiddleware"
	"planetcastdev/openaimiddleware"
	"planetcastdev/replicatemiddleware"
	"planetcastdev/storage"
	"planetcastdev/utils"
	"sort"
	"strings"
	"sync"
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
	Language string    `json:"detected_language"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Id    int64   `json:"id"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
	Words []Word  `json:"words"`
}

type Word struct {
	End   float64 `json:"end"`
	Start float64 `json:"start"`
	Word  string  `json:"word"`
}

type Dubbing struct {
	storage    *storage.Storage
	database   *database.Queries
	logger     *zap.Logger
	ffmpeg     *ffmpegmiddleware.Ffmpeg
	email      *email.Email
	openai     *openaimiddleware.OpenAI
	replicate  *replicatemiddleware.Replicate
	elevenlabs *elevenlabsmiddleware.ElevenLabs
}

type DubbingConnectProps struct {
	Storage    *storage.Storage
	Database   *database.Queries
	Logger     *zap.Logger
	Ffmpeg     *ffmpegmiddleware.Ffmpeg
	Email      *email.Email
	Openai     *openaimiddleware.OpenAI
	Replicate  *replicatemiddleware.Replicate
	ElevenLabs *elevenlabsmiddleware.ElevenLabs
}

func Connect(args DubbingConnectProps) *Dubbing {
	return &Dubbing{
		storage:    args.Storage,
		database:   args.Database,
		logger:     args.Logger,
		ffmpeg:     args.Ffmpeg,
		email:      args.Email,
		openai:     args.Openai,
		replicate:  args.Replicate,
		elevenlabs: args.ElevenLabs,
	}
}

func (d *Dubbing) getTranscript(ctx context.Context, fileName string) (*WhisperOutput, error) {

	fileUrl := d.storage.GetFileLink(fileName)

	retries := 5

	var output any

	for retries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - retries)

		replicateRequestBody := map[string]interface{}{
			"version": "4a60104c44dd709fc08a03dfeca6c6906257633dd03fd58663ec896a4eeba30e",
			"input": map[string]interface{}{
				"audio":           fileUrl,
				"model":           "large-v2",
				"word_timestamps": true,
			},
		}
		jsonBody, err := json.Marshal(replicateRequestBody)
		output, err = d.replicate.MakeRequest(ctx, bytes.NewBuffer(jsonBody))

		if err == nil {
			break
		} else {
			retries -= 1
			d.logger.Error("Whisper request failed, retrying after sleeping", zap.Error(err), zap.Int("sleep_time", sleepTime), zap.Int("retries_left", retries))
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

	if retries <= 0 {
		d.logger.Error("Failed to transcribe whisper request")
		return nil, fmt.Errorf("Failed to transcribe whisper request")
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

	cleanedSegments := cleanSegments(&whisperOutput)
	whisperOutput.Segments = cleanedSegments

	return &whisperOutput, nil
}

func cleanSegments(whisperOutput *WhisperOutput) []Segment {
	segments := whisperOutput.Segments
	var newSegmentArray []Segment
	var idx int64 = 0

	for _, seg := range segments {
		if seg.Start >= seg.End {
			continue
		}
		if len(seg.Text) <= 0 {
			continue
		}
		newSegmentArray = append(
			newSegmentArray,
			Segment{Id: idx, Start: seg.Start, End: seg.End, Text: seg.Text, Words: []Word{}},
		)
		idx += 1
	}

	return newSegmentArray
}

type CreateTransformationParams struct {
	ProjectID int64
	FileName  string
	IsSource  bool
}

func (d *Dubbing) CreateTransformation(
	ctx context.Context,
	args CreateTransformationParams,
) (database.Transformation, error) {

	transcriptPtr, err := d.getTranscript(ctx, args.FileName)

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
		TargetLanguage: strings.ToUpper(transcriptObj.Language),
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

type CreateTranslationProps struct {
	SourceTransformation database.Transformation
	TargetTransformation database.Transformation
	Identifier           string
	LipSync              bool
}

func (d *Dubbing) CreateTranslation(
	ctx context.Context,
	args CreateTranslationProps,
) (*database.Transformation, error) {

	sourceTransformation := args.SourceTransformation
	identifier := args.Identifier
	targetTransformation := args.TargetTransformation

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
		return nil, fmt.Errorf("Error downloading original audio file from S3: %s", err.Error())
	}
	err = os.WriteFile(identifier+".mp4", responseBody, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error writing audio file: %s", err.Error())
	}

	var whisperOutput WhisperOutput
	json.Unmarshal(sourceTransformation.Transcript.RawMessage, &whisperOutput)

	// call chatgpt, convert the source text to target text
	sourceSegments := whisperOutput.Segments

	projectObj, _ := d.database.GetProjectById(ctx, args.SourceTransformation.ProjectID)
	teamObj, _ := d.database.GetTeamById(ctx, projectObj.TeamID)

	userEmail, err := auth.EmailFromContext(ctx)

	if err != nil {
		d.logger.Error("Could not send transformation start alert email to address", zap.Error(err), zap.Int("transformation_id", int(targetTransformation.ID)))
	} else {
		d.email.DubbingStartAlert(email.DubbingAlertProps{
			TargetLanguage: args.TargetTransformation.TargetLanguage,
			ProjectId:      args.SourceTransformation.ProjectID,
			TeamSlug:       teamObj.Slug,
			UserEmail:      userEmail,
		})
	}

	fetchAndDubArgs := fetchAndDubProps{
		segments:               sourceSegments,
		projectId:              sourceTransformation.ProjectID,
		identifier:             identifier,
		targetLanguage:         targetTransformation.TargetLanguage,
		targetTransformationId: targetTransformation.ID,
		lipSync:                args.LipSync,
	}
	translatedSegmentsPtr, err := d.fetchAndDub(ctx, fetchAndDubArgs)
	if err != nil {
		return nil, fmt.Errorf("Error translating and fetching: %s", err.Error())
	}
	translatedSegments := *translatedSegmentsPtr

	if err != nil {
		return nil, fmt.Errorf("Could not process translated segments " + err.Error())
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

	// get the target text, and parse it
	json.Unmarshal(targetTransformation.Transcript.RawMessage, &whisperOutput)
	whisperOutput.Segments = translatedSegments
	whisperOutput.Language = strings.ToLower(string(targetTransformation.TargetLanguage))

	// store the target text in db
	jsonBytesUntimed, err := json.Marshal(whisperOutput)
	transcriptPtr, err := d.getTranscript(ctx, newFileName)
	var transcriptObj WhisperOutput
	var jsonBytesTimed []byte
	if err == nil {
		transcriptObj = *transcriptPtr
		jsonBytesTimed, err = json.Marshal(transcriptObj)
	}

	if err == nil {
		targetTransformation, err = d.database.UpdateTranscriptById(ctx, database.UpdateTranscriptByIdParams{
			ID:         targetTransformation.ID,
			Transcript: pqtype.NullRawMessage{Valid: true, RawMessage: jsonBytesTimed},
		})
	} else {
		targetTransformation, err = d.database.UpdateTranscriptById(ctx, database.UpdateTranscriptByIdParams{
			ID:         targetTransformation.ID,
			Transcript: pqtype.NullRawMessage{Valid: true, RawMessage: jsonBytesUntimed},
		})
	}

	if err != nil {
		return nil, fmt.Errorf("Could not update transformation: " + err.Error())
	}

	targetTransformation, err = d.database.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     targetTransformation.ID,
		Status: "complete",
	})

	targetTransformation, err = d.database.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
		ID:       targetTransformation.ID,
		Progress: 100,
	})

	if err != nil {
		d.logger.Error("Could not send transformation processed alert email to address", zap.Error(err), zap.Int("transformation_id", int(targetTransformation.ID)))
	} else {
		d.email.DubbingEndedAlert(email.DubbingAlertProps{
			TargetLanguage: args.TargetTransformation.TargetLanguage,
			ProjectId:      args.SourceTransformation.ProjectID,
			TeamSlug:       teamObj.Slug,
			UserEmail:      userEmail,
		})
	}

	//Delete any files written
	utils.DeleteFiles([]string{identifier + ".mp4", identifier + "_dubbed.mp4"})

	// return the update transformation
	return &targetTransformation, nil
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

type fetchAndDubProps struct {
	segments               []Segment
	projectId              int64
	identifier             string
	targetLanguage         string
	targetTransformationId int64
	lipSync                bool
}

func (d *Dubbing) fetchAndDub(ctx context.Context, args fetchAndDubProps) (*[]Segment, error) {

	translatedSegments := []Segment{}

	d.database.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     args.targetTransformationId,
		Status: "processing",
	})

	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	wg.Add(len(args.segments))

	for idx := range args.segments {
		go func(idx int) {
			translatedSeg, _ := d.processSegment(ctx, idx, args)
			mutex.Lock()

			translatedSegments = append(translatedSegments, *translatedSeg)

			percentage := 100 * (float64(len(translatedSegments)) / float64(len(args.segments)))
			percentage = math.Round(percentage*100) / 100

			d.database.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
				ID:       args.targetTransformationId,
				Progress: percentage,
			})

			mutex.Unlock()
			wg.Done()
		}(idx)
	}

	wg.Wait()

	// Sort the array by the Id of the segment
	sort.Slice(translatedSegments, func(i, j int) bool {
		return translatedSegments[i].Id < translatedSegments[j].Id
	})

	return &translatedSegments, nil
}

func (d *Dubbing) processSegment(ctx context.Context, idx int, args fetchAndDubProps) (*Segment, error) {

	segments := args.segments
	identifier := args.identifier
	segment := segments[idx]

	logProgress := func(stage string) {
		d.logger.Info(
			stage,
			zap.Int("project_id", int(args.projectId)),
			zap.Int("transformation_id", int(args.targetTransformationId)),
			zap.String("target_language", string(args.targetLanguage)),
			zap.Int("segments_processed", idx+1),
			zap.Int("total_segments", len(segments)),
		)
	}

	beforeOriginalSegments := segments[utils.MaxOf(0, idx-2):idx]
	afterOriginalSegments := segments[idx+1 : utils.MinOf(idx+3, len(segments))]

	beforeOriginalSentences := []string{}
	for _, seg := range beforeOriginalSegments {
		beforeOriginalSentences = append(beforeOriginalSentences, seg.Text)
	}

	afterOriginalSentences := []string{}
	for _, seg := range afterOriginalSegments {
		afterOriginalSentences = append(afterOriginalSentences, seg.Text)
	}

	translatedSegment, err := d.translateSegment(ctx, segment, args.targetLanguage, beforeOriginalSentences, afterOriginalSentences)
	if err != nil {
		return nil, fmt.Errorf("Failed to translated segment %d/%d: %s", idx+1, len(segments), err.Error())
	}
	logProgress("Translation Progress")

	videoSegmentName := getVideoSegmentName(identifier, translatedSegment.Id)
	originalVideoSegmentName := "original_" + videoSegmentName
	originalAudioSegmentName := originalVideoSegmentName + ".mp3"

	generateVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s.mp4' -ss %f -to %f file:'%s'", identifier, translatedSegment.Start, translatedSegment.End, originalVideoSegmentName)
	_, err = d.ffmpeg.Run(ctx, generateVideoClip)

	generateAudioClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -vn -acodec libmp3lame -q:a 4 file:'%s'", originalVideoSegmentName, originalAudioSegmentName)
	_, err = d.ffmpeg.Run(ctx, generateAudioClip)

	if err != nil {
		return nil, fmt.Errorf("Clip %d/%d extraction failed: %s\n%s\n", idx+1, len(segments), err.Error(), generateVideoClip)
	}
	logProgress("Clip Extration")

	err = d.fetchDubbedClip(ctx, *translatedSegment, identifier, args.targetLanguage)
	if err != nil {
		return nil, fmt.Errorf("Could fetch dubbed clip %d/%d: %s\n", idx+1, len(segments), err.Error())
	}
	logProgress("Audio Generation Progress")

	err = d.dubVideoClip(ctx, *translatedSegment, identifier)
	if err != nil {
		return nil, fmt.Errorf("Could not process clip %d/%d: %s\n", idx+1, len(segments), err.Error())
	}
	logProgress("Dubbing Progress")

	if args.lipSync {
		err = d.lipSyncClip(ctx, *translatedSegment, identifier)
		if err != nil {
			return nil, fmt.Errorf("Could not lip sync clip %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		logProgress("Lip Syncing Progress")
	} else {
		videoSegmentName := getVideoSegmentName(identifier, segment.Id)
		dubbedVideoSegmentName := "dubbed_" + videoSegmentName
		syncedVideoSegmentName := "synced_" + videoSegmentName
		copyFileCmd := fmt.Sprintf("cp %s %s", dubbedVideoSegmentName, syncedVideoSegmentName)
		utils.ExecCommand(copyFileCmd)
		utils.DeleteFiles([]string{dubbedVideoSegmentName})
	}

	var beforeSegment *Segment = nil
	if idx > 0 {
		beforeSegment = &segments[idx-1]
	} else if idx == 0 {
		beforeSegment = &Segment{End: 0}
	}

	err = d.addMissingInfo(ctx, addMissingInfoProps{identifier: identifier, currentSegment: *translatedSegment, beforeSegment: beforeSegment})
	if err != nil {
		return nil, fmt.Errorf("Could not add missing info %d/%d: %s\n", idx+1, len(segments), err.Error())
	}
	logProgress("Added Missing Info")

	videoDuration, err := utils.GetAudioFileDuration(identifier + ".mp4")
	if err == nil && idx == len(segments)-1 && videoDuration-translatedSegment.End <= 0.5 {
		flip := true
		err = d.addMissingInfo(ctx, addMissingInfoProps{identifier: identifier, currentSegment: Segment{Id: translatedSegment.Id, Start: videoDuration - 0.1}, beforeSegment: translatedSegment, flip: &flip})
		if err != nil {
			return nil, fmt.Errorf("Could not add missing info %d/%d: %s\n", idx+1, len(segments), err.Error())
		}
		logProgress("Added Missing Info")
	}

	return translatedSegment, nil
}

type addMissingInfoProps struct {
	currentSegment Segment
	identifier     string
	beforeSegment  *Segment
	flip           *bool
}

func (d *Dubbing) addMissingInfo(ctx context.Context, args addMissingInfoProps) error {

	if args.beforeSegment == nil {
		return nil
	}

	start := args.beforeSegment.End
	end := args.currentSegment.Start

	if end-start <= 0.001 {
		return nil
	}

	videoSegmentName := getVideoSegmentName(args.identifier, args.currentSegment.Id)
	beforeSegmentName := "before_" + videoSegmentName
	syncedVideoSegmentName := "synced_" + videoSegmentName

	//extract the middle part with disabled audio
	generateVideoClipCmd := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s.mp4' -ss %f -to %f -af 'volume=0' file:'%s'", args.identifier, start, end, beforeSegmentName)
	d.logger.Info("INTERIM FFMPEG CMD", zap.String("ffmpeg_cmd", generateVideoClipCmd))
	_, err := d.ffmpeg.Run(ctx, generateVideoClipCmd)
	if err != nil {
		return fmt.Errorf("Could not extract interim segment: %s", err.Error())
	}

	batch := []string{beforeSegmentName, syncedVideoSegmentName}

	if args.flip != nil && *args.flip == true {
		batch = []string{syncedVideoSegmentName, beforeSegmentName}
	}

	outputFileName, err := d.concatBatchFiles(ctx, batch, args.identifier, 2)

	if err != nil {
		return fmt.Errorf("Could not concat missing info: %s", err.Error())
	}

	err = os.Rename(outputFileName, syncedVideoSegmentName)

	if err != nil {
		return fmt.Errorf("Could not rename concated missing info output: %s", err.Error())
	}

	return nil
}

func (d *Dubbing) fetchDubbedClip(ctx context.Context, segment Segment, identifier string, language string) error {

	id := segment.Id
	audioFileName := getAudioFileName(identifier, id)

	videoSegmentName := getVideoSegmentName(identifier, id)
	originalVideoSegmentName := "original_" + videoSegmentName

	originalAudioSegmentName := originalVideoSegmentName + ".mp3"

	audioContent, err := d.elevenlabs.ElevenLabsMakeRequest(ctx, elevenlabsmiddleware.ElevenLabsRequestArgs{AudioFileName: originalAudioSegmentName, Text: segment.Text})

	if err != nil {
		return fmt.Errorf("Error reading response from ElevenLabs: %s", err.Error())
	}
	err = os.WriteFile(audioFileName, audioContent, 0644)
	if err != nil {
		return fmt.Errorf("Error writing audio file: %s", err.Error())
	}
	return nil

}

func (d *Dubbing) dubVideoClip(ctx context.Context, segment Segment, identifier string) error {

	id := segment.Id

	audioFileName := getAudioFileName(identifier, id)
	stretchAudioFileName := "stretched_" + audioFileName

	videoSegmentName := getVideoSegmentName(identifier, id)
	originalVideoSegmentName := "original_" + videoSegmentName
	dubbedVideoSegmentName := "dubbed_" + videoSegmentName

	originalAudioSegmentName := originalVideoSegmentName + ".mp3"

	start := segment.Start
	end := segment.End

	videoFileDuration := end - start
	audioFileDuarion, err := utils.GetAudioFileDuration(audioFileName)

	videoStretchRatio := audioFileDuarion / videoFileDuration
	audioStretchRatio := videoFileDuration / audioFileDuarion

	if videoStretchRatio > 1 {

		averageDuration := (videoFileDuration + audioFileDuarion) / 2
		videoStretchRatio = averageDuration / videoFileDuration
		audioStretchRatio = averageDuration / audioFileDuarion

	} else {
		videoStretchRatio = math.Max(videoStretchRatio, 1.0)
		audioStretchRatio = math.Max(audioStretchRatio, 1.0)
	}

	audioStretchRatio = math.Max(1/audioStretchRatio, 0.5)

	if err != nil {
		d.logger.Error("Could not get audio file duration", zap.Error(err))
		videoStretchRatio = 1
		audioStretchRatio = 1
	}

	stretchVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -vf 'setpts=%f*PTS' file:'%s'", originalVideoSegmentName, videoStretchRatio, videoSegmentName)
	_, err = d.ffmpeg.Run(ctx, stretchVideoClip)

	if err != nil {
		return fmt.Errorf("Clip stretching failed: %s\n%s\n", err.Error(), stretchAudioFileName)
	}

	stretchAudioClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -filter:a 'atempo=%f' file:'%s'", audioFileName, audioStretchRatio, stretchAudioFileName)
	_, err = d.ffmpeg.Run(ctx, stretchAudioClip)

	dubVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -i file:'%s' -c:v copy -map 0:v:0 -map 1:a:0 file:'%s'",
		videoSegmentName, stretchAudioFileName, dubbedVideoSegmentName)
	_, err = d.ffmpeg.Run(ctx, dubVideoClip)

	if err != nil {
		return fmt.Errorf("Clip dubbing failed: %s\n%s\n", err.Error(), dubVideoClip)
	}

	utils.DeleteFiles([]string{videoSegmentName, originalVideoSegmentName, audioFileName, stretchAudioFileName, originalAudioSegmentName})

	return nil
}

func (d *Dubbing) lipSyncClip(ctx context.Context, segment Segment, identifier string) error {

	videoSegmentName := getVideoSegmentName(identifier, segment.Id)
	dubbedVideoSegmentName := "dubbed_" + videoSegmentName
	syncedVideoSegmentName := "synced_" + videoSegmentName

	file, err := os.Open(dubbedVideoSegmentName)
	if err != nil {
		d.logger.Error("Error opening file", zap.Error(err))
	}
	defer file.Close()

	d.storage.Upload(dubbedVideoSegmentName, file)
	fileLink := d.storage.GetFileLink(dubbedVideoSegmentName)

	replicateRequestBody := map[string]interface{}{
		"version": "8d65e3f4f4298520e079198b493c25adfc43c058ffec924f2aefc8010ed25eef",
		"input": map[string]string{
			"face":  fileLink,
			"audio": fileLink,
		},
	}

	jsonBody, err := json.Marshal(replicateRequestBody)
	outputUrl, err := d.replicate.MakeRequest(ctx, bytes.NewBuffer(jsonBody))
	d.storage.DeleteFile(dubbedVideoSegmentName)

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

	err = os.WriteFile(syncedVideoSegmentName, responseBody, 0644)

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

func (d *Dubbing) translateSegment(
	ctx context.Context,
	segment Segment,
	targetLang string,
	beforeTranslatedSentences []string,
	afterOriginalSentences []string,
) (*Segment, error) {

	retries := 5
	prompt := generateTranslationPrompt(string(targetLang), segment.Text, beforeTranslatedSentences, afterOriginalSentences)
	chatGptInput := openaimiddleware.ChatRequestInput{
		Model: "gpt-4",
		Messages: []openaimiddleware.ChatCompletionMessage{
			{Role: "user", Content: prompt},
		},
	}

	chatResponse, err := d.openai.MakeAPIRequest(ctx, openaimiddleware.MakeAPIRequestProps{Retries: retries, RequestInput: chatGptInput})
	if err != nil {
		return nil, fmt.Errorf("Open AI Requests Failed: %s", err.Error())
	}

	segment.Text = chatResponse.Choices[0].Message.Content
	return &segment, nil
}

func generateTranslationPrompt(targetLanguage string, targetSentence string, beforeTranslatedSentences []string, afterOriginalSentences []string) string {

	beforeSentence := ""
	afterSentence := ""

	if len(beforeTranslatedSentences) > 0 {
		beforeSentence = fmt.Sprintf(
			"This sentence will come after the following sentences:\n'%s'\nPlease make sure that the translation flows naturally after these sentences and maintains the original meaning and conversational tone.",
			strings.Join(beforeTranslatedSentences, "\n"),
		)
	}

	if len(afterOriginalSentences) > 0 {
		afterSentence = fmt.Sprintf(
			"For additional context, these are the sentences that will come after the sentence that you will be translating:\n'%s'\n",
			strings.Join(afterOriginalSentences, "\n"),
		)
	}

	disclaimer := fmt.Sprintf("Please use vocabulary that is simple, common and even a learner new to %s language would know, please do not use any advanced words, or formal vocabulary. Focus on clarity and simplicity over complex vocabulary", targetLanguage)

	prompt := fmt.Sprintf(
		"Please simplify the meaning following sentence. Then after you simplify it, translate it to everyday, informal, conversational %s, and provide the output in %s Alphabet:\n'%s'\n%s\n%s\n%s\nAgain the sentence that you are supposed to translate is this:\n'%s'\nProvide the output in %s Alphabet.  Simplify the sentence before any translation. Make sure the translation is everyday, conversational, informal, and understandable by a learner of %s. ONLY PROVIDE THE TRANSLATED OUTPUT AND NOTHING ELSE.  Do not surround the output with any quotation marks.",
		targetLanguage, targetLanguage, targetSentence, disclaimer, beforeSentence, afterSentence, targetSentence, targetLanguage, targetLanguage,
	)
	return prompt
}
