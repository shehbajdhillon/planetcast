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

	"github.com/tabbed/pqtype"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

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

	demucsPtr, err := d.runDemucs(ctx, args.FileName)
	if err != nil {
		d.logger.Error("Failed to run demucs", zap.Error(err))
		return database.Transformation{}, err
	}
	demucsObj := *demucsPtr

	demucsFile := []*string{
		demucsObj.Bass,
		demucsObj.Drums,
		demucsObj.Guitar,
		demucsObj.Other,
		demucsObj.Piano,
	}

	demucsFileNames := []string{}

	//Download the files, except vocals. Write files to disk.
	for _, filePtr := range demucsFile {
		if filePtr == nil {
			continue
		}
		fileUrl := *filePtr

		responseBody, err := httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
			Method: "GET",
			Url:    fileUrl,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "audio/mp3",
			},
		})

		demucsFileName := fmt.Sprintf("%s-demucs-%d.mp3", args.FileName, len(demucsFileNames))

		if err != nil {
			return database.Transformation{}, err
		}
		err = os.WriteFile(demucsFileName, responseBody, 0644)
		if err != nil {
			return database.Transformation{}, err
		}
		demucsFileNames = append(demucsFileNames, demucsFileName)
	}

	//Mix files together. Upload to S3.
	fileName := fmt.Sprintf("%s-demucs.mp3", args.FileName)
	ffmpegFiles := []string{}
	for _, fileName := range demucsFileNames {
		ffmpegFiles = append(ffmpegFiles, fmt.Sprintf("-i file:'%s'", fileName))
	}
	inputString := strings.Join(ffmpegFiles, " ")
	ffmpegCmd := fmt.Sprintf("ffmpeg %s -filter_complex 'amix=inputs=%d:duration=longest' file:'%s'", inputString, len(demucsFileNames), fileName)
	d.ffmpeg.Run(ctx, ffmpegCmd)

	file, err := os.Open(fileName)
	if err != nil {
		d.logger.Error("Error opening file", zap.Error(err))
	}
	defer file.Close()

	d.storage.Upload(fileName, file)

	//Delete Files from Disk
	utils.DeleteFiles(demucsFileNames)
	utils.DeleteFiles([]string{fileName})

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
	Gender               string
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

	fileUrl = d.storage.GetFileLink(fmt.Sprintf("%s-demucs.mp3", sourceTransformation.TargetMedia))
	responseBody, err = httpmiddleware.HttpRequest(httpmiddleware.HttpRequestStruct{
		Method: "GET",
		Url:    fileUrl,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "audio/mp3",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Error downloading demucs audio file from S3: %s", err.Error())
	}
	err = os.WriteFile(identifier+"-demucs.mp3", responseBody, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error writing demucs audio file: %s", err.Error())
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
		gender:                 args.Gender,
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
	utils.DeleteFiles([]string{identifier + ".mp4", identifier + "_dubbed.mp4", identifier + "-demucs.mp3"})

	// return the update transformation
	return &targetTransformation, nil
}

type fetchAndDubProps struct {
	segments               []Segment
	projectId              int64
	identifier             string
	targetLanguage         string
	targetTransformationId int64
	lipSync                bool
	gender                 string
}

func (d *Dubbing) fetchAndDub(ctx context.Context, args fetchAndDubProps) (*[]Segment, error) {

	translatedSegments := []Segment{}

	d.database.UpdateTransformationStatusById(ctx, database.UpdateTransformationStatusByIdParams{
		ID:     args.targetTransformationId,
		Status: "processing",
	})

	var wg sync.WaitGroup
	mutex := &sync.Mutex{}
	errChan := make(chan error, 1)

	wg.Add(len(args.segments))

	maxWorkers := 4
	sem := semaphore.NewWeighted(int64(maxWorkers))

	frameRate, err := utils.GetVideoFileFrameRate(args.identifier + ".mp4")

	if err != nil {
		frameRate = 0
	}

	for idx := range args.segments {

		sem.Acquire(ctx, 1)

		go func(idx int) {
			defer sem.Release(1)
			defer wg.Done()

			segmentRetires := 6
			var translatedSeg *Segment
			var err error

			for segmentRetires > 0 {
				sleepTime := utils.GetExponentialDelaySeconds(6 - segmentRetires)

				translatedSeg, err = d.processSegment(ctx, idx, frameRate, args)
				if err == nil {
					break
				} else {
					segmentRetires -= 1
					d.logger.Error(
						"Could not process segment, retrying again after sleeping",
						zap.Int("retries_left", segmentRetires),
						zap.Error(err),
						zap.Int("sleep_time", sleepTime),
					)
					time.Sleep(time.Duration(sleepTime) * time.Second)
				}
			}

			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}

			mutex.Lock()

			translatedSegments = append(translatedSegments, *translatedSeg)

			percentage := 100 * (float64(len(translatedSegments)) / float64(len(args.segments)))
			percentage = math.Round(percentage*100) / 100

			d.database.UpdateTransformationProgressById(ctx, database.UpdateTransformationProgressByIdParams{
				ID:       args.targetTransformationId,
				Progress: percentage,
			})

			mutex.Unlock()

		}(idx)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return nil, err
	}

	// Sort the array by the Id of the segment
	sort.Slice(translatedSegments, func(i, j int) bool {
		return translatedSegments[i].Id < translatedSegments[j].Id
	})

	return &translatedSegments, nil
}

