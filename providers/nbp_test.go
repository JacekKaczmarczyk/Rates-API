package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNbpProvider(t *testing.T) {
	provider := NewNbpProvider()

	if provider.Name != "NBP" {
		t.Errorf("Expected Name to be 'NBP', got %s", provider.Name)
	}

	if provider.APIURL != "https://api.nbp.pl/api/exchangerates/tables/a" {
		t.Errorf("Expected APIURL to be 'https://api.nbp.pl/api/exchangerates/tables/a', got %s", provider.APIURL)
	}

	if provider.DateFormat != "2006-01-02" {
		t.Errorf("Expected DateFormat to be '2006-01-02', got %s", provider.DateFormat)
	}
}

func TestNbpProvider_createGetRequest(t *testing.T) {
	provider := NewNbpProvider()

	tests := []struct {
		name        string
		date        string
		expectErr   bool
		expectedUrl string
	}{
		{
			name:        "valid date",
			date:        "2023-01-01",
			expectErr:   false,
			expectedUrl: "https://api.nbp.pl/api/exchangerates/tables/a/2023-01-01",
		},
		{
			name:        "invalid date",
			date:        "invalid",
			expectErr:   true,
			expectedUrl: "",
		},
		{
			name:        "empty date",
			date:        "",
			expectErr:   false,
			expectedUrl: "https://api.nbp.pl/api/exchangerates/tables/a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := provider.createGetRequest(tt.date)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			expectedURL := tt.expectedUrl
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			if req.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			if req.Header.Get("Accept") != "application/json" {
				t.Errorf("Expected Accept header 'application/json', got %s", req.Header.Get("Accept"))
			}
		})
	}
}

func TestNbpProvider_filterRates(t *testing.T) {
	provider := NewNbpProvider()

	rates := []Rate{
		{Currency: "US Dollar", Code: "USD", Mid: 4.0},
		{Currency: "Euro", Code: "EUR", Mid: 4.5},
		{Currency: "British Pound", Code: "GBP", Mid: 5.0},
	}

	tests := []struct {
		name     string
		codes    []string
		expected []RateValue
	}{
		{
			name:  "single code",
			codes: []string{"USD"},
			expected: []RateValue{
				{Code: "USD", Value: 4.0},
			},
		},
		{
			name:  "multiple codes",
			codes: []string{"USD", "EUR"},
			expected: []RateValue{
				{Code: "USD", Value: 4.0},
				{Code: "EUR", Value: 4.5},
			},
		},
		{
			name:     "no matching codes",
			codes:    []string{"JPY"},
			expected: []RateValue{},
		},
		{
			name:     "empty codes",
			codes:    []string{},
			expected: []RateValue{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.filterRates(rates, tt.codes)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d rates, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected rate %+v, got %+v", expected, result[i])
				}
			}
		})
	}
}

func TestNbpProvider_GetCurrencies(t *testing.T) {
	mockResponse := []NbpResponse{
		{
			Table:         "A",
			No:            "001/A/NBP/2023",
			EffectiveDate: "2023-01-01",
			Rates: []Rate{
				{Currency: "US Dollar", Code: "USD", Mid: 4.0},
				{Currency: "Euro", Code: "EUR", Mid: 4.5},
			},
		},
	}

	mockData, err := json.Marshal(mockResponse)
	if err != nil {
		t.Fatalf("Failed to marshal mock response: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockData)
	}))
	defer server.Close()

	provider := &NbpProvider{
		Name:       "NBP",
		APIURL:     server.URL + "/api/exchangerates/tables/a",
		DateFormat: "2006-01-02",
	}

	response, err, statusCode := provider.GetCurrencies([]string{"USD"}, "2023-01-01")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}

	if response.Provider != "NBP" {
		t.Errorf("Expected provider 'NBP', got %s", response.Provider)
	}

	if response.AsOF != "2023-01-01" {
		t.Errorf("Expected AsOF '2023-01-01', got %s", response.AsOF)
	}

	if len(response.Rates) != 1 {
		t.Errorf("Expected 1 rate, got %d", len(response.Rates))
	}

	if response.Rates[0].Code != "USD" || response.Rates[0].Value != 4.0 {
		t.Errorf("Expected rate {USD 4.0}, got %+v", response.Rates[0])
	}
}

func TestNbpProvider_GetCurrencies_InvalidDate(t *testing.T) {
	provider := NewNbpProvider()

	_, err, statusCode := provider.GetCurrencies([]string{"USD"}, "invalid-date")

	if err == nil {
		t.Error("Expected error for invalid date, got nil")
	}

	if statusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", statusCode)
	}
}

func TestNbpProvider_GetCurrencies_NoRatesFound(t *testing.T) {
	mockResponse := []NbpResponse{
		{
			Table:         "A",
			No:            "001/A/NBP/2023",
			EffectiveDate: "2023-01-01",
			Rates: []Rate{
				{Currency: "US Dollar", Code: "USD", Mid: 4.0},
			},
		},
	}

	mockData, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockData)
	}))
	defer server.Close()

	provider := &NbpProvider{
		Name:       "NBP",
		APIURL:     server.URL + "/api/exchangerates/tables/a",
		DateFormat: "2006-01-02",
	}

	_, err, statusCode := provider.GetCurrencies([]string{"EUR"}, "2023-01-01")

	if err == nil {
		t.Error("Expected error for no rates found, got nil")
	}

	if statusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", statusCode)
	}
}
