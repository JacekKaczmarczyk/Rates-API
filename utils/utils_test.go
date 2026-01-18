package utils

import (
	"testing"
)

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name       string
		date       string
		dateFormat string
		expected   bool
	}{
		{
			name:       "valid date",
			date:       "2023-01-01",
			dateFormat: "2006-01-02",
			expected:   true,
		},
		{
			name:       "invalid date format",
			date:       "01-01-2023",
			dateFormat: "2006-01-02",
			expected:   false,
		},
		{
			name:       "invalid date",
			date:       "2023-13-01",
			dateFormat: "2006-01-02",
			expected:   false,
		},
		{
			name:       "empty date",
			date:       "",
			dateFormat: "2006-01-02",
			expected:   false,
		},
		{
			name:       "different format",
			date:       "01/01/2023",
			dateFormat: "02/01/2006",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateDate(tt.date, tt.dateFormat)
			if result != tt.expected {
				t.Errorf("ValidateDate(%q, %q) = %v; want %v", tt.date, tt.dateFormat, result, tt.expected)
			}
		})
	}
}

func TestValidateCurrencyCodeFormat(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "valid uppercase code",
			code:     "USD",
			expected: true,
		},
		{
			name:     "valid another code",
			code:     "EUR",
			expected: true,
		},
		{
			name:     "lowercase code",
			code:     "usd",
			expected: false,
		},
		{
			name:     "mixed case",
			code:     "UsD",
			expected: false,
		},
		{
			name:     "too short",
			code:     "US",
			expected: false,
		},
		{
			name:     "too long",
			code:     "USDD",
			expected: false,
		},
		{
			name:     "empty string",
			code:     "",
			expected: false,
		},
		{
			name:     "numbers",
			code:     "123",
			expected: false,
		},
		{
			name:     "special characters",
			code:     "U$D",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCurrencyCodeFormat(tt.code)
			if result != tt.expected {
				t.Errorf("ValidateCurrencyCodeFormat(%q) = %v; want %v", tt.code, result, tt.expected)
			}
		})
	}
}