func (d *Dubbing) processSegment(ctx context.Context, idx int, frameRate float64, args fetchAndDubProps) (*Segment, error) {

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

	translatedSegment, err := d.translateSegment(ctx, segment, args.targetLanguage)
	if err != nil {
		return nil, fmt.Errorf("Failed to translated segment %d/%d: %s", idx+1, len(segments), err.Error())
	}
	logProgress("Translation Progress")

	videoSegmentName := getVideoSegmentName(identifier, translatedSegment.Id)
	originalAudioSegmentName := videoSegmentName + ".mp3"
	demucsAudioSegmentName := videoSegmentName + "-demucs.mp3"

	generateVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s.mp4' -ss %f -to %f file:'%s'", identifier, translatedSegment.Start, translatedSegment.End, videoSegmentName)
	_, err = d.ffmpeg.Run(ctx, generateVideoClip)

	generateAudioClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -vn -acodec libmp3lame -q:a 4 file:'%s'", videoSegmentName, originalAudioSegmentName)
	_, err = d.ffmpeg.Run(ctx, generateAudioClip)

	generateDemucsAudioClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -ss %f -to %f -vn -acodec libmp3lame -q:a 4 file:'%s'", identifier+"-demucs.mp3", translatedSegment.Start, translatedSegment.End, demucsAudioSegmentName)
	_, err = d.ffmpeg.Run(ctx, generateDemucsAudioClip)

	if err != nil {
		return nil, fmt.Errorf("Clip %d/%d extraction failed: %s\n%s\n", idx+1, len(segments), err.Error(), generateVideoClip)
	}
	logProgress("Clip Extration")

	err = d.fetchDubbedClip(ctx, *translatedSegment, identifier, args.targetLanguage, args.gender)
	if err != nil {
		return nil, fmt.Errorf("Could fetch dubbed clip %d/%d: %s\n", idx+1, len(segments), err.Error())
	}
	logProgress("Audio Generation Progress")

	err = d.dubVideoClip(ctx, *translatedSegment, identifier, frameRate)
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
	identifier := args.identifier

	if end-start <= 0.001 {
		return nil
	}

	videoSegmentName := getVideoSegmentName(args.identifier, args.currentSegment.Id)
	beforeSegmentName := "before_" + videoSegmentName
	syncedVideoSegmentName := "synced_" + videoSegmentName
	demucsAudioSegmentName := "before_" + videoSegmentName + "-demucs.mp3"

	mixedClipFileName := "mixed_" + beforeSegmentName

	//extract the middle part with disabled audio
	generateVideoClipCmd := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s.mp4' -ss %f -to %f -af 'volume=0' file:'%s'", args.identifier, start, end, beforeSegmentName)
	d.logger.Info("INTERIM FFMPEG CMD", zap.String("ffmpeg_cmd", generateVideoClipCmd))
	_, err := d.ffmpeg.Run(ctx, generateVideoClipCmd)
	if err != nil {
		return fmt.Errorf("Could not extract interim segment: %s", err.Error())
	}

	generateDemucsAudioClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -ss %f -to %f -vn -acodec libmp3lame -q:a 4 file:'%s'", identifier+"-demucs.mp3", start, end, demucsAudioSegmentName)
	_, err = d.ffmpeg.Run(ctx, generateDemucsAudioClip)

	dubVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -i file:'%s' -c:v copy -map 0:v:0 -map 1:a:0 file:'%s'",
		beforeSegmentName, demucsAudioSegmentName, mixedClipFileName)
	_, err = d.ffmpeg.Run(ctx, dubVideoClip)

	batch := []string{mixedClipFileName, syncedVideoSegmentName}

	if args.flip != nil && *args.flip == true {
		batch = []string{syncedVideoSegmentName, mixedClipFileName}
	}

	outputFileName, err := d.concatBatchFiles(ctx, batch, args.identifier, 2)

	if err != nil {
		return fmt.Errorf("Could not concat missing info: %s", err.Error())
	}

	err = os.Rename(outputFileName, syncedVideoSegmentName)
	utils.DeleteFiles([]string{beforeSegmentName, demucsAudioSegmentName})

	if err != nil {
		return fmt.Errorf("Could not rename concated missing info output: %s", err.Error())
	}

	return nil
}

