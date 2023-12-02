package dubbing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
