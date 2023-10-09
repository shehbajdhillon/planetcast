package elevenlabsmiddleware

import (
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"os"
)

type ElevenLabs struct {
	logger      *zap.Logger
	semaphore   *semaphore.Weighted
	apiKey      string
	coquiApiKey string
}

type ElevenLabsConnectProps struct {
	Logger *zap.Logger
}

func Connect(args ElevenLabsConnectProps) *ElevenLabs {
	maxWorkers := 5
	sem := semaphore.NewWeighted(int64(maxWorkers))
	apiKey := os.Getenv("ELEVEN_LABS_KEY")
	coquiApiKey := os.Getenv("COQUI_AI_KEY")
	return &ElevenLabs{logger: args.Logger, semaphore: sem, apiKey: apiKey, coquiApiKey: coquiApiKey}
}
