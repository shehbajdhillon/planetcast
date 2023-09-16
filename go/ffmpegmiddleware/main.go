package ffmpegmiddleware

import (
	"context"
	"fmt"
	"planetcastdev/utils"
	"runtime"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type Ffmpeg struct {
	semaphore *semaphore.Weighted
}

type FfmpegConnectProps struct {
	Logger *zap.Logger
}

func Connect(args FfmpegConnectProps) *Ffmpeg {
	maxWorkers := runtime.GOMAXPROCS(0)
	sem := semaphore.NewWeighted(int64(maxWorkers))
	args.Logger.Info("Setting Up Ffmpeg Middleware", zap.Int("max_semaphore_workers", maxWorkers), zap.Any("semaphore", sem))
	return &Ffmpeg{semaphore: sem}
}

func (f *Ffmpeg) Run(ctx context.Context, ffmpegCmd string) (string, error) {
	if err := f.semaphore.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("Failed to acquire semaphore.")
	}
	defer f.semaphore.Release(1)
	return utils.ExecCommand(ffmpegCmd)
}
