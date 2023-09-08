package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetCurrentDateTimeString() string {
	currentTime := time.Now()
	timeString := strings.ReplaceAll(currentTime.Format("2006-01-02 15:04:05"), " ", "-")
	return timeString
}

func GetAudioFileDuration(fileName string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"format=duration", "-of", "default=noprint_wrappers=1:nokey=1", fileName)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil

}

func DeleteFiles(fileNames []string) error {
	fileNameString := strings.Join(fileNames, " ")
	removeCmd := fmt.Sprintf("rm -rf %s", fileNameString)
	err := exec.Command("sh", "-c", removeCmd).Run()
	return err
}
