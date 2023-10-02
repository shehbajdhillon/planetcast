package youtubemiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"planetcastdev/ffmpegmiddleware"
	"planetcastdev/utils"

	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"go.uber.org/zap"
)

type YoutubeConnectProps struct {
	Logger *zap.Logger
	Ffmpeg *ffmpegmiddleware.Ffmpeg
}

type Youtube struct {
	logger        *zap.Logger
	youtubeClient *youtube.Client
	download      *downloader.Downloader
	ffmpeg        *ffmpegmiddleware.Ffmpeg
}

func Connect(args YoutubeConnectProps) *Youtube {
	client := youtube.Client{}
	download := downloader.Downloader{Client: client}
	return &Youtube{logger: args.Logger, youtubeClient: &client, download: &download, ffmpeg: args.Ffmpeg}
}

func (y *Youtube) getVideoInfo(videoUrl string) (*youtube.Video, error) {
	video, err := y.youtubeClient.GetVideo(videoUrl)
	if err != nil {
		y.logger.Error("Could not fetch video data", zap.Error(err), zap.String("video_url", videoUrl))
		return nil, err
	}
	y.logger.Info("Fetched video data successfully", zap.String("video_url", videoUrl))
	return video, err
}

func (y *Youtube) downloadVideo(video *youtube.Video) (io.ReadSeeker, error) {
	newCtx := context.Background()
	randomString := uuid.NewString()

	randomFileName := randomString + ".mp4"
	encodedFileName := randomString + "_encoded.mp4"

	err := y.download.DownloadComposite(newCtx, randomFileName, video, "hd1080", "")

	if err != nil {
		y.logger.Error("Could not download youtube video", zap.Error(err), zap.String("video_id", video.ID))
		return nil, err
	}
	y.logger.Info("Downloaded youtube file successfully", zap.String("video_id", video.ID), zap.String("file_name", randomFileName))

	ffmpegCmd := fmt.Sprintf("ffmpeg -i file:'%s' -vcodec libx264 -acodec aac -vsync 2 file:'%s'", randomFileName, encodedFileName)
	y.ffmpeg.Run(newCtx, ffmpegCmd)

	file, err := os.Open(encodedFileName)
	if err != nil {
		y.logger.Error("Could not read downloaded youtube file", zap.Error(err), zap.String("video_id", video.ID), zap.String("file_name", randomFileName))
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		y.logger.Error("Could not read downloaded youtube file", zap.Error(err), zap.String("video_id", video.ID), zap.String("file_name", randomFileName))
		return nil, err
	}

	y.logger.Info("Read downloaded youtube file successfully", zap.String("video_id", video.ID), zap.String("file_name", randomFileName))

	readSeeker := bytes.NewReader(fileContent)

	utils.DeleteFiles([]string{randomFileName, encodedFileName})

	return readSeeker, nil
}

func (y *Youtube) Download(videoUrl string) (io.ReadSeeker, string, error) {
	video, err := y.getVideoInfo(videoUrl)
	if err != nil {
		return nil, "", err
	}
	file, err := y.downloadVideo(video)
	if err != nil {
		return nil, "", err
	}
	return file, video.Title, nil
}
