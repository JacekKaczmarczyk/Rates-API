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
