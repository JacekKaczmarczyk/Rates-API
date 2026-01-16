package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		query         string
		expectedCodes []string
		expectedDate  string
	}{
		{
			name:          "no parameters",
			query:         "",
			expectedCodes: defaultCodes,
			expectedDate:  "",
		},
		{
			name:          "with codes",
			query:         "?codes=USD&codes=EUR&codes=GBP",
			expectedCodes: []string{"USD", "EUR", "GBP"},
			expectedDate:  "",
		},
		{
			name:          "with date",
			query:         "?date=2023-01-01",
			expectedCodes: defaultCodes,
			expectedDate:  "2023-01-01",
		},
		{
			name:          "with both codes and date",
			query:         "?codes=USD&date=2023-01-01",
			expectedCodes: []string{"USD"},
			expectedDate:  "2023-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/test"+tt.query, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			codes, date := parseQueryParams(c)

			if len(codes) != len(tt.expectedCodes) {
				t.Errorf("Expected codes %v, got %v", tt.expectedCodes, codes)
				return
			}

			for i, code := range codes {
				if code != tt.expectedCodes[i] {
					t.Errorf("Expected code %s, got %s", tt.expectedCodes[i], code)
				}
			}

			if date != tt.expectedDate {
				t.Errorf("Expected date %s, got %s", tt.expectedDate, date)
			}
		})
	}
}

// Mock provider for testing
// type mockProvider struct{}

// func (m *mockProvider) GetCurrencies(codes []string, date string) (providers.Response, error, int) {
// 	rates := make([]providers.RateValue, len(codes))
// 	for i, code := range codes {
// 		rates[i] = providers.RateValue{Code: code, Value: 1.0}
// 	}

// 	return providers.Response{
// 		AsOF:     date,
// 		Provider: "Mock",
// 		Rates:    rates,
// 	}, nil, http.StatusOK
// }
