package utils

import (
	"time"
)

func ValidateDate(date string, dateFormat string) bool {
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return false
	}

	return true
}
