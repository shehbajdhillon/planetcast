package utils

import (
	"bytes"
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

	cmdString := fmt.Sprintf("ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 file:'%s'", fileName)

	output, err := ExecCommand(cmdString)
	if err != nil {
		return 0, fmt.Errorf("Could not run ffprobe to get duration: %s", err.Error())
	}

	text := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not trim space to get duration: %s", err.Error())
	}

	return duration, nil

}

func DeleteFiles(fileNames []string) (string, error) {
	fileNameString := strings.Join(fileNames, " ")
	removeCmd := fmt.Sprintf("rm -rf %s", fileNameString)
	return ExecCommand(removeCmd)
}

func ExecCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	errMsg := stderr.String()

	if err != nil {
		return "", fmt.Errorf("Command failed: %s, %s", err.Error(), errMsg)
	}
	return string(output), nil
}

func MinOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	return min
}

func MaxOf(vars ...int) int {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}
