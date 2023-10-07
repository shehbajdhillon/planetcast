package ffmpegmiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"planetcastdev/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type Ffmpeg struct {
	semaphore *semaphore.Weighted
	logger    *zap.Logger
}

type FfmpegConnectProps struct {
	Logger *zap.Logger
}

func Connect(args FfmpegConnectProps) *Ffmpeg {
	maxWorkers := 4
	sem := semaphore.NewWeighted(int64(maxWorkers))
	args.Logger.Info("Setting Up Ffmpeg Middleware", zap.Int("max_semaphore_workers", maxWorkers), zap.Any("semaphore", sem))
	return &Ffmpeg{semaphore: sem, logger: args.Logger}
}

func (f *Ffmpeg) Run(ctx context.Context, ffmpegCmd string) (string, error) {
	if err := f.semaphore.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("Failed to acquire semaphore.")
	}
	defer f.semaphore.Release(1)
	return utils.ExecCommand(ffmpegCmd)
}

func (f *Ffmpeg) DownscaleFile(ctx context.Context, fileData io.ReadSeeker) (io.ReadSeeker, error) {

	fileName := uuid.NewString()
	encodedFileName := fileName + "_encoded.mp4"

	body, err := io.ReadAll(fileData)
	if err != nil {
		f.logger.Error("Could not read the file data", zap.Error(err), zap.Any("file_data", fileData))
		return nil, err
	}

	os.WriteFile(fileName, body, 0644)
	fileData.Seek(0, io.SeekStart)

	ffmpegCmd := fmt.Sprintf(`ffmpeg -i file:'%s' -vf 'scale=1280:720' -vcodec libx264 -acodec aac -vsync 2 file:'%s'`, fileName, encodedFileName)
	_, err = f.Run(ctx, ffmpegCmd)
	if err != nil {
		f.logger.Error("Could not execute ffmpeg downscaling command", zap.Error(err), zap.String("file_name", fileName))
		return nil, err
	}
	f.logger.Info("Downscaled file to 720p", zap.String("file_name", fileName))

	fileObj, err := os.Open(encodedFileName)
	if err != nil {
		f.logger.Error("Could not open encoded downscaled file", zap.Error(err), zap.String("file_name", encodedFileName))
		return nil, err
	}
	defer fileObj.Close()

	fileContent, err := io.ReadAll(fileObj)
	if err != nil {
		f.logger.Error("Could not read encoded downscaled file", zap.Error(err), zap.String("file_name", encodedFileName))
		return nil, err
	}

	readSeeker := bytes.NewReader(fileContent)

	utils.DeleteFiles([]string{fileName, encodedFileName})

	return readSeeker, nil
}