func (d *Dubbing) fetchDubbedClip(ctx context.Context, segment Segment, identifier string, language string, gender string) error {

	id := segment.Id
	audioFileName := getAudioFileName(identifier, id)

	videoSegmentName := getVideoSegmentName(identifier, id)

	originalAudioSegmentName := videoSegmentName + ".mp3"

	audioContent, err := d.elevenlabs.ElevenLabsMakeRequest(ctx, elevenlabsmiddleware.ElevenLabsRequestArgs{AudioFileName: originalAudioSegmentName, Text: segment.Text, Gender: gender})

	if err != nil {
		return fmt.Errorf("Error reading response from ElevenLabs: %s", err.Error())
	}
	err = os.WriteFile(audioFileName, audioContent, 0644)
	if err != nil {
		return fmt.Errorf("Error writing audio file: %s", err.Error())
	}
	return nil

}

func (d *Dubbing) dubVideoClip(ctx context.Context, segment Segment, identifier string, frameRate float64) error {

	id := segment.Id

	audioFileName := getAudioFileName(identifier, id)
	stretchAudioFileName := "stretched_" + audioFileName
	mixedAudioFileName := "mixed_" + audioFileName

	videoSegmentName := getVideoSegmentName(identifier, id)
	dubbedVideoSegmentName := "dubbed_" + videoSegmentName

	originalAudioSegmentName := videoSegmentName + ".mp3"
	demucsAudioSegmentName := videoSegmentName + "-demucs.mp3"

	start := segment.Start
	end := segment.End

	videoFileDuration := end - start
	audioFileDuarion, err := utils.GetAudioFileDuration(audioFileName)

	audioStretchRatio := videoFileDuration / audioFileDuarion
	audioStretchRatio = math.Max(1/audioStretchRatio, 0.5)

	if err != nil {
		d.logger.Error("Could not get audio file duration", zap.Error(err))
		audioStretchRatio = 1
	}

	if err != nil {
		return fmt.Errorf("Clip stretching failed: %s\n%s\n", err.Error(), stretchAudioFileName)
	}

	stretchAudioClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -filter:a 'atempo=%f' file:'%s'", audioFileName, audioStretchRatio, stretchAudioFileName)
	_, err = d.ffmpeg.Run(ctx, stretchAudioClip)

	mixAudioClip := fmt.Sprintf(
		"ffmpeg -i file:'%s' -i file:'%s' -filter_complex 'amix=inputs=2:duration=longest' file:'%s'",
		demucsAudioSegmentName,
		stretchAudioFileName,
		mixedAudioFileName,
	)
	_, err = d.ffmpeg.Run(ctx, mixAudioClip)

	dubVideoClip := fmt.Sprintf("ffmpeg -threads 1 -i file:'%s' -i file:'%s' -c:v copy -map 0:v:0 -map 1:a:0 file:'%s'",
		videoSegmentName, mixedAudioFileName, dubbedVideoSegmentName)
	_, err = d.ffmpeg.Run(ctx, dubVideoClip)

	if err != nil {
		return fmt.Errorf("Clip dubbing failed: %s\n%s\n", err.Error(), dubVideoClip)
	}

	utils.DeleteFiles([]string{videoSegmentName, audioFileName, stretchAudioFileName, originalAudioSegmentName, demucsAudioSegmentName, mixedAudioFileName})

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
	url := "https://api.replicate.com/v1/predictions"
	outputUrl, err := d.replicate.MakeRequest(ctx, bytes.NewBuffer(jsonBody), url)
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

func (d *Dubbing) translateSegment(
	ctx context.Context,
	segment Segment,
	targetLang string,
) (*Segment, error) {

	timeTaken := segment.End - segment.Start

	systemPrompt := fmt.Sprintf(
		` You are an expert translator that can translate any text to the %s language.
    You will only use vocabulary that is simple, common and even a new learner to %s language would know.
    You will not use any advanced words, or formal vocabulary.
    You will focus more on clarity and simplicity over complexity of the vocabulary.
    You may simplify the meaning of the sentence first if it means the translation will also use simple, common vocabulary.
    You will translate the input text and will only output the translation.
    Everytime you do a translation, you will first take a deep breath and work on it step-by-step.
  `, targetLang, targetLang)

	userPrompt := fmt.Sprintf(
		` Take a deep breath, and translate the following sentence to %s: '%s'. The original sentence was said in %f seconds, make sure that the translation can also be said in this time.`,
		targetLang, segment.Text, timeTaken)

	retries := 5
	chatGptInput := openaimiddleware.ChatRequestInput{
		Model: "gpt-4",
		Messages: []openaimiddleware.ChatCompletionMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	chatResponse, err := d.openai.MakeAPIRequest(ctx, openaimiddleware.MakeAPIRequestProps{Retries: retries, RequestInput: chatGptInput})
	if err != nil {
		return nil, fmt.Errorf("Open AI Requests Failed: %s", err.Error())
	}

	segment.Text = chatResponse.Choices[0].Message.Content
	return &segment, nil
}
