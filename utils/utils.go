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

func ValidateCurrencyCodeFormat(code string) bool {
	if len(code) != 3 {
		return false
	}
	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
