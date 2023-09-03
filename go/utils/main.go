package utils

import (
	"strings"
	"time"
)

func GetCurrentDateTimeString() string {
	currentTime := time.Now()
	timeString := strings.ReplaceAll(currentTime.Format("2006-01-02 15:04:05"), " ", "-")
	return timeString
}
