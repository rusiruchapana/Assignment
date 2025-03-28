package utils

import (
	"encoding/json"
	"time"
)

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func PrettyPrint(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
