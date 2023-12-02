package dubbing

import (
	"context"
	"fmt"
	"os"
	"planetcastdev/utils"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)


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
