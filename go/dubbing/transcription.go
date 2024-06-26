package dubbing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"planetcastdev/utils"
	"strings"
	"time"

	"go.uber.org/zap"
)

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

func (d *Dubbing) getTranscript(ctx context.Context, fileName string) (*WhisperOutput, error) {

	fileUrl := d.storage.GetFileLink(fileName)

	retries := 5

	var output any

	for retries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - retries)

		replicateRequestBody := map[string]interface{}{
			"input": map[string]interface{}{
				"audio":           fileUrl,
				"model":           "large-v2",
				"word_timestamps": true,
			},
		}
		jsonBody, err := json.Marshal(replicateRequestBody)
		url := "https://api.replicate.com/v1/deployments/shehbajdhillon/whisper-model/predictions"
		output, err = d.replicate.MakeRequest(ctx, bytes.NewBuffer(jsonBody), url)

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

	THRESHOLD_SECONDS := 0.20

	for _, seg := range segments {

		if seg.Start >= seg.End {
			continue
		}
		if len(seg.Text) <= 0 {
			continue
		}

		segmentText := strings.Trim(seg.Text, " ")

		if idx > 0 && (seg.Start-newSegmentArray[idx-1].End <= THRESHOLD_SECONDS) {
			newSegmentArray[idx-1].End = seg.End
			newSegmentArray[idx-1].Text += (" " + segmentText)
		} else {
			newSegmentArray = append(
				newSegmentArray,
				Segment{Id: idx, Start: seg.Start, End: seg.End, Text: segmentText, Words: []Word{}},
			)
			idx += 1
		}
	}

	return newSegmentArray
}

type demucsOutput struct {
	Bass   *string `json:"bass"`
	Drums  *string `json:"drums"`
	Guitar *string `json:"guitar"`
	Other  *string `json:"other"`
	Piano  *string `json:"piano"`
	Vocals *string `json:"vocals"`
}

func (d *Dubbing) runDemucs(ctx context.Context, fileName string) (*demucsOutput, error) {
	fileUrl := d.storage.GetFileLink(fileName)

	retries := 5

	var output any

	for retries > 0 {

		sleepTime := utils.GetExponentialDelaySeconds(5 - retries)

		replicateRequestBody := map[string]interface{}{
			"input": map[string]interface{}{
				"audio": fileUrl,
			},
		}
		jsonBody, err := json.Marshal(replicateRequestBody)
		url := "https://api.replicate.com/v1/deployments/shehbajdhillon/demucs/predictions"
		output, err = d.replicate.MakeRequest(ctx, bytes.NewBuffer(jsonBody), url)

		if err == nil {
			break
		} else {
			retries -= 1
			d.logger.Error("Demucs request failed, retrying after sleeping", zap.Error(err), zap.Int("sleep_time", sleepTime), zap.Int("retries_left", retries))
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

	if retries <= 0 {
		d.logger.Error("Failed to run demucs on input file")
		return nil, fmt.Errorf("Failed to run demucs on input file")
	}

	outputJson, ok := output.(map[string]interface{})

	if !ok {
		d.logger.Error("Could not parse demucs json output")
		return nil, fmt.Errorf("Could not parse demucs json output")
	}

	responseBody, err := json.Marshal(outputJson)
	if err != nil {
		d.logger.Error("Could not parse demucs json body to bytes")
		return nil, fmt.Errorf("Could not parse demucs json body to bytes")
	}

	var demucsOutput demucsOutput
	err = json.Unmarshal(responseBody, &demucsOutput)
	if err != nil {
		d.logger.Error("Could not parse demucs bytes to struct")
		return nil, fmt.Errorf("Could not parse demucs bytes to struct")
	}
	d.logger.Info("Demucs request processes successfully for:", zap.String("fileName", fileName))

	return &demucsOutput, nil
}

func (d *Dubbing) GetTranscriptLength(whisperOutput *WhisperOutput) int {
	segments := whisperOutput.Segments
	length := 0.0
	for _, seg := range segments {
		length += (seg.End - seg.Start)
	}
	length = length / 60
	return int(math.Ceil(length))
}
